<template>
	<dialog ref="streams_dialog" class="modal" @close="dialogGone">
		<div class="modal-box min-w-max">
			<h3 class="m-0">Streams</h3>
			<StreamList q="4K" />
			<StreamList q="1080p" />
			<StreamList q="720p" />
			<StreamList q="SD" />
			<StreamList q="CAM" />
		</div>
	</dialog>
</template>

<script lang="tsx" setup>
	import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

	const StreamList = (props: { q: '4K' | '1080p' | '720p' | 'SD' | 'CAM' }) =>
		streams.value[props.q].length === 0 ? (
			<div></div>
		) : (
			<div>
				{streams.value[props.q].map((stream) => (
					<div>
						{(stream.realdebrid || []).map((s: any) => (
							<div>
								<p>
									<span class={'badge mr-2 ' + badgeClass(props.q)}>{stream.quality}</span>
									<span>
										{stream.info.join(' | ')} | {humanFileSize(s.filesize)}
									</span>
									<br />
									<span class="text-sm">{s.filename}</span>
									<br />
									<span class="btn btn-xs btn-secondary btn-outline mr-3" onClick={() => copyToClipbaord(s.download)}>
										<FontAwesomeIcon icon="clipboard" class="mr-1" />
										Copy to clipboard
									</span>
								</p>

								<div class="divider"></div>
							</div>
						))}
					</div>
				))}
			</div>
		)

	function copyToClipbaord(url: string) {
		navigator.clipboard.writeText(url)
	}

	function badgeClass(q: string) {
		switch (q) {
			case '4K':
				return 'badge-success'
			case '1080p':
				return 'badge-primary'
			case '720p':
				return 'badge-info'
			case 'SD':
				return 'badge-warning'
			case 'CAM':
				return 'badge-error'
			default:
				return 'badge-secondary'
		}
	}
	const streams_dialog = ref<HTMLDialogElement>()

	const data = ref<any[]>([])

	function humanFileSize(size: number): number | string {
		var i = size == 0 ? 0 : Math.floor(Math.log(size) / Math.log(1024))
		return parseFloat((size / Math.pow(1024, i)).toFixed(2)) * 1 + ' ' + ['B', 'kB', 'MB', 'GB', 'TB'][i]
	}

	const streams = ref({
		'4K': data.value.filter((d) => d.quality === '4K'),
		'1080p': data.value.filter((d) => d.quality === '1080p'),
		'720p': data.value.filter((d) => d.quality === '720p'),
		'SD': data.value.filter((d) => d.quality === 'SD'),
		'CAM': data.value.filter((d) => d.quality === 'CAM'),
	})

	watch(data, () => {
		streams.value = {
			'4K': data.value.filter((d) => d.quality === '4K'),
			'1080p': data.value.filter((d) => d.quality === '1080p'),
			'720p': data.value.filter((d) => d.quality === '720p'),
			'SD': data.value.filter((d) => d.quality === 'SD'),
			'CAM': data.value.filter((d) => d.quality === 'CAM'),
		}
	})

	watch(
		() => useStreams().streams,
		() => {
			data.value = useStreams().streams
		}
	)

	watch(
		() => useStreams().triggerModal,
		async () => {
			console.log('dialog open')
			streams_dialog.value?.showModal()
			data.value = (await useStreams().getStreams()) || []
		}
	)

	async function dialogGone() {
		useStreams().mqttClient.unsubscribe(useStreams().topic)
	}
</script>
