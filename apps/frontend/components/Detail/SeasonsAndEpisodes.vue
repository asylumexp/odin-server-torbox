<template>
	<!-- <progress class="progress progress-success" :value="watchedEpisodes" :max="totalEpisodes" /> -->
	<h3>Seasons & Episodes</h3>
	<div class="gap-2 flex flex-wrap">
		<button v-for="s of item.seasons" class="btn btn-sm" :class="s.number == sNumber ? 'btn-primary' : 'btn-outline'" @click="sNumber = s.number">
			{{ s.number == 0 ? 'Specials' : s.number }}
			<!-- <FaIcon icon="check" /> -->
		</button>
	</div>
	<div v-if="tmdbSeasons.length === 0">Loading</div>
	<div v-else>
		<h3>Episodes of {{ getSeasonByNumber().title }}</h3>
		<div class="grid grid-cols-5 gap-5">
			<div
				v-for="(e, k) in getSeasonByNumber().episodes"
				class="cursor-pointer card card-compact shadow-lg bg-base-300 bg-opacity-70 no-underline hover:bg-opacity-10 hover:bg-primary transition-all"
				@click="useStreams().openModal(e, item, getSeasonByNumber().number)"
			>
				<figure class="m-0">
					<img :src="getStill(k)" />
				</figure>
				<div class="card-body">
					<p class="card-title m-0 text-sm">S{{ sNumber }}xE{{ e.number }} - {{ e.title }}</p>
					<p class="text-sm m-0">{{ e.first_aired }}</p>
					<FaIcon v-if="e.watched" icon="check-circle" class="text-success" />
					<!-- <p class="text-sm">{{ e.overview }}</p> -->
				</div>
			</div>
		</div>
	</div>
</template>

<script lang="tsx" setup>
	const props = defineProps({
		item: {
			type: Object,
			required: true,
		},
	})

	if (typeof props.item === 'undefined') {
		useRouter().back()
	}

	const item = props.item

	const sNumber = ref(item.seasons[0].number)

	const getSeasonByNumber = () => {
		return item.seasons.filter((s: any) => s.number === sNumber.value)[0]
	}

	function seasonColor(s: any) {
		const watched = s.episodes.filter((e: any) => e.watched).length
		if (watched > 0) {
			return watched == s.episodes.length ? 'tab-success' : 'tab-secondary'
		} else {
		}
		return 'tab-neutral'
	}

	const totalEpisodes = item.seasons.reduce((acc: number, s: any) => acc + s.episodes.length, 0)

	const watchedEpisodes = item.seasons.reduce((acc: number, s: any) => acc + s.episodes.filter((e: any) => e.watched).length, 0)

	const tmdbSeasons = ref<any[]>([])

	onMounted(async () => {
		if (tmdbSeasons.value.length === 0) {
			tmdbSeasons.value = await usePb().send('/tmdbseasons/' + item.ids.tmdb + '?seasons=' + item.seasons.map((s: any) => s.number).join(','), { method: 'GET' })
			console.log(tmdbSeasons.value)
		}
	})

	function getStill(e: number) {
		const path = tmdbSeasons.value.filter((s: any) => s?.season_number == sNumber.value)[0]?.episodes[e]?.still_path

		if (!path) {
			return 'http://imageipsum.com/1200x675'
		}

		return 'https://image.tmdb.org/t/p/w780' + path
	}
</script>
