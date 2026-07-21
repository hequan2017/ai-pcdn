import service from '@/utils/request'

export const generateBill = (data) => service({ url: '/pcdn/admin/bill/generate', method: 'post', data })
export const getBillList = (params) => service({ url: '/pcdn/admin/bill/list', method: 'get', params })
export const getBill = (params) => service({ url: '/pcdn/admin/bill/find', method: 'get', params })
export const approveBill = (data) => service({ url: '/pcdn/admin/bill/approve', method: 'put', data })
export const rejectBill = (data) => service({ url: '/pcdn/admin/bill/reject', method: 'put', data })
export const payBill = (data) => service({ url: '/pcdn/admin/bill/pay', method: 'put', data })
