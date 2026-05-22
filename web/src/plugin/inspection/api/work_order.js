import service from '@/utils/request'

export const assignBatchTemplate = (data) => {
  return service({ url: '/workOrder/assignBatchTemplate', method: 'post', data })
}

export const assignOrderTemplate = (data) => {
  return service({ url: '/workOrder/assignOrderTemplate', method: 'post', data })
}

export const startInspection = (data) => {
  return service({ url: '/workOrder/startInspection', method: 'post', data })
}

export const startRecheck = (data) => {
  return service({ url: '/workOrder/startRecheck', method: 'post', data })
}

export const getInspectionBatchList = (params) => {
  return service({ url: '/workOrder/getInspectionBatchList', method: 'get', params })
}

export const saveResults = (data) => {
  return service({ url: '/workOrder/saveResults', method: 'post', data })
}

export const saveSingleResult = (data) => {
  return service({ url: '/workOrder/saveSingleResult', method: 'post', data })
}

export const completeInspection = (data) => {
  return service({ url: '/workOrder/completeInspection', method: 'post', data })
}

export const completeRecheck = (data) => {
  return service({ url: '/workOrder/completeRecheck', method: 'post', data })
}

export const returnDevices = (data) => {
  return service({ url: '/workOrder/returnDevices', method: 'post', data })
}

export const getInspectionDetail = (params) => {
  return service({ url: '/workOrder/getInspectionDetail', method: 'get', params })
}

export const getBatchStatusLogs = (params) => {
  return service({ url: '/workOrder/getBatchStatusLogs', method: 'get', params })
}

export const getFlowLogs = (params) => {
  return service({ url: '/workOrder/getFlowLogs', method: 'get', params })
}

export const exportInspectionExcel = (params) => {
  return service({
    url: '/workOrder/exportInspectionExcel',
    method: 'get',
    params,
    responseType: 'blob'
  })
}
