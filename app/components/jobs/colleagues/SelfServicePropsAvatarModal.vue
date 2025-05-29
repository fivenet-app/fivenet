<script lang="ts" setup>
import { useAuthStore } from '~/stores/auth';
import FileUpload from '../../partials/elements/FileUpload.vue';

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('components.jobs.self_service.set_profile_picture') }}
                    </h3>

                    <UButton class="-my-1" color="gray" variant="ghost" icon="i-mdi-window-close" @click="isOpen = false" />
                </div>
            </template>

            <div>
                <FileUpload
                    v-model="activeChar!.avatar"
                    :upload-fn="(opts) => $grpc.citizens.citizens.uploadAvatar(opts)"
                    :delete-fn="(_) => $grpc.citizens.citizens.deleteAvatar({})"
                />
            </div>

            <template #footer>
                <UButtonGroup class="inline-flex w-full">
                    <UButton class="flex-1" color="black" block @click="isOpen = false">
                        {{ $t('common.close', 1) }}
                    </UButton>
                </UButtonGroup>
            </template>
        </UCard>
    </UModal>
</template>
