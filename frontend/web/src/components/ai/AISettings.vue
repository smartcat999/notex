<template>
  <div class="ai-settings-container">
    <div class="settings-header">
      <h2>AI 模型配置</h2>
      <p>配置您的AI模型API密钥和参数</p>
      <el-button v-if="isDev" @click="debugSettings" size="small" type="info">调试设置</el-button>
    </div>

    <el-tabs v-model="activeTab" tab-position="left" class="settings-tabs">
      <el-tab-pane 
        v-for="provider in aiProviders" 
        :key="provider.id" 
        :label="provider.name" 
        :name="provider.id"
      >
        <div class="provider-settings">
          <h3>{{ provider.name }} 设置</h3>
          
          <el-form :model="providerSettings[provider.id]" label-position="top">
            <el-form-item label="API 密钥">
              <el-input 
                v-model="providerSettings[provider.id].apiKey" 
                placeholder="输入您的API密钥"
                show-password
              />
            </el-form-item>
            
            <el-form-item label="API 端点" v-if="provider.hasEndpoint">
              <el-input 
                v-model="providerSettings[provider.id].endpoint" 
                placeholder="输入API端点URL"
              />
            </el-form-item>
            
            <el-form-item label="可用模型">
              <div class="models-list">
                <div 
                  v-for="model in provider.models" 
                  :key="model.id"
                  class="model-item"
                >
                  <div class="model-header">
                    <el-checkbox 
                      v-model="providerSettings[provider.id].enabledModels[model.id]"
                      :label="model.name"
                    />
                    <el-tag size="small" :type="model.isPaid ? 'danger' : 'success'">
                      {{ model.isPaid ? '付费' : '免费' }}
                    </el-tag>
                  </div>
                  
                  <div class="model-description">
                    {{ model.description }}
                  </div>
                  
                  <div class="model-params" v-if="providerSettings[provider.id].enabledModels[model.id]">
                    <el-form-item label="温度">
                      <el-slider 
                        v-model="providerSettings[provider.id].modelParams[model.id].temperature" 
                        :min="0" 
                        :max="1" 
                        :step="0.1"
                        show-stops
                      />
                    </el-form-item>
                    
                    <el-form-item label="最大输出长度">
                      <el-input-number 
                        v-model="providerSettings[provider.id].modelParams[model.id].maxTokens" 
                        :min="100" 
                        :max="8000"
                        :step="100"
                      />
                    </el-form-item>
                  </div>
                </div>
              </div>
            </el-form-item>
          </el-form>
          
          <div class="provider-actions">
            <el-button type="primary" @click="saveProviderSettings(provider.id)">
              保存设置
            </el-button>
            <el-button @click="testConnection(provider.id)" :disabled="!providerSettings[provider.id].apiKey">
              测试连接
            </el-button>
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>

    <el-card class="settings-card">
      <template #header>
        <div class="card-header">
          <h3>默认模型设置</h3>
        </div>
      </template>
      <div class="settings-form">
        <el-form-item label="默认AI模型">
          <el-select v-model="defaultModel" placeholder="选择默认AI模型" @change="handleDefaultModelChange">
            <el-option
              v-for="model in aiStore.availableModels"
              :key="model.id"
              :label="model.name"
              :value="model.id"
            />
          </el-select>
          <div class="form-tip">设置每次打开AI聊天时默认使用的模型</div>
        </el-form-item>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useAIStore } from '@/stores/ai'
import { useUserStore } from '@/stores/user'
import axios from 'axios'

const aiStore = useAIStore()
const userStore = useUserStore()
const activeTab = ref('openai')
const defaultModel = ref('')
const isLoading = ref(false)
const settingsLoaded = ref(false)
const isDev = import.meta.env.DEV
const isLoadingProviders = ref(false)
const isLoadingUserSettings = ref(false)

// API基础URL
const apiBaseUrl = import.meta.env.VITE_API_BASE_URL || '/api'

// 获取认证头
const getAuthHeaders = () => {
  const token = userStore.token
  return token ? { Authorization: `Bearer ${token}` } : {}
}

// AI提供商列表
const aiProviders = ref([])

// 初始化每个提供商的设置
const providerSettings = reactive({})

// 监听用户登录状态变化
watch(() => userStore.isAuthenticated, async (isAuthenticated) => {
  if (isAuthenticated && !settingsLoaded.value) {
    await loadProvidersAndModels()
    await loadUserSettings()
    settingsLoaded.value = true
  }
})

// 监听 AI 存储的初始化状态
watch(() => aiStore.initialized, async (initialized) => {
  if (initialized) {
    defaultModel.value = aiStore.defaultModel
  }
})

// 加载AI提供商和模型
const loadProvidersAndModels = async () => {
  if (isLoadingProviders.value) {
    return
  }
  
  if (aiStore.modelsLoaded) {
    const providers = {}
    aiStore.availableModels.forEach(model => {
      if (!providers[model.provider]) {
        providers[model.provider] = {
          id: model.provider,
          name: model.provider,
          description: '',
          hasEndpoint: model.provider !== 'openai',
          models: []
        }
      }
      
      providers[model.provider].models.push({
        id: model.id,
        name: model.name,
        description: model.description,
        isPaid: model.isPaid
      })
    })
    
    aiProviders.value = Object.values(providers)
    
    aiProviders.value.forEach(provider => {
      initProviderSetting(provider)
    })
    
    if (aiProviders.value.length > 0) {
      activeTab.value = aiProviders.value[0].id
    }
    
    return
  }
  
  try {
    isLoadingProviders.value = true
    isLoading.value = true
    const response = await axios.get(`${apiBaseUrl}/ai/available-models`)
    
    if (response.data && response.data.providers) {
      aiProviders.value = response.data.providers.map(provider => ({
        id: provider.providerId,
        name: provider.name,
        description: provider.description,
        hasEndpoint: provider.hasEndpoint,
        models: provider.models.map(model => ({
          id: model.modelId,
          name: model.name,
          description: model.description,
          isPaid: model.isPaid
        }))
      }))
      
      aiProviders.value.forEach(provider => {
        initProviderSetting(provider)
      })
      
      if (aiProviders.value.length > 0) {
        activeTab.value = aiProviders.value[0].id
      }
    }
  } catch (error) {
    console.error('加载AI提供商和模型失败:', error)
    ElMessage.error('加载AI提供商和模型失败')
  } finally {
    isLoading.value = false
    isLoadingProviders.value = false
  }
}

// 初始化提供商设置
const initProviderSetting = (provider) => {
  const modelParams = {}
  const enabledModels = {}
  
  provider.models.forEach(model => {
    modelParams[model.id] = {
      temperature: 0.7,
      maxTokens: 2000
    }
    enabledModels[model.id] = false
  })
  
  providerSettings[provider.id] = {
    apiKey: '',
    endpoint: provider.hasEndpoint ? '' : undefined,
    enabledModels,
    modelParams
  }
}

// 加载用户设置
const loadUserSettings = async () => {
  if (!userStore.isAuthenticated) {
    return
  }
  
  if (isLoadingUserSettings.value) {
    return
  }

  try {
    isLoadingUserSettings.value = true
    isLoading.value = true
    
    const response = await axios.get(`${apiBaseUrl}/ai/settings`, {
      headers: getAuthHeaders()
    })
    
    if (response.data && response.data.settings) {
      response.data.settings.forEach(setting => {
        if (providerSettings[setting.providerId]) {
          providerSettings[setting.providerId].apiKey = setting.apiKey || ''
          providerSettings[setting.providerId].endpoint = setting.endpoint || ''
          
          if (setting.enabledModels) {
            Object.keys(setting.enabledModels).forEach(modelId => {
              if (providerSettings[setting.providerId].enabledModels.hasOwnProperty(modelId)) {
                providerSettings[setting.providerId].enabledModels[modelId] = setting.enabledModels[modelId]
              }
            })
          }
          
          if (setting.modelParams) {
            Object.keys(setting.modelParams).forEach(modelId => {
              if (providerSettings[setting.providerId].modelParams.hasOwnProperty(modelId)) {
                providerSettings[setting.providerId].modelParams[modelId] = {
                  ...providerSettings[setting.providerId].modelParams[modelId],
                  ...setting.modelParams[modelId]
                }
              }
            })
          }
        }
      })
    }
    
    const defaultResponse = await axios.get(`${apiBaseUrl}/ai/default-setting`, {
      headers: getAuthHeaders()
    })
    if (defaultResponse.data && defaultResponse.data.defaultModel) {
      defaultModel.value = defaultResponse.data.defaultModel
    }
  } catch (error) {
    console.error('加载用户设置失败:', error)
  } finally {
    isLoading.value = false
    isLoadingUserSettings.value = false
  }
}

// 保存提供商设置
const saveProviderSettings = async (providerId) => {
  try {
    isLoading.value = true
    
    const settingData = {
      providerId: providerId,
      apiKey: providerSettings[providerId].apiKey,
      endpoint: providerSettings[providerId].endpoint,
      enabledModels: providerSettings[providerId].enabledModels,
      modelParams: providerSettings[providerId].modelParams
    }
    
    await axios.post(`${apiBaseUrl}/ai/settings`, settingData, {
      headers: getAuthHeaders()
    })
    
    ElMessage.success(`${getProviderName(providerId)} 设置已保存`)
    
    const defaultModelObj = aiStore.availableModels.find(m => m.id === defaultModel.value)
    if (defaultModelObj && defaultModelObj.provider === providerId) {
      await aiStore.initialize()
    }
  } catch (error) {
    console.error('保存设置失败:', error)
    ElMessage.error(`保存设置失败: ${error.response?.data?.error || error.message}`)
  } finally {
    isLoading.value = false
  }
}

// 测试连接
const testConnection = async (providerId) => {
  try {
    isLoading.value = true
    ElMessage.info(`正在测试 ${getProviderName(providerId)} 连接...`)
    
    const defaultModelForProvider = getDefaultModelForProvider(providerId)
    
    const success = await aiStore.testProviderConnection(
      providerId,
      providerSettings[providerId].apiKey,
      providerSettings[providerId].endpoint,
      defaultModelForProvider
    )
    
    if (success) {
      ElMessage.success(`${getProviderName(providerId)} 连接测试成功`)
    } else {
      ElMessage.error(`${getProviderName(providerId)} 连接测试失败`)
    }
  } catch (error) {
    ElMessage.error(`${getProviderName(providerId)} 连接测试失败: ${error.message}`)
  } finally {
    isLoading.value = false
  }
}

// 获取提供商名称
const getProviderName = (providerId) => {
  const provider = aiProviders.value.find(p => p.id === providerId)
  return provider ? provider.name : providerId
}

// 获取提供商的默认模型
const getDefaultModelForProvider = (providerId) => {
  const provider = aiProviders.value.find(p => p.id === providerId)
  if (!provider || !provider.models || provider.models.length === 0) {
    return null
  }
  
  const enabledModels = Object.entries(providerSettings[providerId].enabledModels)
    .filter(([_, enabled]) => enabled)
    .map(([modelId]) => modelId)
  
  if (enabledModels.length > 0) {
    return enabledModels[0]
  }
  
  return provider.models[0].id
}

const handleDefaultModelChange = async (modelId) => {
  try {
    isLoading.value = true
    
    const defaultData = { defaultModel: modelId }
    
    await axios.post(`${apiBaseUrl}/ai/default-setting`, defaultData, {
      headers: getAuthHeaders()
    })
    
    await aiStore.saveDefaultModel(modelId)
    
    ElMessage.success('默认模型设置已保存')
    
    await aiStore.initialize()
  } catch (error) {
    console.error('保存默认模型失败:', error)
    ElMessage.error(`保存默认模型失败: ${error.response?.data?.error || error.message}`)
  } finally {
    isLoading.value = false
  }
}

// 调试设置
const debugSettings = () => {
  if (isDev) {
    console.error('当前设置状态:', {
      providerSettings,
      defaultModel: defaultModel.value,
      aiStore: {
        initialized: aiStore.initialized,
        defaultModel: aiStore.defaultModel,
        currentModel: aiStore.currentModel,
        currentProvider: aiStore.currentProvider,
        availableModels: aiStore.availableModels
      },
      userAuth: userStore.isAuthenticated,
      settingsLoaded: settingsLoaded.value
    })
  }
}

// 初始化
onMounted(async () => {
  if (settingsLoaded.value) {
    return
  }
  
  await loadProvidersAndModels()
  
  if (userStore.isAuthenticated) {
    await loadUserSettings()
    settingsLoaded.value = true
  }
  
  if (aiStore.initialized) {
    defaultModel.value = aiStore.defaultModel
  }
})
</script>

<style scoped lang="scss">
.ai-settings-container {
  padding: 32px;
  max-width: 1200px;
  margin: 0 auto;
  background-color: #f5f7fa;
  min-height: 100vh;
}

.settings-header {
  margin-bottom: 40px;
  text-align: center;
  
  h2 {
    font-size: 28px;
    font-weight: 600;
    margin-bottom: 16px;
    background: linear-gradient(135deg, #8b5cf6 0%, #6366f1 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    letter-spacing: 0.5px;
  }
  
  p {
    color: #64748b;
    font-size: 16px;
    line-height: 1.6;
  }
}

.settings-tabs {
  background: white;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05);
  border: 1px solid #e2e8f0;
  
  :deep(.el-tabs__header) {
    margin-right: 0;
    background-color: #f8fafc;
    border-right: 1px solid #e2e8f0;
    width: 240px;
  }
  
  :deep(.el-tabs__item) {
    height: 56px;
    line-height: 56px;
    text-align: left;
    padding: 0 24px;
    font-size: 15px;
    color: #64748b;
    transition: all 0.3s ease;
    
    &:hover {
      color: #8b5cf6;
      background-color: rgba(139, 92, 246, 0.04);
    }
    
    &.is-active {
      background-color: rgba(139, 92, 246, 0.08);
      color: #8b5cf6;
      font-weight: 500;
      border-right: 3px solid #8b5cf6;
    }
  }
  
  :deep(.el-tabs__content) {
    padding: 32px;
    background-color: #ffffff;
  }
}

.provider-settings {
  h3 {
    font-size: 22px;
    font-weight: 600;
    margin-bottom: 32px;
    color: #1e293b;
    position: relative;
    padding-bottom: 12px;
    
    &::after {
      content: '';
      position: absolute;
      bottom: 0;
      left: 0;
      width: 48px;
      height: 3px;
      background: linear-gradient(90deg, #8b5cf6 0%, #6366f1 100%);
      border-radius: 2px;
    }
  }
}

.models-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 24px;
  margin-top: 24px;
  
  .model-item {
    border: 1px solid #e2e8f0;
    border-radius: 12px;
    padding: 24px;
    transition: all 0.3s ease;
    background-color: #ffffff;
    
    &:hover {
      border-color: #8b5cf6;
      box-shadow: 0 8px 16px rgba(139, 92, 246, 0.08);
      transform: translateY(-2px);
    }
    
    .model-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 16px;
      
      .el-checkbox {
        font-weight: 500;
        font-size: 15px;
      }
    }
    
    :deep(.el-tag) {
      border-radius: 6px;
      padding: 4px 8px;
      font-size: 12px;
      font-weight: 500;
    }
    
    :deep(.el-tag--success) {
      background-color: rgba(34, 197, 94, 0.1);
      border-color: rgba(34, 197, 94, 0.2);
      color: #16a34a;
    }
    
    :deep(.el-tag--danger) {
      background-color: rgba(239, 68, 68, 0.1);
      border-color: rgba(239, 68, 68, 0.2);
      color: #dc2626;
    }
    
    .model-description {
      color: #64748b;
      font-size: 14px;
      line-height: 1.6;
      margin-bottom: 24px;
    }
    
    .model-params {
      background-color: #f8fafc;
      border-radius: 10px;
      padding: 24px;
      margin-top: 24px;
      border: 1px solid #e2e8f0;
      
      :deep(.el-form-item__label) {
        color: #475569;
        font-weight: 500;
        font-size: 14px;
      }
    }
  }
}

.provider-actions {
  margin-top: 40px;
  display: flex;
  gap: 16px;
  
  .el-button {
    padding: 12px 28px;
    font-size: 15px;
    border-radius: 8px;
    transition: all 0.3s ease;
    font-weight: 500;
    
    &--primary {
      background: linear-gradient(135deg, #8b5cf6 0%, #6366f1 100%);
      border: none;
      
      &:hover {
        transform: translateY(-1px);
        box-shadow: 0 8px 16px rgba(139, 92, 246, 0.16);
      }
    }
    
    &:not(.el-button--primary) {
      border-color: #e2e8f0;
      color: #475569;
      
      &:hover {
        border-color: #8b5cf6;
        color: #8b5cf6;
      }
      
      &[disabled] {
        border-color: #e2e8f0;
        color: #94a3b8;
        background-color: #f8fafc;
      }
    }
  }
}

.settings-card {
  margin-top: 32px;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05);
  border: 1px solid #e2e8f0;
  
  :deep(.el-card__header) {
    background: linear-gradient(135deg, #8b5cf6 0%, #6366f1 100%);
    padding: 20px 24px;
    border-bottom: none;
    
    .card-header {
      h3 {
        margin: 0;
        color: #ffffff;
        font-size: 18px;
        font-weight: 600;
      }
    }
  }
  
  .settings-form {
    padding: 24px;
    
    :deep(.el-form-item) {
      margin-bottom: 28px;
      
      .el-form-item__label {
        color: #1e293b;
        font-weight: 500;
        font-size: 15px;
      }
      
      .form-tip {
        margin-top: 8px;
        font-size: 13px;
        color: #64748b;
      }
    }
    
    :deep(.el-select) {
      width: 100%;
      max-width: 320px;
      
      .el-input__wrapper {
        background-color: #ffffff;
        border: 1px solid #e2e8f0;
        border-radius: 8px;
        padding: 0 12px;
        transition: all 0.3s ease;
        box-shadow: none;
        
        &:hover {
          border-color: #8b5cf6;
        }
        
        &.is-focus {
          border-color: #8b5cf6;
          box-shadow: 0 0 0 3px rgba(139, 92, 246, 0.1);
        }
        
        .el-input__inner {
          color: #1e293b;
          height: 40px;
          line-height: 40px;
          font-size: 14px;
        }
      }
    }
  }
}

:deep(.el-slider) {
  .el-slider__runway {
    height: 6px;
    background-color: #e2e8f0;
    border-radius: 3px;
  }
  
  .el-slider__bar {
    background: linear-gradient(90deg, #8b5cf6 0%, #6366f1 100%);
    border-radius: 3px;
  }
  
  .el-slider__button-wrapper {
    top: -16px;
  }
  
  .el-slider__button {
    width: 16px;
    height: 16px;
    border: 2px solid #8b5cf6;
    background-color: #ffffff;
    transition: transform 0.3s ease;
    
    &:hover {
      transform: scale(1.2);
    }
  }
  
  .el-slider__stop {
    background-color: #e2e8f0;
  }
}

:deep(.el-input-number) {
  .el-input__wrapper {
    background-color: #ffffff;
    border: 1px solid #e2e8f0;
    border-radius: 8px;
    padding: 0 12px;
    transition: all 0.3s ease;
    box-shadow: none;
    
    &:hover {
      border-color: #8b5cf6;
    }
    
    &.is-focus {
      border-color: #8b5cf6;
      box-shadow: 0 0 0 3px rgba(139, 92, 246, 0.1);
    }
    
    .el-input__inner {
      color: #1e293b;
      height: 40px;
      line-height: 40px;
      font-size: 14px;
    }
  }
  
  .el-input-number__decrease,
  .el-input-number__increase {
    border-color: #e2e8f0;
    color: #64748b;
    
    &:hover {
      color: #8b5cf6;
    }
  }
}
</style> 