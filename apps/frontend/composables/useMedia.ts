export const useMedia = defineStore('useMedia', () => {
	const detail = ref<{ [id: string]: any }>({})

	const lists = ref<{ [id: string]: [] }>({})

	async function getDetail(id: string, type: string) {
		const trakt = id.split('-').at(-1) || ''
		if (!detail.value[id] && id) {
			const res = await usePb().send(`/_trakt/search/trakt/${trakt}?type=${type}`, { method: 'GET' })
			if (res.length > 0) {
				detail.value[id] = res.map((m: any) => fixItem(m))[0]
			}
		}
		return detail.value[id]
	}

	async function search(term: string) {
		if (term.length < 2) {
			return []
		}
		const movies = (await usePb().send(`/_trakt/search/movie?query=${term}`, { method: 'GET' })).map((m: any) => fixItem(m))
		const shows = (await usePb().send(`/_trakt/search/show?query=${term}`, { method: 'GET' })).map((m: any) => fixItem(m))

		return { movies, shows }
	}

	const getId = (item: any) => {
		if (item['type'] === 'episode') {
			return `${item['show']['ids']['slug']}-${item['show']['ids']['trakt']}`
		}
		return `${item['ids']['slug']}-${item['ids']['trakt']}`
	}

	const getLink = (item: any) => {
		if (item['type'] === 'episode') {
			return `/show/${getId(item)}`
		}
		return `/${item['type']}/${getId(item)}`
	}

	const setDetail = (item: any) => {
		if (!detail.value[getId(item)]) {
			if (item.type === 'episode') {
				console.log('episode')
				detail.value[getId(item)] = item.show
			} else {
				detail.value[getId(item)] = item
			}
		}
	}

	const getList = async (url: string) => {
		if (!lists.value[url]) {
			lists.value[url] = (
				await usePb().send(`/_trakt${url}`, {
					method: 'GET',
					cache: 'no-cache',
				})
			).map((m: any) => fixItem(m))
		}
		return lists.value[url]
	}

	function fixItem(m: any) {
		let item = null
		for (const k of Object.keys(m)) {
			if (k === 'movie' || k === 'show' || k === 'episode') {
				item = k
				break
			}
		}

		if (!item) {
			return m
		}

		m[item].type = item
		if (item === 'episode') {
			m[item].tmdb = m[item].show.tmdb
		}

		return m[item]
	}

	return {
		setDetail,
		getDetail,
		getList,
		getLink,
		detail,
		search,
	}
})

if (import.meta.hot) {
	import.meta.hot.accept(acceptHMRUpdate(useMedia, import.meta.hot))
}
