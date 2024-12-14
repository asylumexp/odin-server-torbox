<template>
	<div>
		<h1>Settings</h1>
		<form @submit="submit">
			<div class="grid grid-cols-3">
				<div>
					<h2>Trakt</h2>
					<label class="label">Client ID: </label>
					<input type="text" v-model="useSettings().settings.trakt.clientId" placeholder="trakt clientId" class="input input-sm input-bordered mr-2" />
					<label class="label">Client Secret: </label>
					<input type="text" v-model="useSettings().settings.trakt.clientSecret" placeholder="trakt clientSecret" class="input input-sm input-bordered mr-2" />
				</div>
				<div>
					<h2>TMDB</h2>
					<label class="label">Key: </label>
					<input type="text" v-model="useSettings().settings.tmdb.key" placeholder="tmdb key" class="input input-sm input-bordered mr-2" />
				</div>
				<div>
					<h2>Scraper</h2>
					<label class="label">URL: </label>
					<input type="text" v-model="useSettings().settings.scraper_url" placeholder="http://odin-scraper:6969" class="input input-sm input-bordered mr-2" />
				</div>

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
