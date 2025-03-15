<template>
  <div class="ai-settings-container">
    <div class="settings-header">
      <h2>AI 模型配置</h2>
      <p>配置您的AI模型API密钥和参数</p>
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
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useAIStore } from '@/stores/ai'

const aiStore = useAIStore()
const activeTab = ref('openai')
const defaultModel = ref('')

// AI提供商列表
const aiProviders = [
  {
    id: 'openai',
    name: 'OpenAI',
    hasEndpoint: false,
    models: [
      {
        id: 'gpt-3.5-turbo',
        name: 'GPT-3.5 Turbo',
        description: '强大的语言模型，适合大多数任务，响应速度快。',
        isPaid: true
      },
      {
        id: 'gpt-4',
        name: 'GPT-4',
        description: '最先进的语言模型，具有更强的推理能力和更广泛的知识。',
        isPaid: true
      }
    ]
  },
  {
    id: 'anthropic',
    name: 'Anthropic',
    hasEndpoint: false,
    models: [
      {
        id: 'claude-3-opus',
        name: 'Claude 3 Opus',
        description: 'Anthropic的顶级模型，具有强大的推理和创作能力。',
        isPaid: true
      },
      {
        id: 'claude-3-sonnet',
        name: 'Claude 3 Sonnet',
        description: '平衡性能和速度的模型，适合大多数任务。',
        isPaid: true
      }
    ]
  },
  {
    id: 'google',
    name: 'Google AI',
    hasEndpoint: false,
    models: [
      {
        id: 'gemini-pro',
        name: 'Gemini Pro',
        description: 'Google的多模态AI模型，具有强大的理解和生成能力。',
        isPaid: true
      }
    ]
  },
  {
    id: 'deepseek',
    name: 'DeepSeek',
    hasEndpoint: false,
    models: [
      {
        id: 'deepseek-chat',
        name: 'DeepSeek Chat',
        description: 'DeepSeek的通用对话模型，擅长自然语言理解和生成。',
        isPaid: true
      },
      {
        id: 'deepseek-coder',
        name: 'DeepSeek Coder',
        description: 'DeepSeek的代码生成模型，专注于编程和开发任务。',
        isPaid: true
      }
    ]
  },
  {
    id: 'custom',
    name: '自定义模型',
    hasEndpoint: true,
    models: [
      {
        id: 'custom-model',
        name: '自定义模型',
        description: '配置您自己的AI模型API端点。',
        isPaid: false
      }
    ]
  }
]

// 初始化每个提供商的设置
const initProviderSettings = () => {
  const settings = {}
  
  aiProviders.forEach(provider => {
    const modelParams = {}
    const enabledModels = {}
    
    provider.models.forEach(model => {
      modelParams[model.id] = {
        temperature: 0.7,
        maxTokens: 2000
      }
      enabledModels[model.id] = false
    })
    
    settings[provider.id] = {
      apiKey: '',
      endpoint: provider.hasEndpoint ? '' : undefined,
      enabledModels,
      modelParams
    }
  })
  
  return settings
}

const providerSettings = reactive(initProviderSettings())

// 加载保存的设置
onMounted(() => {
  // 从localStorage或后端API加载设置
  const savedSettings = localStorage.getItem('aiProviderSettings')
  if (savedSettings) {
    try {
      const parsed = JSON.parse(savedSettings)
      Object.keys(parsed).forEach(providerId => {
        if (providerSettings[providerId]) {
          Object.assign(providerSettings[providerId], parsed[providerId])
        }
      })
    } catch (error) {
      console.error('Failed to load AI settings:', error)
    }
  }

  // 初始化时读取默认模型
  defaultModel.value = aiStore.defaultModel
})

// 保存提供商设置
const saveProviderSettings = (providerId) => {
  // 保存到localStorage或发送到后端API
  const allSettings = JSON.parse(localStorage.getItem('aiProviderSettings') || '{}')
  allSettings[providerId] = providerSettings[providerId]
  localStorage.setItem('aiProviderSettings', JSON.stringify(allSettings))
  
  ElMessage.success(`${getProviderName(providerId)} 设置已保存`)
}

// 测试连接
const testConnection = async (providerId) => {
  try {
    ElMessage.info(`正在测试 ${getProviderName(providerId)} 连接...`)
    
    // 使用AI store测试连接
    const success = await aiStore.testProviderConnection(providerId)
    
    if (success) {
      ElMessage.success(`${getProviderName(providerId)} 连接测试成功`)
    } else {
      ElMessage.error(`${getProviderName(providerId)} 连接测试失败`)
    }
  } catch (error) {
    ElMessage.error(`${getProviderName(providerId)} 连接测试失败: ${error.message}`)
  }
}

// 获取提供商名称
const getProviderName = (providerId) => {
  const provider = aiProviders.find(p => p.id === providerId)
  return provider ? provider.name : providerId
}

const handleDefaultModelChange = (modelId) => {
  aiStore.saveDefaultModel(modelId)
  ElMessage.success('默认模型设置已保存')
}
</script>

<style scoped lang="scss">
.ai-settings-container {
  padding: 16px 0;
}

.settings-header {
  margin-bottom: 24px;
  
  h2 {
    font-size: 20px;
    font-weight: 600;
    margin-bottom: 8px;
    color: #2B5876;
  }
  
  p {
    color: #6b7280;
    font-size: 14px;
  }
}

.settings-tabs {
  border: 1px solid rgba(0, 0, 0, 0.06);
  border-radius: 12px;
  overflow: hidden;
  
  :deep(.el-tabs__header) {
    margin-right: 0;
    background-color: #f9fafb;
    border-right: 1px solid rgba(0, 0, 0, 0.06);
  }
  
  :deep(.el-tabs__item) {
    height: 50px;
    line-height: 50px;
    text-align: left;
    padding: 0 20px;
    
    &.is-active {
      background-color: rgba(43, 88, 118, 0.08);
      color: #2B5876;
      font-weight: 500;
    }
  }
  
  :deep(.el-tabs__content) {
    padding: 20px;
  }
}

.provider-settings {
  h3 {
    font-size: 18px;
    font-weight: 600;
    margin-bottom: 20px;
    color: #2B5876;
  }
}

.models-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
  
  .model-item {
    border: 1px solid rgba(0, 0, 0, 0.08);
    border-radius: 10px;
    padding: 16px;
    
    .model-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 8px;
      
      .el-checkbox {
        font-weight: 500;
      }
    }
    
    .model-description {
      color: #6b7280;
      font-size: 14px;
      margin-bottom: 16px;
    }
    
    .model-params {
      background-color: #f9fafb;
      border-radius: 8px;
      padding: 16px;
      margin-top: 16px;
    }
  }
}

.provider-actions {
  margin-top: 24px;
  display: flex;
  gap: 12px;
}

.settings-card {
  margin-bottom: 24px;
  background-color: #1f1f1f;
  border: 1px solid #2a2a2a;
  
  :deep(.el-card__header) {
    background: linear-gradient(180deg, #252525 0%, #1f1f1f 100%);
    border-bottom: 1px solid #2a2a2a;
    padding: 16px 20px;
    
    .card-header {
      h3 {
        margin: 0;
        color: #e0e0e0;
        font-size: 16px;
        font-weight: 600;
      }
    }
  }
  
  .settings-form {
    padding: 20px;
    
    .el-form-item {
      margin-bottom: 24px;
      
      :deep(.el-form-item__label) {
        color: #e0e0e0;
      }
      
      .form-tip {
        margin-top: 8px;
        font-size: 12px;
        color: #a0a0a0;
      }
    }
    
    :deep(.el-select) {
      width: 100%;
      max-width: 300px;
      
      .el-input__wrapper {
        background-color: #252525;
        border: 1px solid #2a2a2a;
        
        &:hover {
          border-color: #409EFF;
        }
        
        &.is-focus {
          border-color: #409EFF;
          box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.1);
        }
        
        .el-input__inner {
          color: #e0e0e0;
        }
      }
    }
  }
}
</style> 