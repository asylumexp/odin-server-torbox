<template>
	<div>
		<div class="grid grid-cols-3">
			<div>
				<h2>RealDebrid</h2>
				<RealDebrid @success="setRealDebrid" />
			</div>

			<div>
				<h2>AllDebrid</h2>
				<AllDebrid />
			</div>
		</div>
		<div class="divider"></div>
		<UserList />
	</div>
</template>
<script lang="ts" setup>
	definePageMeta({
		layout: 'admin',
	})

	async function setRealDebrid(data: any) {
		useSettings().settings.real_debrid = data
		await useSettings().save()
	}
</script>
