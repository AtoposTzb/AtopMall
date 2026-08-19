<template>
  <div class="user-profile">
    <h3 class="mb-md">个人信息</h3>

    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="100px"
      v-loading="loading"
      style="max-width: 500px"
    >
      <el-form-item label="用户名" prop="name">
        <el-input v-model="form.name" placeholder="请输入用户名" />
      </el-form-item>
      <el-form-item label="性别" prop="gender">
        <el-radio-group v-model="form.gender">
          <el-radio value="male">男</el-radio>
          <el-radio value="female">女</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="生日" prop="birthday">
        <el-date-picker
          v-model="form.birthday"
          type="date"
          placeholder="选择生日"
          format="YYYY-MM-DD"
          value-format="YYYY-MM-DD"
          style="width: 100%"
        />
      </el-form-item>
      <el-form-item label="手机号">
        <el-input :model-value="userStore.userInfo?.mobile" disabled />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="handleSave" :loading="saving"
          >保存修改</el-button
        >
        <el-button @click="handleReset">重置</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from "vue";
import { ElMessage, type FormInstance, type FormRules } from "element-plus";
import { getUserDetail, updateUser } from "@/api/user";
import { useUserStore } from "@/store/user";

const userStore = useUserStore();
const loading = ref(false);
const saving = ref(false);
const formRef = ref<FormInstance>();

const form = reactive({
  name: "",
  gender: "male",
  birthday: "",
});

const rules: FormRules = {
  name: [
    { required: true, message: "请输入用户名", trigger: "blur" },
    { min: 3, max: 10, message: "长度在 3 到 10 个字符", trigger: "blur" },
  ],
  gender: [{ required: true, message: "请选择性别", trigger: "change" }],
  birthday: [{ required: true, message: "请选择生日", trigger: "change" }],
};

const loadUserInfo = async () => {
  loading.value = true;
  try {
    const res = await getUserDetail();
    form.name = res.name || "";
    form.gender = res.gender || "male";
    form.birthday = res.birthday || "";
  } catch (error) {
    console.error("加载用户信息失败", error);
  } finally {
    loading.value = false;
  }
};

const handleSave = async () => {
  if (!formRef.value) return;

  const valid = await formRef.value.validate().catch(() => false);
  if (!valid) return;

  saving.value = true;
  try {
    await updateUser(form);
    ElMessage.success("更新成功");
    await userStore.initUser();
  } catch (error) {
    // 错误已在拦截器中处理
  } finally {
    saving.value = false;
  }
};

const handleReset = () => {
  formRef.value?.resetFields();
  loadUserInfo();
};

onMounted(() => {
  loadUserInfo();
});
</script>

<style lang="scss" scoped>
.user-profile {
  padding: 8px 0;

  h3 {
    font-size: 20px;
    font-weight: 600;
    color: $text-primary;
    margin-bottom: 28px;
    padding-bottom: 16px;
    border-bottom: 2px solid #f0f5ff;
    position: relative;

    &::after {
      content: '';
      position: absolute;
      bottom: -2px;
      left: 0;
      width: 60px;
      height: 2px;
      background: $primary-color;
      border-radius: 1px;
    }
  }

  :deep(.el-form) {
    .el-form-item__label {
      font-weight: 500;
      color: $text-primary;
    }

    .el-input__wrapper {
      border-radius: 8px;
      transition: all 0.3s;
      box-shadow: 0 0 0 1px $border-light inset;

      &:hover {
        box-shadow: 0 0 0 1px $primary-color inset;
      }

      &.is-focus {
        box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.2) inset;
      }
    }

    .el-radio-group {
      .el-radio {
        margin-right: 24px;
      }
    }

    .el-button--primary {
      border-radius: 8px;
      padding: 10px 28px;
      font-weight: 500;
    }

    .el-button:not(.el-button--primary) {
      border-radius: 8px;
      padding: 10px 28px;
    }
  }
}

@media (max-width: $bp-mobile) {
  .user-profile {
    :deep(.el-form-item__label) {
      width: 60px !important;
    }

    :deep(.el-form) {
      max-width: 100% !important;
    }
  }
}
</style>