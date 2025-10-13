const state = {
    balance: 0
}

const mutation = {
    SET_BALANCE(state, balance){
        state.balance = balance
    }
}

const actions = {
    updateBalance({commit}, balance) {
        commit('SET_BALANCE', balance)
    }
}

export default {
    namespaced: true,
    state,
    mutation,
    actions,
}