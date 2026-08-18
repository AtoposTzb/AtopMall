<template>
  <div class="user-message">
    <h3 class="page-title">我的留言</h3>

    <!-- 新建留言按钮 -->
    <div class="action-bar">
      <el-button type="primary" @click="showCreateDialog = true">
        <el-icon><Edit /></el-icon>
        写留言
      </el-button>
    </div>

    <!-- 留言列表 -->
    <div v-loading="loading" class="message-list">
      <div
        v-for="item in messages"
        :key="item.id"
        class="message-item card"
        @click="openDetail(item)"
      >
        <div class="message-header">
          <el-tag :type="getTypeTag(item.type)" size="small" effect="plain">
            {{ getTypeText(item.type) }}
          </el-tag>
          <h4 class="message-subject">{{ item.subject }}</h4>
          <el-icon class="arrow-icon"><ArrowRight /></el-icon>
        </div>
        <div class="message-body">
          <p>{{ item.message }}</p>
        </div>
        <div class="message-file" v-if="item.file">
          <el-icon><Link /></el-icon>
          <span>有附件</span>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-if="!loading && messages.length === 0" class="empty">
      <el-icon class="empty-icon"><ChatLineSquare /></el-icon>
      <p>暂无留言记录</p>
      <el-button type="primary" @click="showCreateDialog = true">写留言</el-button>
    </div>

    <!-- 新建留言弹窗 -->
    <el-dialog
      v-model="showCreateDialog"
      title="写留言"
      width="560px"
      :close-on-click-modal="false"
      destroy-on-close
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="80px"
        @submit.prevent
      >
        <el-form-item label="留言类型" prop="type">
          <el-select v-model="form.type" placeholder="请选择留言类型" style="width: 100%">
            <el-option label="留言" :value="1" />
            <el-option label="投诉" :value="2" />
            <el-option label="询问" :value="3" />
            <el-option label="售后" :value="4" />
            <el-option label="求购" :value="5" />
          </el-select>
        </el-form-item>
        <el-form-item label="主题" prop="subject">
          <el-input
            v-model="form.subject"
            placeholder="请输入留言主题"
            maxlength="100"
            show-word-limit
          />
        </el-form-item>
        <el-form-item label="内容" prop="message">
          <el-input
            v-model="form.message"
            type="textarea"
            :rows="5"
            placeholder="请输入留言内容"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>
        <el-form-item label="上传附件">
          <div class="upload-area">
            <input
              ref="fileInputRef"
              type="file"
              class="file-input-hidden"
              @change="handleFileChange"
            />
            <el-button
              :loading="uploading"
              :disabled="uploading"
              @click="fileInputRef?.click()"
            >
              <el-icon><Upload /></el-icon>
              {{ uploading ? '上传中...' : '选择文件' }}
            </el-button>
            <span class="upload-hint">支持图片、文档等格式，大小不超过 10MB</span>
          </div>
          <!-- 上传预览 -->
          <div class="upload-preview" v-if="form.file">
            <template v-if="isImageFile(form.file)">
              <div class="preview-image-box">
                <img :src="form.file" class="preview-image" />
                <button class="preview-remove" @click="removeFile">&times;</button>
              </div>
            </template>
            <template v-else>
              <div class="preview-file-box">
                <el-icon class="file-icon"><Document /></el-icon>
                <span class="file-name">{{ getFileName(form.file) }}</span>
                <el-button type="danger" link size="small" @click="removeFile">移除</el-button>
              </div>
            </template>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          提交留言
        </el-button>
      </template>
    </el-dialog>

    <!-- 留言详情弹窗 -->
    <el-dialog
      v-model="showDetail"
      title="留言详情"
      width="580px"
      :close-on-click-modal="true"
      destroy-on-close
    >
      <template v-if="currentMessage">
        <div class="detail-header">
          <el-tag :type="getTypeTag(currentMessage.type)" size="default" effect="plain">
            {{ getTypeText(currentMessage.type) }}
          </el-tag>
        </div>
        <div class="detail-subject">
          <label>主题</label>
          <p>{{ currentMessage.subject }}</p>
        </div>
        <div class="detail-body">
          <label>内容</label>
          <p>{{ currentMessage.message }}</p>
        </div>
        <div class="detail-file" v-if="currentMessage.file">
          <label>附件</label>
          <template v-if="isImageFile(currentMessage.file)">
            <div class="file-image-box">
              <img :src="currentMessage.file" class="file-image" />
            </div>
            <div class="file-actions">
              <el-button type="primary" size="small" @click="downloadFile(currentMessage.file)">
                <el-icon><Download /></el-icon>
                下载原图
              </el-button>
              <a :href="currentMessage.file" target="_blank" class="file-link">新窗口打开</a>
            </div>
          </template>
          <template v-else>
            <div class="file-other-box">
              <el-icon class="file-type-icon"><Document /></el-icon>
              <span class="file-name-text">{{ getFileName(currentMessage.file) }}</span>
            </div>
            <div class="file-actions">
              <el-button type="primary" size="small" @click="downloadFile(currentMessage.file)">
                <el-icon><Download /></el-icon>
                下载附件
              </el-button>
              <a :href="currentMessage.file" target="_blank" class="file-link">新窗口打开</a>
            </div>
          </template>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from "vue";
import { ElMessage, type FormInstance, type FormRules } from "element-plus";
import { getMessageList, createMessage, type MessageItem } from "@/api/message";
import { uploadFile } from "@/api/oss";

const messages = ref<MessageItem[]>([]);
const loading = ref(false);
const submitting = ref(false);
const uploading = ref(false);
const showCreateDialog = ref(false);
const showDetail = ref(false);
const currentMessage = ref<MessageItem | null>(null);
const formRef = ref<FormInstance>();
const fileInputRef = ref<HTMLInputElement>();

const form = reactive({
  type: 1,
  subject: "",
  message: "",
  file: "",
});

const rules: FormRules = {
  type: [{ required: true, message: "请选择留言类型", trigger: "change" }],
  subject: [
    { required: true, message: "请输入留言主题", trigger: "blur" },
    { min: 2, max: 100, message: "主题长度在 2 到 100 个字符", trigger: "blur" },
  ],
  message: [
    { required: true, message: "请输入留言内容", trigger: "blur" },
    { min: 5, max: 500, message: "内容长度在 5 到 500 个字符", trigger: "blur" },
  ],
};

const getTypeText = (type: number): string => {
  const map: Record<number, string> = {
    1: "留言",
    2: "投诉",
    3: "询问",
    4: "售后",
    5: "求购",
  };
  return map[type] || "留言";
};

const getTypeTag = (type: number): string => {
  const map: Record<number, string> = {
    1: "",
    2: "danger",
    3: "warning",
    4: "info",
    5: "success",
  };
  return map[type] || "";
};

const isImageFile = (url: string): boolean => {
  return /\.(jpg|jpeg|png|gif|webp|bmp|svg)(\?|$)/i.test(url);
};

const getFileName = (url: string): string => {
  const parts = url.split("/");
  const last = parts[parts.length - 1];
  return decodeURIComponent(last.split("?")[0]) || "附件";
};

const handleFileChange = async (event: Event) => {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  if (!file) return;

  const maxSize = 10 * 1024 * 1024;
  if (file.size > maxSize) {
    ElMessage.warning("文件大小不能超过 10MB");
    return;
  }

  uploading.value = true;
  try {
    const url = await uploadFile(file);
    form.file = url;
    ElMessage.success("文件上传成功");
  } catch {
    ElMessage.error("文件上传失败，请重试");
  } finally {
    uploading.value = false;
    if (target) target.value = "";
  }
};

const removeFile = () => {
  form.file = "";
  if (fileInputRef.value) {
    fileInputRef.value.value = "";
  }
};

const downloadFile = (url: string) => {
  const a = document.createElement("a");
  a.href = url;
  a.download = getFileName(url);
  a.target = "_blank";
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
};

const openDetail = (item: MessageItem) => {
  currentMessage.value = item;
  showDetail.value = true;
};

const loadMessages = async () => {
  loading.value = true;
  try {
    const res = await getMessageList();
    messages.value = (res as any).data || [];
  } catch {
    // 错误已在拦截器中处理
  } finally {
    loading.value = false;
  }
};

const handleSubmit = async () => {
  if (!formRef.value) return;
  const valid = await formRef.value.validate().catch(() => false);
  if (!valid) return;

  submitting.value = true;
  try {
    await createMessage({
      type: form.type,
      subject: form.subject,
      message: form.message,
      file: form.file || undefined,
    });
    ElMessage.success("留言提交成功");
    showCreateDialog.value = false;
    form.type = 1;
    form.subject = "";
    form.message = "";
    form.file = "";
    await loadMessages();
  } catch {
    // 错误已在拦截器中处理
  } finally {
    submitting.value = false;
  }
};

onMounted(() => {
  loadMessages();
});
</script>

<style lang="scss" scoped>
.user-message {
  .page-title {
    font-size: 20px;
    font-weight: 600;
    color: $text-primary;
    margin-bottom: 28px;
    padding-bottom: 16px;
    border-bottom: 2px solid #f0f5ff;
    position: relative;

    &::after {
      content: "";
      position: absolute;
      bottom: -2px;
      left: 0;
      width: 60px;
      height: 2px;
      background: $primary-color;
      border-radius: 1px;
    }
  }
}

.action-bar {
  margin-bottom: 20px;
}

.message-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.message-item {
  padding: 20px;
  border-radius: 12px;
  border: 1px solid #f0f0f0;
  transition: all 0.3s;
  cursor: pointer;

  &:hover {
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
    border-color: transparent;
    transform: translateY(-2px);
  }

  .message-header {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 12px;

    .message-subject {
      font-size: 16px;
      font-weight: 600;
      color: $text-primary;
      margin: 0;
      flex: 1;
    }

    .arrow-icon {
      color: #c0c4cc;
      font-size: 14px;
      transition: color 0.2s;
    }
  }

  &:hover .arrow-icon {
    color: $primary-color;
  }

  .message-body {
    p {
      font-size: 14px;
      color: $text-regular;
      line-height: 1.8;
      margin: 0;
    }
  }

  .message-file {
    margin-top: 12px;
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 13px;
    color: $primary-color;

    a {
      color: $primary-color;
      text-decoration: none;

      &:hover {
        text-decoration: underline;
      }
    }
  }
}

.empty {
  text-align: center;
  padding: 80px 0;
  color: $text-secondary;

  .empty-icon {
    font-size: 56px;
    margin-bottom: 16px;
    color: #dcdfe6;
  }

  p {
    font-size: 15px;
    margin-bottom: 20px;
  }
}

// ========== 上传区域 ==========
.upload-area {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;

  .file-input-hidden {
    display: none;
  }

  .upload-hint {
    font-size: 12px;
    color: #999;
  }
}

.upload-preview {
  margin-top: 12px;
}

.preview-image-box {
  position: relative;
  display: inline-block;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #e8e8e8;

  .preview-image {
    max-width: 200px;
    max-height: 150px;
    object-fit: cover;
    display: block;
  }

  .preview-remove {
    position: absolute;
    top: 4px;
    right: 4px;
    width: 22px;
    height: 22px;
    border-radius: 50%;
    border: none;
    background: rgba(0, 0, 0, 0.55);
    color: #fff;
    font-size: 14px;
    line-height: 1;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;

    &:hover {
      background: rgba(0, 0, 0, 0.75);
    }
  }
}

.preview-file-box {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: #f5f7fa;
  border-radius: 6px;
  font-size: 13px;

  .file-icon {
    font-size: 20px;
    color: $primary-color;
  }

  .file-name {
    color: $text-regular;
    max-width: 200px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}

// ========== 详情弹窗 ==========
.detail-header {
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.detail-subject,
.detail-body,
.detail-file {
  margin-bottom: 20px;

  label {
    display: block;
    font-size: 13px;
    font-weight: 600;
    color: $text-secondary;
    margin-bottom: 8px;
  }
}

.detail-subject p {
  font-size: 16px;
  font-weight: 500;
  color: $text-primary;
  margin: 0;
}

.detail-body p {
  font-size: 14px;
  color: $text-regular;
  line-height: 1.8;
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
}

.file-image-box {
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #f0f0f0;
  background: #fafafa;
  display: inline-block;

  .file-image {
    max-width: 100%;
    max-height: 400px;
    display: block;
    object-fit: contain;
  }
}

.file-other-box {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 16px;
  background: #f5f7fa;
  border-radius: 8px;
  border: 1px solid #f0f0f0;

  .file-type-icon {
    font-size: 28px;
    color: $primary-color;
  }

  .file-name-text {
    font-size: 14px;
    color: $text-regular;
    word-break: break-all;
  }
}

.file-actions {
  margin-top: 12px;
  display: flex;
  align-items: center;
  gap: 12px;
}

.file-link {
  font-size: 13px;
  color: $primary-color;
  text-decoration: none;

  &:hover {
    text-decoration: underline;
  }
}
</style>