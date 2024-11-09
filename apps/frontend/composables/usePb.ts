import pocketbase, { LocalAuthStore } from 'pocketbase'

export const usePb = () => {
	console.log('USEPB', useSettings().config?.pbUrl)
	return new pocketbase(useSettings().config?.pbUrl, new LocalAuthStore('odin-frontend'))
	// return new pocketbase(useRuntimeConfig().public.pbUrl)
}
