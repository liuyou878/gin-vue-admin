const fullScreenMenuNames = ['inspectionWorkOrder', 'InspectWorkOrder', 'work_order']

export const isFullScreenMenu = (name) => fullScreenMenuNames.includes(name)

export const openFullScreenMenuIfNeeded = (name) => {
  if (!isFullScreenMenu(name)) return false
  const url = `${window.location.origin}${window.location.pathname}#/inspectWorkOrder`
  window.open(url, '_blank')
  return true
}
