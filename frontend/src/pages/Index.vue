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
                    v-if="$store.state.auth.isAuthinticatedw"
                />
            </div>
        </div>

        <gofer-filters />

        <q-inner-loading :showing="loading">
            <q-spinner-gears size="50px" color="primary"/>
        </q-inner-loading>

        <div v-if="!loading && filteredListings.lenght===0" class="text-center q-my-xl">
            <q-icon name="mdi-package-variant" size="100px" color="grey-4"/>
            <div class="text-h6 q-mt-md text-grey">Нет доступных листингов</div>
        </div>

        <div class="row q-col-gutter-md">
            <div v-for="listing in filteredListings" :key="listing.id" class="col-12 col-sm-6 col-md-4 col-lg-3">
                <app-gofer-card
                    :listing="listing"
                    @click="$router.push(`listing/${listing.id}`)"
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
import {computed, ref, onMounted} from 'vue'
import {useStore} from 'vuex'

export default {
    name: 'PageIndex',
    setup() {
        const store = useStore
        const showCreateModal = ref(false)

        const loading = computed(() => store.state.listing.loading)
        const filteredListings = computed(() => store.getters['listings/filteredListings'])

        onMounted(() => {
            store.dispatch('listings/fetchListings')
        })

        const handleListingCreated = () => {
            store.dispatch('listings/fetchListings')            
        }

        return {
            loading,
            filteredListings,
            showCreateModal,
            handleListingCreated
        }
    }
}
</script>