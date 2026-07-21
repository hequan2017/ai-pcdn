// 个人门户登录态（localStorage 持久化，简化版无 pinia）
const TOKEN_KEY = 'pcdn_portal_token'
const USER_KEY = 'pcdn_portal_user'

export const getToken = () => localStorage.getItem(TOKEN_KEY)
export const setToken = (t) => localStorage.setItem(TOKEN_KEY, t)
export const clearToken = () => localStorage.removeItem(TOKEN_KEY)

export const getUser = () => JSON.parse(localStorage.getItem(USER_KEY) || 'null')
export const setUser = (u) => localStorage.setItem(USER_KEY, JSON.stringify(u || {}))

export const clearAuth = () => {
  clearToken()
  localStorage.removeItem(USER_KEY)
}
