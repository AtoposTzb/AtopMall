<template>
  <div class="message-list-page">
    <h2>留言管理</h2>

    <el-card>
      <el-table :data="messages" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="user_id" label="用户ID" width="100" />
        <el-table-column label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getTypeTag(row.type)" size="small" effect="plain">
              {{ getTypeText(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="subject" label="主题" min-width="160" show-overflow-tooltip />
        <el-table-column prop="message" label="内容" min-width="200">
          <template #default="{ row }">
            <span class="message-preview" @click="openDetail(row)">{{ row.message }}</span>
          </template>
        </el-table-column>
        <el-table-column label="附件" width="100" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.file" type="success" size="small" effect="plain">有附件</el-tag>
            <span v-else class="no-file">-</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="80" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="openDetail(row)">查看</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div v-if="!loading && messages.length === 0" class="empty-hint">
        暂无留言
      </div>
    </el-card>

    <!-- 留言详情弹窗 -->
    <el-dialog
      v-model="showDetail"
      title="留言详情"
      width="640px"
      :close-on-click-modal="true"
      destroy-on-close
    >
      <template v-if="currentMessage">
        <div class="detail-header">
          <el-tag :type="getTypeTag(currentMessage.type)" size="default" effect="plain">
            {{ getTypeText(currentMessage.type) }}
          </el-tag>
          <span class="detail-user">用户ID：{{ currentMessage.user_id }}</span>
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
import { ref, onMounted } from "vue";
import { getMessageList, type MessageItem } from "@/api/message";

const messages = ref<MessageItem[]>([]);
const loading = ref(false);
const showDetail = ref(false);
const currentMessage = ref<MessageItem | null>(null);

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

const downloadFile = (url: string) => {
  const a = document.createElement("a");
  a.href = url;
  a.download = getFileName(url);
  a.target = "_blank";
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
};

const openDetail = (row: MessageItem) => {
  currentMessage.value = row;
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

onMounted(() => {
  loadMessages();
});
</script>

<style lang="scss" scoped>
.message-list-page {
  h2 {
    font-size: 20px;
    font-weight: 600;
    margin-bottom: 20px;
    color: $text-primary;
  }
}

.message-preview {
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  cursor: pointer;
  color: $text-regular;
  transition: color 0.2s;

  &:hover {
    color: $primary-color;
  }
}

.no-file {
  color: #c0c4cc;
}

.empty-hint {
  text-align: center;
  padding: 60px 0;
  color: #c0c4cc;
  font-size: 14px;
}

// ========== 详情弹窗 ==========
.detail-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f0f0f0;

  .detail-user {
    font-size: 13px;
    color: $text-secondary;
  }
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