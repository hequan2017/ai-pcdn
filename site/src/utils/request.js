import axios from 'axios'
import { ElMessage } from 'element-plus'
import { getToken, clearAuth } from '../store/auth'

const service = axios.create({
  baseURL: import.meta.env.VITE_API_BASE || '',
  timeout: 15000
})

// 请求携带 x-token（与 GVA 主前端一致）
service.interceptors.request.use((cfg) => {
  const t = getToken()
  if (t) {
    cfg.headers['x-token'] = t
  }
  return cfg
})

// 统一处理后台响应 {code,data,msg}：code===0 成功
service.interceptors.response.use(
  (resp) => {
    const data = resp.data
    if (data && data.code === 0) {
      return data
    }
    // 鉴权类失败清理登录态
    if (data && data.code === 7) {
      clearAuth()
      if (location.hash !== '#/login') {
        location.hash = '#/login'
      }
    }
    ElMessage.error((data && data.msg) || '请求失败')
    return Promise.reject(data)
  },
  (err) => {
    ElMessage.error(err.message || '网络错误')
    return Promise.reject(err)
  }
)

export default service
