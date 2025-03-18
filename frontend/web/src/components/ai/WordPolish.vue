<template>
  <div class="ai-document-container">
    <div class="document-header">
      <div class="title-area">
        <el-icon class="header-icon"><Document /></el-icon>
        <h2>AI 文档</h2>
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

    <div class="document-content">
      <div class="document-intro">
        <h3>让AI解读文档、论文、网页，轻松又专业</h3>
        <p>上传文档获取AI优化、解读与分析</p>
      </div>

      <div class="upload-section">
        <el-upload
          class="upload-area"
          drag
          action="#"
          :auto-upload="false"
          :on-change="handleFileChange"
          accept=".doc,.docx,.pdf,.txt"
          :show-file-list="false"
        >
          <div class="upload-inner">
            <el-icon class="upload-icon"><upload-filled /></el-icon>
            <div class="upload-text">
              拖拽文件到此处或 <em>点击上传</em>
            </div>
            <div class="upload-tip">
              支持 .doc, .docx, .pdf, .txt 格式的文档
            </div>
          </div>
        </el-upload>
      </div>

      <div class="features-section">
        <div class="features-title">
          <h3>使用场景</h3>
        </div>
        <div class="feature-cards">
          <div class="feature-card">
            <el-icon class="feature-icon"><Reading /></el-icon>
            <div class="feature-text">
              <h4>学术论文解读</h4>
              <p>帮你深入了解和把握论文和学术文章</p>
            </div>
          </div>
          <div class="feature-card">
            <el-icon class="feature-icon"><DataAnalysis /></el-icon>
            <div class="feature-text">
              <h4>数据报表解读</h4>
              <p>轻松获取数据报表洞察，效率加倍</p>
            </div>
          </div>
          <div class="feature-card">
            <el-icon class="feature-icon"><Connection /></el-icon>
            <div class="feature-text">
              <h4>网页内容解析</h4>
              <p>深入解读网页链接，一键洞悉核心内容</p>
            </div>
          </div>
        </div>
      </div>

      <div class="document-processing" v-if="fileContent">
        <el-card class="process-card">
          <template #header>
            <div class="card-header">
              <div class="file-info">
                <el-icon><Document /></el-icon>
                <span class="file-name">{{ currentFileName }}</span>
              </div>
              <div class="card-actions">
                <el-button type="text" @click="resetForm">
                  <el-icon><Delete /></el-icon>
                </el-button>
              </div>
            </div>
          </template>

          <div class="card-tabs">
            <el-tabs v-model="activeTab">
              <el-tab-pane label="优化配置" name="options">
                <el-form :model="polishOptions" label-width="100px" class="options-form">
                  <el-form-item label="优化类型">
                    <el-radio-group v-model="polishOptions.type">
                      <el-radio-button value="polish">文档优化</el-radio-button>
                      <el-radio-button value="summarize">内容摘要</el-radio-button>
                      <el-radio-button value="analyze">深度分析</el-radio-button>
                    </el-radio-group>
                  </el-form-item>

                  <template v-if="polishOptions.type === 'polish'">
                    <el-form-item label="优化风格">
                      <el-select v-model="polishOptions.style" placeholder="选择优化风格">
                        <el-option label="学术论文" value="academic" />
                        <el-option label="商务文档" value="business" />
                        <el-option label="创意写作" value="creative" />
                        <el-option label="技术文档" value="technical" />
                        <el-option label="自定义" value="custom" />
                      </el-select>
                    </el-form-item>

                    <el-form-item v-if="polishOptions.style === 'custom'" label="自定义风格">
                      <el-input
                        v-model="polishOptions.customStyle"
                        type="textarea"
                        :rows="2"
                        placeholder="请输入自定义优化风格要求"
                      />
                    </el-form-item>

                    <el-form-item label="优化重点">
                      <el-checkbox-group v-model="polishOptions.focus">
                        <el-checkbox value="grammar">语法优化</el-checkbox>
                        <el-checkbox value="vocabulary">词汇提升</el-checkbox>
                        <el-checkbox value="structure">结构优化</el-checkbox>
                        <el-checkbox value="tone">语气调整</el-checkbox>
                      </el-checkbox-group>
                    </el-form-item>

                    <el-form-item label="语言风格">
                      <el-radio-group v-model="polishOptions.language">
                        <el-radio value="formal">正式</el-radio>
                        <el-radio value="casual">随意</el-radio>
                        <el-radio value="neutral">中性</el-radio>
                      </el-radio-group>
                    </el-form-item>
                  </template>

                  <template v-if="polishOptions.type === 'summarize'">
                    <el-form-item label="摘要长度">
                      <el-slider v-model="polishOptions.summaryLength" :min="1" :max="5" show-stops />
                    </el-form-item>
                    <el-form-item label="关注点">
                      <el-checkbox-group v-model="polishOptions.summaryFocus">
                        <el-checkbox value="main">主要观点</el-checkbox>
                        <el-checkbox value="arguments">论证过程</el-checkbox>
                        <el-checkbox value="conclusion">结论</el-checkbox>
                        <el-checkbox value="methodology">方法论</el-checkbox>
                      </el-checkbox-group>
                    </el-form-item>
                  </template>

                  <template v-if="polishOptions.type === 'analyze'">
                    <el-form-item label="分析维度">
                      <el-checkbox-group v-model="polishOptions.analysisDimensions">
                        <el-checkbox value="structure">文档结构</el-checkbox>
                        <el-checkbox value="arguments">论点有效性</el-checkbox>
                        <el-checkbox value="evidence">证据支持</el-checkbox>
                        <el-checkbox value="logic">逻辑连贯性</el-checkbox>
                        <el-checkbox value="language">语言表达</el-checkbox>
                      </el-checkbox-group>
                    </el-form-item>
                    <el-form-item label="分析深度">
                      <el-radio-group v-model="polishOptions.analysisDepth">
                        <el-radio value="basic">基础</el-radio>
                        <el-radio value="detailed">详细</el-radio>
                        <el-radio value="expert">专家</el-radio>
                      </el-radio-group>
                    </el-form-item>
                  </template>
                </el-form>

                <div class="form-actions">
                  <el-button 
                    type="primary" 
                    :disabled="!fileContent || isPolishing"
                    @click="startPolish"
                    :loading="isPolishing"
                  >
                    {{ getActionButtonText() }}
                  </el-button>
                  <el-button @click="resetForm" :disabled="isPolishing">重置</el-button>
                </div>
              </el-tab-pane>

              <el-tab-pane label="原文内容" name="original" v-if="fileContent">
                <div class="preview-box" v-html="fileContent"></div>
              </el-tab-pane>

              <el-tab-pane label="处理结果" name="result" v-if="polishedContent">
                <div class="result-content" v-html="polishedContent"></div>
                <div class="result-actions">
                  <el-button type="primary" @click="downloadResult">
                    <el-icon><Download /></el-icon>
                    下载处理结果
                  </el-button>
                  <el-button @click="copyToClipboard">
                    <el-icon><CopyDocument /></el-icon>
                    复制内容
                  </el-button>
                </div>
              </el-tab-pane>
            </el-tabs>
          </div>
        </el-card>
      </div>

      <div class="sample-documents" v-if="!fileContent">
        <div class="section-title">
          <h3>试试解读</h3>
        </div>
        <div class="document-examples">
          <div class="example-row">
            <el-card class="example-card" shadow="hover">
              <div class="example-icon pdf-icon">
                <el-icon><Document /></el-icon>
              </div>
              <div class="example-info">
                <div class="example-name">论文样例.pdf</div>
                <div class="example-size">2.7MB</div>
              </div>
            </el-card>
            <el-card class="example-card" shadow="hover">
              <div class="example-icon excel-icon">
                <el-icon><Grid /></el-icon>
              </div>
              <div class="example-info">
                <div class="example-name">数据报表.xlsx</div>
                <div class="example-size">1.5MB</div>
              </div>
            </el-card>
            <el-card class="example-card" shadow="hover">
              <div class="example-icon web-icon">
                <el-icon><Link /></el-icon>
              </div>
              <div class="example-info">
                <div class="example-name">网页链接</div>
                <div class="example-url">https://example.com/article</div>
              </div>
            </el-card>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, watch, onMounted, computed } from 'vue'
import { 
  Document, UploadFilled, Download, Delete, Reading, 
  DataAnalysis, Connection, Grid, Link, CopyDocument 
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { useAIStore } from '@/stores/ai'
import { useUserStore } from '@/stores/user'
import { marked } from 'marked'
import mammoth from 'mammoth'
import { Document as DocxDocument, Packer, Paragraph, TextRun } from 'docx'
import { saveAs } from 'file-saver'

const aiStore = useAIStore()
const userStore = useUserStore()
const fileContent = ref('')
const polishedContent = ref('')
const isPolishing = ref(false)
const selectedModel = ref('')
const currentFileName = ref('')
const activeTab = ref('options')

const polishOptions = reactive({
  type: 'polish',
  style: 'academic',
  focus: ['grammar', 'vocabulary'],
  language: 'formal',
  customStyle: '',
  summaryLength: 3,
  summaryFocus: ['main', 'conclusion'],
  analysisDimensions: ['structure', 'arguments', 'evidence'],
  analysisDepth: 'detailed'
})

// 监听 AI 存储的初始化状态
watch(() => aiStore.initialized, (initialized) => {
  // console.log('WordPolish: AI存储初始化状态变化:', initialized)
  if (initialized) {
    selectedModel.value = aiStore.currentModel
    // console.log('WordPolish: 从AI存储加载当前模型:', selectedModel.value)
  }
})

// 监听 AI 存储的当前模型变化
watch(() => aiStore.currentModel, (newModel) => {
  if (newModel && newModel !== selectedModel.value) {
    // console.log('WordPolish: AI存储当前模型变化:', newModel)
    selectedModel.value = newModel
  }
})

// 监听用户登录状态变化
watch(() => userStore.isAuthenticated, async (isAuthenticated) => {
  console.log('WordPolish: 用户认证状态变化:', isAuthenticated)
  if (isAuthenticated && !aiStore.initialized) {
    await aiStore.initialize()
  }
})

// 初始化时设置默认模型
onMounted(async () => {
  // 如果 AI 存储已初始化，使用其默认模型
  if (aiStore.initialized) {
    selectedModel.value = aiStore.currentModel
    // console.log('WordPolish: 已从 AI 存储加载当前模型:', selectedModel.value)
  } else if (userStore.isAuthenticated) {
    // 如果用户已登录但 AI 存储未初始化，初始化 AI 存储
    if (!aiStore.isInitializing) { // 检查是否正在初始化
      if (!aiStore.modelsLoaded) {
        // console.log('WordPolish: 用户已认证，AI存储未初始化，开始初始化')
        await aiStore.initialize()
      } else {
        // console.log('WordPolish: 用户已认证，模型已加载，仅加载默认模型')
        await aiStore.loadDefaultModel()
        aiStore.initialized = true
      }
      selectedModel.value = aiStore.currentModel
      // console.log('WordPolish: 已初始化 AI 存储并加载当前模型:', selectedModel.value)
    } else {
      // console.log('WordPolish: AI存储正在初始化中，跳过重复初始化')
    }
  } else {
    // 如果用户未登录，加载可用模型
    if (!aiStore.modelsLoaded) {
      // console.log('WordPolish: 用户未认证，加载可用模型')
      await aiStore.loadAvailableModels()
      aiStore.modelsLoaded = true
    } else {
      // console.log('WordPolish: 用户未认证，模型已加载，跳过重复加载')
    }
    if (aiStore.availableModels.length > 0) {
      selectedModel.value = aiStore.availableModels[0].id
      // console.log('WordPolish: 已加载第一个可用模型:', selectedModel.value)
    }
  }
  
  if (selectedModel.value) {
    aiStore.setCurrentModel(selectedModel.value)
  }
})

const handleFileChange = async (file) => {
  if (!file) return
  
  try {
    const arrayBuffer = await file.raw.arrayBuffer()
    const result = await mammoth.convertToHtml({ arrayBuffer })
    fileContent.value = result.value
    polishedContent.value = ''
    currentFileName.value = file.name
  } catch (error) {
    console.error('Error reading file:', error)
    ElMessage.error('文件读取失败，请确保上传的是有效的Word文档')
    fileContent.value = ''
    polishedContent.value = ''
  }
}

const cancelPolish = async () => {
  try {
    await aiStore.cancelCurrentRequest()
    isPolishing.value = false
    ElMessage.success('已取消优化')
  } catch (error) {
    console.error('Cancel error:', error)
    ElMessage.error('取消失败，请重试')
  }
}

const startPolish = async () => {
  if (!fileContent.value) {
    ElMessage.warning('请先上传文档')
    return
  }

  const content = fileContent.value.trim()
  if (!content) {
    ElMessage.warning('文档内容为空，请重新上传')
    return
  }

  // 检查是否选择了模型
  if (!aiStore.currentModel) {
    ElMessage.warning('请先选择AI模型')
    return
  }

  // 检查用户是否已登录
  if (!userStore.isAuthenticated) {
    ElMessage.warning('请先登录后再使用优化功能')
    return
  }

  isPolishing.value = true
  polishedContent.value = '' // 清空之前的结果
  
  try {
    const prompt = generatePolishPrompt()
    if (!prompt || !prompt.trim()) {
      ElMessage.error('生成提示词失败，请重试')
      return
    }

    // 使用流式响应
    await aiStore.sendMessage(prompt.trim(), -1, (response) => {
      // 确保content是字符串
      const content = typeof response === 'object' ? response.fullContent || '' : response
      // 实时更新优化结果
      polishedContent.value = marked(content)
    })
  } catch (error) {
    if (error.name === 'AbortError') {
      return
    }
    ElMessage.error('优化失败，请重试')
  } finally {
    isPolishing.value = false
  }
}

const generatePolishPrompt = () => {
  const content = fileContent.value.trim() // 确保内容不为空
  if (!content) {
    ElMessage.error('文档内容为空，请重新上传')
    return null
  }

  // 处理文档优化类型
  if (polishOptions.type === 'polish') {
    const styleMap = {
      academic: '学术论文',
      business: '商务文档',
      creative: '创意写作',
      technical: '技术文档',
      custom: polishOptions.customStyle || '自定义风格'
    }

    const focusMap = {
      grammar: '语法优化',
      vocabulary: '词汇提升',
      structure: '结构优化',
      tone: '语气调整'
    }

    const languageMap = {
      formal: '正式',
      casual: '随意',
      neutral: '中性'
    }

    // 验证自定义风格
    if (polishOptions.style === 'custom' && !polishOptions.customStyle.trim()) {
      ElMessage.error('请输入自定义优化风格要求')
      return null
    }

    return `请对以下文本进行优化，要求：
1. 优化风格：${styleMap[polishOptions.style]}
2. 优化重点：${polishOptions.focus.map(f => focusMap[f]).join('、')}
3. 语言风格：${languageMap[polishOptions.language]}
4. 保持原文的核心意思不变
5. 输出优化后的完整文本，保持原有的格式（包括标题、段落、列表等）
6. 使用Markdown格式输出，包括：
   - 使用 # 表示标题
   - 使用 * 或 _ 表示强调
   - 使用 > 表示引用
   - 使用 - 或 1. 表示列表
   - 使用 | 表示表格
7. 请逐段输出优化结果，每段之间保持适当的空行

原文：
${content}`
  }
  
  // 处理摘要生成
  else if (polishOptions.type === 'summarize') {
    const lengthMap = {
      1: '非常简短（约100字）',
      2: '简短（约200字）',
      3: '中等（约400字）',
      4: '详细（约800字）',
      5: '非常详细（约1500字）'
    }
    
    const focusMap = {
      main: '主要观点',
      arguments: '论证过程',
      conclusion: '结论',
      methodology: '方法论'
    }
    
    return `请对以下文本生成摘要，要求：
1. 摘要长度：${lengthMap[polishOptions.summaryLength]}
2. 关注点：${polishOptions.summaryFocus.map(f => focusMap[f]).join('、')}
3. 保持原文的核心意思不变，提取最重要的信息
4. 使用Markdown格式输出，包括：
   - 使用 # 表示摘要标题
   - 使用 ## 表示各部分小标题（如有必要）
   - 使用 * 或 _ 表示强调重点
   - 使用 > 表示引用原文中的关键句
5. 摘要结构应包含：
   - 开头简述文档主旨
   - 中间部分概括核心内容
   - 结尾总结主要结论或观点

原文：
${content}`
  }
  
  // 处理深度分析
  else if (polishOptions.type === 'analyze') {
    const dimensionsMap = {
      structure: '文档结构',
      arguments: '论点有效性',
      evidence: '证据支持',
      logic: '逻辑连贯性',
      language: '语言表达'
    }
    
    const depthMap = {
      basic: '基础（概括性分析）',
      detailed: '详细（深入分析主要方面）',
      expert: '专家（全面且批判性地分析）'
    }
    
    return `请对以下文本进行深度分析，要求：
1. 分析维度：${polishOptions.analysisDimensions.map(d => dimensionsMap[d]).join('、')}
2. 分析深度：${depthMap[polishOptions.analysisDepth]}
3. 提供客观、专业的分析评价，指出优缺点
4. 使用Markdown格式输出，包括：
   - 使用 # 表示分析标题
   - 使用 ## 表示各个分析维度的小标题
   - 使用 * 或 _ 表示强调重点
   - 使用 > 表示引用原文中的关键句作为分析依据
5. 分析结构应包含：
   - 开头简述文档整体评价
   - 中间按照每个分析维度分节进行详细分析
   - 结尾提供改进建议或总体评价
   - 可选：如有必要，使用评分（1-5分）表示各维度的表现

原文：
${content}`
  }
  
  // 默认返回优化模式的提示词
  else {
    ElMessage.warning('未选择处理类型，默认使用文档优化')
    return `请对以下文本进行语言优化，保持原文意思不变，使表达更加清晰流畅。

原文：
${content}`
  }
}

const downloadResult = async () => {
  try {
    // 将 HTML 转换为纯文本
    const tempDiv = document.createElement('div')
    tempDiv.innerHTML = polishedContent.value
    const textContent = tempDiv.textContent || tempDiv.innerText

    // 创建 Word 文档
    const doc = new DocxDocument({
      sections: [{
        properties: {},
        children: textContent.split('\n').map(line => {
          // 处理标题
          if (line.startsWith('# ')) {
            return new Paragraph({
              text: line.replace('# ', ''),
              heading: HeadingLevel.HEADING_1,
              spacing: { after: 200 }
            })
          } else if (line.startsWith('## ')) {
            return new Paragraph({
              text: line.replace('## ', ''),
              heading: HeadingLevel.HEADING_2,
              spacing: { after: 200 }
            })
          } else if (line.startsWith('### ')) {
            return new Paragraph({
              text: line.replace('### ', ''),
              heading: HeadingLevel.HEADING_3,
              spacing: { after: 200 }
            })
          }
          // 处理普通段落
          else if (line.trim()) {
            return new Paragraph({
              children: [
                new TextRun({
                  text: line,
                  size: 24 // 12pt
                })
              ],
              spacing: { after: 200 }
            })
          }
          // 处理空行
          else {
            return new Paragraph({
              spacing: { after: 200 }
            })
          }
        })
      }]
    })

    // 生成文档并下载
    const blob = await Packer.toBlob(doc)
    saveAs(blob, currentFileName.value)
    ElMessage.success('文档下载成功')
  } catch (error) {
    console.error('下载文档失败:', error)
    ElMessage.error('下载文档失败，请重试')
  }
}

const copyToClipboard = async () => {
  try {
    const markdownText = polishedContent.value.replace(/<[^>]+>/g, '')
    await navigator.clipboard.writeText(markdownText)
    ElMessage.success('已复制到剪贴板')
  } catch (error) {
    ElMessage.error('复制失败，请手动复制')
  }
}

const resetForm = () => {
  fileContent.value = ''
  polishedContent.value = ''
  polishOptions.style = 'academic'
  polishOptions.focus = ['grammar', 'vocabulary']
  polishOptions.language = 'formal'
  polishOptions.customStyle = ''
  currentFileName.value = ''
}

// 根据当前选择的处理类型，返回按钮文本
const getActionButtonText = () => {
  switch(polishOptions.type) {
    case 'polish': return '开始优化';
    case 'summarize': return '生成摘要';
    case 'analyze': return '开始分析';
    default: return '开始处理';
  }
}
</script>

<style lang="scss" scoped>
.ai-document-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: #f8fafc;
  
  .document-header {
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
  
  .document-content {
    flex: 1;
    padding: 24px;
    overflow-y: auto;
    
    .document-intro {
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
    
    .upload-section {
      margin-bottom: 40px;
      
      .upload-area {
        background: #ffffff;
        border: 2px dashed #e2e8f0;
        border-radius: 16px;
        transition: all 0.3s ease;
        
        &:hover {
          border-color: #8b5cf6;
        }
        
        .upload-inner {
          padding: 60px 0;
          text-align: center;
          
          .upload-icon {
            font-size: 64px;
            color: #8b5cf6;
            margin-bottom: 24px;
          }
          
          .upload-text {
            color: #475569;
            font-size: 18px;
            margin-bottom: 12px;
            
            em {
              color: #8b5cf6;
              font-style: normal;
              font-weight: 500;
            }
          }
          
          .upload-tip {
            color: #94a3b8;
            font-size: 15px;
          }
        }
      }
    }
    
    .features-section {
      margin-bottom: 40px;
      
      .features-title {
        margin-bottom: 24px;
        
        h3 {
          font-size: 20px;
          font-weight: 600;
          color: #1e293b;
        }
      }
      
      .feature-cards {
        display: flex;
        gap: 24px;
        justify-content: space-between;
        
        .feature-card {
          flex: 1;
          background: #ffffff;
          border-radius: 16px;
          padding: 24px;
          display: flex;
          flex-direction: column;
          align-items: center;
          text-align: center;
          box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
          transition: transform 0.3s ease, box-shadow 0.3s ease;
          
          &:hover {
            transform: translateY(-5px);
            box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
          }
          
          .feature-icon {
            font-size: 40px;
            color: #8b5cf6;
            margin-bottom: 16px;
          }
          
          .feature-text {
            h4 {
              font-size: 18px;
              font-weight: 600;
              color: #1e293b;
              margin-bottom: 8px;
            }
            
            p {
              font-size: 14px;
              color: #64748b;
              line-height: 1.5;
            }
          }
        }
      }
    }
    
    .document-processing {
      .process-card {
        margin-bottom: 40px;
        border-radius: 16px;
        overflow: hidden;
        box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
        
        :deep(.el-card__header) {
          padding: 16px 24px;
          background: linear-gradient(135deg, #8b5cf6 0%, #6366f1 100%);
          
          .card-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            
            .file-info {
              display: flex;
              align-items: center;
              gap: 8px;
              
              .el-icon {
                color: #ffffff;
                font-size: 20px;
              }
              
              .file-name {
                color: #ffffff;
                font-size: 16px;
                font-weight: 500;
              }
            }
            
            .card-actions {
              .el-button {
                color: #ffffff;
                
                &:hover {
                  color: rgba(255, 255, 255, 0.8);
                }
              }
            }
          }
        }
        
        :deep(.el-card__body) {
          padding: 0;
          
          .card-tabs {
            .el-tabs__header {
              margin: 0;
              padding: 0 16px;
              background-color: #f8fafc;
              border-bottom: 1px solid #e2e8f0;
            }
            
            .el-tabs__content {
              padding: 24px;
            }
            
            .el-tabs__nav-wrap::after {
              background-color: transparent;
            }
            
            .el-tabs__item {
              height: 50px;
              line-height: 50px;
              color: #64748b;
              
              &.is-active {
                color: #8b5cf6;
                font-weight: 500;
              }
            }
            
            .options-form {
              max-width: 600px;
              margin: 0 auto;
              
              .el-form-item {
                margin-bottom: 20px;
              }
            }
            
            .form-actions {
              margin-top: 32px;
              display: flex;
              justify-content: center;
              gap: 16px;
              
              .el-button {
                padding: 12px 24px;
                font-size: 16px;
                min-width: 120px;
                
                &--primary {
                  background: linear-gradient(135deg, #8b5cf6 0%, #6366f1 100%);
                  border: none;
                  
                  &:hover:not(:disabled) {
                    transform: translateY(-1px);
                    box-shadow: 0 4px 8px rgba(139, 92, 246, 0.2);
                  }
                }
              }
            }
            
            .preview-box, .result-content {
              background-color: #ffffff;
              border: 1px solid #e2e8f0;
              border-radius: 8px;
              padding: 20px;
              min-height: 300px;
              color: #1e293b;
              font-size: 15px;
              line-height: 1.6;
            }
            
            .result-actions {
              margin-top: 24px;
              display: flex;
              justify-content: flex-end;
              gap: 16px;
              
              .el-button {
                &--primary {
                  background: linear-gradient(135deg, #8b5cf6 0%, #6366f1 100%);
                  border: none;
                  
                  &:hover {
                    transform: translateY(-1px);
                    box-shadow: 0 4px 8px rgba(139, 92, 246, 0.2);
                  }
                }
              }
            }
          }
        }
      }
    }
    
    .sample-documents {
      margin-top: 40px;
      
      .section-title {
        margin-bottom: 24px;
        
        h3 {
          font-size: 20px;
          font-weight: 600;
          color: #1e293b;
        }
      }
      
      .document-examples {
        .example-row {
          display: flex;
          gap: 24px;
          
          .example-card {
            flex: 1;
            cursor: pointer;
            transition: transform 0.3s ease, box-shadow 0.3s ease;
            border-radius: 12px;
            
            &:hover {
              transform: translateY(-5px);
              box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
            }
            
            :deep(.el-card__body) {
              padding: 20px;
              display: flex;
              align-items: center;
              gap: 16px;
              
              .example-icon {
                width: 48px;
                height: 48px;
                display: flex;
                align-items: center;
                justify-content: center;
                border-radius: 12px;
                color: #ffffff;
                font-size: 24px;
                
                &.pdf-icon {
                  background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
                }
                
                &.excel-icon {
                  background: linear-gradient(135deg, #22c55e 0%, #16a34a 100%);
                }
                
                &.web-icon {
                  background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
                }
              }
              
              .example-info {
                flex: 1;
                
                .example-name {
                  font-size: 16px;
                  font-weight: 500;
                  color: #1e293b;
                  margin-bottom: 4px;
                }
                
                .example-size, .example-url {
                  font-size: 13px;
                  color: #64748b;
                }
              }
            }
          }
        }
      }
    }
  }
}

@media (max-width: 768px) {
  .document-content {
    .feature-cards {
      flex-direction: column;
    }
    
    .document-examples {
      .example-row {
        flex-direction: column;
      }
    }
  }
}
</style> 