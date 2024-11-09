export default defineNuxtRouteMiddleware(async (to, from) => {
	console.log('AUTH', usePb().authStore.isValid)
	if (process.server) return
	if (usePb().authStore.isValid && ['/login'].includes(to.path)) {
		console.log('SHOULD LOGIN')
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
