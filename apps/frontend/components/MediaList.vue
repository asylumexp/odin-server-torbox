<template>
	<div v-if="data.length === 0 && loading" class="text-center h-80">
		<span class="loading loading-bars text-primary loading-lg"></span>
	</div>
	<div v-if="data.length === 0 && !loading" class="text-center h-80">
		<div class="border opacity-30 border-dashed p-3 rounded-lg">No items</div>
	</div>
	<div v-else class="flex overflow-y-scroll gap-4 h-80">
		<NuxtLink :to="useMedia().getLink(item)" v-for="item in data" @click="useMedia().setDetail(item)" class="no-underline">
			<div class="card w-32 bg-base-300 shadow-md flex-shrink-0 rounded-md h-80">
				<figure class="m-0 p-0">
					<img v-if="item.tmdb && item.tmdb.poster_path" :src="'https://image.tmdb.org/t/p/w780' + item.tmdb.poster_path" />
					<img v-else src="https://placehold.co/160x250/11111b/ffffff/png?text=no\nimage" />
				</figure>
				<div class="card-body p-3 pb-5">
					<p class="card-title m-0 text-sm">
						<span class="w-96 break-words overflow-hidden max-h-10">{{ item['title'] }}</span>
						<FaIcon v-if="item['watched']" icon="check-circle" class="text-success" />
					</p>
					<span v-if="item.year" class="badge">{{ item.year }}</span>
				</div>
			</div>
		</NuxtLink>
	</div>
</template>
<script lang="ts" setup>
	const props = defineProps({
		url: String,
	})

	const loading = ref(true)
	const data = ref<any[]>([])

	onMounted(async () => {
		data.value = await useMedia().getList(props.url!)
		loading.value = false
	})
</script>
