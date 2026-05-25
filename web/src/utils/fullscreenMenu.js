export const openFullScreenMenuIfNeeded = (name) => {
  if (name !== 'inspectionWorkOrder' && name !== 'InspectWorkOrder') {
    return false
  }
  const url = `${window.location.origin}${window.location.pathname}#/inspectWorkOrder`
  window.open(url, '_blank')
  return true
}
