<template>
	<div>
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
</template>

<script lang="ts" setup>
	const emit = defineEmits(['success'])

	let profile: null | any = null

	try {
		profile = await usePb().send('/realdebrid/user', { method: 'get' })
	} catch (e) {
		console.log(e)
	}

	const host = 'https://api.real-debrid.com/oauth/v2'

	const login_dialog = ref<HTMLDialogElement>()
	const user_code = ref<string>()
	const url = ref<string>()
	const device_code = ref<string>()
	async function realDebridLogin() {
		login_dialog.value?.showModal()
		const res = await useFetch(`${host}/device/code?client_id=X245A4XAIBGVM&=new_credentials=yes`)

		const data = res.data.value as any
		url.value = data.verification_url
		user_code.value = data.user_code
		device_code.value = data.device_code

		const poll = setInterval(async () => {
			const res2 = await useFetch(`${host}/device/credentials?client_id=X245A4XAIBGVM&code=${device_code.value}`, { method: 'get', cache: 'no-cache' })
			if (res2.data.value !== null) {
				clearInterval(poll)
				const data = res2.data.value as any
				const formData = new FormData()
				formData.append('client_id', data.client_id)
				formData.append('client_secret', data.client_secret)
				formData.append('code', device_code.value || '')
				formData.append('grant_type', 'http://oauth.net/grant_type/device/1.0')

				const res3 = await useFetch(`${host}/token`, {
					method: 'post',
					cache: 'no-cache',
					body: formData,
				})
				if (res3.data.value !== null) {
					emit('success', {
						...res3.data.value,
						...res2.data.value,
					})
					login_dialog.value?.close()
				}
			}
		}, 5000)
	}
</script>
