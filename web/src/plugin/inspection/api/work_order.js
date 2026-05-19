import service from '@/utils/request'

export const startInspection = (data) => {
  return service({ url: '/workOrder/startInspection', method: 'post', data })
}

export const saveResults = (data) => {
  return service({ url: '/workOrder/saveResults', method: 'post', data })
}

export const completeInspection = (data) => {
  return service({ url: '/workOrder/completeInspection', method: 'post', data })
}

export const getInspectionDetail = (params) => {
  return service({ url: '/workOrder/getInspectionDetail', method: 'get', params })
}
