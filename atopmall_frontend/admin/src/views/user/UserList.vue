<template>
  <div class="user-list-page">
    <h2>用户管理</h2>

    <el-card class="mb-md">
      <el-form :inline="true">
        <el-form-item label="搜索">
          <el-input v-model="searchKey" placeholder="用户名或手机号" clearable @keyup.enter="loadUsers" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadUsers">搜索</el-button>
          <el-button @click="searchKey = ''; loadUsers()">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card>
      <el-table :data="users" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="用户名" width="150" />
        <el-table-column prop="nick_name" label="昵称" width="150" />
        <el-table-column prop="mobile" label="手机号" width="150" />
        <el-table-column label="性别" width="80">
          <template #default="{ row }">
            <el-tag :type="row.gender === 'male' ? 'primary' : 'danger'" size="small">
              {{ row.gender === 'male' ? '男' : row.gender === 'female' ? '女' : '-' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="生日" width="150">
          <template #default="{ row }">
            {{ row.birthday || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="角色" width="100">
          <template #default="{ row }">
            <el-tag :type="row.role === 1 ? 'danger' : 'info'" size="small">
              {{ row.role === 1 ? '普通用户' : row.role === 2 ? '管理员' : '普通用户' }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination mt-md">
        <el-pagination
          v-model:current-page="currentPage"
          :page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="loadUsers"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getUserList, type UserInfo } from '@/api/user'

const users = ref<UserInfo[]>([])
const loading = ref(false)
const total = ref(0)
const currentPage = ref(1)
const pageSize = 10
const searchKey = ref('')

const loadUsers = async () => {
  loading.value = true
  try {
    const params: any = { pn: currentPage.value, psize: pageSize }
    if (searchKey.value) {
      params.search = searchKey.value
    }
    const res = await getUserList(params)
    // 后端返回的可能是数组或分页对象
    if (Array.isArray(res)) {
      users.value = res
      total.value = res.length
    } else {
      users.value = (res as any).data || []
      total.value = (res as any).total || users.value.length
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadUsers()
})
</script>

<style lang="scss" scoped>
h2 {
  margin-bottom: 20px;
}

.mb-md {
  margin-bottom: 16px;
}

.mt-md {
  margin-top: 16px;
}
</style>