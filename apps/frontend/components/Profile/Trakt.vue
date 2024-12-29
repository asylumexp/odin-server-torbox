<template>
	<div v-if="traktProfile !== null" class="card w-96 bg-base-100 shadow-xl">
		<div class="card-body">
			<h2 class="card-title">Trakt</h2>
			<p>{{ traktProfile.user.username }}<br />{{ traktProfile.user.name }}</p>
			<div class="card-actions justify-end">
				<button class="btn btn-primary">Refresh</button>
			</div>
		</div>
	</div>
	<div v-else class="card w-96 bg-base-100 shadow-xl">
		<dialog ref="login_dialog" class="modal">
			<div class="modal-box">
				<h3 class="font-bold text-lg">Login to Trakt</h3>
				<p class="py-4">
					Go to: <a :href="url">{{ url }}</a>
				</p>
				<p class="py-4">Enter code:</p>
				<p>{{ user_code }}</p>
			</div>
		</dialog>
		<div class="card-body">
			<h2 class="card-title">Trakt</h2>
			<p>Click below to login into Trakt</p>
			<div class="card-actions justify-end">
				<button class="btn" @click="traktLogin()">Login</button>
			</div>
		</div>
	</div>
</template>

<script lang="ts" setup>
	const login_dialog = ref<HTMLDialogElement>()
	const user_code = ref<string>()
	const url = ref<string>()
	const device_code = ref<string>()
	let traktProfile: any = null
	try {
		traktProfile = await usePb().send('/-/trakt/users/settings', {
			method: 'GET',
		})
	} catch (e) {
		console.log(e)
	}
	async function traktLogin() {
		const settings = useSettings()

		login_dialog.value?.showModal()
		const res = await usePb().send('/-/trakt/oauth/device/code?fresh=true', {
			method: 'POST',
			body: {
				client_id: 'd0ba20c3bb7de7c8108d02f2b2c1eb1b85f74cff5c11dd17554ac063dce9ab12',
			},
		})

		url.value = res.verification_url
		user_code.value = res.user_code
		device_code.value = res.device_code

		const poll = setInterval(async () => {
			const res = await usePb().send('/-/trakt/oauth/device/token?fresh=true', {
				method: 'POST',
				body: {
					client_id: 'd0ba20c3bb7de7c8108d02f2b2c1eb1b85f74cff5c11dd17554ac063dce9ab12',
					client_secret: '1643cc9159f628ff5c42bf023732b0c1f21f1e1e618ab95f68c9e3d5fb1f7186',
					code: device_code.value,
				},
			})
			if (res !== null) {
				console.log(usePb().authStore.model?.id, res)
				await usePb()
					.collection('users')
					.update(usePb().authStore.model?.id, { trakt_token: { ...res, device_code: device_code.value } })
				clearInterval(poll)
				login_dialog.value?.close()
			}
		}, 5000)
	}
</script>
