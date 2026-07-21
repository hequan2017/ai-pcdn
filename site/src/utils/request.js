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

// 用 HTTP 状态码区分鉴权失败：
//   GVA NoAuth 返回 HTTP 401 + code=7（鉴权失败，清 token 跳登录）
//   FailWithMessage 返回 HTTP 200 + code=7（业务错误，仅弹错，不清 token）
service.interceptors.response.use(
  (resp) => {
    const data = resp.data
    if (data && data.code === 0) {
      return data
    }
    // HTTP 200 但 code!=0：业务错误
    ElMessage.error((data && data.msg) || '请求失败')
    return Promise.reject(data)
  },
  (err) => {
    if (err.response && err.response.status === 401) {
      clearAuth()
      if (location.hash !== '#/login') {
        location.hash = '#/login'
      }
      ElMessage.error('登录已过期，请重新登录')
    } else {
      ElMessage.error(err.message || '网络错误')
    }
    return Promise.reject(err)
  }
)

export default service
