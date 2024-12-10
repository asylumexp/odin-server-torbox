import pocketbase, { LocalAuthStore } from 'pocketbase'

export const usePb = () => {
	return new pocketbase(useRuntimeConfig().public.pbUrl)
}
