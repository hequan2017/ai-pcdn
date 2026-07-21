import service from '@/utils/request'

export const getReleaseList = (params) => service({ url: '/pcdn/admin/release/list', method: 'get', params })
export const getRelease = (params) => service({ url: '/pcdn/admin/release/find', method: 'get', params })
export const createRelease = (data) => service({ url: '/pcdn/admin/release/create', method: 'post', data })
export const updateRelease = (data) => service({ url: '/pcdn/admin/release/update', method: 'put', data })
export const deleteRelease = (params) => service({ url: '/pcdn/admin/release/delete', method: 'delete', params })
