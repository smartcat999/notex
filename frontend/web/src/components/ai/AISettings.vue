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
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
  background-color: #ffffff;
}

.settings-header {
  margin-bottom: 32px;
  text-align: center;
  
  h2 {
    font-size: 24px;
    font-weight: 600;
    margin-bottom: 12px;
    color: #1a1a1a;
    background: linear-gradient(135deg, #2B5876 0%, #4E4376 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }
  
  p {
    color: #666;
    font-size: 16px;
  }
}

.settings-tabs {
  border: 1px solid #e5e7eb;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  
  :deep(.el-tabs__header) {
    margin-right: 0;
    background-color: #f8fafc;
    border-right: 1px solid #e5e7eb;
    width: 240px;
  }
  
  :deep(.el-tabs__item) {
    height: 56px;
    line-height: 56px;
    text-align: left;
    padding: 0 24px;
    font-size: 15px;
    color: #4b5563;
    transition: all 0.3s ease;
    
    &:hover {
      color: #2B5876;
    }
    
    &.is-active {
      background-color: rgba(43, 88, 118, 0.08);
      color: #2B5876;
      font-weight: 500;
    }
  }
  
  :deep(.el-tabs__content) {
    padding: 32px;
    background-color: #ffffff;
  }
}

.provider-settings {
  h3 {
    font-size: 20px;
    font-weight: 600;
    margin-bottom: 24px;
    color: #1a1a1a;
    position: relative;
    padding-bottom: 12px;
    
    &::after {
      content: '';
      position: absolute;
      bottom: 0;
      left: 0;
      width: 40px;
      height: 3px;
      background: linear-gradient(90deg, #2B5876 0%, #4E4376 100%);
      border-radius: 2px;
    }
  }
}

.models-list {
  display: flex;
  flex-direction: column;
  gap: 24px;
  
  .model-item {
    border: 1px solid #e5e7eb;
    border-radius: 12px;
    padding: 20px;
    transition: all 0.3s ease;
    background-color: #ffffff;
    
    &:hover {
      border-color: #2B5876;
      box-shadow: 0 4px 12px rgba(43, 88, 118, 0.1);
    }
    
    .model-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 12px;
      
      .el-checkbox {
        font-weight: 500;
        font-size: 15px;
      }
    }
    
    .model-description {
      color: #6b7280;
      font-size: 14px;
      line-height: 1.6;
      margin-bottom: 20px;
    }
    
    .model-params {
      background-color: #f8fafc;
      border-radius: 10px;
      padding: 20px;
      margin-top: 20px;
      border: 1px solid #e5e7eb;
    }
  }
}

.provider-actions {
  margin-top: 32px;
  display: flex;
  gap: 16px;
  
  .el-button {
    padding: 12px 24px;
    font-size: 15px;
    border-radius: 8px;
    transition: all 0.3s ease;
    
    &--primary {
      background: linear-gradient(135deg, #2B5876 0%, #4E4376 100%);
      border: none;
      
      &:hover {
        transform: translateY(-1px);
        box-shadow: 0 4px 12px rgba(43, 88, 118, 0.2);
      }
    }
  }
}

.settings-card {
  margin-top: 32px;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  
  :deep(.el-card__header) {
    background: linear-gradient(135deg, #2B5876 0%, #4E4376 100%);
    padding: 20px 24px;
    
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
    
    .el-form-item {
      margin-bottom: 28px;
      
      :deep(.el-form-item__label) {
        color: #1a1a1a;
        font-weight: 500;
      }
      
      .form-tip {
        margin-top: 8px;
        font-size: 13px;
        color: #6b7280;
      }
    }
    
    :deep(.el-select) {
      width: 100%;
      max-width: 320px;
      
      .el-input__wrapper {
        background-color: #ffffff;
        border: 1px solid #e5e7eb;
        border-radius: 8px;
        padding: 0 12px;
        transition: all 0.3s ease;
        
        &:hover {
          border-color: #2B5876;
        }
        
        &.is-focus {
          border-color: #2B5876;
          box-shadow: 0 0 0 3px rgba(43, 88, 118, 0.1);
        }
        
        .el-input__inner {
          color: #1a1a1a;
          height: 40px;
          line-height: 40px;
        }
      }
    }
  }
}

:deep(.el-slider) {
  .el-slider__runway {
    height: 6px;
    background-color: #e5e7eb;
    border-radius: 3px;
  }
  
  .el-slider__bar {
    background: linear-gradient(90deg, #2B5876 0%, #4E4376 100%);
    border-radius: 3px;
  }
  
  .el-slider__button-wrapper {
    top: -16px;
  }
  
  .el-slider__button {
    width: 16px;
    height: 16px;
    border: 2px solid #2B5876;
    background-color: #ffffff;
    
    &:hover {
      transform: scale(1.2);
    }
  }
}

:deep(.el-input-number) {
  .el-input__wrapper {
    background-color: #ffffff;
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    padding: 0 12px;
    transition: all 0.3s ease;
    
    &:hover {
      border-color: #2B5876;
    }
    
    &.is-focus {
      border-color: #2B5876;
      box-shadow: 0 0 0 3px rgba(43, 88, 118, 0.1);
    }
    
    .el-input__inner {
      color: #1a1a1a;
      height: 40px;
      line-height: 40px;
    }
  }
}
</style> 