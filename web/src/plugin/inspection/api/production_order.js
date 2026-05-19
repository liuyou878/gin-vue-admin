import service from '@/utils/request'

export const createProductionOrder = (data) => {
  return service({ url: '/productionOrder/createProductionOrder', method: 'post', data })
}

export const deleteProductionOrder = (params) => {
  return service({ url: '/productionOrder/deleteProductionOrder', method: 'delete', params })
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
