<template>
    <q-page class="q-pa-md">
        <div v-if="loading" class="flex flex-center">
            <q-spinner-gears size="50px" color="primary"/>
        </div>

        <div v-else-if="listing" class="row q-col-gutter-xl">
            <div class="col-12 col-md-6">
                <q-card class="q-pa-md">
                    <q-img :src="getImageURL(listing)" :alt="listing.gofer.name" ratio="1" class="rounded-borders" placeholder-src="../../image/gofer-placeholder.png">
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

                    <div class="q-mb-lg">
                        <div class="text-h6 q-mb-sm">Описание</div>
                        <p class="text-body1">{{ listing.description || 'Нет описания' }}</p>
                    </div>

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

export default {
    name: 'PageListingDetail',
    setup() {
        const store = useStore()
        const route = useRoute()
        const router = useRouter()
        const $q = useQuasar()

        const loading = ref(true)
        const listing = computed(() => store.state.listing.currentListing)
        const isAuthenticated = computed(() => store.state.auth.isAuthenticated)

        const getImageURL = (listing) => {
            return listing.image?.source_url || `/api/listings/${listing.id}/image` ||'../../image/gofer-placeholder.png'
        }

        const handleBuy = async () => {
            if (!isAuthenticated.value) {
                router.push('/login')
                return
            }

            try {
                await store.dispatch('listings/buyListing', listing.value.id)
                $q.notify({
                    type: 'positive',
                    message: 'Успешная покупка!'
                })
                router.push('/')
            } catch (error) {
                $q.notify({
                    type: 'negative',
                    message: 'Ошибка при покупке'
                })
            }
        }

        onMounted(async () => {
            try {
                await store.dispatch('listings/fetchListing', route.params.id)
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
            handleBuy
        }
    }
}
</script>