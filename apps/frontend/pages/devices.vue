<template>
	<div>
		<h1>Devices</h1>
		<button class="btn btn-primary" @click="openModal"><FaIcon icon="computer" /> Connect</button>
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
				<div v-if="loading" class="text-center">
					<h3 class="text-success">Linking devices</h3>
					<span class="loading loading-infinity loading-lg text-success"></span>
					<p>{{ url }}</p>
				</div>
				<div v-else class="text-center">
					<h1><FaIcon icon="tv" class="mr-5" />Link Device</h1>
					<p>Please enter the code shown on your TV app</p>
					<input type="text" class="input input-bordered mr-5" v-model="id" />
					<button class="btn btn-md btn-primary" v-if="!verified && !loading" @click="linkDevice">Link Device</button>
					<p class="text-sm m-0 mt-10">The following URL will be sent to the device:</p>
					<p class="text-success text-sm m-0">{{ url }}</p>
				</div>
			</div>
		</dialog>
	</div>
</template>

<script lang="ts" setup>
	const devices = ref()
	const loading = ref(false)
	const device_dialog = ref<HTMLDialogElement>()
	const id = ref('')
	const url = ref('')
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
		const deviceToken = await generateToken()

		const d = await usePb().collection('devices').create({
			user: usePb().authStore.model?.id,
			token: deviceToken,
			verified: false,
			name: 'My Device',
		})
		const topic = `odinmovieshow-${id.value}`
		const deviceId = d.id
		const data = {
			url: url.value,
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
		verified.value = false
		device_dialog.value?.close()
		devices.value = await getDevices()
	}
	async function getDevices() {
		return await usePb()
			.collection('devices')
			.getFullList({ filter: `(user='${usePb().authStore.model?.id}')` })
	}
	onMounted(async () => {
		url.value = (await usePb().send('/backendurl', { method: 'GET' })).url
		devices.value = await getDevices()
	})
</script>
