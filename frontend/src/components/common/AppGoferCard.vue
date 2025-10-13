<template>
    <q-card
        class="gofer-card"
        :class="`rarity-${listing.gofer.rarity}`"
        flat
        bordered
    >
        <q-img
            :src="getImageURL(listing)"
            :alt="listing.gofer.name"
            ratio="1"
            class="gofer-image"
            placeholder-src="../../image/gofer-placeholder.png"
        >
            <template v-slot:error>
                <div class="absolute-full flex flex-center bg-grey-3 text-grey-8">
                    <q-icon name="mdi-image-off" size="xl"></q-icon>
                </div>
            </template>

            <div v-if="listing.is_sold" class="absolute-full flex flex-center bg-dak overlay-sold">
                <q-icon name="mdi-cancel" size="xl" color="white"></q-icon>
                <div class="text-h6 text-white q-mt-md">Продано</div>
            </div>
        </q-img>

        <q-card-section>
            <div class="text-h6 text-weight-bold">{{ listing.gofer.name }}</div>
            <div class="text-caption text-grey-7 q-mt-xs">{{ truncateDescription(listing.description) }}</div>
        </q-card-section>

        <q-card-section class="qt-pt-none">
            <div class="row items-center justify-between">
                <div class="price-section">
                    <div class="text-caption text-grey">Цена</div>
                    <div class="text-h6 text-weight-bold text-primary">{{ formatPrice(listing.price) }} горутиг</div>
                </div>
                <q-btn
                color="primary"
                label="Купить"
                :disabled="listing.is_sold || !isAuthenticated"
                @click="handleBuy"
                >
                </q-btn>
            </div>
        </q-card-section>
    </q-card>
</template>

<script>
import {computed} from 'vue'
import {useStore} from 'vuex'
import {useRouter} from 'vue-router'

import { RARITY_COLORS, RARITY_NAMES } from '../../utils/constants'
import { formatPrice, truncateText } from '../../utils/formatters'

export default {
    name:'AppGoferCard',
    props: {
        listing: {
            type: Object,
            required: true
        }
    },
    setup(props) {
        const store = useStore
        const router = useRouter

        const isAuthenticated = computed(() => store.state.auth.isAuthenticated)
        const getRarityColor = (rarity) => RARITY_COLORS[rarity] || RARITY_COLORS[1]
        const getRarityName = (rarity) => RARITY_NAMES[rarity] || RARITY_NAMES[1]

        const truncateDescription = (description) => {
            return truncateText(description || 'Нет описания', 80)
        }

        const getImageURL = (listing) => {
            return listing.image?.source_url || `/api/listings/${listing.id}/image` || '/placeholder-gofer.png'
        }

        const handleBuy = async () => {
            if (!isAuthenticated) {
                router.push('/login')
                return
            }
            try {
                await store.dispatch('listings/buyListing', props.listing.id)
                store.dispatch('auth/fetchProfile')
            } catch (error) {
                console.error('Purchase failed:', error)
            }
        }
        return {
            isAuthenticated,
            getRarityColor,
            getRarityName,
            truncateDescription,
            getImageURL,
            formatPrice,
            handleBuy,
        }
    }
}
</script>

<style scoped>
.gofer-card {
    max-width: 300px;
    transition: transform 0.3s ease, box-shadow 0.3s ease;
    border-left: 4px solid transparent;
}

.gofer-card:hover {
    transform: translateY(-5px);
    box-shadow: 0 10px 20px rgba(0,0,0,0.1);
}

.rarity-1 {
    border-left-color: #E6D5B8;
}

.rarity-2 {
    border-left-color: #D4B896;
}

.rarity-3 {
    border-left-color: #8B4513;
}

.overlay-sold {
    background: rgba(0, 0, 0, 0.7);
    backdrop-filter: blur(2px);
}
</style>