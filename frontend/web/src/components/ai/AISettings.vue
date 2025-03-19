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
    <div v-else>
      <el-tabs v-model="activeTab" tab-position="left" class="settings-tabs">
        <!-- 提供商设置 -->
        <el-tab-pane 
          v-for="provider in aiProviders" 
          :key="provider.providerId" 
          :label="provider.name" 
          :name="provider.providerId"
        >
          <div class="provider-settings">
            <h3>{{ provider.name }} 设置</h3>
            
            <el-form :model="providerSettings[provider.providerId]" label-position="top">
              <el-form-item label="API 密钥">
                <div class="api-key-input">
                  <el-input 
                    v-model="providerSettings[provider.providerId].apiKey" 
                    placeholder="输入您的API密钥"
                    :type="apiKeyVisibility[provider.providerId] ? 'text' : 'password'"
                  />
                  <el-button 
                    @click="toggleApiKeyVisibility(provider.providerId)" 
                    :icon="apiKeyVisibility[provider.providerId] ? 'View' : 'Hide'"
                    type="text"
                    class="visibility-toggle"
                  ></el-button>
                </div>
                <!-- 调试显示 -->
                <div v-if="isDev" class="debug-info">
                  密钥长度: {{ providerSettings[provider.providerId]?.apiKey?.length || 0 }}
                  <br>密钥开头: {{ providerSettings[provider.providerId]?.apiKey?.substring(0, 5) || '' }}
                </div>
              </el-form-item>
              
              <el-form-item label="API 端点" v-if="provider.hasEndpoint">
                <el-input 
                  v-model="providerSettings[provider.providerId].endpoint" 
                  placeholder="输入API端点URL"
                />
              </el-form-item>

              <div class="api-actions">
                <el-button type="primary" @click="saveSettings" :loading="isSubmitting">
                  保存设置
                </el-button>
                <el-button @click="testModelConnection(provider.providerId)" :loading="isSubmitting">
                  测试连接
                </el-button>
              </div>
              
              <el-divider>模型列表</el-divider>
              
              <!-- 统一模型配置区域 -->
              <div class="unified-models-list">
                <!-- 文本模型部分 -->
                <div class="model-type-section" v-if="getProviderTextModels(provider.providerId).length > 0">
                  <h4>文本模型</h4>
                  <el-table :data="getProviderTextModels(provider.providerId)" style="width: 100%" :fit="false" :row-class-name="getRowClassName">
                    <el-table-column label="模型" min-width="180" width="180">
                      <template #default="scope">
                        <div class="model-name">
                          <span>{{ scope.row.name }}</span>
                          <el-tag size="small" v-if="scope.row.isPaid" type="danger">付费</el-tag>
                        </div>
                      </template>
                    </el-table-column>
                    <el-table-column label="启用" width="60" align="center">
                      <template #default="scope">
                        <el-checkbox 
                          v-model="providerSettings[provider.providerId].enabledModels[scope.row.modelId]"
                        ></el-checkbox>
                      </template>
                    </el-table-column>
                    <el-table-column label="默认" width="100" align="center">
                      <template #default="scope">
                        <el-radio
                          v-model="defaultModel"
                          :label="scope.row.modelId"
                          @change="handleModelSelection(scope.row.modelId)"
                          :disabled="!providerSettings[provider.providerId].enabledModels[scope.row.modelId]"
                        ></el-radio>
                      </template>
                    </el-table-column>
                    <el-table-column label="参数配置" min-width="400">
                      <template #default="scope">
                        <div class="model-params" v-if="providerSettings[provider.providerId].enabledModels[scope.row.modelId]">
                          <div class="param-item">
                            <span class="param-label">温度：</span>
                            <el-slider 
                              v-model="providerSettings[provider.providerId].modelParams[scope.row.modelId].temperature" 
                              :min="0" 
                              :max="1" 
                              :step="0.1"
                              show-stops
                            ></el-slider>
                            <span class="param-value">{{ providerSettings[provider.providerId].modelParams[scope.row.modelId].temperature }}</span>
                          </div>
                          <div class="param-item">
                            <span class="param-label">最大输出长度：</span>
                            <el-input-number 
                              v-model="providerSettings[provider.providerId].modelParams[scope.row.modelId].maxTokens" 
                              :min="100" 
                              :max="8000"
                              :step="100"
                              size="small"
                            ></el-input-number>
                          </div>
                        </div>
                        <div v-else class="model-disabled-notice">
                          启用模型以配置参数
                        </div>
                      </template>
                    </el-table-column>
                    <el-table-column label="描述" min-width="300">
                      <template #default="scope">
                        <span class="model-description">{{ scope.row.description || '无描述' }}</span>
                      </template>
                    </el-table-column>
                  </el-table>
                </div>
                
                <!-- 图像模型部分 -->
                <div class="model-type-section" v-if="getProviderImageModels(provider.providerId).length > 0">
                  <h4>图像模型</h4>
                  <el-table :data="getProviderImageModels(provider.providerId)" style="width: 100%" :fit="false" :row-class-name="getRowClassName">
                    <el-table-column label="模型" min-width="180" width="180">
                      <template #default="scope">
                        <div class="model-name">
                          <span>{{ scope.row.name }}</span>
                          <el-tag size="small" v-if="scope.row.isPaid" type="danger">付费</el-tag>
                        </div>
                      </template>
                    </el-table-column>
                    <el-table-column label="启用" width="60" align="center">
                      <template #default="scope">
                        <el-checkbox 
                          v-model="providerSettings[provider.providerId].enabledModels[scope.row.modelId]"
                        ></el-checkbox>
                      </template>
                    </el-table-column>
                    <el-table-column label="默认" width="100" align="center">
                      <template #default="scope">
                        <el-radio
                          v-model="imageDefaultModel"
                          :label="scope.row.modelId"
                          @change="handleImageModelSelection(scope.row.modelId)"
                          :disabled="!providerSettings[provider.providerId].enabledModels[scope.row.modelId]"
                        ></el-radio>
                      </template>
                    </el-table-column>
                    <el-table-column label="描述" min-width="400">
                      <template #default="scope">
                        <span class="model-description">{{ scope.row.description || '无描述' }}</span>
                      </template>
                    </el-table-column>
                  </el-table>
                </div>
              </div>
            </el-form>
          </div>
        </el-tab-pane>
      </el-tabs>
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
const apiKeyVisibility = ref({})
const defaultModel = ref('')
const imageDefaultModel = ref('')

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
      
      // 初始化API密钥可见性
      if (apiKeyVisibility[provider.providerId] === undefined) {
        apiKeyVisibility[provider.providerId] = false;
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

// 切换API密钥可见性
const toggleApiKeyVisibility = (providerId) => {
  apiKeyVisibility[providerId] = !apiKeyVisibility[providerId]
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
  min-height: 600px;
  
  :deep(.el-tabs__header) {
    margin-right: 0;
    background-color: #f8fafc;
    border-right: 1px solid #e2e8f0;
    width: 200px;
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
    padding: 24px;
    background-color: #ffffff;
    overflow-x: auto;
  }
}

.provider-settings {
  min-width: 800px;

  h3 {
    font-size: 22px;
    font-weight: 600;
    margin-bottom: 24px;
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

.unified-models-list, .all-models-view {
  margin-top: 20px;
  width: 100%;
  
  .model-type-section {
    margin-bottom: 30px;
    
    h4 {
      margin-bottom: 15px;
      font-size: 16px;
      font-weight: 600;
      color: #409EFF;
    }
  }
  
  .model-name {
    display: flex;
    align-items: center;
    gap: 8px;
  }
  
  :deep(.selected-row) {
    background-color: #f0f9ff;
  }
}

.api-key-input {
  display: flex;
  align-items: center;
  
  .el-input {
    flex-grow: 1;
  }
  
  .el-button {
    margin-left: 10px;
  }
}

.api-actions {
  display: flex;
  gap: 10px;
  margin: 20px 0;
}

.model-params {
  padding: 12px;
  background-color: #f8f9fa;
  border-radius: 6px;
  
  .param-item {
    display: flex;
    align-items: center;
    margin-bottom: 12px;
    
    &:last-child {
      margin-bottom: 0;
    }
    
    .param-label {
      min-width: 110px;
      color: #5c6b77;
      font-size: 14px;
    }
    
    .param-value {
      margin-left: 10px;
      color: #1f2937;
      font-weight: 500;
      min-width: 30px;
    }
    
    .el-slider {
      flex: 1;
      min-width: 150px;
    }
  }
}

.model-disabled-notice {
  color: #909399;
  font-size: 13px;
  font-style: italic;
  padding: 8px 0;
}

.model-description {
  display: block;
  word-break: break-word;
  white-space: normal;
  line-height: 1.5;
  color: #475569;
  font-size: 14px;
  max-width: 380px;
  padding-right: 10px;
}

.el-radio {
  width: 100%;
  display: flex;
  justify-content: center;
  height: auto;
  padding: 0;
  margin-right: 0;
  
  :deep(.el-radio__input) {
    white-space: nowrap;
  }
  
  :deep(.el-radio__label) {
    display: none;
  }
}

.debug-info {
  margin-top: 8px;
  padding: 6px 10px;
  background-color: #f0f9ff;
  border-radius: 4px;
  font-size: 12px;
  color: #0369a1;
  font-family: monospace;
}

.loading-container {
  padding: 32px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05);
  border: 1px solid #e2e8f0;
  margin-bottom: 24px;
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

.visibility-toggle {
  padding: 4px 8px;
  margin-left: 8px;
  
  &:hover {
    background-color: var(--el-fill-color-light);
  }
}

@media (max-width: 1024px) {
  .ai-settings-container {
    padding: 20px;
  }
  
  .settings-tabs {
    :deep(.el-tabs__header) {
      width: 160px;
    }
    
    :deep(.el-tabs__item) {
      padding: 0 16px;
      font-size: 14px;
    }
  }
  
  .provider-settings {
    min-width: 600px;
  }
}
</style> 
