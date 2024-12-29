import pocketbase, { LocalAuthStore } from 'pocketbase'

export const usePb = () => {
	return new pocketbase('/', new LocalAuthStore('odin'))
}
