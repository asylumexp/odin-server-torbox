import pocketbase from 'pocketbase'

export const usePb = () => {
	console.log('USEPB', useSettings().config?.pbUrl)
	return new pocketbase(useSettings().config?.pbUrl)
	// return new pocketbase(useRuntimeConfig().public.pbUrl)
}
