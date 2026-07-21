<template>
  <div class="auth-wrap">
    <el-card class="auth-card">
      <h2>注册账号</h2>
      <el-form :model="form" label-width="72px">
        <el-form-item label="用户名"><el-input v-model="form.username" placeholder="登录用户名" /></el-form-item>
        <el-form-item label="密码"><el-input v-model="form.password" type="password" show-password /></el-form-item>
        <el-form-item label="昵称"><el-input v-model="form.nickName" /></el-form-item>
        <el-form-item label="手机"><el-input v-model="form.phone" /></el-form-item>
        <el-form-item label="邮箱"><el-input v-model="form.email" /></el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="onSubmit">注册</el-button>
          <router-link to="/login" class="link">已有账号？去登录</router-link>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
  import { reactive, ref } from 'vue'
  import { useRouter } from 'vue-router'
  import { ElMessage } from 'element-plus'
  import { register } from '../api/portal'

  const router = useRouter()
  const loading = ref(false)
  const form = reactive({ username: '', password: '', nickName: '', phone: '', email: '' })

  const onSubmit = async () => {
    if (!form.username || !form.password) {
      ElMessage.warning('请填写用户名和密码')
      return
    }
    loading.value = true
    try {
      await register(form)
      ElMessage.success('注册成功，请登录')
      router.push('/login')
    } finally {
      loading.value = false
    }
  }
</script>

<style scoped>
  .auth-wrap { display: flex; justify-content: center; padding: 60px 20px; }
  .auth-card { width: 420px; }
  .link { margin-left: 12px; color: #2563eb; }
</style>
