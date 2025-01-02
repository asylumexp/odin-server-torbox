<template>
	<div>
		<div v-if="rdloading">Loading...</div>
		<div v-else>
			<div v-if="profile">
				<p>{{ profile.username }}</p>
				<p>{{ profile.email }}</p>
				<p>{{ profile.expiration }}</p>
			</div>
			<div v-else>
				<dialog ref="login_dialog" class="modal">
					<div class="modal-box">
						<h3 class="font-bold text-lg">Login to RealDebrid</h3>
						<p class="py-4">Go to: {{ url }}</p>
						<p class="py-4">Enter code:</p>
						<p>{{ user_code }}</p>
					</div>
				</dialog>
				<button class="btn" @click="realDebridLogin()">Login</button>
			</div>
		</div>
	</div>
</template>

<script lang="ts" setup>
	const emit = defineEmits(['success'])

	const profile = ref<{ username: string; email: string; expiration: string } | null>(null)
	const rdloading = ref(true)

	onMounted(async () => {
		try {
			profile.value = await usePb().send('/-/realdebrid/user', { method: 'get' })
		} catch (e) {
			console.error(e)
		}
		rdloading.value = false
	})

	const host = 'https://api.real-debrid.com/oauth/v2'

	const login_dialog = ref<HTMLDialogElement>()
	const user_code = ref<string>()
	const url = ref<string>()
	const device_code = ref<string>()
	async function realDebridLogin() {
		login_dialog.value?.showModal()
		const data = await usePb().send(`/-/realdebrid/isAuth/device/code?client_id=X245A4XAIBGVM&=new_credentials=yes`, { method: 'get', cache: 'no-cache' })
		url.value = data.verification_url
		user_code.value = data.user_code
		device_code.value = data.device_code

		const poll = setInterval(async () => {
			const data2 = await usePb().send(`/-/realdebrid/isAuth/device/credentials?client_id=X245A4XAIBGVM&code=${device_code.value}`, { method: 'get', cache: 'no-cache' })
			console.log(data2)
			if (data2 !== null) {
				clearInterval(poll)
				const formData = new FormData()
				formData.append('client_id', data2.client_id)
				formData.append('client_secret', data2.client_secret)
				formData.append('code', device_code.value || '')
				formData.append('grant_type', 'http://oauth.net/grant_type/device/1.0')

				console.log(formData)

				const data3 = await usePb().send(`/-/realdebrid/isAuth/token`, {
					method: 'post',
					cache: 'no-cache',
					body: { client_id: data2.client_id, client_secret: data2.client_secret, code: device_code.value, grant_type: 'http://oauth.net/grant_type/device/1.0' },
				})
				if (data3 !== null) {
					emit('success', {
						...data3,
						...data2,
					})
					login_dialog.value?.close()
					profile.value = await usePb().send('/-/realdebrid/user', { method: 'get' })
				}
			}
		}, 5000)
	}
</script>
