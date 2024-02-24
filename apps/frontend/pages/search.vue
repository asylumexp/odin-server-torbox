<template>
	<div>
		<h2>Movies</h2>
		<MediaList :url="`/search/movie?limit=10&query=${useRoute().query['q']}`" />
		<h2>Shows</h2>
		<MediaList :url="`/search/show?limit=10&query=${useRoute().query['q']}`" />
	</div>
</template>
<script lang="ts" setup></script>
