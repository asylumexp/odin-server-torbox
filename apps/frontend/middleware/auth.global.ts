export default defineNuxtRouteMiddleware(async (to, _) => {
	if (process.server) return
	const id = usePb().authStore.model?.id || null
	console.log('ID', id)
	console.log('valid', usePb().authStore.isValid)
	console.log('admin', usePb().authStore.isAdmin)
	console.log('paths', to.path)
	if (
		id !== null &&
		usePb().authStore.isValid &&
		((usePb().authStore.isAdmin && !(await usePb().admins.getOne(id))) || (!usePb().authStore.isAdmin && !(await usePb().collection('users').getOne(id))))
	) {
		console.log('NO USER')
		usePb().authStore.clear()
	}
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
