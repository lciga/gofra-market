import { createApp } from 'vue'
import { Quasar, Notify, Dialog } from 'quasar'
import quasarLang from 'quasar/lang/ru'
import quasarIconSet from 'quasar/icon-set/material-icons'

import '@quasar/extras/material-icons/material-icons.css'
import 'quasar/src/index.sass'

import App from './App.vue'
import router from './router'
import store from './store'

const app = createApp(App)

app.use(Quasar, {
    plugins: { Notify, Dialog },
    lang: quasarLang,
    iconSet: quasarIconSet,
})

app.use(store)
app.use(router)

app.mount('#app')