<template>
    <q-page class="flex flex-center by-grey-1">
        <q-card class="q-pa-lg" style="width: 400px;">
            <q-card-section>
                <div class="text-h4 text-center text-primary q-mb-md">Вход</div>
                <q-form @submit="handleLogin" class="q-gutter-md">
                    <q-input
                        v-model="form.login"
                        label="Логин"
                        :rules="[val => !!val || 'Логин обязателен!']"
                    />

                    <q-input
                        v-model="form.password"
                        label="Пароль"
                        type="password"
                        :rules="[val => !!val || 'Пароль обязателен']"
                    />

                    <div>
                        <q-btn
                            label="Войти"
                            type="submit"
                            color="primary"
                            :loading="loading"
                            class="full-width"
                        />
                        <q-btn
                            label="Регистрация"
                            flat
                            to="/register"
                            class="full-width q-mt-sm"
                        />
                    </div>
                </q-form>
            </q-card-section>
        </q-card>
    </q-page>
</template>

<script>
import {ref} from 'vue'
import {useStore} from 'vuex'
import {useRouter} from 'vue-router'
import {useQuasar} from 'quasar'

export default {
    name: 'PageLogin',
    setup() {
        const store = useStore()
        const router = useRouter()
        const $q = useQuasar()

        const loading = ref(false)
        const form = ref({
            login: '',
            password: ''
        })

        const handleLogin = async () => {
            loading.value = true
            try {
                await store.dispatch('auth/login', form.value)
                $q.notify({
                    type: 'positive',
                    message: 'Успешеый вход!'
                })
                router.push('/')
            } catch (error) {
                $q.notify({
                    type: 'negative',
                    message: 'Ошибка входа!'
                }) 
            } finally {
                loading.value = false
            }
        }

        return {
            form,
            loading,
            handleLogin
        }
    }
}
</script>