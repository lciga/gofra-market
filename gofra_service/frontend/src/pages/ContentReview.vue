<template>
    <q-page class="q-pa-lg bg-grey-1">
        <div class="row justify-center">
            <div class="col-12 col-md-8 col-lg-6">
                <q-card bordered flat>
                    <q-card-section>
                        <div class="text-h4 text-primary q-mb-sm">Материал редактору</div>
                        <div class="text-grey-7">
                            Заполните форму, а кнопка подготовит письмо редактору сайта. Почтовый сервер тут не притворяется героем, поэтому открывается обычный почтовый клиент.
                        </div>
                    </q-card-section>

                    <q-separator />

                    <q-card-section>
                        <q-form @submit.prevent="handleSubmit" class="column q-gutter-md">
                            <q-input
                                v-model="form.title"
                                label="Заголовок материала"
                                outlined
                                :rules="[val => !!val || 'Введите заголовок']"
                            />
                            <q-input
                                v-model="form.text"
                                type="textarea"
                                label="Текст материала"
                                outlined
                                autogrow
                                :rules="[val => !!val || 'Введите текст материала']"
                            />
                            <q-btn
                                label="Отправить редактору"
                                color="primary"
                                type="submit"
                                :loading="loading"
                                class="full-width"
                            />
                        </q-form>
                    </q-card-section>

                    <q-card-section v-if="review.editor_email" class="bg-grey-2">
                        <div class="text-caption text-grey-7">Адрес редактора</div>
                        <div class="text-body1">{{ review.editor_email }}</div>
                    </q-card-section>
                </q-card>
            </div>
        </div>
    </q-page>
</template>

<script>
import { defineComponent, reactive, ref } from 'vue'
import { useQuasar } from 'quasar'

import { contentAPI } from '../utils/api'

export default defineComponent({
    name: 'PageContentReview',
    setup() {
        const $q = useQuasar()
        const loading = ref(false)
        const form = reactive({
            title: '',
            text: '',
        })
        const review = reactive({
            editor_email: '',
            subject: '',
            body: '',
            mailto: '',
        })

        const handleSubmit = async () => {
            loading.value = true
            try {
                const { data } = await contentAPI.submit({
                    title: form.title,
                    text: form.text,
                })
                Object.assign(review, data)
                window.location.href = data.mailto
                $q.notify({ type: 'positive', message: 'Письмо подготовлено для отправки редактору.' })
            } catch (error) {
                const message = error.response?.data?.error || error.message || 'Не удалось подготовить письмо'
                $q.notify({ type: 'negative', message })
            } finally {
                loading.value = false
            }
        }

        return {
            form,
            review,
            loading,
            handleSubmit,
        }
    },
})
</script>
