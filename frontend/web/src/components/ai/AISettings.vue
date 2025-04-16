<template>
  <div class="ai-settings-container">
    <div class="settings-header">
      <h2>AI 模型配置</h2>
      <p>配置您的AI模型API密钥和参数</p>
      <el-button v-if="isDev" @click="debugSettings" size="small" type="info">调试设置</el-button>
    </div>

    <!-- 加载状态 -->
    <div v-if="isLoading" class="loading-container">
      <el-skeleton :rows="5" animated />
    </div>

    <!-- 内容区域 -->
    <div v-else class="settings-main">
      <!-- 左侧提供商列表 -->
      <div class="provider-list">
        <h3 class="list-title">AI 提供商</h3>
        <ul class="provider-menu">
          <li 
            v-for="provider in aiProviders" 
            :key="provider.providerId"
            :class="{ active: activeTab === provider.providerId }"
            @click="activeTab = provider.providerId"
          >
            <div class="provider-icon" :class="provider.providerId">
              <span>{{ provider.name.substring(0, 2) }}</span>
            </div>
            <span class="provider-name">{{ provider.name }}</span>
          </li>
        </ul>
      </div>

      <!-- 右侧内容区域 -->
      <div class="settings-content">
        <div v-for="provider in aiProviders" 
             :key="provider.providerId" 
             v-show="activeTab === provider.providerId"
             class="provider-settings"
        >
          <!-- API设置卡片 -->
          <div class="settings-card api-settings">
            <div class="card-header">
              <h3>{{ provider.name }} API 设置</h3>
              <p class="card-description">配置您的API密钥和服务端点以使用{{ provider.name }}的AI模型</p>
            </div>
            
            <div class="card-content">
              <el-form :model="providerSettings[provider.providerId]" label-position="top">
                <el-form-item label="API 密钥">
                  <div class="api-key-input">
                    <el-input 
                      v-model="providerSettings[provider.providerId].apiKey" 
                      placeholder="输入您的API密钥"
                      :type="passwordVisible[provider.providerId] ? 'text' : 'password'"
                    />
                    <el-button 
                      @click="togglePasswordVisibility(provider.providerId)" 
                      class="visibility-toggle"
                    >
                      {{ passwordVisible[provider.providerId] ? '隐藏' : '显示' }}
                    </el-button>
                  </div>
                </el-form-item>
                
                <el-form-item label="API 端点" v-if="provider.hasEndpoint">
                  <el-input 
                    v-model="providerSettings[provider.providerId].endpoint" 
                    placeholder="输入API端点URL"
                  />
                </el-form-item>

                <div class="form-actions">
                  <el-button type="primary" @click="saveSettings" :loading="isSubmitting">
                    保存设置
                  </el-button>
                  <el-button @click="testModelConnection(provider.providerId)" :loading="isSubmitting">
                    测试连接
                  </el-button>
                </div>
              </el-form>
            </div>
          </div>

          <!-- 模型选择与配置卡片 -->
          <div class="settings-card models-settings">
            <div class="card-header">
              <h3>模型管理</h3>
              <p class="card-description">启用或禁用模型，设置默认模型和调整参数</p>
            </div>
            
            <div class="card-content">
              <!-- 模型类型选择器 -->
              <div class="model-tabs">
                <el-radio-group v-model="modelTypeFilter" size="large" button-style="solid">
                  <el-radio-button label="text">文本模型</el-radio-button>
                  <el-radio-button label="image" v-if="getProviderImageModels(provider.providerId).length > 0">图像模型</el-radio-button>
                </el-radio-group>
              </div>
              
              <!-- 文本模型列表 -->
              <div v-if="modelTypeFilter === 'text'" class="model-list text-models">
                <div 
                  v-for="model in getProviderTextModels(provider.providerId)" 
                  :key="model.modelId" 
                  class="model-card"
                  :class="{ 'selected': defaultModel === model.modelId, 'disabled': !providerSettings[provider.providerId].enabledModels[model.modelId] }"
                >
                  <div class="model-header">
                    <div class="model-name-section">
                      <h4 class="model-name">{{ model.name }}</h4>
                      <el-tag size="small" v-if="model.isPaid" type="danger">付费</el-tag>
                    </div>
                    <el-switch 
                      v-model="providerSettings[provider.providerId].enabledModels[model.modelId]"
                      inline-prompt
                      active-text="启用"
                      inactive-text="禁用"
                    />
                  </div>
                  
                  <div class="model-description">
                    {{ model.description || '无描述' }}
                  </div>
                  
                  <div v-if="providerSettings[provider.providerId].enabledModels[model.modelId]" class="model-settings">
                    <div class="param-item">
                      <div class="param-header">
                        <span class="param-label">温度：</span>
                        <span class="param-value">{{ providerSettings[provider.providerId].modelParams[model.modelId].temperature }}</span>
                      </div>
                      <el-slider 
                        v-model="providerSettings[provider.providerId].modelParams[model.modelId].temperature" 
                        :min="0" 
                        :max="1" 
                        :step="0.1"
                        show-stops
                      />
                    </div>
                    <div class="param-item">
                      <div class="param-header">
                        <span class="param-label">最大输出长度：</span>
                      </div>
                      <el-input-number 
                        v-model="providerSettings[provider.providerId].modelParams[model.modelId].maxTokens" 
                        :min="100" 
                        :max="8000"
                        :step="100"
                        size="small"
                      />
                    </div>
                  </div>
                  
                  <div class="model-actions">
                    <el-button 
                      type="primary" 
                      :disabled="!providerSettings[provider.providerId].enabledModels[model.modelId]"
                      @click="handleModelSelection(model.modelId)"
                      :class="{ 'is-default': defaultModel === model.modelId }"
                    >
                      {{ defaultModel === model.modelId ? '当前默认' : '设为默认' }}
                    </el-button>
                  </div>
                </div>
              </div>
              
              <!-- 图像模型列表 -->
              <div v-if="modelTypeFilter === 'image'" class="model-list image-models">
                <div 
                  v-for="model in getProviderImageModels(provider.providerId)" 
                  :key="model.modelId" 
                  class="model-card"
                  :class="{ 'selected': imageDefaultModel === model.modelId, 'disabled': !providerSettings[provider.providerId].enabledModels[model.modelId] }"
                >
                  <div class="model-header">
                    <div class="model-name-section">
                      <h4 class="model-name">{{ model.name }}</h4>
                      <el-tag size="small" v-if="model.isPaid" type="danger">付费</el-tag>
                    </div>
                    <el-switch 
                      v-model="providerSettings[provider.providerId].enabledModels[model.modelId]"
                      inline-prompt
                      active-text="启用"
                      inactive-text="禁用"
                    />
                  </div>
                  
                  <div class="model-description">
                    {{ model.description || '无描述' }}
                  </div>
                  
                  <div class="model-actions">
                    <el-button 
                      type="primary" 
                      :disabled="!providerSettings[provider.providerId].enabledModels[model.modelId]"
                      @click="handleImageModelSelection(model.modelId)"
                      :class="{ 'is-default': imageDefaultModel === model.modelId }"
                    >
                      {{ imageDefaultModel === model.modelId ? '当前默认' : '设为默认' }}
                    </el-button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch, computed, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { useAIStore } from '@/stores/ai'
import { useUserStore } from '@/stores/user'
import axios from 'axios'
import { getProviderEndpoint } from '@/services/aiService'

const aiStore = useAIStore()
const userStore = useUserStore()
const apiBaseUrl = import.meta.env.VITE_API_BASE_URL || '/api'
const isDev = import.meta.env.DEV

// 状态
const activeTab = ref('openai')
const isLoading = ref(false)
const isSubmitting = ref(false)
const settingsLoaded = ref(false)
const passwordVisible = reactive({})
const defaultModel = ref('')
const imageDefaultModel = ref('')
const modelTypeFilter = ref('text')

// 提供商数据
const aiProviders = ref([])

// 提供商设置
const providerSettings = reactive({})

// 获取认证头
const getAuthHeaders = () => {
  const token = userStore.token
  return token ? { Authorization: `Bearer ${token}` } : {}
}

// 获取分类后的模型
const textModels = computed(() => {
  return aiStore.textModels;
});

const imageModels = computed(() => {
  return aiStore.imageModels;
});

// 处理模型选择
const handleModelSelection = async (modelId) => {
  try {
    if (!userStore.isAuthenticated) {
      ElMessage.warning('请先登录');
      return;
    }
    
    // 更新前端状态
    defaultModel.value = modelId;
    aiStore.setCurrentModel(modelId);
    
    // 保存到后端
    const result = await saveDefaultModelSetting(modelId);
    if (result) {
      ElMessage.success('已设置为默认文本模型');
    }
  } catch (error) {
    console.error('设置默认模型失败:', error);
    defaultModel.value = aiStore.currentModel; // 恢复原值
    ElMessage.error(`设置默认模型失败: ${error.response?.data?.error || error.message}`);
  }
};

// 处理图像模型选择
const handleImageModelSelection = async (modelId) => {
  try {
    if (!userStore.isAuthenticated) {
      ElMessage.warning('请先登录');
      return;
    }
    
    // 更新本地状态
    imageDefaultModel.value = modelId;
    aiStore.setCurrentImageModel(modelId);
    
    // 保存到后端
    const result = await saveDefaultModelSetting(null, modelId);
    if (result) {
      ElMessage.success('已设置为默认图像模型');
    }
  } catch (error) {
    console.error('设置默认图像模型失败:', error);
    imageDefaultModel.value = aiStore.currentImageModel; // 恢复原值
    ElMessage.error(`设置默认图像模型失败: ${error.response?.data?.error || error.message}`);
  }
};

// 保存默认模型设置
const saveDefaultModelSetting = async (textModelId = null, imageModelId = null) => {
  try {
    if (!userStore.isAuthenticated) {
      return false;
    }
    
    // 准备请求数据
    const requestData = {};
    
    // 只包含有变化的字段
    if (textModelId) {
      requestData.defaultModel = textModelId;
    }
    
    if (imageModelId) {
      requestData.defaultImageModel = imageModelId;
    }
    
    // 发送请求
    await axios.post(`${apiBaseUrl}/ai/default-setting`, requestData, {
      headers: getAuthHeaders()
    });
    
    return true;
  } catch (error) {
    console.error('保存默认模型设置失败:', error);
    
    // 处理常见错误情况
    if (error.response?.status === 404) {
      ElMessage.error('API端点未找到，请确认后端服务是否正确配置');
    } else if (error.response?.status === 401) {
      ElMessage.error('未授权访问，请确认您已登录');
    } else if (error.response?.status === 400) {
      const errorMsg = error.response?.data?.error || '未知错误';
      ElMessage.error(`请求参数错误: ${errorMsg}`);
    } else {
      ElMessage.error(`保存默认模型设置失败: ${error.response?.data?.error || error.message}`);
    }
    
    throw error;
  }
};

// 获取行的类名（用于高亮当前选中模型）
const getRowClassName = ({ row }) => {
  // 根据模型类型确定使用哪个当前选中的模型ID
  const currentActiveModel = row.type === 'image' 
    ? aiStore.currentImageModel 
    : aiStore.currentModel;
  
  return row.id === currentActiveModel ? 'selected-row' : '';
};

// 获取提供商名称
const getProviderName = (providerId) => {
  const provider = aiProviders.value.find(p => p.providerId === providerId);
  return provider ? provider.name : providerId;
};

// 获取特定提供商的文本模型
const getProviderTextModels = (providerId) => {
  const provider = aiProviders.value.find(p => p.providerId === providerId);
  if (!provider || !provider.models) {
    return [];
  }
  
  return provider.models.filter(model => model.type === 'text' || !model.type);
};

// 获取特定提供商的图像模型
const getProviderImageModels = (providerId) => {
  const provider = aiProviders.value.find(p => p.providerId === providerId);
  if (!provider || !provider.models) {
    return [];
  }
  
  return provider.models.filter(model => model.type === 'image');
};

// 加载提供商和模型数据
const loadProvidersAndModels = async () => {
  try {
    isLoading.value = true
    
    // 获取提供商和模型
    const response = await axios.get(`${apiBaseUrl}/ai/available-models`, {
      headers: getAuthHeaders()
    })
    
    if (!response.data || !response.data.providers) {
      ElMessage.error('加载提供商和模型失败')
      return
    }
    
    // 更新提供商
    aiProviders.value = response.data.providers
    
    // 预先初始化所有提供商的设置对象（重要：确保响应式）
    aiProviders.value.forEach(provider => {
      if (!providerSettings[provider.providerId]) {
        // 创建提供商设置，确保对象是响应式的
        providerSettings[provider.providerId] = {
          apiKey: '',
          endpoint: '',
          enabledModels: {},
          modelParams: {}
        };
      }
      
      // 初始化密码可见性
      if (passwordVisible[provider.providerId] === undefined) {
        passwordVisible[provider.providerId] = false;
      }
    });
    
    // 然后为每个提供商初始化模型设置
    aiProviders.value.forEach(provider => {
      // 为每个模型初始化设置
      if (provider.models) {
        provider.models.forEach(model => {
          const modelId = model.modelId;
          
          // 确保每个模型都有启用状态设置
          if (providerSettings[provider.providerId].enabledModels[modelId] === undefined) {
            providerSettings[provider.providerId].enabledModels[modelId] = true;
          }
          
          // 确保每个模型都有参数设置
          if (!providerSettings[provider.providerId].modelParams[modelId]) {
            providerSettings[provider.providerId].modelParams[modelId] = {
              temperature: 0.7,
              maxTokens: 2000
            };
          }
        });
      }
    });
    
    // 如果提供商列表不为空且未设置活动标签，设置第一个为活动标签
    if (aiProviders.value.length > 0 && !activeTab.value) {
      activeTab.value = aiProviders.value[0].providerId;
    }
  } catch (error) {
    console.error('加载提供商和模型失败:', error)
    ElMessage.error('加载提供商和模型失败')
  } finally {
    isLoading.value = false
  }
}

// 加载用户设置
const loadUserSettings = async () => {
  try {
    // 获取用户设置
    const response = await axios.get(`${apiBaseUrl}/ai/settings`, {
      headers: getAuthHeaders()
    });
    
    if (!response.data || !response.data.settings) {
      return;
    }
    
    // 遍历用户设置并创建临时对象
    const tempSettings = {};
    
    response.data.settings.forEach(setting => {
      if (setting && setting.providerId) {
        const providerId = setting.providerId;
        
        // 创建或更新该提供商的设置
        tempSettings[providerId] = {
          apiKey: setting.apiKey || '',
          endpoint: setting.endpoint || '',
          enabledModels: { ...setting.enabledModels || {} },
          modelParams: { ...setting.modelParams || {} }
        };
      }
    });
    
    // 合并已有设置和新设置
    Object.keys(tempSettings).forEach(providerId => {
      // 创建一个完整的设置对象
      if (!providerSettings[providerId]) {
        providerSettings[providerId] = {
          apiKey: '',
          endpoint: '',
          enabledModels: {},
          modelParams: {}
        };
      }
      
      // 更新设置
      providerSettings[providerId].apiKey = tempSettings[providerId].apiKey;
      providerSettings[providerId].endpoint = tempSettings[providerId].endpoint;
      
      // 更新模型启用状态和参数
      Object.assign(providerSettings[providerId].enabledModels, tempSettings[providerId].enabledModels);
      Object.assign(providerSettings[providerId].modelParams, tempSettings[providerId].modelParams);
    });
    
    // 确保所有模型的设置都存在
    aiProviders.value.forEach(provider => {
      if (provider.models && providerSettings[provider.providerId]) {
        provider.models.forEach(model => {
          const modelId = model.modelId;
          
          // 设置模型默认值（如果尚未设置）
          if (providerSettings[provider.providerId].enabledModels[modelId] === undefined) {
            providerSettings[provider.providerId].enabledModels[modelId] = false;
          }
          
          if (!providerSettings[provider.providerId].modelParams[modelId]) {
            providerSettings[provider.providerId].modelParams[modelId] = {
              temperature: 0.7,
              maxTokens: 2000
            };
          }
        });
      }
    });
    
    // 获取默认设置
    try {
      const defaultResponse = await axios.get(`${apiBaseUrl}/ai/default-setting`, {
        headers: getAuthHeaders()
      });
      
      if (defaultResponse.data) {
        if (defaultResponse.data.defaultModel) {
          defaultModel.value = defaultResponse.data.defaultModel;
        }
        
        if (defaultResponse.data.defaultImageModel) {
          imageDefaultModel.value = defaultResponse.data.defaultImageModel;
        }
      }
    } catch (error) {
      console.error('获取默认设置失败:', error);
    }
    
    // 触发视图更新
    setTimeout(() => {
      aiProviders.value = [...aiProviders.value];
    }, 50);
  } catch (error) {
    console.error('加载用户设置失败:', error);
    ElMessage.error('加载用户设置失败');
  }
}

// 切换密码可见性
const togglePasswordVisibility = (providerId) => {
  console.log('Toggle password visibility for provider:', providerId)
  // 使用双重取反确保是布尔值
  passwordVisible[providerId] = !!(passwordVisible[providerId]) ? false : true
  console.log('New visibility state:', passwordVisible[providerId])
}

// 测试连接
const testModelConnection = async (providerId) => {
  try {
    const settings = providerSettings[providerId]
    
    if (!settings?.apiKey) {
      ElMessage.warning('请先输入API密钥')
      return
    }
    
    isSubmitting.value = true
    
    // 找到第一个启用的模型
    let modelId = null
    for (const [id, enabled] of Object.entries(settings.enabledModels)) {
      if (enabled) {
        modelId = id
        break
      }
    }
    
    if (!modelId) {
      ElMessage.warning('请至少启用一个模型')
      return
    }
    
    // 调用API测试连接
    const testResult = await aiStore.testProviderConnection(
      providerId, 
      settings.apiKey.trim(), 
      settings.endpoint?.trim() || '', 
      modelId
    )
    
    if (testResult === true) {
      ElMessage.success('连接测试成功')
    } else {
      ElMessage.error('连接测试失败')
    }
  } catch (error) {
    console.error('连接测试失败:', error)
    ElMessage.error(`连接测试失败: ${error.message}`)
  } finally {
    isSubmitting.value = false
  }
}

// 保存设置
const saveSettings = async () => {
  try {
    isSubmitting.value = true
    
    // 检查用户是否已登录
    if (!userStore.isAuthenticated) {
      ElMessage.warning('请先登录')
      return
    }
    
    // 获取当前激活的提供商
    const providerId = activeTab.value
    const settings = providerSettings[providerId]
    
    if (!settings?.apiKey) {
      ElMessage.warning('请先输入API密钥')
      return
    }
    
    // 创建请求数据
    const requestData = {
      providerId: providerId,
      apiKey: settings.apiKey.trim(),
      endpoint: settings.endpoint?.trim() || '',
      enabledModels: settings.enabledModels,
      modelParams: settings.modelParams
    }
    
    // 发送请求保存设置
    try {
      const response = await axios.post(`${apiBaseUrl}/ai/settings`, requestData, {
        headers: getAuthHeaders()
      });
    } catch (error) {
      // 处理常见错误情况
      if (error.response?.status === 404) {
        ElMessage.error('API端点未找到，请确认后端服务是否正确配置');
        return;
      } else if (error.response?.status === 401) {
        ElMessage.error('未授权访问，请确认您已登录');
        return;
      } else if (error.response?.status === 400) {
        ElMessage.error(`请求参数错误: ${error.response?.data?.error || '未知错误'}`);
        return;
      }
      
      // 重新抛出错误，让外部的catch块处理
      throw error;
    }
    
    // 保存设置到localStorage
    const allSettings = JSON.parse(localStorage.getItem('aiProviderSettings') || '{}')
    allSettings[providerId] = {
      apiKey: settings.apiKey,
      endpoint: settings.endpoint,
      modelParams: settings.modelParams
    }
    localStorage.setItem('aiProviderSettings', JSON.stringify(allSettings))
    
    // 更新AI Store状态
    if (aiStore.initialized) {
      await aiStore.initialize()
    }
    
    ElMessage.success('设置已保存')
  } catch (error) {
    console.error('保存设置失败:', error)
    ElMessage.error(`保存设置失败: ${error.response?.data?.error || error.message}`)
  } finally {
    isSubmitting.value = false
  }
}

// 调试设置
const debugSettings = () => {
  if (isDev) {
    console.error('当前设置状态:', {
      providerSettings,
      defaultModel: defaultModel.value,
      imageDefaultModel: imageDefaultModel.value,
      aiStore: {
        initialized: aiStore.initialized,
        defaultModel: aiStore.defaultModel,
        currentModel: aiStore.currentModel,
        currentImageModel: aiStore.currentImageModel,
        currentProvider: aiStore.currentProvider,
        availableModels: aiStore.availableModels,
        textModels: aiStore.textModels,
        imageModels: aiStore.imageModels
      },
      userAuth: userStore.isAuthenticated,
      settingsLoaded: settingsLoaded.value
    })
  }
}

// 监听活动标签的变化
watch(activeTab, (newTab) => {
  if (isDev) {
    // console.log('切换到提供商:', newTab);
    // console.log('该提供商设置:', providerSettings[newTab]);
  }
});

// 初始化
onMounted(async () => {
  try {
    if (settingsLoaded.value) {
      return;
    }
    
    isLoading.value = true;
    
    // 先加载提供商和模型数据
    await loadProvidersAndModels();
    
    // 然后加载用户设置
    if (userStore.isAuthenticated) {
      await loadUserSettings();
    }
    
    // 如果已经初始化AI Store，加载默认模型
    if (aiStore.initialized) {
      defaultModel.value = aiStore.defaultModel || '';
      imageDefaultModel.value = aiStore.currentImageModel || '';
    }
    
    settingsLoaded.value = true;
    
    // 使用nextTick确保DOM更新
    await nextTick();
  } catch (error) {
    console.error('初始化AI设置失败:', error);
    ElMessage.error('初始化AI设置失败');
  } finally {
    isLoading.value = false;
  }
});
</script>

<style scoped lang="scss">
.ai-settings-container {
  padding: 32px;
  max-width: 1200px;
  margin: 0 auto;
  background-color: #f8fafc;
  min-height: calc(100vh - 64px);
  color: #1e293b;
}

.settings-header {
  margin-bottom: 40px;
  text-align: center;
  
  h2 {
    font-size: 28px;
    font-weight: 600;
    margin-bottom: 16px;
    background: linear-gradient(135deg, #0f172a 0%, #334155 100%);
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

.loading-container {
  padding: 32px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
  border: 1px solid #e2e8f0;
  margin-bottom: 24px;
}

.settings-main {
  display: flex;
  gap: 24px;
  min-height: 700px;
  
  .provider-list {
    width: 220px;
    flex-shrink: 0;
    background: white;
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
    border: 1px solid #e2e8f0;
    
    .list-title {
      padding: 20px;
      font-size: 18px;
      font-weight: 600;
      color: #0f172a;
      border-bottom: 1px solid #e2e8f0;
      margin: 0;
    }
    
    .provider-menu {
      list-style-type: none;
      padding: 0;
      margin: 0;
      
      li {
        display: flex;
        align-items: center;
        padding: 16px 20px;
        cursor: pointer;
        transition: all 0.2s ease;
        gap: 12px;
        color: #64748b;
        
        &:hover {
          background-color: #f8fafc;
          color: #0f172a;
        }
        
        &.active {
          background-color: #f1f5f9;
          color: #0f172a;
          position: relative;
          font-weight: 500;
          
          &::after {
            content: '';
            position: absolute;
            right: 0;
            top: 50%;
            transform: translateY(-50%);
            height: 60%;
            width: 3px;
            background: #0f172a;
            border-radius: 2px 0 0 2px;
          }
        }
        
        .provider-icon {
          width: 36px;
          height: 36px;
          border-radius: 8px;
          display: flex;
          align-items: center;
          justify-content: center;
          font-size: 14px;
          font-weight: 600;
          color: white;
          
          &.openai {
            background: linear-gradient(135deg, #10a37f, #0a8a6c);
          }
          
          &.anthropic {
            background: linear-gradient(135deg, #b73999, #8f2d79);
          }
          
          &.google {
            background: linear-gradient(135deg, #4285f4, #34a853);
          }
          
          &.deepseek {
            background: linear-gradient(135deg, #ff5c35, #d64c2d);
          }
          
          &.custom {
            background: linear-gradient(135deg, #334155, #1e293b);
          }
        }
        
        .provider-name {
          font-size: 15px;
        }
      }
    }
  }
  
  .settings-content {
    flex: 1;
    overflow: hidden;
    min-width: 0;
    
    .provider-settings {
      display: flex;
      flex-direction: column;
      gap: 24px;
    }
  }
}

.settings-card {
  background: white;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.08);
  border: 1px solid #e2e8f0;
  
  .card-header {
    padding: 24px 28px;
    border-bottom: 1px solid #e2e8f0;
    background-color: #fafafa;
    
    h3 {
      margin: 0 0 8px 0;
      font-size: 18px;
      font-weight: 600;
      color: #0f172a;
    }
    
    .card-description {
      margin: 0;
      color: #64748b;
      font-size: 14px;
    }
  }
  
  .card-content {
    padding: 28px;
  }
}

.form-actions {
  display: flex;
  gap: 12px;
  margin-top: 24px;
}

.api-key-input {
  display: flex;
  align-items: stretch;
  gap: 8px;
  
  .el-input {
    flex-grow: 1;
    
    :deep(.el-input__wrapper) {
      border-radius: 8px;
    }
  }
  
  .visibility-toggle {
    min-width: 70px;
    border: 1px solid #e2e8f0;
    border-radius: 8px;
    background: #f8fafc;
    color: #334155;
    font-size: 14px;
    transition: all 0.2s ease;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0 16px;
    
    &:hover {
      color: #0f172a;
      border-color: #94a3b8;
      background-color: #f1f5f9;
    }
    
    &:active {
      background-color: #e2e8f0;
    }
  }
}

.model-tabs {
  margin-bottom: 28px;
  
  :deep(.el-radio-group) {
    display: flex;
    width: 100%;
    
    .el-radio-button {
      flex: 1;
      
      &:first-child {
        .el-radio-button__inner {
          border-radius: 8px 0 0 8px;
        }
      }
      
      &:last-child {
        .el-radio-button__inner {
          border-radius: 0 8px 8px 0;
        }
      }
      
      .el-radio-button__inner {
        width: 100%;
        padding: 12px 20px;
        border-color: #e2e8f0;
        color: #64748b;
        font-size: 15px;
        transition: all 0.2s ease;
        
        &:hover {
          color: #0f172a;
        }
      }
      
      &.is-active {
        .el-radio-button__inner {
          background: #0f172a;
          color: white;
          border-color: #0f172a;
          box-shadow: 0 0 0 1px #0f172a;
        }
      }
    }
  }
}

.model-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}

.model-card {
  background: #f8fafc;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  transition: all 0.2s ease;
  
  &:hover {
    transform: translateY(-3px);
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.08);
  }
  
  &.selected {
    border-color: #0f172a;
    box-shadow: 0 0 0 1px #0f172a, 0 10px 30px rgba(15, 23, 42, 0.1);
    background-color: #f8fafc;
    
    .model-name {
      color: #0f172a;
    }
  }
  
  &.disabled {
    opacity: 0.6;
    
    .model-description {
      color: #94a3b8;
    }
  }
  
  .model-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    
    .model-name-section {
      display: flex;
      align-items: center;
      gap: 8px;
      
      .model-name {
        margin: 0;
        font-size: 16px;
        font-weight: 600;
        color: #1e293b;
      }
    }
  }
  
  .model-description {
    color: #64748b;
    font-size: 14px;
    line-height: 1.5;
    min-height: 42px;
  }
  
  .model-settings {
    display: flex;
    flex-direction: column;
    gap: 16px;
    background-color: rgba(255, 255, 255, 0.7);
    padding: 16px;
    border-radius: 8px;
    border: 1px solid #e2e8f0;
    
    .param-item {
      display: flex;
      flex-direction: column;
      gap: 8px;
      
      .param-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        
        .param-label {
          color: #64748b;
          font-size: 14px;
        }
        
        .param-value {
          font-weight: 600;
          color: #1e293b;
        }
      }
    }
  }
  
  .model-actions {
    margin-top: auto;
    
    .el-button {
      width: 100%;
      border-radius: 8px;
      padding: 10px 0;
      
      &.is-default {
        background-color: #f1f5f9;
        color: #0f172a;
        border-color: #cbd5e1;
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
    background: #334155;
    border-radius: 3px;
  }
  
  .el-slider__button-wrapper {
    top: -16px;
  }
  
  .el-slider__button {
    width: 16px;
    height: 16px;
    border: 2px solid #334155;
    background-color: #ffffff;
    transition: transform 0.2s ease;
    
    &:hover {
      transform: scale(1.2);
    }
  }
  
  .el-slider__stop {
    background-color: #cbd5e1;
  }
}

:deep(.el-switch) {
  &.is-checked {
    .el-switch__core {
      background-color: #334155;
      border-color: #1e293b;
    }
  }
  
  .el-switch__core {
    border-color: #cbd5e1;
  }
}

:deep(.el-button--primary) {
  background: #0f172a;
  border: none;
  box-shadow: 0 4px 16px rgba(15, 23, 42, 0.16);
  
  &:hover:not(:disabled) {
    background: #1e293b;
    transform: translateY(-1px);
    box-shadow: 0 8px 20px rgba(15, 23, 42, 0.2);
  }
  
  &:active:not(:disabled) {
    transform: translateY(0);
    box-shadow: 0 4px 12px rgba(15, 23, 42, 0.12);
  }
  
  &:disabled {
    background: #cbd5e1;
    opacity: 0.7;
    box-shadow: none;
  }
}

:deep(.el-input-number) {
  .el-input__wrapper {
    background-color: #ffffff;
    border: 1px solid #e2e8f0;
    border-radius: 8px;
    padding: 0 12px;
    transition: all 0.2s ease;
    box-shadow: none;
    
    &:hover {
      border-color: #94a3b8;
    }
    
    &.is-focus {
      border-color: #334155;
      box-shadow: 0 0 0 3px rgba(51, 65, 85, 0.1);
    }
    
    .el-input__inner {
      color: #1e293b;
    }
  }
  
  .el-input-number__decrease, 
  .el-input-number__increase {
    color: #64748b;
    
    &:hover {
      color: #0f172a;
    }
  }
}

@media (max-width: 1024px) {
  .settings-main {
    flex-direction: column;
    
    .provider-list {
      width: 100%;
      
      .provider-menu {
        display: flex;
        flex-wrap: wrap;
        
        li {
          flex: 1;
          min-width: 120px;
        }
      }
    }
  }
  
  .model-list {
    grid-template-columns: 1fr;
  }
}
</style> 
