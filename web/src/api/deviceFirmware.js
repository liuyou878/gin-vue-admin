import service from '@/utils/request'

export const createDeviceCategory = (data) => {
  return service({
    url: '/deviceCategory/createDeviceCategory',
    method: 'post',
    data
  })
}

export const deleteDeviceCategory = (params) => {
  return service({
    url: '/deviceCategory/deleteDeviceCategory',
    method: 'delete',
    params
  })
}

export const deleteDeviceCategoryByIds = (data) => {
  return service({
    url: '/deviceCategory/deleteDeviceCategoryByIds',
    method: 'delete',
    data
  })
}

export const updateDeviceCategory = (data) => {
  return service({
    url: '/deviceCategory/updateDeviceCategory',
    method: 'put',
    data
  })
}

export const findDeviceCategory = (params) => {
  return service({
    url: '/deviceCategory/findDeviceCategory',
    method: 'get',
    params
  })
}

export const getDeviceCategoryList = (params) => {
  return service({
    url: '/deviceCategory/getDeviceCategoryList',
    method: 'get',
    params
  })
}

export const createDeviceModel = (data) => {
  return service({
    url: '/deviceModel/createDeviceModel',
    method: 'post',
    data
  })
}

export const deleteDeviceModel = (params) => {
  return service({
    url: '/deviceModel/deleteDeviceModel',
    method: 'delete',
    params
  })
}

export const deleteDeviceModelByIds = (data) => {
  return service({
    url: '/deviceModel/deleteDeviceModelByIds',
    method: 'delete',
    data
  })
}

export const updateDeviceModel = (data) => {
  return service({
    url: '/deviceModel/updateDeviceModel',
    method: 'put',
    data
  })
}

export const findDeviceModel = (params) => {
  return service({
    url: '/deviceModel/findDeviceModel',
    method: 'get',
    params
  })
}

export const getDeviceModelList = (params) => {
  return service({
    url: '/deviceModel/getDeviceModelList',
    method: 'get',
    params
  })
}

export const createFirmwareVersion = (data) => {
  return service({
    url: '/firmwareVersion/createFirmwareVersion',
    method: 'post',
    data
  })
}

export const deleteFirmwareVersion = (params) => {
  return service({
    url: '/firmwareVersion/deleteFirmwareVersion',
    method: 'delete',
    params
  })
}

export const deleteFirmwareVersionByIds = (data) => {
  return service({
    url: '/firmwareVersion/deleteFirmwareVersionByIds',
    method: 'delete',
    data
  })
}

export const updateFirmwareVersion = (data) => {
  return service({
    url: '/firmwareVersion/updateFirmwareVersion',
    method: 'put',
    data
  })
}

export const findFirmwareVersion = (params) => {
  return service({
    url: '/firmwareVersion/findFirmwareVersion',
    method: 'get',
    params
  })
}

export const getFirmwareVersionList = (params) => {
  return service({
    url: '/firmwareVersion/getFirmwareVersionList',
    method: 'get',
    params
  })
}

export const changeFirmwareVersionStatus = (data) => {
  return service({
    url: '/firmwareVersion/changeFirmwareVersionStatus',
    method: 'post',
    data
  })
}

export const createModelFirmwareRel = (data) => {
  return service({
    url: '/modelFirmware/createModelFirmwareRel',
    method: 'post',
    data
  })
}

export const deleteModelFirmwareRel = (params) => {
  return service({
    url: '/modelFirmware/deleteModelFirmwareRel',
    method: 'delete',
    params
  })
}

export const deleteModelFirmwareRelByIds = (data) => {
  return service({
    url: '/modelFirmware/deleteModelFirmwareRelByIds',
    method: 'delete',
    data
  })
}

export const updateModelFirmwareRel = (data) => {
  return service({
    url: '/modelFirmware/updateModelFirmwareRel',
    method: 'put',
    data
  })
}

export const findModelFirmwareRel = (params) => {
  return service({
    url: '/modelFirmware/findModelFirmwareRel',
    method: 'get',
    params
  })
}

export const getModelFirmwareRelList = (params) => {
  return service({
    url: '/modelFirmware/getModelFirmwareRelList',
    method: 'get',
    params
  })
}

export const setModelFirmwareRecommended = (data) => {
  return service({
    url: '/modelFirmware/setModelFirmwareRecommended',
    method: 'post',
    data
  })
}

export const setModelFirmwareTestResult = (data) => {
  return service({
    url: '/modelFirmware/setModelFirmwareTestResult',
    method: 'post',
    data
  })
}

export const createFirmwareTag = (data) => {
  return service({
    url: '/firmwareTag/createFirmwareTag',
    method: 'post',
    data
  })
}

export const updateFirmwareTag = (data) => {
  return service({
    url: '/firmwareTag/updateFirmwareTag',
    method: 'put',
    data
  })
}

export const deleteFirmwareTag = (params) => {
  return service({
    url: '/firmwareTag/deleteFirmwareTag',
    method: 'delete',
    params
  })
}

export const getFirmwareTagList = (params) => {
  return service({
    url: '/firmwareTag/getFirmwareTagList',
    method: 'get',
    params
  })
}

export const setFirmwareTags = (data) => {
  return service({
    url: '/firmwareTag/setFirmwareTags',
    method: 'post',
    data
  })
}

export const getFirmwareVersionLogList = (params) => {
  return service({
    url: '/firmwareVersionLog/getFirmwareVersionLogList',
    method: 'get',
    params
  })
}

export const findFirmwareVersionLog = (params) => {
  return service({
    url: '/firmwareVersionLog/findFirmwareVersionLog',
    method: 'get',
    params
  })
}
