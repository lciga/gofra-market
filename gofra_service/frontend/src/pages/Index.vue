<template>
    <q-page class="q-pa-md">
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
                    <div class="col-auto" style="padding-bottom: 22px;">
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
            <q-spinner-gears size="50px" color="primary"/>
        </q-inner-loading>

        <div v-if="!loading && filteredListings.length === 0" class="text-center q-my-xl">
            <q-icon name="mdi-package-variant" size="100px" color="grey-4"/>
            <div class="text-h6 q-mt-md text-grey">Нет доступных листингов</div>
        </div>

        <div class="row q-col-gutter-md">
            <div v-for="listing in filteredListings" :key="listing.id" class="col-12 col-sm-6 col-md-4 col-lg-3">
                <app-gofer-card
                    :listing="listing"
                    @click="$router.push(`/listing/${listing.id}`)"
                />
            </div>
        </div>

        <create-listing-modal
            v-model="showCreateModal"
            @created="handleListingCreated"
        />
    </q-page>
</template>

<script>
import { computed, ref, onMounted, defineComponent } from 'vue'
import { useStore } from 'vuex'

import AppGoferCard from '../components/common/AppGoferCard.vue'
import GoferFilters from '../components/marketplace/GoferFilters.vue'
import CreateListingModal from '../components/marketplace/CreateListingModal.vue'

export default defineComponent({
    name: 'PageIndex',
    components: {
        AppGoferCard,
        GoferFilters,
        CreateListingModal,
    },
    setup() {
        const store = useStore()
        const showCreateModal = ref(false)
        const searchQuery = ref('')

        const loading = computed(() => store.state.listing.loading)
        const filteredListings = computed(() => store.getters['listing/filteredListings'])

        onMounted(() => {
            store.dispatch('listing/fetchListings')
        })

        const handleListingCreated = () => {
            store.dispatch('listing/fetchListings')
        }

        const handleSearch = () => {
            store.commit('listing/SET_SEARCH_QUERY', searchQuery.value)
            store.dispatch('listing/fetchListings')
        }

        return {
            loading,
            filteredListings,
            showCreateModal,
            searchQuery,
            handleListingCreated,
            handleSearch,
        }
    },
})
</script>