import service from '@/utils/request'

/**
 * 新增检测项
 */
export const createItem = (data) => {
  return service({
    url: '/inspectionItem/createItem',
    method: 'post',
    data
  })
}

/**
 * 删除检测项
 */
export const deleteItem = (params) => {
  return service({
    url: '/inspectionItem/deleteItem',
    method: 'delete',
    params
  })
}

/**
 * 批量删除检测项
 */
export const deleteItemByIds = (params) => {
  return service({
    url: '/inspectionItem/deleteItemByIds',
    method: 'delete',
    params
  })
}

/**
 * 更新检测项
 */
export const updateItem = (data) => {
  return service({
    url: '/inspectionItem/updateItem',
    method: 'put',
    data
  })
}

/**
 * 根据ID获取检测项
 */
export const findItem = (params) => {
  return service({
    url: '/inspectionItem/findItem',
    method: 'get',
    params
  })
}

/**
 * 获取检测项列表
 */
export const getItemList = (params) => {
  return service({
    url: '/inspectionItem/getItemList',
    method: 'get',
    params
  })
}
