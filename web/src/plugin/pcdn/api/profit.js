import service from '@/utils/request'

export const getProfitSummary = (params) => service({ url: '/pcdn/admin/profit/summary', method: 'get', params })
export const getRevenueByPlatform = (params) => service({ url: '/pcdn/admin/profit/revenueByPlatform', method: 'get', params })
export const getCostByOwner = (params) => service({ url: '/pcdn/admin/profit/costByOwner', method: 'get', params })
export const getProfitTrend = (params) => service({ url: '/pcdn/admin/profit/trend', method: 'get', params })
