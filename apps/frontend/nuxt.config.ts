// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
	ssr: false,
	devtools: { enabled: true },
	modules: ['nuxt-primevue', '@nuxtjs/tailwindcss', '@pinia/nuxt'],

	runtimeConfig: {
		public: {
			pbUrl: process.env.PB_URL || 'http://localhost:8090',
		},
	},

	css: ['primevue/resources/themes/lara-dark-indigo/theme.css'],

	app: {
		pageTransition: { name: 'page', mode: 'out-in' },
		layoutTransition: { name: 'layout', mode: 'out-in' },
	},
})
