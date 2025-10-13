import { listingAPI } from "../../utils/api"

const state = {
	listings: [],
	currentListing: null,
	filters: {
		rarity: null,
		minPrice: 0,
		maxPrice: 1000,
	},
	loading: false,
}

const mutations = {
	SET_LISTING(state, listing) {
		state.listing = listing
	},
	SET_CURRENT_LISTING(state, listing) {
		state.currentListing = listing
	},
	ADD_LISTING(state, listing){
		state.listing.unshift(listing)
	},
	UPDATE_LISTING(state, updatedListing) {
		const index = state.listing.findIndex(l => l.id === updatedListing.id)
		if (index !== -1) {
			state.listing.splice(index, 1, updatedListing)
		}
	},
	SET_LOADING(state, loading){
		state.loading = loading
	},
	SET_FILTERS(state, filters){
		state.filters = {...state.filters, ...filters}
	},
}

const getters = {
	filteredListing: (state) => {
		return state.listing.filter(listing =>{
			const matchesRarity = !state.filters.rarity || listing.gofer.rarity === state.filter.rarity
			const matchesPrice = listing.price >= state.filter.minPrice && listing.price <= listing.filter.maxPrice
			return matchesRarity, matchesPrice, !listing.is_sold
		})
	}
}

const actions = {
	async fetchListings({commit, state}){
		commit('SET_LOADING', true)
		try {
			const response = await listingAPI.getMarket(state.filters)
			commit('SET_LISTINGS', response.data)
		} catch (error) {
			console.error('Error fetching listings:', error)
			throw error
		} finally {
			commit('SET_LOADING', false)
		}
	},

	async fetchListing({commit}, id){
		try{
			const response = await listingAPI.getListing(id)
			commit('SET_CURRENT_LISTING', response.data)
		} catch(error) {
			console.error('Error fetching listing:', error);
			throw error
		}
	},

	async createListing({commit, dispatch}, listingData){
		try {
			const response = await listingAPI.createListing(listingData)
			commit('ADD_LISTING', response.data)
			await dispatch('auth/fetchProfile', null, {root:true})
			return response
		} catch (error) {
			console.error('Error create listing:', error);
			throw error
		}
	},

	async buyListing({commit, dispatch}, listingID){
		try {
			await listingAPI.buy({listing_id: listingID})
			await dispatch('fetchListing')
			await dispatch('auth/fetchProfile', null, {root: true})
		} catch (error) {
			console.error('Error buying listing:', error)
			throw error
		}
	},

	async uploadImage({commit}, {listingID, imageURL}) {
		try {
			await listingAPI.uploadImage(listingID, {url: imageURL})
		} catch (error) {
			console.error('Error upload image:', error);
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