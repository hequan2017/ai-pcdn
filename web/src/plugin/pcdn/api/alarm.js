import service from '@/utils/request'

// 告警规则
export const getAlarmRuleList = (params) => service({ url: '/pcdn/admin/alarm/rule/list', method: 'get', params })
export const getAlarmRule = (params) => service({ url: '/pcdn/admin/alarm/rule/find', method: 'get', params })
export const createAlarmRule = (data) => service({ url: '/pcdn/admin/alarm/rule/create', method: 'post', data })
export const updateAlarmRule = (data) => service({ url: '/pcdn/admin/alarm/rule/update', method: 'put', data })
export const deleteAlarmRule = (params) => service({ url: '/pcdn/admin/alarm/rule/delete', method: 'delete', params })
export const deleteAlarmRuleByIds = (data) => service({ url: '/pcdn/admin/alarm/rule/deleteByIds', method: 'delete', data })

// 告警记录
export const getAlarmRecordList = (params) => service({ url: '/pcdn/admin/alarm/record/list', method: 'get', params })
