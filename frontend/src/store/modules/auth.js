import { authAPI } from "../../utils/api"

const status = {
    user: null,
    token: localStorage.getItem('token'),
    isAuthenticated: !!localStorage.getItem('token')
}

const mutations = {
    SET_USER(state, user) {
        state.user = user
        state.isAuthenticated = true
    },
    SET_TOKEN(state, token) {
        state.token = token
        localStorage.setItem('token', token)
    },
    LOGOUT(state) {
        state.user = null
        state.token = null
        state.isAuthenticated = false
        localStorage.removeItem('token')
    }
}

const actions = {
    async login({commit, dispatch}, credenials) {
        try {
            const response = await authAPI.login(credenials)
            commit('SET_TOKEN', response.data.token)
            await dispatch('fetchProfile')
            return response
        } catch(error) {
            commit('LOGOUT')
            throw error
        }
    },

    async register({commit, dispatch}, userData) {
        try {
            const response = await authAPI.register(userData)
            commit('SET_TOKEN', response.data.token)
            await dispatch('fetchProfile')
            return response
        } catch (error) {
            commit('LOGOUT')
            throw error
        }
    },

    async fetchProfile({commit}) {
        try {
            const response = await authAPI.getPorfile()
            commit('SET_USER', response.data)
            return response
        } catch (error) {
            commit('LOGOUT')
            throw error
        }
    },

    logout({commit}) {
        commit('LOGOUT')
    }
}

export default {
    namespaced: true,
    state,
    mutations,
    actions,
}