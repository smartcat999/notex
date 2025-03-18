<template>
  <div class="ai-writing-container">
    <div class="writing-header">
      <div class="title-area">
        <el-icon class="header-icon"><Edit /></el-icon>
        <h2>AI 写作</h2>
      </div>
      <div class="header-actions">
        <el-select v-model="selectedModel" placeholder="选择AI模型" class="model-selector">
          <el-option
            v-for="model in aiStore.availableModels"
            :key="model.id"
            :label="model.name"
            :value="model.id"
          />
        </el-select>
      </div>
    </div>

    <div class="writing-content">
      <div class="writing-intro">
        <h3>让灵感帮我创作，才思如泉涌</h3>
        <p>选择类型，输入主题，AI帮你轻松创作优质内容</p>
      </div>

      <div class="writing-input-section">
        <div class="theme-input">
          <el-input
            v-model="writingTheme"
            placeholder="输入创作主题"
            class="theme-input-box"
            @keydown.enter="generateContent"
          >
            <template #prefix>
              <el-icon><EditPen /></el-icon>
            </template>
          </el-input>
          <el-button 
            type="primary" 
            class="generate-button" 
            @click="generateContent"
            :disabled="!writingTheme || isGenerating"
            :loading="isGenerating"
          >
            <el-icon><Promotion /></el-icon>
            生成内容
          </el-button>
        </div>

        <div class="type-tabs">
          <el-tabs v-model="activeTab" @tab-click="handleTabChange">
            <el-tab-pane label="全部" name="all"></el-tab-pane>
            <el-tab-pane label="学习教育" name="education"></el-tab-pane>
            <el-tab-pane label="工作" name="work"></el-tab-pane>
            <el-tab-pane label="营销" name="marketing"></el-tab-pane>
            <el-tab-pane label="回复" name="reply"></el-tab-pane>
          </el-tabs>
        </div>
      </div>

      <div class="template-section" v-if="!generatedContent">
        <div class="template-grid">
          <div v-for="template in filteredTemplates" :key="template.id" class="template-card" @click="selectTemplate(template)">
            <div class="template-icon" :class="template.category">
              <el-icon><component :is="template.icon" /></el-icon>
            </div>
            <div class="template-info">
              <h4>{{ template.name }}</h4>
              <p>{{ template.description }}</p>
            </div>
          </div>
        </div>
      </div>

      <div class="result-section" v-if="generatedContent">
        <el-card class="result-card">
          <template #header>
            <div class="result-header">
              <div class="result-title">
                <span>{{ currentTemplate ? currentTemplate.name : '生成内容' }}</span>
                <span class="theme-tag">{{ writingTheme }}</span>
              </div>
              <div class="result-actions">
                <el-button type="text" @click="newWriting">
                  <el-icon><CirclePlus /></el-icon>
                  新建
                </el-button>
                <el-button type="text" @click="copyContent">
                  <el-icon><CopyDocument /></el-icon>
                  复制
                </el-button>
              </div>
            </div>
          </template>
          <div class="result-content" v-html="renderedContent"></div>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { 
  Edit, EditPen, Promotion, CirclePlus, CopyDocument, Document, 
  Reading, Briefcase, Histogram, ChatLineRound, Timer, 
  GoldMedal, PriceTag, Collection, Opportunity, List, School, Money
} from '@element-plus/icons-vue'
import { useAIStore } from '@/stores/ai'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import { marked } from 'marked'

const aiStore = useAIStore()
const userStore = useUserStore()
const selectedModel = ref('')
const writingTheme = ref('')
const activeTab = ref('all')
const generatedContent = ref('')
const isGenerating = ref(false)
const currentTemplate = ref(null)

// 所有写作模板
const templates = [
  {
    id: 1,
    name: '作文',
    description: '输入标题，一键生成分段作文',
    icon: 'Reading',
    category: 'education',
    prompt: '请根据"{theme}"这个标题，写一篇800字左右的分段作文。要求：结构清晰，开头引人入胜，中间过渡自然，结尾有深意，适当使用比喻、排比等修辞手法。请使用Markdown格式输出。'
  },
  {
    id: 2,
    name: '读后感',
    description: '一键生成对应书籍的读后感',
    icon: 'Reading',
    category: 'education',
    prompt: '请为我写一篇关于《{theme}》的读后感，800字左右。要求：包含对作品主题的理解，对作者写作风格的赏析，阅读过程中的感受，以及从作品中获得的启示。请使用Markdown格式输出。'
  },
  {
    id: 3,
    name: '议论文',
    description: '支持自定义论点的议论文',
    icon: 'Reading',
    category: 'education',
    prompt: '请以"{theme}"为主题，写一篇议论文，800字左右。要求：有明确的论点，三个论据支撑，论据包含事实论据和道理论据，有适当的引用。请使用Markdown格式输出。'
  },
  {
    id: 4,
    name: '课程大纲',
    description: '一键生成结构化课程大纲',
    icon: 'School',
    category: 'education',
    prompt: '请为"{theme}"这门课程设计一个详细的教学大纲，包括课程目标、课程内容（分为若干模块和课时）、教学方法、学习资源和考核方式。使用Markdown格式输出，保持层次清晰。'
  },
  {
    id: 5,
    name: '周报',
    description: '快速生成工作周报',
    icon: 'Timer',
    category: 'work',
    prompt: '请帮我生成一篇关于"{theme}"项目的工作周报。包括本周工作内容概述、工作成果、遇到的问题与解决方案、下周工作计划等部分。使用Markdown格式输出，保持专业简洁的语言风格。'
  },
  {
    id: 6,
    name: '工作计划',
    description: '一键生成详细工作计划',
    icon: 'Briefcase',
    category: 'work',
    prompt: '请为"{theme}"制定一份详细的工作计划。要包含项目背景介绍、目标设定、具体工作内容、时间节点安排、资源需求以及预期成果等部分。使用Markdown格式输出，语言要专业规范。'
  },
  {
    id: 7,
    name: '会议纪要',
    description: '快速生成会议纪要模板',
    icon: 'List',
    category: 'work',
    prompt: '请生成一篇关于"{theme}"的会议纪要。包括会议基本信息(时间、地点、参会人员）、会议议程、讨论内容、决议事项、后续行动计划等部分。使用Markdown格式输出，保持条理清晰。'
  },
  {
    id: 8,
    name: '产品方案',
    description: '生成产品方案文档',
    icon: 'Opportunity',
    category: 'work',
    prompt: '请为"{theme}"产品编写一份产品方案文档。包括产品背景与定位、目标用户分析、功能需求描述、技术实现方案、项目计划与里程碑、风险评估等部分。使用Markdown格式输出，语言专业简洁。'
  },
  {
    id: 9,
    name: '营销文案',
    description: '生成吸引人的营销文案',
    icon: 'PriceTag',
    category: 'marketing',
    prompt: '请为"{theme}"产品/服务编写一篇吸引人的营销文案。文案需包含引人注目的标题、产品/服务价值点描述、情感共鸣部分、用户痛点分析与解决方案、号召行动等要素。使用Markdown格式输出，语言要生动有感染力。'
  },
  {
    id: 10,
    name: '公众号文章',
    description: '生成吸引读者的公众号内容',
    icon: 'Collection',
    category: 'marketing',
    prompt: '请以"{theme}"为主题，创作一篇适合公众号发布的文章。要求：标题吸引人，开头引人入胜，内容分3-5个小标题展开，每个小标题下内容围绕主题展开，配以适当的例子、数据或故事，结尾有号召行动。使用Markdown格式输出，全文1500字左右。'
  },
  {
    id: 11,
    name: '电子邮件',
    description: '快速生成邮件内容',
    icon: 'ChatLineRound',
    category: 'reply',
    prompt: '请帮我起草一封关于"{theme}"的电子邮件。邮件需包含适当的称呼、清晰的主题陈述、详细的内容说明、礼貌的结束语和签名。使用Markdown格式输出，语言要正式、专业。'
  },
  {
    id: 12,
    name: '感谢信',
    description: '生成诚挚的感谢信',
    icon: 'GoldMedal',
    category: 'reply',
    prompt: '请以"{theme}"为主题，帮我写一封真诚的感谢信。信中应包含具体的感谢原因、对方帮助的价值与意义、个人感受以及未来的祝福或期许。使用Markdown格式输出，语言要真诚、温暖而有礼貌。'
  }
]

// 根据当前选择的标签过滤模板
const filteredTemplates = computed(() => {
  if (activeTab.value === 'all') {
    return templates
  } else {
    return templates.filter(template => template.category === activeTab.value)
  }
})

// 将markdown内容渲染为HTML
const renderedContent = computed(() => {
  if (!generatedContent.value) return ''
  return marked(generatedContent.value)
})

// 初始化
onMounted(async () => {
  // 如果AI存储已初始化，直接使用当前模型
  if (aiStore.initialized) {
    selectedModel.value = aiStore.currentModel
    console.log('AIWriting: AIStore已初始化，使用当前模型:', selectedModel.value)
  } 
  // 如果AI存储有可用模型但未初始化
  else if (aiStore.availableModels.length > 0) {
    if (aiStore.defaultModel) {
      selectedModel.value = aiStore.defaultModel
      console.log('AIWriting: 使用默认模型:', selectedModel.value)
    } else {
      selectedModel.value = aiStore.availableModels[0].id
      console.log('AIWriting: 使用第一个可用模型:', selectedModel.value)
    }
  }
  // 如果AI存储未初始化，等待初始化
  else {
    console.log('AIWriting: AIStore未初始化，等待初始化完成')
    // 用户已登录但AI存储未初始化，手动初始化
    if (userStore.isAuthenticated && !aiStore.initialized && !aiStore.isInitializing) {
      await aiStore.initialize()
    }
  }
})

// 监听AI存储初始化状态变化
watch(() => aiStore.initialized, (isInitialized) => {
  if (isInitialized) {
    console.log('AIWriting: AIStore初始化状态变化，重新设置模型')
    selectedModel.value = aiStore.currentModel
  }
})

// 监听AI存储默认模型变化
watch(() => aiStore.defaultModel, (newDefaultModel) => {
  if (newDefaultModel && (!selectedModel.value || selectedModel.value !== newDefaultModel)) {
    console.log('AIWriting: 默认模型变化:', newDefaultModel)
    selectedModel.value = newDefaultModel
  }
})

// 监听AI存储当前模型变化
watch(() => aiStore.currentModel, (newCurrentModel) => {
  if (newCurrentModel && selectedModel.value !== newCurrentModel) {
    console.log('AIWriting: 当前模型变化:', newCurrentModel)
    selectedModel.value = newCurrentModel
  }
})

// 处理标签切换
const handleTabChange = (tab) => {
  console.log('切换到标签:', tab.props.name)
}

// 选择模板
const selectTemplate = (template) => {
  currentTemplate.value = template
  if (!writingTheme.value) {
    ElMessage.warning('请先输入创作主题')
    return
  }
  generateContent()
}

// 生成内容
const generateContent = async () => {
  if (!writingTheme.value.trim()) {
    ElMessage.warning('请输入创作主题')
    return
  }
  
  if (!currentTemplate.value) {
    ElMessage.warning('请选择写作模板')
    return
  }
  
  if (isGenerating.value) return
  
  try {
    isGenerating.value = true
    
    // 使用用户输入的主题替换提示词中的{theme}
    const prompt = currentTemplate.value.prompt.replace('{theme}', writingTheme.value)
    
    // 使用AI Store发送消息
    await aiStore.sendMessage(prompt.trim(), -1, (response) => {
      // 确保content是字符串
      const content = typeof response === 'object' ? response.fullContent || '' : response
      // 实时更新生成结果
      generatedContent.value = content
    })
    
    ElMessage.success('内容生成完成')
  } catch (error) {
    console.error('生成内容失败:', error)
    ElMessage.error('生成内容失败，请重试')
  } finally {
    isGenerating.value = false
  }
}

// 新建写作
const newWriting = () => {
  writingTheme.value = ''
  generatedContent.value = ''
  currentTemplate.value = null
}

// 复制内容
const copyContent = () => {
  if (!generatedContent.value) {
    ElMessage.warning('没有内容可复制')
    return
  }
  
  // 创建一个纯文本版本，移除markdown标记
  const textToCopy = generatedContent.value
  
  navigator.clipboard.writeText(textToCopy)
    .then(() => {
      ElMessage.success('内容已复制到剪贴板')
    })
    .catch(err => {
      console.error('复制失败:', err)
      ElMessage.error('复制失败，请手动选择文本并复制')
    })
}
</script>

<style lang="scss" scoped>
.ai-writing-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: #f8fafc;
  
  .writing-header {
    padding: 16px 24px;
    background: #ffffff;
    border-bottom: 1px solid #e2e8f0;
    display: flex;
    align-items: center;
    justify-content: space-between;
    position: sticky;
    top: 0;
    z-index: 10;
    
    .title-area {
      display: flex;
      align-items: center;
      gap: 8px;
      
      .header-icon {
        font-size: 24px;
        color: #8b5cf6;
      }
      
      h2 {
        font-size: 22px;
        font-weight: 600;
        background: linear-gradient(135deg, #8b5cf6 0%, #6366f1 100%);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        margin: 0;
      }
    }
    
    .model-selector {
      width: 240px;
      
      :deep(.el-input__wrapper) {
        background-color: #ffffff;
        border: 1px solid #e2e8f0;
        border-radius: 8px;
        box-shadow: none;
        
        &:hover {
          border-color: #8b5cf6;
        }
        
        &.is-focus {
          border-color: #8b5cf6;
          box-shadow: 0 0 0 3px rgba(139, 92, 246, 0.1);
        }
      }
    }
  }
  
  .writing-content {
    flex: 1;
    padding: 24px;
    overflow-y: auto;
    
    .writing-intro {
      text-align: center;
      margin-bottom: 32px;
      
      h3 {
        font-size: 28px;
        font-weight: 600;
        color: #1e293b;
        margin-bottom: 12px;
      }
      
      p {
        font-size: 16px;
        color: #64748b;
      }
    }
    
    .writing-input-section {
      margin-bottom: 32px;
      
      .theme-input {
        display: flex;
        gap: 16px;
        margin-bottom: 24px;
        
        .theme-input-box {
          flex: 1;
          
          :deep(.el-input__wrapper) {
            background-color: #ffffff;
            border: 1px solid #e2e8f0;
            border-radius: 8px;
            box-shadow: none;
            height: 48px;
            
            &:hover {
              border-color: #8b5cf6;
            }
            
            &.is-focus {
              border-color: #8b5cf6;
              box-shadow: 0 0 0 3px rgba(139, 92, 246, 0.1);
            }
            
            .el-input__prefix {
              color: #8b5cf6;
            }
          }
          
          :deep(.el-input__inner) {
            font-size: 16px;
          }
        }
        
        .generate-button {
          height: 48px;
          background: linear-gradient(135deg, #8b5cf6 0%, #6366f1 100%);
          border: none;
          border-radius: 8px;
          font-weight: 500;
          min-width: 140px;
          
          .el-icon {
            margin-right: 8px;
          }
          
          &:hover:not(:disabled) {
            transform: translateY(-1px);
            box-shadow: 0 4px 8px rgba(139, 92, 246, 0.2);
          }
        }
      }
      
      .type-tabs {
        :deep(.el-tabs__header) {
          margin-bottom: 24px;
        }
        
        :deep(.el-tabs__nav-wrap::after) {
          background-color: #e2e8f0;
          height: 1px;
        }
        
        :deep(.el-tabs__active-bar) {
          background-color: #8b5cf6;
          height: 3px;
          border-radius: 3px;
        }
        
        :deep(.el-tabs__item) {
          color: #64748b;
          font-size: 15px;
          height: 40px;
          
          &.is-active {
            color: #8b5cf6;
            font-weight: 500;
          }
          
          &:hover {
            color: #8b5cf6;
          }
        }
      }
    }
    
    .template-section {
      .template-grid {
        display: grid;
        grid-template-columns: repeat(3, 1fr);
        gap: 24px;
        
        .template-card {
          background: #ffffff;
          border-radius: 16px;
          padding: 24px;
          display: flex;
          align-items: center;
          gap: 16px;
          cursor: pointer;
          transition: all 0.3s ease;
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
          
          &:hover {
            transform: translateY(-5px);
            box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
          }
          
          .template-icon {
            width: 48px;
            height: 48px;
            display: flex;
            align-items: center;
            justify-content: center;
            border-radius: 12px;
            color: #ffffff;
            font-size: 24px;
            
            &.education {
              background: linear-gradient(135deg, #8b5cf6 0%, #6366f1 100%);
            }
            
            &.work {
              background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
            }
            
            &.marketing {
              background: linear-gradient(135deg, #f97316 0%, #ea580c 100%);
            }
            
            &.reply {
              background: linear-gradient(135deg, #22c55e 0%, #16a34a 100%);
            }
          }
          
          .template-info {
            flex: 1;
            
            h4 {
              font-size: 16px;
              font-weight: 600;
              color: #1e293b;
              margin: 0 0 4px;
            }
            
            p {
              font-size: 14px;
              color: #64748b;
              margin: 0;
            }
          }
        }
      }
    }
    
    .result-section {
      .result-card {
        border-radius: 16px;
        overflow: hidden;
        
        :deep(.el-card__header) {
          padding: 16px 24px;
          background: linear-gradient(135deg, #8b5cf6 0%, #6366f1 100%);
          
          .result-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            
            .result-title {
              color: #ffffff;
              font-size: 16px;
              font-weight: 500;
              display: flex;
              align-items: center;
              gap: 12px;
              
              .theme-tag {
                background: rgba(255, 255, 255, 0.2);
                padding: 4px 12px;
                border-radius: 100px;
                font-size: 14px;
                font-weight: normal;
              }
            }
            
            .result-actions {
              display: flex;
              gap: 16px;
              
              .el-button {
                color: #ffffff;
                font-size: 14px;
                
                .el-icon {
                  margin-right: 4px;
                }
                
                &:hover {
                  color: rgba(255, 255, 255, 0.8);
                }
              }
            }
          }
        }
        
        :deep(.el-card__body) {
          padding: 24px;
          
          .result-content {
            max-height: 600px;
            overflow-y: auto;
            line-height: 1.8;
            color: #1e293b;
            font-size: 15px;
            
            :deep(h1) {
              font-size: 24px;
              margin: 24px 0 16px;
              color: #1e293b;
            }
            
            :deep(h2) {
              font-size: 20px;
              margin: 20px 0 14px;
              color: #1e293b;
            }
            
            :deep(h3) {
              font-size: 18px;
              margin: 18px 0 12px;
              color: #1e293b;
            }
            
            :deep(p) {
              margin: 12px 0;
              line-height: 1.8;
            }
            
            :deep(ul), :deep(ol) {
              margin: 12px 0;
              padding-left: 24px;
            }
            
            :deep(blockquote) {
              margin: 16px 0;
              padding: 8px 16px;
              border-left: 4px solid #8b5cf6;
              background-color: #f8fafc;
            }
            
            :deep(pre) {
              background-color: #f1f5f9;
              padding: 16px;
              border-radius: 8px;
              margin: 16px 0;
              overflow-x: auto;
            }
            
            :deep(code) {
              font-family: 'Fira Code', monospace;
              background-color: #f1f5f9;
              padding: 2px 4px;
              border-radius: 4px;
              font-size: 90%;
            }
            
            :deep(a) {
              color: #8b5cf6;
              text-decoration: none;
              
              &:hover {
                text-decoration: underline;
              }
            }
            
            &::-webkit-scrollbar {
              width: 8px;
            }
            
            &::-webkit-scrollbar-track {
              background: #f1f5f9;
              border-radius: 4px;
            }
            
            &::-webkit-scrollbar-thumb {
              background: #cbd5e1;
              border-radius: 4px;
              
              &:hover {
                background: #94a3b8;
              }
            }
          }
        }
      }
    }
  }
}

@media (max-width: 1200px) {
  .writing-content {
    .template-section {
      .template-grid {
        grid-template-columns: repeat(2, 1fr);
      }
    }
  }
}

@media (max-width: 768px) {
  .writing-content {
    .template-section {
      .template-grid {
        grid-template-columns: 1fr;
      }
    }
    
    .writing-input-section {
      .theme-input {
        flex-direction: column;
      }
    }
  }
}
</style> 