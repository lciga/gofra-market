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
            :placeholder-src="goferPlaceholder"
            img-class="crossorigin-anonymous"
        >
            <template v-slot:error>
                <div class="absolute-full flex flex-center bg-grey-3 text-grey-8">
                    <q-icon name="mdi-image-off" size="xl"></q-icon>
                </div>
            </template>

            <div v-if="listing.is_sold && !hideActions" class="absolute-full flex flex-center bg-dak overlay-sold">
                <q-icon name="mdi-cancel" size="xl" color="white"></q-icon>
                <div class="text-h6 text-white q-mt-md">Продано</div>
            </div>
        </q-img>

        <q-card-section>
            <div class="text-h6 text-weight-bold">{{ listing.gofer.name }}</div>
            <div v-if="showDescription && listing.description" class="text-caption text-grey q-mt-sm">
                {{ truncateDescription(listing.description) }}
            </div>
        </q-card-section>

        <q-card-section class="qt-pt-none">
            <div class="row items-center justify-between">
                <div class="price-section">
                    <div class="text-caption text-grey">Цена</div>
                    <div class="text-h6 text-weight-bold text-primary">{{ formatPrice(listing.price) }} горутин</div>
                </div>
                <q-btn
                v-if="!hideActions"
                color="primary"
                label="Купить"
                :disabled="listing.is_sold || !isAuthenticated"
                @click.stop="handleBuy"
                >
                </q-btn>
            </div>
        </q-card-section>
    </q-card>
</template>

<script>
import { computed } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'

import goferPlaceholder from 'assets/gofer-placeholder.png'
import { RARITY_COLORS, RARITY_NAMES } from '../../utils/constants'
import { formatPrice, truncateText } from '../../utils/formatters'
import { API_URL } from '../../utils/api'

export default {
	name: 'AppGoferCard',
	props: {
		listing: {
			type: Object,
			required: true,
		},
		hideActions: {
			type: Boolean,
			default: false,
		},
		showDescription: {
			type: Boolean,
			default: false,
		},
	},
	setup(props) {
		const store = useStore()
		const router = useRouter()

		const isAuthenticated = computed(() => store.state.auth.isAuthenticated)
		const getRarityColor = (rarity) => RARITY_COLORS[rarity] || RARITY_COLORS[1]
		const getRarityName = (rarity) => RARITY_NAMES[rarity] || RARITY_NAMES[1]

		const truncateDescription = (description) => truncateText(description, 80)

		const getImageURL = (listing) => {
			if (listing.image && (listing.image.source_url || listing.image.content_type)) {
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
				await store.dispatch('listing/buyListing', props.listing.id)
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
            goferPlaceholder,
		}
	},
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