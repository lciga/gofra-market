<template>
    <q-header elevated class="bg-primary text-white">
        <q-toolbar>
            <q-toolbar-title class="cursor-pointer" @click="$router.push('/')">
                <q-avatar>
                    <img src="../../../image/gofer-placeholder.png" alt="Gofra Market"></img>
                </q-avatar>
                Gofra Market
            </q-toolbar-title>

            <div v-if="isAuthenticated" class="row items-center q-gutter-sm">
                <q-btn flat icon="mdi-account" :label="user.login" @click="$router.push('/profile')"></q-btn>
                <q-bange color="accent" class="q-pa-sm">Баланс: {{ formatPrice(user.balance) }} горутин</q-bange>
                <q-btn flat icon="mdi-logout" @click="handleLogout"></q-btn>
            </div>
            <div v-else>
                <q-btn flat label="Войти" to="/login"></q-btn>
                <q-btn flat label="Регистрация" to="/register"></q-btn>
            </div>
        </q-toolbar>
    </q-header>
</template>

<script>
import {computed} from 'vue'
import {useStore} from 'vuex'
import {useRouter} from 'veu-router'
import {formatPrice} from '../../utils/formatters'

export default {
    name: 'AppHeader',
    setup() {
        const store = useStore
        const router = useRouter

        const isAuthenticated = computed(() => store.state.auth.isAuthenticated)
        const user = computed(() => store.state.auth.user || {})

        const handleLogout = () => {
            store.dispatch('auth/logout')
            router.push('/')
        }

        return {
            isAuthenticated, 
            user,
            formatPrice,
            handleLogout
        }
    }
}
</script>