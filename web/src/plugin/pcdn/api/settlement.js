import service from '@/utils/request'

export const importSettlement = (data) => service({ url: '/pcdn/admin/settlement/import', method: 'post', data })
export const recheckSettlement = (params) => service({ url: '/pcdn/admin/settlement/recheck', method: 'put', params })
export const getSettlementList = (params) => service({ url: '/pcdn/admin/settlement/list', method: 'get', params })
export const getRevenueSummary = (params) => service({ url: '/pcdn/admin/settlement/revenue', method: 'get', params })
export const deleteSettlement = (params) => service({ url: '/pcdn/admin/settlement/delete', method: 'delete', params })
