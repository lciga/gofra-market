import { authAPI } from "../../utils/api"

const state = {
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
    async login({commit, dispatch}, credentials) {
        try {
            const response = await authAPI.login(credentials)
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

    async fetchProfile({commit, state}) {
        try {
            const response = await authAPI.getProfile()
            commit('SET_USER', response.data)
            if (!state.isAuthenticated && state.token) {
                state.isAuthenticated = true
            }
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