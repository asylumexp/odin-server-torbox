<template>
	<div>
		<h1 class="mb-0">Users</h1>

		<table class="table table-zebra mt-0">
			<thead>
				<tr>
					<td>ID</td>
					<td>Username</td>
					<td>Email</td>
					<td></td>
				</tr>
			</thead>
			<tr v-for="user in users" :key="user.id">
				<td>{{ user.id }}</td>
				<td>{{ user.username }}</td>
				<td>{{ user.email }}</td>
				<td>
					<button class="btn btn-sm">
						<FaIcon icon="fa-trash" @click="askDelete(user)" />
					</button>
				</td>
			</tr>
		</table>
		<button class="btn btn-primary btn-sm" @click="createUser()"><FaIcon icon="fa-plus" />Add user</button>
		<dialog ref="user_delete_dialog" class="modal">
			<div class="modal-box">
				<h3 class="font-bold text-lg">Delete {{ userToDelete?.username }}?</h3>
				<p class="py-4">Do you really want to delete {{ userToDelete?.email }}?</p>
				<div class="modal-action">
					<form method="dialog">
						<button class="btn btn-primary" @click="confirmDelete()">Yes</button>
						<button class="btn">Cancel</button>
					</form>
				</div>
			</div>
		</dialog>
		<dialog ref="user_create_dialog" class="modal">
			<div class="modal-box">
				<h3 class="font-bold text-lg">Create user</h3>

				<div class="modal-action">
					<form method="dialog" class="md:container" @submit="confirmCreate">
						<div>
							<label class="label w-full max-w-xs">
								<span class="label-text"><FaIcon icon="fa-user" />Username</span>
							</label>
							<input type="text" v-model="newUser.username" placeholder="" class="input input-bordered input-sm w-full max-w-xs" />
						</div>
						<div>
							<label class="label w-full max-w-xs">
								<span class="label-text"><FaIcon icon="fa-envelope" />Email</span>
							</label>
							<input type="email" v-model="newUser.email" placeholder="" class="input input-bordered input-sm w-full max-w-xs" />
						</div>

						<div>
							<label class="label w-full max-w-xs">
								<span class="label-text"><FaIcon icon="fa-key" />Password</span>
							</label>
							<input type="password" v-model="newUser.password" placeholder="" class="input input-bordered input-sm w-full max-w-xs" />
						</div>

						<div>
							<label class="label w-full max-w-xs">
								<span class="label-text"><FaIcon icon="fa-key" />Confirm Password</span>
							</label>
							<input type="password" v-model="newUser.passwordConfirm" placeholder="" class="input input-bordered input-sm w-full max-w-xs" />
						</div>
						<button class="btn btn-primary">Create User</button>
					</form>
				</div>
			</div>
		</dialog>
	</div>
</template>

<script lang="ts" setup>
	let users = ref((await usePb().collection('users').getList()).items)
	const defaultUser = {
		username: '',
		email: '',
		password: '',
		passwordConfirm: '',
		verified: true,
		trakt_sections: {
			home: [
				{
					big: true,
					paginate: true,
					title: 'Trending movies',
					url: '/movies/trending',
				},
				{
					big: true,
					paginate: true,
					title: 'Trending shows',
					url: '/shows/trending',
				},
				{
					big: true,
					paginate: false,
					title: 'Your Todays episodes',
					url: '/calendars/my/shows/::year::-::month:-1:-::day::/::monthdays::',
				},
				{
					big: true,
					paginate: false,
					title: 'Your tomorrows episodes',
					url: '/calendars/my/shows/::year::-::month::-::day:+1:/1',
				},
			],
			movies: [
				{
					big: false,
					paginate: true,
					title: 'Most Watched Today',
					url: '/movies/watched/daily',
				},
				{
					big: false,
					paginate: true,
					title: 'Popular ::year::/::year:-1: releases',
					url: '/movies/popular?years=::year::,::year:-1:',
				},
				{
					big: false,
					paginate: false,
					title: 'Box Office',
					url: '/movies/boxoffice',
				},
				{
					big: false,
					paginate: false,
					title: 'Your watchlist',
					url: '/sync/watchlist/movies/title',
				},
				{
					big: false,
					paginate: true,
					title: 'Highly anticipated',
					url: '/movies/anticipated',
				},
			],
			shows: [
				{
					big: false,
					paginate: true,
					title: 'Most Watched today',
					url: '/shows/watched/daily',
				},
				{
					big: false,
					paginate: true,
					title: 'Popular ::year::/::year:-1: releases',
					url: '/shows/popular/?years=::year::,::year:-1:',
				},
				{
					big: false,
					paginate: false,
					title: 'Your watchlist',
					url: '/sync/watchlist/shows/title',
				},
				{
					big: false,
					paginate: true,
					title: 'Highly anticipated',
					url: '/shows/anticipated',
				},
				{
					big: true,
					paginate: false,
					title: 'Recently watched',
					url: '/sync/watched/shows',
				},
			],
		},
	}
	let newUser = defaultUser

	let userToDelete = ref<any>()
	let user_delete_dialog = ref<HTMLDialogElement>()
	let user_create_dialog = ref<HTMLDialogElement>()
	async function askDelete(user: any) {
		userToDelete.value = user
		user_delete_dialog.value?.showModal()
	}
	async function confirmDelete() {
		if (userToDelete.value) {
			await usePb().collection('users').delete(userToDelete.value.id)
			users.value = users.value.filter((u: any) => u.id !== userToDelete.value.id)
			userToDelete = ref(null)
		}
	}

	async function createUser() {
		user_create_dialog.value?.showModal()
	}

	async function confirmCreate(e: Event) {
		e.preventDefault()
		try {
			const user = await usePb().collection('users').create(newUser)
			users.value.push(user as any)
			user_create_dialog.value?.close()
			newUser = defaultUser
		} catch (e) {}
	}
</script>

<style scoped>
	span svg {
		margin-right: 0.5rem;
	}
</style>
