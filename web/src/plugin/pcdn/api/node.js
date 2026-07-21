import service from '@/utils/request'

// 查询节点列表
export const getNodeList = (params) => {
  return service({ url: '/pcdn/admin/node/list', method: 'get', params })
}

// 查询节点详情
export const getNode = (params) => {
  return service({ url: '/pcdn/admin/node/find', method: 'get', params })
}

// 创建节点
export const createNode = (data) => {
  return service({ url: '/pcdn/admin/node/create', method: 'post', data })
}

// 更新节点
export const updateNode = (data) => {
  return service({ url: '/pcdn/admin/node/update', method: 'put', data })
}

// 删除节点
export const deleteNode = (params) => {
  return service({ url: '/pcdn/admin/node/delete', method: 'delete', params })
}

// 批量删除节点
export const deleteNodeByIds = (data) => {
  return service({ url: '/pcdn/admin/node/deleteByIds', method: 'delete', data })
}

// 查询节点流量曲线
export const getNodeTraffic = (params) => {
  return service({ url: '/pcdn/admin/node/traffic', method: 'get', params })
}

// 查询节点95值
export const getNode95 = (params) => {
  return service({ url: '/pcdn/admin/node/n95', method: 'get', params })
}
