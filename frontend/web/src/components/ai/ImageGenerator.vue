<template>
  <div class="image-generator-container">
    <div class="image-generator-header">
      <h3>AI图像生成</h3>
      <div class="model-selection">
        <el-select
          v-model="selectedModel"
          placeholder="选择图像模型"
          class="model-selector"
          :loading="isLoadingModels"
          :disabled="isGenerating"
        >
          <el-option
            v-for="model in aiStore.imageModels"
            :key="model.id"
            :label="`${model.name} (${model.providerName || model.provider})`"
            :value="model.id"
          />
        </el-select>
      </div>
    </div>
    
    <div class="image-generator-content">
      <div class="prompt-container">
        <el-input
          v-model="prompt"
          type="textarea"
          :rows="3"
          placeholder="描述您想要生成的图像，越详细越好..."
          resize="none"
          :disabled="isGenerating"
        />
        
        <div class="image-options">
          <el-select 
            v-model="imageSize" 
            placeholder="图像尺寸" 
            :disabled="isGenerating"
          >
            <el-option label="256 x 256" value="256x256" />
            <el-option label="512 x 512" value="512x512" />
            <el-option label="1024 x 1024" value="1024x1024" />
          </el-select>
          
          <el-select 
            v-model="imageCount" 
            placeholder="生成数量" 
            :disabled="isGenerating"
          >
            <el-option v-for="i in 4" :key="i" :label="`${i}张`" :value="i" />
          </el-select>
        </div>
      </div>
      
      <div v-if="showApiKeyPrompt" class="api-key-prompt">
        <div class="prompt-content">
          <el-icon class="warning-icon"><Warning /></el-icon>
          <div class="prompt-text">
            <h4>需要配置API密钥</h4>
            <p>请先在AI设置中配置 <strong>{{ currentProvider }}</strong> 的API密钥</p>
          </div>
        </div>
        <div class="prompt-actions">
          <el-button type="primary" @click="goToAISettings">
            前往配置
          </el-button>
          <el-button @click="showApiKeyPrompt = false">
            取消
          </el-button>
        </div>
      </div>
      
      <div class="generation-area" v-if="!showApiKeyPrompt">
        <div 
          class="generate-button-container"
          :class="{ 'generating': isGenerating }"
        >
          <el-button
            type="primary"
            :loading="isGenerating"
            :disabled="!prompt.trim() || isGenerating"
            @click="generateImage"
            class="generate-button"
          >
            {{ isGenerating ? '生成中...' : '生成图像' }}
          </el-button>
        </div>
        
        <div v-if="error" class="error-message">
          {{ error }}
        </div>
      </div>
      
      <div v-if="generatedImages.length > 0" class="generated-images">
        <h4>生成结果</h4>
        <div class="images-grid">
          <div 
            v-for="(image, index) in generatedImages" 
            :key="index"
            class="image-item"
          >
            <img :src="image" :alt="`生成的图像 ${index + 1}`" />
            <div class="image-actions">
              <el-button type="primary" size="small" @click="sendImageToChat(image)">
                发送到聊天
              </el-button>
              <el-button type="default" size="small" @click="downloadImage(image, index)">
                <el-icon><Download /></el-icon>
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useAIStore } from '@/stores/ai'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Warning, Download } from '@element-plus/icons-vue'

const props = defineProps({
  onSend: {
    type: Function,
    default: null
  },
  onClose: {
    type: Function,
    default: null
  }
})

const emit = defineEmits(['send', 'close'])

const aiStore = useAIStore()
const router = useRouter()

// 状态变量
const prompt = ref('')
const imageSize = ref('1024x1024')
const imageCount = ref(1)
const isLoadingModels = ref(false)
const isGenerating = ref(false)
const generatedImages = ref([])
const error = ref('')
const showApiKeyPrompt = ref(false)
const selectedModel = ref('')

// 计算属性
const currentProvider = computed(() => {
  if (!selectedModel.value) return ''
  const model = aiStore.imageModels.find(m => m.id === selectedModel.value)
  return model ? model.provider : ''
})

// 设置默认模型
watch(
  () => aiStore.imageModels,
  (models) => {
    if (models.length > 0 && !selectedModel.value) {
      selectedModel.value = models[0].id
    }
  },
  { immediate: true }
)

// 监听模型变化，重置错误并保存为默认值
watch(selectedModel, async (newValue, oldValue) => {
  error.value = ''
  showApiKeyPrompt.value = false
  
  // 仅当选择发生实际变化时保存
  if (newValue && newValue !== oldValue) {
    try {
      await aiStore.saveDefaultImageModel(newValue)
    } catch (err) {
      console.error('保存默认图像模型失败:', err)
    }
  }
})

// 初始化
onMounted(async () => {
  isLoadingModels.value = true
  try {
    // 确保图像模型已加载
    if (aiStore.imageModels.length === 0) {
      await aiStore.loadImageModels()
    }
    
    // 设置默认选择的模型
    if (aiStore.currentImageModel) {
      selectedModel.value = aiStore.currentImageModel
    } else if (aiStore.imageModels.length > 0) {
      selectedModel.value = aiStore.imageModels[0].id
    }
  } catch (err) {
    console.error('加载图像模型失败:', err)
    error.value = '无法加载图像模型：' + (err.message || '未知错误')
  } finally {
    isLoadingModels.value = false
  }
})

// 生成图像
const generateImage = async () => {
  if (!prompt.value.trim()) {
    ElMessage.warning('请输入图像描述')
    return
  }
  
  if (!selectedModel.value) {
    ElMessage.warning('请选择图像生成模型')
    return
  }
  
  error.value = ''
  showApiKeyPrompt.value = false
  
  try {
    isGenerating.value = true
    generatedImages.value = []
    
    // 更新并保存AI Store中的当前图像模型
    await aiStore.saveDefaultImageModel(selectedModel.value)
    
    // 生成图像
    const images = await aiStore.generateImage(prompt.value, {
      n: imageCount.value,
      size: imageSize.value
    })
    
    if (images && images.length > 0) {
      generatedImages.value = images
      ElMessage.success(`成功生成${images.length}张图像`)
    } else {
      error.value = '未能生成图像，请重试'
    }
  } catch (err) {
    console.error('生成图像失败:', err)
    if (err.message.includes('未配置API密钥')) {
      showApiKeyPrompt.value = true
    } else {
      error.value = '生成图像失败：' + (err.message || '未知错误')
    }
  } finally {
    isGenerating.value = false
  }
}

// 发送图像到聊天
const sendImageToChat = (imageUrl) => {
  const imageMarkdown = `![AI生成的图像](${imageUrl})\n\n*基于提示: ${prompt.value}*`
  
  if (props.onSend) {
    props.onSend(imageMarkdown)
  } else {
    emit('send', imageMarkdown)
  }
  
  // 如果有关闭回调，则调用
  if (props.onClose) {
    props.onClose()
  } else {
    emit('close')
  }
}

// 下载图像
const downloadImage = async (imageUrl, index) => {
  try {
    // 获取图像数据
    const response = await fetch(imageUrl)
    const blob = await response.blob()
    
    // 创建下载链接
    const downloadUrl = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = downloadUrl
    a.download = `ai-image-${Date.now()}-${index}.png`
    document.body.appendChild(a)
    a.click()
    
    // 清理
    document.body.removeChild(a)
    URL.revokeObjectURL(downloadUrl)
    
    ElMessage.success('图像下载成功')
  } catch (error) {
    console.error('下载图像失败:', error)
    ElMessage.error('下载图像失败')
  }
}

// 前往AI设置页面
const goToAISettings = () => {
  router.push('/ai/settings')
}

// 当选择的模型改变时保存为默认值
const handleModelChange = async (model) => {
  selectedModel.value = model
  // 更新AI store的当前图像模型
  await aiStore.saveDefaultImageModel(selectedModel.value)
}
</script>

<style lang="scss" scoped>
.image-generator-container {
  width: 100%;
  background-color: #ffffff;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  
  .image-generator-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px 20px;
    border-bottom: 1px solid #f0f0f0;
    
    h3 {
      margin: 0;
      font-size: 18px;
      font-weight: 600;
      color: #333;
    }
    
    .model-selection {
      .model-selector {
        width: 220px;
      }
    }
  }
  
  .image-generator-content {
    padding: 20px;
    display: flex;
    flex-direction: column;
    gap: 20px;
    
    .prompt-container {
      display: flex;
      flex-direction: column;
      gap: 12px;
      
      .image-options {
        display: flex;
        gap: 12px;
        
        .el-select {
          flex: 1;
        }
      }
    }
    
    .generation-area {
      display: flex;
      flex-direction: column;
      gap: 12px;
      
      .generate-button-container {
        display: flex;
        justify-content: center;
        margin: 10px 0;
        
        &.generating {
          opacity: 0.7;
        }
        
        .generate-button {
          width: 200px;
          height: 48px;
          font-size: 16px;
          background: linear-gradient(135deg, #8b5cf6 0%, #6366f1 100%);
          border: none;
          
          &:hover:not(:disabled) {
            background: linear-gradient(135deg, #7c4ddb 0%, #5459e3 100%);
            transform: translateY(-1px);
            box-shadow: 0 4px 8px rgba(99, 102, 241, 0.3);
          }
        }
      }
      
      .error-message {
        color: #dc2626;
        text-align: center;
        padding: 10px;
        background-color: #fee2e2;
        border-radius: 6px;
      }
    }
    
    .api-key-prompt {
      margin: 10px 0;
      padding: 16px;
      border-radius: 8px;
      border: 1px solid #f0f0f0;
      background-color: #fffbeb;
      
      .prompt-content {
        display: flex;
        align-items: flex-start;
        gap: 12px;
        margin-bottom: 16px;
        
        .warning-icon {
          font-size: 20px;
          color: #f59e0b;
          margin-top: 2px;
        }
        
        .prompt-text {
          h4 {
            margin: 0 0 8px 0;
            color: #b45309;
          }
          
          p {
            margin: 0;
            color: #78350f;
          }
        }
      }
      
      .prompt-actions {
        display: flex;
        justify-content: flex-end;
        gap: 8px;
      }
    }
    
    .generated-images {
      margin-top: 20px;
      
      h4 {
        margin: 0 0 16px 0;
        font-size: 16px;
        font-weight: 600;
        color: #333;
      }
      
      .images-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
        gap: 16px;
        
        .image-item {
          border-radius: 8px;
          overflow: hidden;
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
          background-color: #f8fafc;
          
          img {
            width: 100%;
            display: block;
            aspect-ratio: 1 / 1;
            object-fit: cover;
          }
          
          .image-actions {
            display: flex;
            justify-content: space-between;
            padding: 8px;
            
            .el-button {
              flex: 1;
              
              &:first-child {
                margin-right: 8px;
              }
            }
          }
        }
      }
    }
  }
}
</style> 