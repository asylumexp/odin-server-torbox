export default defineEventHandler((event) => {
	return { pbUrl: process.env.NUXT_PB_URL }
})
