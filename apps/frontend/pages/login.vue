<template>
	<div><button @click="login">Login</button></div>
</template>
<script setup>
	import { getActivePinia } from 'pinia'
	getActivePinia()._s.forEach((s) => {
		s.$dispose()
	})
	definePageMeta({
		layout: 'empty',
	})

	async function login() {
		try {
			await usePb().admins.authWithPassword('admin@odin.local', 'odinAdmin1')
			window.location.reload(true)
		} catch (e) {}

		// await usePb().collection('users').authWithPassword('adis', 'odinAdis1')
		// window.location.reload(true)
	}
</script>
