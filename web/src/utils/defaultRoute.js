function isExternalUrl(val) {
  return typeof val === 'string' && /^(https?:)?\/\//.test(val)
}

function isAccessibleRoute(route) {
  if (!route) return false
  if (route.hidden || route.meta?.hidden) return false
  if (route.name === 'Reload' || route.name === 'layout') return false
  if (isExternalUrl(route.path) || isExternalUrl(route.name) || isExternalUrl(route.component)) {
    return false
  }
  return true
}

function findRouteByName(routes, targetName) {
  if (!Array.isArray(routes) || !targetName) return null

  for (const route of routes) {
    if (route?.name === targetName && isAccessibleRoute(route)) {
      return route
    }
    const found = findRouteByName(route?.children || [], targetName)
    if (found) {
      return found
    }
  }

  return null
}

export function findFirstAccessibleRouteName(routes) {
  if (!Array.isArray(routes)) return null

  for (const route of routes) {
    if (!isAccessibleRoute(route)) {
      continue
    }

    if (route.children && route.children.length > 0) {
      const childRouteName = findFirstAccessibleRouteName(route.children)
      if (childRouteName) {
        return childRouteName
      }
    }

    if (route.name) {
      return route.name
    }
  }

  return null
}

export function resolveDefaultRouterName(_configuredDefaultRouter, routes) {
  const routerName = String(_configuredDefaultRouter || '')
  if (routerName) {
    const matchedRoute = findRouteByName(routes, routerName)
    if (matchedRoute) {
      return routerName
    }
  }

  const firstRouteName = findFirstAccessibleRouteName(routes)
  if (firstRouteName) {
    return firstRouteName
  }

  return null
}
