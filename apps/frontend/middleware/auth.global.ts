export default defineNuxtRouteMiddleware(async (to, from) => {
  if (process.server) return
  if (usePb().authStore.isValid && ['/login'].includes(to.path)) {
    if (usePb().authStore.isAdmin) {
      return navigateTo('/admin')
    }
    return navigateTo('/')
  }
  if (!usePb().authStore.isValid && !['/login'].includes(to.path)) {
    return navigateTo('/login')
  }
  // if (usePb().authStore.isValid && from.path !== "/" && to.path === "/") {
  //   return navigateTo(from.fullPath);
  // }
})
