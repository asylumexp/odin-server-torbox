export const useProfile = defineStore('useProfile', () => {
	let me = ref()
	async function init() {
		if (usePb().authStore.isAdmin) {
			me.value = await usePb().admins.getOne(usePb().authStore.model?.id)
			return
		}
		me.value = await usePb()
			.collection('users')
			.getOne(usePb().authStore.model?.id)
	}
	return {
		me,
		init,
	}
})
