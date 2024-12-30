<template>
	<div class="pt-20">
		<div class="card w-96 bg-base-200 text-neutral-content m-auto shadow-lg">
			<form class="card-body items-center text-center" @submit="login">
				<img src="/logo.svg" alt="logo" class="w-24 mb-5" />
				<h2 class="card-title mb-5">Enjoy movies & TV</h2>
				<input class="input input-sm input-bordered" type="text" v-model="email" placeholder="Username/Email" />
				<input class="input input-sm input-bordered" type="password" v-model="password" placeholder="**********" />
				<div class="card-actions justify-end">
					<button class="btn btn-primary btn-sm mt-5" type="submit">Login</button>
				</div>
			</form>
		</div>
	</div>
</template>
<script setup>
	import { getActivePinia } from 'pinia'
	getActivePinia()._s.forEach((s) => {
		s.$dispose()
	})
	definePageMeta({
		layout: 'empty',
	})

	const email = ref('')
	const password = ref('')

	async function login(e) {
		e.preventDefault()
		try {
			await usePb().admins.authWithPassword(email.value, password.value)
			// await usePb().admins.authWithPassword('admin@odin.local', 'odinAdmin1')
			window.location.reload(true)
		} catch (_) {}

		try {
			await usePb().collection('users').authWithPassword(email.value, password.value)
			window.location.reload(true)
		} catch (e) {}
	}
</script>
