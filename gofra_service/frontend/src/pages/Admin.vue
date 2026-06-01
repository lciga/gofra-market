<template>
    <q-page class="q-pa-lg bg-grey-1">
        <div class="row items-center q-mb-lg">
            <div class="col">
                <div class="text-h4 text-primary">Админская панель</div>
                <div class="text-grey-7">Пользователи и статистика посещений</div>
            </div>
            <q-btn color="primary" icon="refresh" label="Обновить" :loading="loading" @click="fetchDashboard" />
        </div>

        <q-banner v-if="error" class="bg-red-1 text-red-9 q-mb-md" rounded>
            {{ error }}
        </q-banner>

        <div class="row q-col-gutter-md q-mb-lg">
            <div class="col-12 col-md-4">
                <q-card bordered flat>
                    <q-card-section>
                        <div class="text-caption text-grey-7">Всего пользователей</div>
                        <div class="text-h4 text-primary">{{ stats.users_total }}</div>
                    </q-card-section>
                </q-card>
            </div>
            <div class="col-12 col-md-4">
                <q-card bordered flat>
                    <q-card-section>
                        <div class="text-caption text-grey-7">Активные пользователи</div>
                        <div class="text-h4 text-positive">{{ stats.active_users }}</div>
                    </q-card-section>
                </q-card>
            </div>
            <div class="col-12 col-md-4">
                <q-card bordered flat>
                    <q-card-section>
                        <div class="text-caption text-grey-7">Счётчик входов</div>
                        <div class="text-h4 text-accent">{{ stats.total_visits }}</div>
                    </q-card-section>
                </q-card>
            </div>
        </div>

        <q-card bordered flat>
            <q-card-section class="row items-center">
                <div class="text-h6">Пользователи в базе данных</div>
                <q-space />
                <q-badge color="primary" class="q-pa-sm">{{ users.length }}</q-badge>
            </q-card-section>
            <q-separator />
            <q-table
                :rows="users"
                :columns="columns"
                row-key="id"
                :loading="loading"
                flat
            >
                <template #body-cell-role="props">
                    <q-td :props="props">
                        <q-badge :color="roleColor(props.row.role)">
                            {{ roleName(props.row.role) }}
                        </q-badge>
                    </q-td>
                </template>
                <template #body-cell-created_at="props">
                    <q-td :props="props">
                        {{ formatDate(props.row.created_at) }}
                    </q-td>
                </template>
            </q-table>
        </q-card>
    </q-page>
</template>

<script>
import { computed, defineComponent, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'

import { adminAPI } from '../utils/api'
import { formatDate, formatPrice } from '../utils/formatters'

export default defineComponent({
    name: 'PageAdmin',
    setup() {
        const store = useStore()
        const router = useRouter()
        const loading = ref(false)
        const error = ref('')
        const users = ref([])
        const stats = ref({
            users_total: 0,
            active_users: 0,
            total_visits: 0,
        })

        const currentUser = computed(() => store.state.auth.user || {})

        const columns = [
            { name: 'login', field: 'login', label: 'Логин', align: 'left', sortable: true },
            { name: 'role', field: 'role', label: 'Роль', align: 'left', sortable: true },
            { name: 'balance', field: row => formatPrice(row.balance), label: 'Баланс', align: 'right', sortable: true },
            { name: 'created_at', field: 'created_at', label: 'Создан', align: 'left', sortable: true },
        ]

        const roleName = (role) => {
            if (role === 'admin') return 'Администратор'
            if (role === 'system') return 'Системный'
            return 'Редактор'
        }

        const roleColor = (role) => {
            if (role === 'admin') return 'negative'
            if (role === 'system') return 'grey'
            return 'primary'
        }

        const fetchDashboard = async () => {
            loading.value = true
            error.value = ''
            try {
                const { data } = await adminAPI.dashboard()
                users.value = data.users || []
                stats.value = data.stats || stats.value
            } catch (err) {
                if (err.response?.status === 403) {
                    error.value = 'Эта страница доступна только администраторам.'
                    return
                }
                error.value = err.response?.data?.error || err.message || 'Не удалось загрузить админскую панель'
            } finally {
                loading.value = false
            }
        }

        onMounted(async () => {
            if (!store.state.auth.isAuthenticated) {
                router.push('/login')
                return
            }
            if (!currentUser.value.role) {
                await store.dispatch('auth/fetchProfile')
            }
            if ((store.state.auth.user || {}).role !== 'admin') {
                error.value = 'Эта страница доступна только администраторам.'
                return
            }
            fetchDashboard()
        })

        return {
            loading,
            error,
            users,
            stats,
            columns,
            roleName,
            roleColor,
            formatDate,
            fetchDashboard,
        }
    },
})
</script>
