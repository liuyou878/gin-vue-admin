import service from '@/utils/request'

export const createProductionOrder = (data) => {
  return service({ url: '/productionOrder/createProductionOrder', method: 'post', data })
}

export const deleteProductionOrder = (params) => {
  return service({ url: '/productionOrder/deleteProductionOrder', method: 'delete', params })
}

export const forceDeleteProductionOrder = (params) => {
  return service({ url: '/productionOrder/forceDeleteProductionOrder', method: 'delete', params })
}

export const updateProductionOrder = (data) => {
  return service({ url: '/productionOrder/updateProductionOrder', method: 'put', data })
}

export const findProductionOrder = (params) => {
  return service({ url: '/productionOrder/findProductionOrder', method: 'get', params })
}

export const getProductionOrderList = (params) => {
  return service({ url: '/productionOrder/getProductionOrderList', method: 'get', params })
}

export const getSubmittedDeviceList = (params) => {
  return service({ url: '/productionOrder/getSubmittedDeviceList', method: 'get', params })
}

export const findSubmittedDevice = (params) => {
  return service({ url: '/productionOrder/findSubmittedDevice', method: 'get', params })
}

export const deleteSubmittedDevice = (params) => {
  return service({ url: '/productionOrder/deleteSubmittedDevice', method: 'delete', params })
}

export const confirmReworkDone = (data) => {
  return service({ url: '/productionOrder/confirmReworkDone', method: 'post', data })
}

export const confirmReworkReceived = (data) => {
  return service({ url: '/productionOrder/confirmReworkReceived', method: 'post', data })
}

export const scanAssignBatch = (data) => {
  return service({ url: '/productionOrder/scanAssignBatch', method: 'post', data })
}

export const getDeviceStatusLogs = (params) => {
  return service({ url: '/productionOrder/getDeviceStatusLogs', method: 'get', params })
}
