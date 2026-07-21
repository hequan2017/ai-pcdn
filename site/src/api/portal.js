import service from '../utils/request'

// 个人注册
export const register = (data) => service({ url: '/pcdn/portal/register', method: 'post', data })

// 个人登录
export const login = (data) => service({ url: '/pcdn/portal/login', method: 'post', data })

// 我的节点列表
export const myNodes = (params) => service({ url: '/pcdn/portal/myNodes', method: 'get', params })

// 我的节点流量
export const myNodeTraffic = (params) => service({ url: '/pcdn/portal/myNodeTraffic', method: 'get', params })

// 添加节点（自助上机）
export const addNode = (data) => service({ url: '/pcdn/portal/addNode', method: 'post', data })

// 我的账单
export const myBills = (params) => service({ url: '/pcdn/portal/myBills', method: 'get', params })
