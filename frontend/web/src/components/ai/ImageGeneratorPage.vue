<template>
  <div class="image-generator-page">
    <div class="page-header">
      <h1>AI 图像生成</h1>
      <el-button @click="navigateToSettings" v-if="showApiKeyPrompt" type="primary" plain>
        配置API密钥
      </el-button>
    </div>

    <div v-if="showApiKeyPrompt" class="api-key-prompt">
      <el-alert
        title="需要配置API密钥"
        type="warning"
        description="要使用图像生成功能，请先配置有效的OpenAI或Stability AI API密钥。"
        show-icon
        :closable="false"
      />
    </div>

    <div class="generator-container">
      <!-- 加载状态 -->
      <div v-if="isInitializing" class="loading-state">
        <el-skeleton :rows="3" animated />
      </div>

      <!-- 生成器控件 -->
      <div v-else class="generator-controls">
        <div class="model-selector">
          <el-form-item label="选择模型">
            <el-select 
              v-model="selectedModel" 
              placeholder="选择图像生成模型" 
              style="width: 100%"
              :loading="aiStore.imageModels.length === 0"
            >
              <el-option
                v-for="model in aiStore.imageModels"
                :key="model.id"
                :label="model.name"
                :value="model.id"
              >
                <div class="model-option">
                  <span>{{ model.name }}</span>
                  <span class="model-provider">{{ getProviderName(model.provider) }}</span>
                </div>
              </el-option>
            </el-select>
          </el-form-item>
        </div>

        <div class="prompt-input">
          <el-form-item label="图像描述">
            <el-input
              v-model="prompt"
              type="textarea"
              :rows="4"
              placeholder="请描述您想要生成的图像内容..."
              resize="none"
            />
          </el-form-item>
        </div>

        <div class="generation-options">
          <el-form-item label="图像数量">
            <el-slider v-model="imageCount" :min="1" :max="4" :step="1" show-stops />
          </el-form-item>

          <el-form-item label="图像尺寸">
            <el-select v-model="imageSize" style="width: 200px">
              <el-option
                v-for="option in imageSizeOptions"
                :key="option.value"
                :label="option.label"
                :value="option.value"
              />
            </el-select>
          </el-form-item>
        </div>

        <div class="action-buttons">
          <el-button
            type="primary"
            @click="generateImage"
            :loading="aiStore.isGeneratingImage"
            :disabled="!prompt.trim() || !selectedModel || showApiKeyPrompt"
          >
            {{ aiStore.isGeneratingImage ? '生成中...' : '生成图像' }}
          </el-button>
        </div>
      </div>

      <div class="generated-images" v-if="aiStore.generatedImages.length > 0">
        <h3>生成的图像</h3>
        <div class="images-grid">
          <div v-for="(image, index) in aiStore.generatedImages" :key="index" class="image-item">
            <img :src="image" alt="生成的图像" />
            <div class="image-actions">
              <el-button size="small" @click="downloadImage(image, index)">
                下载
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useAIStore } from '@/stores/ai';
import { useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';

const aiStore = useAIStore();
const router = useRouter();

// 状态
const prompt = ref('');
const selectedModel = ref('');
const imageCount = ref(1);
const imageSize = ref('512x512');
const isInitializing = ref(true);

// 监听模型变化
watch(() => aiStore.imageModels, (newModels) => {
  if (newModels.length > 0 && !selectedModel.value) {
    // 如果有默认模型则使用默认模型
    if (aiStore.currentImageModel) {
      selectedModel.value = aiStore.currentImageModel;
    } else {
      // 否则使用第一个模型
      selectedModel.value = newModels[0].id;
    }
  }
}, { immediate: true });

// 监听 currentImageModel 变化
watch(() => aiStore.currentImageModel, (newModel) => {
  if (newModel && !selectedModel.value) {
    selectedModel.value = newModel;
  }
}, { immediate: true });

// 计算属性
const showApiKeyPrompt = computed(() => {
  if (!selectedModel.value) return true;
  
  // 找到当前选择的模型
  const model = aiStore.imageModels.find(m => m.id === selectedModel.value);
  if (!model) return true;
  
  // 检查该模型提供商的设置
  const settings = JSON.parse(localStorage.getItem('aiProviderSettings') || '{}');
  const providerSettings = settings[model.provider];
  
  return !providerSettings?.apiKey;
});

const imageSizeOptions = computed(() => {
  if (!selectedModel.value) return [];
  
  const model = aiStore.imageModels.find(m => m.id === selectedModel.value);
  if (!model) return [];
  
  // Stability AI SDXL 模型的特殊尺寸选项
  if (model.provider === 'stabilityai' && model.id.includes('xl')) {
    return [
      { label: '1024 x 1024', value: '1024x1024' },
      { label: '1152 x 896', value: '1152x896' },
      { label: '1216 x 832', value: '1216x832' },
      { label: '1344 x 768', value: '1344x768' },
      { label: '1536 x 640', value: '1536x640' },
      { label: '640 x 1536', value: '640x1536' },
      { label: '768 x 1344', value: '768x1344' },
      { label: '832 x 1216', value: '832x1216' },
      { label: '896 x 1152', value: '896x1152' }
    ];
  }
  
  // 其他模型的标准尺寸选项
  return [
    { label: '小 (256x256)', value: '256x256' },
    { label: '中 (512x512)', value: '512x512' },
    { label: '大 (1024x1024)', value: '1024x1024' }
  ];
});

// 监听模型变化时更新尺寸
watch(selectedModel, (newModel) => {
  if (newModel) {
    const model = aiStore.imageModels.find(m => m.id === newModel);
    if (model?.provider === 'stabilityai' && model.id.includes('xl')) {
      // 对于 Stability AI SDXL 模型，默认使用 1024x1024
      imageSize.value = '1024x1024';
    } else {
      // 其他模型使用默认的中等尺寸
      imageSize.value = '512x512';
    }
  }
});

// 初始化函数
const initializeImageGenerator = async () => {
  try {
    isInitializing.value = true;
    
    // 确保 AI Store 已初始化
    if (!aiStore.initialized) {
      await aiStore.initialize();
    }
    
  } catch (error) {
    console.error('初始化图像生成器失败:', error);
    ElMessage.error('加载模型失败，请刷新页面重试');
  } finally {
    isInitializing.value = false;
  }
};

// 方法
const generateImage = async () => {
  if (!prompt.value.trim() || !selectedModel.value) {
    ElMessage.warning('请填写完整的图像描述并选择模型');
    return;
  }
  
  try {
    await aiStore.generateImage(prompt.value, {
      n: imageCount.value,
      size: imageSize.value
    });
  } catch (error) {
    ElMessage.error('生成图像失败: ' + (error.message || '未知错误'));
  }
};

const downloadImage = (imageUrl, index) => {
  const link = document.createElement('a');
  link.href = imageUrl;
  link.download = `generated-image-${new Date().getTime()}-${index}.png`;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
};

const navigateToSettings = () => {
  router.push('/ai/settings');
};

const getProviderName = (providerId) => {
  const providers = {
    'openai': 'OpenAI',
    'stabilityai': 'Stability AI'
  };
  
  return providers[providerId] || providerId;
};

// 生命周期钩子
onMounted(() => {
  initializeImageGenerator();
});
</script>

<style lang="scss" scoped>
.image-generator-page {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
  
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
    
    h1 {
      font-size: 24px;
      font-weight: 600;
      margin: 0;
    }
  }
  
  .api-key-prompt {
    margin-bottom: 24px;
  }
  
  .generator-container {
    display: flex;
    flex-direction: column;
    gap: 32px;
    
    .loading-state {
      background-color: #fff;
      border-radius: 12px;
      padding: 24px;
      box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
    }
    
    .generator-controls {
      background-color: #fff;
      border-radius: 12px;
      padding: 24px;
      box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
      display: flex;
      flex-direction: column;
      gap: 20px;
      
      .model-option {
        display: flex;
        justify-content: space-between;
        align-items: center;
        width: 100%;
        
        .model-provider {
          color: #909399;
          font-size: 12px;
        }
      }
    }
    
    .generated-images {
      h3 {
        font-size: 18px;
        font-weight: 600;
        margin-bottom: 16px;
      }
      
      .images-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
        gap: 16px;
        
        .image-item {
          background-color: #fff;
          border-radius: 12px;
          overflow: hidden;
          box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
          
          img {
            width: 100%;
            height: auto;
            display: block;
          }
          
          .image-actions {
            padding: 12px;
            display: flex;
            justify-content: center;
          }
        }
      }
    }
  }
}
</style> 