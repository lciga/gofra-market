<template>
    <q-page class="flex flex-center text-center q-pa-md">
        <div>
            <div class="text-h1 text-primary">404</div>
            <div class="text-h4 q-mt-md">Страница не найдена</div>
            <div class="text-subtitle1 q-mt-md q-mb-xl">
                Переадресуем на главную через {{ redirectCountdown }} секунд(ы).
            </div>
            <q-btn label="На главную" to="/" color="primary"/>
        </div>
    </q-page>
</template>

<script>
import { defineComponent, onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'

export default defineComponent({
    name: 'Error404',
    setup() {
        const router = useRouter()
        const redirectCountdown = ref(5)
        let timerId = null

        onMounted(() => {
            timerId = window.setInterval(() => {
                redirectCountdown.value -= 1
                if (redirectCountdown.value <= 0) {
                    window.clearInterval(timerId)
                    router.replace('/')
                }
            }, 1000)
        })

        onUnmounted(() => {
            if (timerId) {
                window.clearInterval(timerId)
            }
        })

        return {
            redirectCountdown,
        }
    }
})
</script>