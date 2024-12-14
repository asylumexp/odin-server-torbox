<template>
	<div>
		<div v-if="rdloading">Loading...</div>
		<div v-else>
			<div v-if="profile">
				<p>{{ profile.data.user.username }}</p>
				<p>{{ profile.data.user.email }}</p>
				<p>{{ profile.data.user.premiumUntil }}</p>
			</div>
			<div v-else>
				<dialog ref="login_dialog" class="modal">
					<div class="modal-box">
						<h3 class="font-bold text-lg">Login to RealDebrid</h3>
						<p class="py-4">Go to: {{ url }}</p>
						<p class="py-4">Enter code:</p>
						<p>{{ pin }}</p>
					</div>
				</dialog>
				<button class="btn" @click="allDebridLogin()">Login</button>
			</div>
		</div>
	</div>
</template>

<script lang="ts" setup>
	const emit = defineEmits(['success'])

	const profile = ref<{ data: { user: { username: string; email: string; premiumUntil: string } } } | null>(null)
	const rdloading = ref(true)

	onMounted(async () => {
		try {
			profile.value = await usePb().send('/_alldebrid/user', { method: 'get' })
			console.log(profile.value)
			console.log('DONE')
		} catch (e) {
			console.error(e)
		}
		rdloading.value = false
	})

	const host = 'https://api.alldebrid.com/v4'

	const login_dialog = ref<HTMLDialogElement>()
	const pin = ref<string>()
	const url = ref<string>()
	const check = ref<string>()
	async function allDebridLogin() {
		login_dialog.value?.showModal()
		const res = await useFetch(`${host}/pin/get?agent=odinMovieShow`)
		const data = res.data.value as any
		console.log(data)
		url.value = data.data.user_url
		pin.value = data.data.pin
		check.value = data.data.check

		const poll = setInterval(async () => {
			const res2 = await useFetch(`${host}/pin/check?agent=odinMovieShow&check=${check.value}&pin=${pin.value}`)
			if (res2.data.value !== null) {
				const data2 = res2.data.value as any
				console.log(data2)
				if (data2.data.activated !== true) {
					clearInterval(poll)
					emit('success', {
						...res2.data.value,
					})
				}
				login_dialog.value?.close()
			}
		}, 5000)
	}
</script>
