<template>
    <q-page class="q-pa-md">
        <div class="row items-center q-mb-xl">
            <div class="col">
                <div class="text-h4">Профиль</div>
                <div class="text-grey">Управление вашим аккаунтом</div>
            </div>
        </div>

        <div class="row q-col-gutter-xl">
            <div class="col-12 col-md-4">
                <q-card class="q-pa-md">
                    <div class="text-center">
                        <q-avatar size="100px" color="primary" text-color="white" class="q-mb-md">
                            {{ user.login?.charAt(0).toUpperCase() }}
                        </q-avatar>
                        <div class="text-h6">{{ user.login }}</div>
                        <div class="text-caption text-grey">ID: {{ user.user_id }}</div>
                    </div>

                    <q-separator class="q-my-md"/>

                    <div class="text-center">
                        <div class="text-h4 text-primary">{{ formatPrice(user.balance) }}</div>
                        <div class="text-caption">горутин на балансе</div>
                    </div>
                </q-card>
            </div>

            <div class="col-12 col-md-8">
                <q-card>
                    <q-tabs v-model="tab" dense class="text-grey" active-color="primary" indicator-color="primary" align="justify" narrow-indicator>
                        <q-tab name="listings" label="Мои листинги"/>
                        <q-tab name="purchases" label="Мои покупки"/>
                    </q-tabs>

                    <q-separator/>

                    <q-tab-panels v-model="tab" animated>
                        <q-tab-panel name="listings">
                            <div v-if="userListings.length === 0" class="text-center q-py-xl">
                                <q-icon name="mdi-package-variant" size="60" color="grey-4"/>
                                <div class="text-h6 q-mt-md text-grey">У вас нет листингов</div>
                                <q-btn 
                                    label="Создать листинг"
                                    color="primary"
                                    to="/"
                                    class="q-mt-md"
                                />
                            </div>

                            <div class="row q-col-gutter-md">
                                <div v-for="listing in userListings" :key="listing.id" class="col-12 col-sm-6">
                                    <app-gofer-card :listing="listing" :show-description="true"/>
                                </div>
                            </div>
                        </q-tab-panel>

                        <q-tab-panel name="purchases">
                            <div v-if="userPurchases.length === 0" class="text-center q-py-xl">
                                <q-icon name="mdi-cart-outline" size="60" color="grey-4"/>
                                <div class="text-h6 q-mt-md text-grey">У вас пока нет покупок</div>
                                <q-btn 
                                    label="Перейти к покупкам"
                                    color="primary"
                                    to="/"
                                    class="q-mt-md"
                                />
                            </div>

                            <div class="row q-col-gutter-md">
                                <div v-for="listing in userPurchases" :key="listing.id" class="col-12 col-sm-6">
                                    <app-gofer-card :listing="listing" :hide-actions="true" :show-description="true"/>
                                </div>
                            </div>
                        </q-tab-panel>
                    </q-tab-panels>
                </q-card>
            </div>
        </div>
    </q-page>
</template>

<script>
import { defineComponent, ref, computed, onMounted, watch } from 'vue'
import { useStore } from 'vuex'
import { useRoute } from 'vue-router'

import AppGoferCard from '../components/common/AppGoferCard.vue'
import { formatPrice } from '../utils/formatters'
import { listingAPI } from '../utils/api'

export default defineComponent({
    name: 'PageProfile',
    components: {
        AppGoferCard,
    },
    setup() {
        const store = useStore()
        const route = useRoute()
        const tab = ref('listings')
        const myListings = ref([])

        const user = computed(() => store.state.auth.user || {})

        const userListings = computed(() => {
            if (!user.value.user_id) return []
            return myListings.value.filter((listing) => listing.seller_id === user.value.user_id && !listing.is_sold)
        })

        const userPurchases = computed(() => {
            if (!user.value.user_id) return []
            return myListings.value.filter((listing) => listing.buyer_id === user.value.user_id && listing.is_sold)
        })

        const fetchMyListings = async () => {
            try {
                const response = await listingAPI.getMyListings()
                myListings.value = response.data.listings || []
            } catch (error) {
                console.error('Failed to fetch my listings:', error)
            }
        }

        onMounted(() => {
            fetchMyListings()
        })

        watch(
            () => route.path,
            (newPath) => {
                if (newPath === '/profile') {
                    fetchMyListings()
                }
            }
        )

        return {
            tab,
            user,
            userListings,
            userPurchases,
            formatPrice,
            fetchMyListings,
        }
    },
})
</script>