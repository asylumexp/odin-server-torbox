<template>
	<div>
		<button class="btn btn-sm btn-primary mt-10" @click="openModal"><FaIcon icon="computer" /> Link Device</button>
		<table class="table">
			<tbody>
				<tr v-for="device in devices" :key="device.id">
					<td><FaIcon icon="check-circle" class="mr-2" :class="device.verified ? 'text-success' : ''" />{{ device.id }}</td>
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
					<p class="text-sm m-0 mt-10">The following URL will be sent to the device:</p>
					<p class="text-success text-sm m-0 mt-5">{{ url }}</p>
					<p class="text-xs m-0 mt-1 opacity-50">Make sure the device can access it.</p>
					<p>Please enter the code shown on your TV app</p>
					<input type="text" class="input input-bordered mr-5" v-model="id" />
					<button class="btn btn-md btn-primary" v-if="!verified && !loading" @click="linkDevice">Connect</button>
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
		id.value = ''
	}

	async function linkDevice() {
		loading.value = true

		const d = await usePb().collection('devices').create({
			user: usePb().authStore.model?.id,
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
			.getFullList({ filter: `(user='${usePb().authStore.model?.id}')`, sort: '-created' })
	}
	onMounted(async () => {
		url.value = `${location.protocol}//${location.host}`
		devices.value = await getDevices()
	})
</script>
