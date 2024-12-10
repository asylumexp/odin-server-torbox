// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  ssr: false,
  devtools: { enabled: true },
  modules: ['nuxt-primevue', '@nuxtjs/tailwindcss', '@pinia/nuxt'],
  css: ['primevue/resources/themes/lara-dark-indigo/theme.css'],
  runtimeConfig: {
    public: {
      pbUrl: 'http://127.0.0.1:8090',
    },
  },
  compatibilityDate: '2024-07-03',
  devServer: {
    host: '0.0.0.0'
  },
  app: {
    pageTransition: { name: 'page', mode: 'out-in' },
    layoutTransition: { name: 'layout', mode: 'out-in' },
  },
})
