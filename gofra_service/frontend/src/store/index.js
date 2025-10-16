import {createStore} from 'vuex'
import auth from './modules/auth'
import listing from './modules/listing'
import user from './modules/user'

export default createStore({
    modules:{
        auth,
        listing,
        user,
    }
})