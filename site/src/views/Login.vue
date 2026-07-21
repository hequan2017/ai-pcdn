<template>
  <div class="auth-wrap">
    <el-card class="auth-card">
      <h2>登录</h2>
      <el-form :model="form" label-width="72px">
        <el-form-item label="用户名"><el-input v-model="form.username" /></el-form-item>
        <el-form-item label="密码"><el-input v-model="form.password" type="password" show-password @keyup.enter="onSubmit" /></el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="onSubmit">登录</el-button>
          <router-link to="/register" class="link">没有账号？注册</router-link>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
  import { reactive, ref } from 'vue'
  import { useRouter } from 'vue-router'
  import { ElMessage } from 'element-plus'
  import { login } from '../api/portal'
  import { setToken, setUser } from '../store/auth'

  const router = useRouter()
  const loading = ref(false)
  const form = reactive({ username: '', password: '' })

  const onSubmit = async () => {
    if (!form.username || !form.password) {
      ElMessage.warning('请填写用户名和密码')
      return
    }
    loading.value = true
    try {
      const res = await login(form)
      setToken(res.data.token)
      setUser(res.data.user)
      ElMessage.success('登录成功')
      router.push('/dashboard')
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
