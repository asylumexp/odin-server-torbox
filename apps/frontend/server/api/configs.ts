export default defineEventHandler((event) => {
	return { pbUrl: process.env.NUXT_PB_URL, mqtt: { url: process.env.MQTT_URL, user: process.env.MQTT_USER, pass: process.env.MQTT_PASSWORD } }
})
