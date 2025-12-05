import { createApp } from 'vue'
import { Quasar, Notify, Dialog } from 'quasar'
import quasarLang from 'quasar/lang/ru'
import quasarIconSet from 'quasar/icon-set/material-icons'

import '@quasar/extras/material-icons/material-icons.css'
import 'quasar/src/index.sass'

import App from './App.vue'
import router from './router'
import store from './store'

if (typeof window !== 'undefined') {
    const isInteractiveElement = (target) => {
        if (!target || !(target instanceof Element)) {
            return false
        }
        return Boolean(target.closest('input, textarea, [contenteditable="true"], .allow-copy'))
    }

    const blockEvent = (event) => {
        if (isInteractiveElement(event.target)) {
            return
        }
        event.preventDefault()
    }

    ['copy', 'cut', 'contextmenu'].forEach((evt) => {
        document.addEventListener(evt, blockEvent)
    })

    document.addEventListener('selectstart', (event) => {
        if (isInteractiveElement(event.target)) {
            return
        }
        event.preventDefault()
    })
}

const app = createApp(App)

app.use(Quasar, {
    plugins: { Notify, Dialog },
    lang: quasarLang,
    iconSet: quasarIconSet,
})

app.use(store)
app.use(router)

app.mount('#app')