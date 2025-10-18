<template>
    <q-page class="q-pa-md">
        <div v-if="loading" class="flex flex-center">
            <q-spinner-gears size="50px" color="primary"/>
        </div>

        <div v-else-if="listing" class="row q-col-gutter-xl">
            <div class="col-12 col-md-6">
                <q-card class="q-pa-md">
                    <q-img :src="getImageURL(listing)" :alt="listing.gofer.name" ratio="1" class="rounded-borders" :placeholder-src="goferPlaceholder">
                        <template v-slot="error">
                            <div class="absolute-full flex flex-center bg-grey-3 text-grey">
                                <q-icon name="mdi-image-off" size="xl"/>
                            </div>
                        </template>
                    </q-img>
                </q-card>
            </div>
            <div class="col-12 col-md-6">
                <q-card class="q-pa-md">
                    <div class="row items-center q-mb-md">
                        <div class="col">
                            <div class="text-h4 text-weight-bold">{{ listing.gofer.name }}</div>
                            <app-rarity-badge :rarity="listing.gofer.rarity"></app-rarity-badge>
                        </div>
                        <div class="col-auto">
                            <div class="text-h5 text-primary text-weight-bold">
                                {{ formatPrice(listing.price) }} горутин
                            </div>
                        </div>
                    </div>

                    <q-separator class="q-mb-md"/>

                    <!-- Description скрыт, но подтягивается с бэка для NoSQL injection -->

                    <q-btn
                        color="primary"
                        size="lg"
                        label="Купить сейчас"
                        :disabled="listing.is_sold || !isAuthenticated"
                        @click="handleBuy"
                        class="full-width"
                    />
                    <div v-if="listing.is_sold" class="text-negative text-center q-mt-sm">Этот гофер уже продан</div>
                </q-card>
            </div>
        </div>
        <div v-else class="text-center q-my-xl">
            <q-icon name="mdi-alert-circle" size="100px" color="negative"/>
            <div class="text-h6 q-mt-md">Листинг не найден</div>
            <q-btn label="На главную" to="/" color="primary" class="q-mt-md"/>
        </div>
    </q-page>
</template>

<script>
import {ref, onMounted, computed} from 'vue'
import { useStore } from 'vuex'
import { useRoute, useRouter } from 'vue-router'
import { useQuasar } from 'quasar'
import { formatPrice } from '../utils/formatters'
import goferPlaceholder from '../assets/gofer-placeholder.png'
import AppRarityBadge from '../components/common/AppRarityBadge.vue'
import { API_URL } from '../utils/api'

export default {
    name: 'PageListingDetail',
    components: {
        AppRarityBadge,
    },
    setup() {
        const store = useStore()
        const route = useRoute()
        const router = useRouter()
        const $q = useQuasar()

        const loading = ref(true)
        const listing = computed(() => store.state.listing.currentListing)
        const isAuthenticated = computed(() => store.state.auth.isAuthenticated)

        const getImageURL = (listing) => {
            // Используем API эндпоинт для получения изображения
            // Если есть изображение (source_url или uploaded file), бэкенд вернет его
            // Иначе вернется 404 и покажется placeholder через slot:error
            if (listing?.image && (listing.image.source_url || listing.image.content_type)) {
                // Добавляем случайный параметр для обхода кеша браузера
                const cacheBuster = Math.random().toString(36).substring(7)
                return `${API_URL}/listings/${listing.id}/image?v=${cacheBuster}`
            }
            return goferPlaceholder
        }

        const handleBuy = async () => {
            if (!isAuthenticated.value) {
                router.push('/login')
                return
            }

            try {
                await store.dispatch('listing/buyListing', listing.value.id)
                $q.notify({
                    type: 'positive',
                    message: 'Успешная покупка!'
                })
                // Navigate to home page (listings will be refreshed automatically)
                await router.push('/')
            } catch (error) {
                const errorMessage = error.response?.data?.error || error.message || 'Ошибка при покупке'
                $q.notify({
                    type: 'negative',
                    message: errorMessage
                })
            }
        }

        onMounted(async () => {
            try {
                await store.dispatch('listing/fetchListing', route.params.id)
            } catch (error) {
                console.error('Error loading listing', error)
            } finally {
                loading.value = false
            }
        })

        return {
            loading,
            listing,
            isAuthenticated,
            getImageURL,
            formatPrice,
            handleBuy,
            goferPlaceholder
        }
    }
}
</script>