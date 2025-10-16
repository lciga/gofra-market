<template>
    <q-header elevated class="bg-primary text-white">
        <q-toolbar>
            <q-toolbar-title class="cursor-pointer" @click="$router.push('/')">
                <q-avatar>
                    <img :src="goferPlaceholder" alt="Gofra Market" />
                </q-avatar>
                Gofra Market
            </q-toolbar-title>

            <div v-if="isAuthenticated && user" class="row items-center q-gutter-sm">
                <q-btn flat icon="person" :label="user.login || 'Профиль'" @click="$router.push('/profile')" />
                <q-badge color="accent" class="q-pa-sm">Баланс: {{ formatPrice(user.balance || 0) }} горутин</q-badge>
                <q-btn flat icon="logout" label="Выход" @click="handleLogout" />
            </div>
            <div v-else>
                <q-btn flat label="Войти" to="/login" />
                <q-btn flat label="Регистрация" to="/register" />
            </div>
        </q-toolbar>
    </q-header>
</template>

<script>
import { computed, defineComponent } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import { formatPrice } from '../../utils/formatters'
import goferPlaceholder from 'assets/gofer-placeholder.png'

export default defineComponent({
    name: 'AppHeader',
    setup() {
        const store = useStore()
        const router = useRouter()

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
            handleLogout,
            goferPlaceholder,
        }
    },
})
</script>