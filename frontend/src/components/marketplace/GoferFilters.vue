<template>
    <q-card class="q-mb-md" flat bordered>
        <q-card-section>
            <div class="text-h6">Фильтры</div>
            <div class="row q-col-gutter-md q-mt-sm">
                <div class="col-12 col-sm-4">
                    <q-select
                        v-model="filters.rarity"
                        :options="rarityOptions"
                        lable="Редкость"
                        clearable
                        emit-value
                        map-options
                    />
                </div>
                <div class="col-12 col-sm-4">
                    <q-input
                        v-model.number="filters.minPrice"
                        type="number"
                        label="Макс. цена"
                        :min="filters.minPrice"
                    />
                </div>
            </div>
            <div class="row justify-end q-mt-md">
                <q-btn color="primary" label="Применить" @click="applyFilters"></q-btn>
                <q-btn flat label="Сбросить" @click="resetFilters" class="q-ml-sm"></q-btn>
            </div>
        </q-card-section>
    </q-card>
</template>

<script>
import {ref, watch} from 'vue'
import {useStore} from 'vuex'
import { RARITY_NAMES } from '../../utils/constants'

export default {
    name: 'GoferFilters',
    setup() {
        const store = useStore

        const rarityOptions = [
            { label: RARITY_NAMES[1], value: 1 },
            { label: RARITY_NAMES[2], value: 2 },
            { label: RARITY_NAMES[3], value: 3 }
        ]

        const filters = ref({
            rarity: null,
            minPrice: 0,
            maxPrice: 1000
        })

        const applyFilters = () => {
            store.commit('listings/SET_FILTERS', filters.value)
            store.dispatch('listings/fetchListings')
        }

        const resetFilters = () => {
            filters.value = {
                rarity: null,
                minPrice: 0,
                maxPrice: 1000
            }
            applyFilters()
        }

        watch(filters, applyFilters, { deep: true, immediate: false })

        return {
            filters,
            rarityOptions,
            applyFilters,
            resetFilters,
        }
    }
}
</script>