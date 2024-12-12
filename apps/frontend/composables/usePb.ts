import pocketbase, { LocalAuthStore } from 'pocketbase'

export const usePb = () => {
	return new pocketbase(useSettings().config?.pbUrl, new LocalAuthStore('odin'))
}
