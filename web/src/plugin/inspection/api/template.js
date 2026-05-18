import service from '@/utils/request'

export const createTemplate = (data) => {
  return service({ url: '/inspectionTemplate/createTemplate', method: 'post', data })
}

export const deleteTemplate = (params) => {
  return service({ url: '/inspectionTemplate/deleteTemplate', method: 'delete', params })
}

export const updateTemplate = (data) => {
  return service({ url: '/inspectionTemplate/updateTemplate', method: 'put', data })
}

export const findTemplate = (params) => {
  return service({ url: '/inspectionTemplate/findTemplate', method: 'get', params })
}

export const getTemplateList = (params) => {
  return service({ url: '/inspectionTemplate/getTemplateList', method: 'get', params })
}
