<template>
	<dialog class="modal" ref="video_dialog">
		<VideoPlayer v-if="dialogOpen" :src="useStreams().videoUrl" autoplay controls type="video/mp4" />
		<div v-else></div>
	</dialog>
</template>

<style></style>

<script lang="ts" setup>
	import { VideoPlayer } from '@videojs-player/vue'
	import 'video.js/dist/video-js.css'
	const video_dialog = ref<HTMLDialogElement>()
	const dialogOpen = ref(false)
	watch(
		() => useStreams().triggerVideoModal,
		async () => {
			dialogOpen.value = true
			console.log('video', useStreams().videoUrl)
			video_dialog.value?.showModal()
		}
	)

	onMounted(() => {
		// player = videojs(video_player.value!, {
		// 	autoplay: true,
		// 	controls: true,
		// 	sources: [
		// 		{
		// 			src: 'https://vjs.zencdn.net/v/oceans.mp4',
		// 			type: 'video/mp4',
		// 		},
		// 	],
		// })

		video_dialog.value!.onclose = () => {
			console.log('close')
			dialogOpen.value = false
		}
	})
</script>
