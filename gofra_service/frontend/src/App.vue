<template>
    <q-layout view="hHh Lpr fFf">
        <app-header/>

        <q-page-container>
            <router-view/>
        </q-page-container>

        <app-footer/>
    </q-layout>
</template>

<script>
import { defineComponent, onMounted } from 'vue'
import { useStore } from 'vuex'

import AppHeader from './components/layout/AppHeader.vue'
import AppFooter from './components/layout/AppFooter.vue'

export default defineComponent({
    name: 'App',
    components: {
        AppHeader,
        AppFooter,
    },
    setup() {
        const store = useStore()

        onMounted(() => {
            if (localStorage.getItem('token')) {
                store.dispatch('auth/fetchProfile')
            }
        })
    }
})
</script>