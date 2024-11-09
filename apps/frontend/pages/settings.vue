<template>
	<div>
		<h1>Settings</h1>

		<button class="btn btn-primary" @click="saveSections()">Save sections</button>
		<div class="grid grid-cols-2 gap-10">
			<div>
				<h2>Sections</h2>
				<section v-for="(_, k) in sections" class="border border-dashed border-slate-700 mb-5 p-4">
					<h3>{{ k }}</h3>
					<draggable v-model="sections[k]" :group="k" item-key="id" handle=".handle" ghost-class="ghost">
						<template #item="{ element, index }">
							<div class="flex mb-2">
								<FaIcon icon="list" class="handle mt-2 mr-2 opacity-30" size="sm" />
								<input class="input input-sm input-bordered mr-2 flex-1" v-model="element.title" placeholder="Title" />
								<input class="input input-sm input-bordered flex-grow" placeholder="URL" v-model="element.url" />
								<label class="cursor-pointer label">
									<span class="label-text mr-2">Big</span>
									<input type="checkbox" v-model="element.big" class="checkbox checkbox-secondary checkbox-sm" />
								</label>
								<label class="cursor-pointer label">
									<span class="label-text mr-2">Paginate</span>
									<input type="checkbox" v-model="element.paginate" class="checkbox checkbox-secondary checkbox-sm" />
								</label>
								<button class="btn btn-sm ml-2" @click="deleteSection(k, index)">
									<FaIcon icon="trash" />
								</button>
							</div>
						</template>
					</draggable>

					<button class="btn" @click="addSection(k)">Add</button>
				</section>
			</div>
			<div>
				<h2>Templates</h2>
				<draggable v-model="templates" :group="{ name: 'movies', pull: 'clone', put: false }" item-key="id" handle=".handle" ghost-class="ghost">
					<template #item="{ element, index }">
						<div class="flex mb-2">
							<FaIcon icon="list" class="handle mt-2 mr-2 opacity-30" size="sm" />
							<span>Title: {{ element.title }} - URL: {{ element.url }}</span>
						</div>
					</template>
				</draggable>
			</div>
		</div>
	</div>
</template>

<style scoped>
	.ghost {
		@apply bg-primary bg-opacity-10;
	}
</style>

<script lang="ts" setup>
	const me = useProfile().me
	let sections = me['trakt_sections']
	import draggable from 'vuedraggable'

	const templates = [
		{ title: 'Movies', url: '/movies', big: false, paginate: false },
		{ title: 'Shows', url: '/movies', big: false, paginate: false },
	]

	if (sections === null) {
		sections = {
			home: [],
			movies: [],
			shows: [],
		}
		usePb().collection('users').update(me.id, { trakt_sections: sections })
	}
	function addSection(place: any) {
		sections[place].push({
			title: '',
			url: '',
			big: false,
			paginate: false,
		})
	}

	function deleteSection(place: any, index: number) {
		console.log(place, index)
		sections[place].splice(index, 1)
	}

	function saveSections() {
		usePb().collection('users').update(me.id, { trakt_sections: sections })
	}
</script>
