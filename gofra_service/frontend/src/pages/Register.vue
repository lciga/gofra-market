<template>
    <q-page class="flex flex-center bg-grey-1">
        <q-card class="q-pa-lg" style="width: 400px;">
            <q-card-section>
                <div class="text-h4 text-center text-primary q-mb-md">Регистрация</div>
                <q-form @submit="handleRegister" class="q-gutter-md">
                    <q-input
                        v-model="form.login"
                        label="Логин"
                        :rules="[
                            val => !!val || 'Логин обязателен!',
                            val => /^[a-zA-Z0-9]+$/.test(val) || 'Только буквы и цифры'
                        ]"
                    />

                    <q-input
                        v-model="form.password"
                        label="Пароль"
                        type="password"
                        :rules="[
                            val => !!val || 'Пароль обязателен',
                            val => val.length >= 6 || 'Минимум 6 символов'
                        ]"
                    />

                    <div>
                        <q-btn
                            label="Войти"
                            flat
                            to="/login"
                            class="full-width q-mt-sm"
                        />
                        <q-btn
                            label="Регистрация"
                            type="submit"
                            color="primary"
                            :loading="loading"
                            class="full-width"
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
    name: 'PageRegister',
    setup() {
        const store = useStore()
        const router = useRouter()
        const $q = useQuasar()

        const loading = ref(false)
        const form = ref({
            login: '',
            password: ''
        })

        const handleRegister = async () => {
            loading.value = true
            try {
                await store.dispatch('auth/register', form.value)
                $q.notify({
                    type: 'positive',
                    message: 'Регистрация успешна! Вам зачислено 100 горутин'
                }) 
                router.push('/content-review')
            } catch (error) {
                const errorMessage = error.response?.data?.error || error.message || 'Ошибка регистрации'
                $q.notify({
                    type: 'negative',
                    message: errorMessage
                })
            } finally {
                loading.value = false
            }
        }

        return {
            form,
            loading,
            handleRegister
        }
    }
}  
</script>