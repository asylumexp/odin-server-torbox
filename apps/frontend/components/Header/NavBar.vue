<template>
	<div class="navbar md:container md:mx-auto">
		<div class="navbar-start">
			<img src="/logo.svg" alt="logo" class="w-16" />
			<form class="join" action="/search" method="GET">
				<input type="text" name="q" placeholder="Search movies or tv shows ..." class="ml-5 input input-sm input-bordered w-64 join-item placeholder-slate-600" />
				<button type="submit" class="join-item btn btn-sm input-bordered"><FaIcon icon="search" /></button>
			</form>
		</div>
		<div class="navbar-center hidden lg:flex">
			<ul class="menu menu-horizontal px-1">
				<li>
					<NuxtLink to="/"><FaIcon icon="fa-home" />Home</NuxtLink>
				</li>
				<li>
					<NuxtLink to="/movies"><FaIcon icon="fa-film" />Movies</NuxtLink>
				</li>
				<li>
					<NuxtLink to="/shows"><FaIcon icon="fa-tv" />TV Shows</NuxtLink>
				</li>
			</ul>
		</div>
		<div class="navbar-end">
			<div class="dropdown dropdown-end">
				<div tabindex="0" role="button" class="btn btn-ghost btn-circle avatar">
					<div class="avatar placeholder">
						<div class="bg-neutral text-neutral-content w-8 rounded-full">
							<span class="text-xs">{{ initials() }}</span>
						</div>
					</div>
				</div>
				<ul tabindex="0" class="menu menu-sm dropdown-content mt-3 z-[50] p-2 shadow bg-base-100 rounded-box w-52">
					<li>
						<NuxtLink to="/" class="disabled">Logged in as {{ useProfile().me?.username }}</NuxtLink>
					</li>
					<li>
						<NuxtLink to="/profile"><FaIcon icon="user" /> Profile </NuxtLink>
					</li>
					<li>
						<NuxtLink to="/settings"><FaIcon icon="gears" />Settings</NuxtLink>
					</li>
					<li @click="linkDevice">
						<NuxtLink to="/devices"><FaIcon icon="computer" />Devices</NuxtLink>
					</li>
					<li class="disabled cursor-pointer">
						<div class="divider"></div>
					</li>
					<li @click="logout" class="text-error">
						<span><FaIcon icon="fa-right-from-bracket" /> Logout </span>
					</li>
				</ul>
			</div>
		</div>
	</div>
</template>

<script lang="ts" setup>
	import AutoComplete from 'primevue/autocomplete'
	const found = ref()

	async function logout() {
		usePb().authStore.clear()
		return navigateTo('/login')
	}

	const initials = () => useProfile().me?.username?.slice(0, 1).toUpperCase()

	async function linkDevice() {}
</script>
