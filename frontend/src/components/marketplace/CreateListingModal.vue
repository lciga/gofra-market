<template>
    <q-dialog v-model="showModal" persistent>
        <q-card style="min-width: 500px;">
            <q-card-section>
                <div class="text-h6">Создать листинг</div>
            </q-card-section>
            <q-card-section class="a-pt-none">
                <q-form @submit="handleSubmit" class="q-gutter-md">
                    <q-input
                        v-model="form.gofer_id"
                        label="ID Гофера"
                        :rules="[val => !!val || 'ID Гофера обязателен']"
                    />

                    <q-input
                        v-model.number="form.price"
                        type="nimber"
                        label="Цена"
                        :rules="[
                            val => !!val || 'Цена обязательна',
                            val => val > 0 || 'Цена должна быть больше 0'
                        ]"
                    />

                    <q-input
                        v-model="form.description"
                        label="Описание"
                        type="textarea"
                        :rules="[val => !!val || 'Описание обязательно']"
                    />

                    <q-input
                        v-model="form.image_url"
                        label="URL изображения"
                        hint="Необязательное поле"
                    />

                    <div class="row justify-end q-gutter-sm">
                        <q-btn label="Отмена" color="negative" flat @click="closeModal"></q-btn>
                        <q-btn label="Создать" type="submit" color="primary" :loading="loading"></q-btn>
                    </div>
                </q-form>
            </q-card-section>
        </q-card>
    </q-dialog>
</template> 

<script>
import {ref, watch} from 'vue'
import {useStore} from 'vuex'
import {useQuasar} from 'quasar'

export default {
    name: 'CreateListingModal',
    props: {
        modelValue: {
            type: Boolean,
            default: false
        }
    },
    emits: ['update:modelValue', 'created'],
    setup(props, {emit}) {
        const store = useStore
        const $q = useQuasar

        const showModal = ref(false)
        const loading = ref(false)

        const form = ref({
            gofer_id: '',
            price: 0,
            description: '',
            image_url: ''
        })
        
        watch(() => props.modelValue, (val) => {
            showModal.value = val
        })

        const handleSubmit = async () => {
            loading.value = true
            try {
                await store.dispatch('listing/createListing', form.value)

                if (form.value.image_url) {
                    const listingID = store.state.listings[0]?.id
                    if (listingID) {
                        await store.dispatch('listing/uploadImage', {
                            listingID,
                            imageURL: form.value.image_url
                        })
                    }
                }

                $q.notify({
                    type: 'positive',
                    message: 'Листинг успешно создан'
                })

                closeModal()
                emit('created')
            } catch (error) {
                $q.notify({
                    type: 'negative',
                    message: 'Ошибка при создании листинга'
                })
            } finally {
                loading.value = false
            }
        }
        return {
            showModal, 
            loading,
            form,
            closeModal,
            handleSubmit
        }
    }
}
</script>