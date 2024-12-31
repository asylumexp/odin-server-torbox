export default defineNuxtRouteMiddleware(async (to, _) => {
	if (process.server) return
	const id = usePb().authStore.model?.id || null
	console.log('ID', id)
	console.log('valid', usePb().authStore.isValid)
	console.log('admin', usePb().authStore.isAdmin)
	console.log('paths', to.path)
	if (id !== null && usePb().authStore.isValid) {
		let a = null
		try {
			a = await usePb().admins.getOne(id)
		} catch (_) {}
		if (usePb().authStore.isAdmin) {
			if (a !== null) {
				return
			}
			console.log('NO ADMIN')
			usePb().authStore.clear()
			return navigateTo('/login')
		} else {
			let u = null
			try {
				u = await usePb().collection('users').getOne(id)
			} catch (_) {}
			if (u !== null) {
				return
			}
			console.log('NO USER')
			usePb().authStore.clear()
			return navigateTo('/login')
		}
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
