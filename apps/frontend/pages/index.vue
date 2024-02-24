<template>
	<div>
		<div v-for="e in list">
			<h2>{{ parseTitle(e.title) }}</h2>
			<MediaList :url="e.url" />
		</div>
	</div>
</template>
<script lang="ts" setup>
	const list = useProfile().me['trakt_sections']['home']
</script>
