<template>
    <q-bage :color="rarityColor" class="text-white">{{ rarityName }}</q-bage>
</template>

<script>
import {computed} from 'vue'
import { RARITY_COLORS, RARITY_NAMES } from '../../utils/constants';

export default {
    name: 'AppRarityBage',
    props: {
        rarity: {
            type: Number,
            requred: true,
            validator: (value) => value >= 1 && value <= 3
        }
    },
    setup(props) {
        const rarityColor = computed(() => RARITY_COLORS[props.rarity])
        const rarityName = computed(() => RARITY_NAMES[props.rarity])

        return {
            rarityColor,
            rarityName,
        }
    }
}
</script>