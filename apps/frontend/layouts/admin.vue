<template>
  <div class="md:container md:mx-auto">
    <div class="navbar bg-base-100">
      <div class="navbar-start">
        <div class="dropdown">
          <div tabindex="0" role="button" class="btn btn-ghost btn-circle">
            <h1>ODIN</h1>
          </div>
        </div>
      </div>
      <div class="navbar-center hidden lg:flex">
        <ul class="menu menu-horizontal px-1">
          <li>
            <NuxtLink to="/admin/users"
              ><FaIcon icon="fa-users" />Users</NuxtLink
            >
          </li>
          <li>
            <NuxtLink to="/admin/settings"
              ><FaIcon icon="fa-gears" />Settings</NuxtLink
            >
          </li>
        </ul>
      </div>
      <div class="navbar-end">
        <div class="dropdown dropdown-end">
          <div tabindex="0" class="avatar placeholder cursor-pointer">
            <div class="bg-neutral text-neutral-content rounded-full w-16">
              <span class="text-xl">AD</span>
            </div>
          </div>
          <ul
            tabindex="0"
            class="mt-3 z-[1] p-2 shadow menu menu-sm dropdown-content bg-base-100 rounded-box w-52"
          >
            <li @click="logout">
              <span><FaIcon icon="fa-right-from-bracket" /> Logout </span>
            </li>
          </ul>
        </div>
      </div>
    </div>
    <div class="prose md:container md:mx-auto"><slot /></div>
  </div>
</template>

<script lang="ts" setup>
await useSettings().init();
async function logout() {
  usePb().authStore.clear();
  return navigateTo("/login");
}
</script>
