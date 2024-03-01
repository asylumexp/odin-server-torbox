import pocketbase from 'pocketbase'

export const usePb = () => {
	return new pocketbase(useRuntimeConfig().public.pbUrl)
}
