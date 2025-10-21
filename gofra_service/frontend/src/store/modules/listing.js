import { listingAPI } from "../../utils/api"

const state = {
	listings: [],
	currentListing: null,
	filters: {
		rarity: null,
		minPrice: 0,
		maxPrice: 1000000,
	},
	searchQuery: '',
	loading: false,
}

const mutations = {
	SET_LISTINGS(state, listings) {
		state.listings = listings
	},
	SET_CURRENT_LISTING(state, listing) {
		state.currentListing = listing
	},
	ADD_LISTING(state, listing) {
		state.listings.unshift(listing)
	},
	UPDATE_LISTING(state, updatedListing) {
		const index = state.listings.findIndex(l => l.id === updatedListing.id)
		if (index !== -1) {
			state.listings.splice(index, 1, updatedListing)
		}
	},
	SET_LOADING(state, loading) {
		state.loading = loading
	},
	SET_FILTERS(state, filters) {
		state.filters = { ...state.filters, ...filters }
	},
	SET_SEARCH_QUERY(state, query) {
		state.searchQuery = query
	},
}

const getters = {
	filteredListings: (state) => {
		return state.listings.filter(listing => {
			const matchesRarity = !state.filters.rarity || listing.gofer?.rarity === state.filters.rarity
			const matchesPrice = listing.price >= state.filters.minPrice && listing.price <= state.filters.maxPrice
			return matchesRarity && matchesPrice && !listing.is_sold
		})
	}
}

const actions = {
	async fetchListings({ commit, state }) {
		commit('SET_LOADING', true)
		try {
			let params = { ...state.filters }
			
			if (state.searchQuery && state.searchQuery.trim()) {
				const query = state.searchQuery.trim()
				if (query.startsWith('{') || query.startsWith('[')) {
					params.filter = query
				} else {
					params.filter = JSON.stringify({ 
						"gofer.name": { "$regex": query, "$options": "i" } 
					})
				}
			}
			
			const response = await listingAPI.getMarket(params)
			commit('SET_LISTINGS', response.data.items || [])
		} catch (error) {
			console.error('Error fetching listings:', error)
			throw error
		} finally {
			commit('SET_LOADING', false)
		}
	},

	async fetchListing({ commit }, id) {
		try {
			const response = await listingAPI.getListing(id)
			commit('SET_CURRENT_LISTING', response.data)
		} catch (error) {
			console.error('Error fetching listing:', error)
			throw error
		}
	},

	async createListing({ commit, dispatch }, listingData) {
		try {
			const response = await listingAPI.createListing(listingData)
			await dispatch('fetchListings')
			await dispatch('auth/fetchProfile', null, { root: true })
			return response
		} catch (error) {
			console.error('Error create listing:', error)
			throw error
		}
	},

	async buyListing({ commit, dispatch }, listingID) {
		try {
			await listingAPI.buy({ listing_id: listingID })
			await dispatch('fetchListing', listingID)
			await dispatch('fetchListings')
			await dispatch('auth/fetchProfile', null, { root: true })
		} catch (error) {
			console.error('Error buying listing:', error)
			throw error
		}
	},

	async uploadImageFromUrl({ commit }, { listingID, imageURL }) {
		try {
			await listingAPI.uploadImageFromUrl(listingID, { url: imageURL })
		} catch (error) {
			console.error('Error uploading image from URL:', error)
			throw error
		}
	},

	async uploadImageFile({ commit }, { listingID, file }) {
		try {
			const formData = new FormData()
			formData.append('image', file)
			await listingAPI.uploadImageFile(listingID, formData)
		} catch (error) {
			console.error('Error uploading image file:', error)
			throw error
		}
	},
}

export default {
	namespaced: true,
	state,
	mutations,
	actions,
	getters,
}