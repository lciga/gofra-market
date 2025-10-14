import { listingAPI } from "../../utils/api"

const state = {
	listings: [],
	currentListing: null,
	filters: {
		rarity: null,
		minPrice: 0,
		maxPrice: 1000000,
	},
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
			const response = await listingAPI.getMarket(state.filters)
			// API returns { items: [], total: N }
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
			commit('ADD_LISTING', response.data)
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
			await dispatch('auth/fetchProfile', null, { root: true })
		} catch (error) {
			console.error('Error buying listing:', error)
			throw error
		}
	},

	async uploadImage({ commit }, { listingID, imageURL }) {
		try {
			await listingAPI.uploadImage(listingID, { url: imageURL })
		} catch (error) {
			console.error('Error upload image:', error)
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