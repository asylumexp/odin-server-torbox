<template>
	<div id="app" data-theme="odin">
		<NuxtLayout>
			<NuxtPage />
		</NuxtLayout>
	</div>
</template>

<style>
	.page-enter-active,
	.page-leave-active {
		transition: all 0.2s;
	}
	.page-enter-from,
	.page-leave-to {
		opacity: 0;
		filter: blur(1rem);
	}

	.layout-enter-active,
	.layout-leave-active {
		transition: all 0.2s;
	}
	.layout-enter-from,
	.layout-leave-to {
		opacity: 0;
		filter: blur(1rem);
	}

	#app {
		min-height: 100vh;
	}
</style>
