import { boot } from 'quasar/wrappers'
import api from '../utils/api'

export default boot(({ app }) => {
    app.config.globalProperties.$api = api
})