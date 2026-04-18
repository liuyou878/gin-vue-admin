import service from '@/utils/request'

export const getPublicFirmwareDownloadPage = (params) => {
  return service({
    url: '/firmwarePublic/getPublicFirmwareDownloadPage',
    method: 'get',
    params
  })
}
