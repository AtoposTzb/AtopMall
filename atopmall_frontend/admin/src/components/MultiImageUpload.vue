<template>
  <div class="multi-image-upload">
    <div class="upload-area">
      <el-upload
        :http-request="handleUpload"
        :show-file-list="false"
        :before-upload="beforeUpload"
        accept="image/*"
        :disabled="uploading"
      >
        <el-button :loading="uploading" :disabled="disabled">
          {{ uploading ? '上传中...' : '点击上传' }}
        </el-button>
      </el-upload>
      <span class="upload-tip">支持 jpg、png、gif、bmp，单张不超过 10MB，可多张上传</span>
    </div>
    <div v-if="modelValue.length" class="preview-list">
      <div v-for="(url, index) in modelValue" :key="index" class="preview-item">
        <el-image :src="url" fit="cover" style="width: 120px; height: 80px;" />
        <el-button type="danger" link size="small" @click="handleRemove(index)">移除</el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { uploadFile } from '@/api/oss'

const props = defineProps<{
  modelValue: string[]
  disabled?: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', val: string[]): void
}>()

const uploading = ref(false)

const beforeUpload = (file: File) => {
  const isImage = file.type.startsWith('image/')
  if (!isImage) {
    ElMessage.error('只能上传图片文件')
    return false
  }
  const isLt10M = file.size / 1024 / 1024 < 10
  if (!isLt10M) {
    ElMessage.error('图片大小不能超过 10MB')
    return false
  }
  return true
}

const handleUpload = async (options: any) => {
  uploading.value = true
  try {
    const res = await uploadFile(options.file) as any
    if (res.code === 200) {
      emit('update:modelValue', [...props.modelValue, res.data.url])
      ElMessage.success('上传成功')
    } else {
      ElMessage.error(res.msg || '上传失败')
    }
  } catch (e: any) {
    ElMessage.error(e?.message || '上传失败，请检查网络连接')
  } finally {
    uploading.value = false
  }
}

const handleRemove = (index: number) => {
  const list = [...props.modelValue]
  list.splice(index, 1)
  emit('update:modelValue', list)
}
</script>

<style lang="scss" scoped>
.multi-image-upload {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.upload-area {
  display: flex;
  align-items: center;
  gap: 12px;
}

.upload-tip {
  font-size: 12px;
  color: #999;
}

.preview-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.preview-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}
</style>