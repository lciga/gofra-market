<template>
    <q-page class="q-pa-md">
        <q-banner class="q-mb-lg bg-primary text-white" dense>
            <div class="row items-center no-wrap">
                <div class="col">
                    <div class="text-subtitle1">{{ activeUsersText }}</div>
                    <div class="text-caption q-mt-xs">Обновляется каждые 5 секунд</div>
                </div>
                <div class="col-auto">
                    <q-spinner-gears v-if="activeUsersLoading" size="24px" color="white" />
                </div>
            </div>
        </q-banner>

        <div class="row q-col-gutter-xl">
            <div class="col-12 col-lg-8">
                <div class="row items-center q-mb-md">
                    <div class="col">
                        <h4 class="q-ma-none text-primary">Gofra Market - магазин гоферов</h4>
                        <p class="text-grey">Покупайте и продавайте уникальных гоферов - зарабатывайте горутины</p>
                    </div>
                    <div class="col-auto">
                        <q-btn
                            color="primary"
                            icon="add"
                            label="Создать листинг"
                            @click="showCreateModal = true"
                            v-if="$store.state.auth.isAuthenticated"
                        />
                    </div>
                </div>

                <q-card class="q-mb-md" flat bordered>
                    <q-card-section>
                        <div class="row q-col-gutter-md items-center">
                            <div class="col">
                                <q-input
                                    v-model="searchQuery"
                                    label="Поиск по имени гофера"
                                    outlined
                                    clearable
                                    @keyup.enter="handleSearch"
                                    hint="Введите имя гофера для поиска"
                                >
                                    <template v-slot:prepend>
                                        <q-icon name="search" />
                                    </template>
                                </q-input>
                            </div>
                            <div class="col-auto search-button-wrapper">
                                <q-btn
                                    color="primary"
                                    label="Найти"
                                    @click="handleSearch"
                                    unelevated
                                />
                            </div>
                        </div>
                    </q-card-section>
                </q-card>

                <gofer-filters />

                <q-inner-loading :showing="loading">
                    <q-spinner-gears size="50px" color="primary" />
                </q-inner-loading>

                <div v-if="!loading && filteredListings.length === 0" class="text-center q-my-xl">
                    <q-icon name="mdi-package-variant" size="100px" color="grey-4" />
                    <div class="text-h6 q-mt-md text-grey">Нет доступных листингов</div>
                </div>

                <div class="row q-col-gutter-md">
                    <div
                        v-for="listing in filteredListings"
                        :key="listing.id"
                        class="col-12 col-sm-6 col-md-4 col-lg-6"
                    >
                        <app-gofer-card
                            :listing="listing"
                            @click="$router.push(`/listing/${listing.id}`)"
                        />
                    </div>
                </div>
            </div>

            <div class="col-12 col-lg-4">
                <div class="sticky-forms">
                    <q-card bordered flat>
                        <q-card-section class="bg-light-blue-1 text-primary text-weight-medium">
                            Google форма обратной связи
                        </q-card-section>
                        <q-separator />
                        <q-card-section>
                            <q-form @submit.prevent="handleGoogleFormSubmit" class="column q-gutter-md">
                                <q-input v-model="googleForm.name" label="Имя" outlined required />
                                <q-input v-model="googleForm.email" label="Email" type="email" outlined required />
                                <q-input
                                    v-model="googleForm.message"
                                    type="textarea"
                                    label="Сообщение"
                                    outlined
                                    autogrow
                                />
                                <q-btn label="Отправить" color="primary" type="submit" />
                            </q-form>
                        </q-card-section>
                    </q-card>

                    <q-card bordered flat>
                        <q-card-section class="bg-amber-2 text-warning text-weight-medium">
                            Yandex форма обратной связи
                        </q-card-section>
                        <q-separator />
                        <q-card-section>
                            <q-form @submit.prevent="handleYandexFormSubmit" class="column q-gutter-md">
                                <q-input v-model="yandexForm.company" label="Компания" outlined required />
                                <q-input v-model="yandexForm.phone" label="Телефон" outlined required />
                                <q-input
                                    v-model="yandexForm.comment"
                                    type="textarea"
                                    label="Комментарий"
                                    outlined
                                    autogrow
                                />
                                <q-btn label="Отправить" color="warning" text-color="white" type="submit" />
                            </q-form>
                        </q-card-section>
                    </q-card>
                </div>
            </div>
        </div>

        <create-listing-modal
            v-model="showCreateModal"
            @created="handleListingCreated"
        />
    </q-page>
</template>

<script>
import { computed, ref, onMounted, onUnmounted, defineComponent, reactive } from 'vue'
import { useStore } from 'vuex'
import { useQuasar } from 'quasar'

import AppGoferCard from '../components/common/AppGoferCard.vue'
import GoferFilters from '../components/marketplace/GoferFilters.vue'
import CreateListingModal from '../components/marketplace/CreateListingModal.vue'
import { statsAPI } from '../utils/api'

export default defineComponent({
    name: 'PageIndex',
    components: {
        AppGoferCard,
        GoferFilters,
        CreateListingModal,
    },
    setup() {
        const store = useStore()
        const $q = useQuasar()

        const showCreateModal = ref(false)
        const searchQuery = ref('')

        const loading = computed(() => store.state.listing.loading)
        const filteredListings = computed(() => store.getters['listing/filteredListings'])

        const activeUsers = ref(0)
        const activeUsersLoading = ref(false)
        const activeUsersError = ref('')
        const activeUsersTimer = ref(null)

        const googleForm = reactive({
            name: '',
            email: '',
            message: '',
        })

        const yandexForm = reactive({
            company: '',
            phone: '',
            comment: '',
        })

        const fetchActiveUsers = async () => {
            activeUsersLoading.value = true
            try {
                const { data } = await statsAPI.getActiveUsers()
                activeUsers.value = data?.active_users ?? 0
                activeUsersError.value = ''
            } catch (error) {
                console.error('Unable to load active users', error)
                activeUsersError.value = 'Не удалось обновить статистику'
            } finally {
                activeUsersLoading.value = false
            }
        }

        onMounted(() => {
            store.dispatch('listing/fetchListings')
            fetchActiveUsers()
            activeUsersTimer.value = window.setInterval(fetchActiveUsers, 5000)
        })

        onUnmounted(() => {
            if (activeUsersTimer.value) {
                window.clearInterval(activeUsersTimer.value)
            }
        })

        const activeUsersText = computed(() => {
            return activeUsersError.value || `Сейчас на сервере ${activeUsers.value} клиент(ов) в онлайне`
        })

        const handleListingCreated = () => {
            store.dispatch('listing/fetchListings')
        }

        const handleSearch = () => {
            store.commit('listing/SET_SEARCH_QUERY', searchQuery.value)
            store.dispatch('listing/fetchListings')
        }

        const handleGoogleFormSubmit = () => {
            if (!googleForm.name || !googleForm.email) {
                $q.notify({ type: 'negative', message: 'Заполните имя и email для отправки формы.' })
                return
            }

            $q.notify({ type: 'positive', message: 'Заявка отправлена в Google формы.' })
            googleForm.name = ''
            googleForm.email = ''
            googleForm.message = ''
        }

        const handleYandexFormSubmit = () => {
            if (!yandexForm.company || !yandexForm.phone) {
                $q.notify({ type: 'negative', message: 'Введите название компании и телефон.' })
                return
            }

            $q.notify({ type: 'positive', message: 'Заявка отправлена в Yandex формы.' })
            yandexForm.company = ''
            yandexForm.phone = ''
            yandexForm.comment = ''
        }

        return {
            loading,
            filteredListings,
            showCreateModal,
            searchQuery,
            handleListingCreated,
            handleSearch,
            activeUsers,
            activeUsersLoading,
            activeUsersText,
            googleForm,
            yandexForm,
            handleGoogleFormSubmit,
            handleYandexFormSubmit,
        }
    },
})
</script>

<style scoped>
.sticky-forms {
    position: sticky;
    top: 88px;
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.search-button-wrapper {
    padding-bottom: 22px;
}

@media (max-width: 1023px) {
    .sticky-forms {
        position: static;
    }
}
</style>