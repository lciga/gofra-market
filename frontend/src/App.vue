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
import {defineComponent, onMounted} from 'vue'
import {useStore} from 'vue'

export default defineComponent({
    name: 'App',
    setup() {
        const store = useStore

        onMounted(() => {
            if (localStorage.getItem('token')) {
                store.dispatch('auth/fetchProfile')
            }
        })
    }
})
</script>