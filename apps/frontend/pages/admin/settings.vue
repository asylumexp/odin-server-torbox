<template>
	<div>
		<h1>Settings</h1>
		<form @submit="submit">
			<div class="grid grid-cols-3">
				<div>
					<h2>RealDebrid</h2>
					<RealDebrid @success="setRealDebrid" />
				</div>

				<div>
					<h2>AllDebrid</h2>
					<AllDebrid @success="setAllDebrid" />
				</div>
			</div>
			<div class="divider"></div>
			<button type="submit" class="btn btn-primary">Save</button>
		</form>
	</div>
</template>

<script lang="ts" setup>
	definePageMeta({
		layout: 'admin',
	})

	async function submit(e: Event) {
		e.preventDefault()
		useSettings().save()
	}

	async function setAllDebrid(data: any) {}

	async function setRealDebrid(data: any) {
		useSettings().settings.real_debrid = data
		await useSettings().save()
	}
</script>
