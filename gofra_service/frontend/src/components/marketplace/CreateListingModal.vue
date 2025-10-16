<template>
    <q-dialog v-model="showModal" persistent>
        <q-card style="min-width: 500px;">
            <q-card-section>
                <div class="text-h6">Создать листинг</div>
            </q-card-section>
            <q-card-section class="q-pt-none">
                <q-form @submit="handleSubmit" class="q-gutter-md">
                    <q-input
                        v-model="form.gofer_name"
                        label="Имя Гофера"
                        :rules="[val => !!val || 'Имя Гофера обязательно']"
                        hint="Создайте нового гофера для продажи"
                    />

                    <q-select
                        v-model="form.gofer_rarity"
                        :options="rarityOptions"
                        label="Редкость Гофера"
                        emit-value
                        map-options
                        :rules="[val => val !== null || 'Редкость обязательна']"
                    />

                    <q-input
                        v-model.number="form.price"
                        type="number"
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

                    <div class="q-mt-md">
                        <div class="text-subtitle2 q-mb-sm">Изображение (необязательно)</div>
                        <q-tabs v-model="imageUploadMethod" dense class="text-grey" active-color="primary">
                            <q-tab name="url" label="По ссылке" />
                            <q-tab name="file" label="Загрузить файл" />
                        </q-tabs>

                        <q-tab-panels v-model="imageUploadMethod" animated class="q-mt-sm">
                            <q-tab-panel name="url" class="q-pa-none">
                                <q-input
                                    v-model="form.image_url"
                                    label="URL изображения"
                                    hint="Укажите ссылку на изображение"
                                    clearable
                                />
                            </q-tab-panel>

                            <q-tab-panel name="file" class="q-pa-none">
                                <q-file
                                    v-model="form.image_file"
                                    label="Выберите файл"
                                    accept="image/*"
                                    max-file-size="5242880"
                                    @rejected="onFileRejected"
                                    clearable
                                >
                                    <template v-slot:prepend>
                                        <q-icon name="attach_file" />
                                    </template>
                                </q-file>
                            </q-tab-panel>
                        </q-tab-panels>
                    </div>

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
import { ref, watch, computed } from 'vue'
import { useStore } from 'vuex'
import { useQuasar } from 'quasar'

export default {
    name: 'CreateListingModal',
    props: {
        modelValue: {
            type: Boolean,
            default: false,
        },
    },
    emits: ['update:modelValue', 'created'],
    setup(props, { emit }) {
        const store = useStore()
        const $q = useQuasar()

        const showModal = ref(false)
        const loading = ref(false)
        const imageUploadMethod = ref('url')
        const form = ref({
            gofer_name: '',
            gofer_rarity: 1,
            price: 0,
            description: '',
            image_url: '',
            image_file: null,
        })

        const rarityOptions = [
            { label: 'Обычный (1)', value: 1 },
            { label: 'Редкий (2)', value: 2 },
            { label: 'Эпический (3)', value: 3 },
            { label: 'Легендарный (4)', value: 4 },
            { label: 'Мифический (5)', value: 5 },
        ]

        const listingsState = computed(() => store.state.listing.listings)

        const resetForm = () => {
            form.value = {
                gofer_name: '',
                gofer_rarity: 1,
                price: 0,
                description: '',
                image_url: '',
                image_file: null,
            }
            imageUploadMethod.value = 'url'
        }

        const onFileRejected = () => {
            $q.notify({
                type: 'negative',
                message: 'Файл слишком большой (макс. 5MB) или неверный формат',
            })
        }

        const closeModal = () => {
            showModal.value = false
            emit('update:modelValue', false)
            resetForm()
        }

        watch(
            () => props.modelValue,
            (val) => {
                showModal.value = val
            },
            { immediate: true }
        )

        watch(showModal, (val) => {
            if (!val) {
                resetForm()
            }
            emit('update:modelValue', val)
        })

        const handleSubmit = async () => {
            loading.value = true
            try {
                const response = await store.dispatch('listing/createListing', form.value)
                const listingID = response.data.id

                // Upload image if provided
                if (imageUploadMethod.value === 'url' && form.value.image_url) {
                    // SSRF vulnerability: URL загрузка без валидации
                    await store.dispatch('listing/uploadImageFromUrl', {
                        listingID,
                        imageURL: form.value.image_url,
                    })
                } else if (imageUploadMethod.value === 'file' && form.value.image_file) {
                    // Upload file
                    await store.dispatch('listing/uploadImageFile', {
                        listingID,
                        file: form.value.image_file,
                    })
                }

                $q.notify({
                    type: 'positive',
                    message: 'Листинг успешно создан',
                })

                // Refresh listings after image upload completes
                await store.dispatch('listing/fetchListings')
                
                emit('created')
                closeModal()
            } catch (error) {
                console.error('Error creating listing:', error)
                const errorMessage = error.response?.data?.error || error.message || 'Ошибка при создании листинга'
                $q.notify({
                    type: 'negative',
                    message: errorMessage,
                })
            } finally {
                loading.value = false
            }
        }

        return {
            showModal,
            loading,
            imageUploadMethod,
            form,
            rarityOptions,
            handleSubmit,
            closeModal,
            onFileRejected,
        }
    },
}
</script>