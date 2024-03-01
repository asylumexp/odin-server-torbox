<template>
	<div>
		<h1>Devices</h1>
		<button class="btn btn-primary" @click="openModal"><FaIcon icon="computer" /> Link Device</button>
		<table class="table">
			<tbody>
				<tr v-for="device in devices" :key="device.id">
					<td><FaIcon icon="check-circle" class="mr-2" :class="device.verified ? 'text-success' : ''" />{{ device.token }}</td>
					<td>{{ device.name }}</td>
					<td>{{ device.created }}</td>
				</tr>
			</tbody>
		</table>
		<dialog class="modal" ref="device_dialog">
			<div class="modal-box">
				<h1>Dialog open</h1>
				<input type="text" v-model="id" />

				<button v-if="!verified && !loading" @click="linkDevice">Link Device</button>
				<div v-if="loading">Loading</div>
			</div>
		</dialog>
	</div>
</template>

<script lang="ts" setup>
	const devices = ref()
	const loading = ref(false)
	const device_dialog = ref<HTMLDialogElement>()
	const id = ref('')
	const verified = ref(false)
	const dialogOpen = ref(false)

	function openModal() {
		dialogOpen.value = false
		device_dialog.value?.showModal()
	}

	async function generateToken() {
		const id = Math.random().toString(36).substring(2).toUpperCase()

		const existing = await usePb()
			.collection('devices')
			.getFullList({ filter: `(token='${id}')` })
		if (existing.length > 0) {
			return generateToken()
		} else {
			return id
		}
	}
	async function linkDevice() {
		loading.value = true
		// const url = useRequestURL().origin
		const url = 'https://local-8090.add.dnmc.in'
		const deviceToken = await generateToken()

		const d = await usePb()
			.collection('devices')
			.create({
				user: usePb().authStore.model?.id,
				token: deviceToken,
				verified: false,
				name: 'My Device',
			})
		const topic = `odinmovieshows-${id.value}`
		const deviceId = d.id
		const data = {
			url,
			deviceId,
		}

		while (!verified.value) {
			verified.value = (await usePb().collection('devices').getOne(deviceId)).verified === true
			if (verified.value) {
				break
			}
			console.log(deviceId, `https://ntfy.sh/${topic}`)
			await useFetch(`https://ntfy.sh/${topic}`, { method: 'POST', body: JSON.stringify(data) })
			await new Promise((resolve) => setTimeout(resolve, 5000))
		}
		loading.value = false
		device_dialog.value?.close()
		devices.value = await getDevices()
	}
	async function getDevices() {
		return await usePb()
			.collection('devices')
			.getFullList({ filter: `(user='${usePb().authStore.model?.id}')` })
	}
	onMounted(async () => {
		devices.value = await getDevices()
	})
</script>
