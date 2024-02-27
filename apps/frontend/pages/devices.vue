<template>
	<div>
		<h1>Devices</h1>
		<button class="btn btn-primary" @click="linkDevice"><FaIcon icon="computer" /> Link Device</button>
		<table class="table">
			<tbody>
				<tr v-for="device in devices" :key="device.id">
					<td><FaIcon icon="check-circle" class="mr-2" :class="device.verified ? 'text-success' : ''" />{{ device.token }}</td>
					<td>{{ device.name }}</td>
					<td>{{ device.created }}</td>
				</tr>
			</tbody>
		</table>
	</div>
</template>

<script lang="ts" setup>
	const devices = ref()

	function generateNtfyId() {
		return Math.random().toString(36).substring(2).toUpperCase().substring(0, 6)
	}

	async function generateID() {
		const id = Math.random().toString(36).substring(2).toUpperCase()
		const existing = await usePb()
			.collection('devices')
			.getFullList({ filter: `(token='${id}')` })
		if (existing.length > 0) {
			return generateID()
		} else {
			return id
		}
	}
	async function linkDevice() {
		const url = useRequestURL().origin
		const id = generateNtfyId()
		const deviceToken = await generateID()

		const d = await usePb()
			.collection('devices')
			.create({
				user: usePb().authStore.model?.id,
				token: deviceToken,
				verified: false,
				name: 'My Device',
			})
		const topic = `odinmovieshows-${id}`
		const deviceId = d.id
		const data = {
			url,
			deviceId,
		}

		let verified = false

		while (!verified) {
			verified = (await usePb().collection('devices').getOne(deviceId)).verified === true
			if (verified) {
				return
			}
			console.log(deviceId)
			// await useFetch(`https://ntfy.sh/${topic}`, { method: 'POST', body: JSON.stringify(data) })
			await new Promise((resolve) => setTimeout(resolve, 5000))
		}

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
