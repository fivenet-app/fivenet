<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import NotSupportedTabletBlock from '~/components/partials/NotSupportedTabletBlock.vue';
import { useNotificatorStore } from '~/stores/notificator';
import { useSettingsStore } from '~/stores/settings';
import type { File, FileInfo } from '~~/gen/ts/resources/filestore/file';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { UploadFileResponse } from '~~/gen/ts/services/settings/filestore';

const emit = defineEmits<{
    (e: 'uploaded', file: FileInfo): void;
}>();

const { $grpc } = useNuxtApp();

const { fileUpload } = useAppConfig();

const { isOpen } = useModal();

const notifications = useNotificatorStore();

const settingsStore = useSettingsStore();
const { nuiEnabled } = storeToRefs(settingsStore);

const schema = z.object({
    category: z.string().min(3).max(255),
    name: z.string().min(3).max(255),
    file: zodFileSingleSchema(fileUpload.fileSizes.fileStore, fileUpload.types.images),
});

type Schema = z.output<typeof schema>;

const state = reactive({
    category: '',
    name: '',
    file: undefined,
});

const categories = ['jobassets'];

async function uploadFile(values: Schema): Promise<UploadFileResponse | undefined> {
    const file = {} as File;

    if (!values.file[0]) {
        return;
    }

    file.data = new Uint8Array(await values.file[0].arrayBuffer());

    try {
        const call = $grpc.settings.filestore.uploadFile({
            prefix: values.category,
            name: values.name,
            file: file,
        });
        const { response } = await call;

        if (response.file) {
            emit('uploaded', response.file);

            notifications.add({
                title: { key: 'notifications.action_successfull.title', parameters: {} },
                description: { key: 'notifications.action_successfull.content', parameters: {} },
                type: NotificationType.SUCCESS,
            });
        }

        isOpen.value = false;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await uploadFile(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl font-semibold leading-6">
                            {{ $t('common.upload') }}
                        </h3>

                        <UButton class="-my-1" color="gray" variant="ghost" icon="i-mdi-window-close" @click="isOpen = false" />
                    </div>
                </template>

                <div>
                    <UFormGroup class="flex-1" name="category" :label="$t('common.category')" required>
                        <ClientOnly>
                            <USelectMenu
                                v-model="state.category"
                                :options="categories"
                                :searchable-placeholder="$t('common.search_field')"
                            />
                        </ClientOnly>
                    </UFormGroup>

                    <UFormGroup class="flex-1" name="name" :label="$t('common.name')" required>
                        <UInput v-model="state.name" type="text" name="name" :placeholder="$t('common.name')" />
                    </UFormGroup>

                    <UFormGroup class="flex-1" name="file" :label="$t('common.file')">
                        <NotSupportedTabletBlock v-if="nuiEnabled" />
                        <template v-else>
                            <UInput
                                type="file"
                                name="file"
                                :accept="fileUpload.types.images.join(',')"
                                :placeholder="$t('common.image')"
                                @change="state.file = $event"
                            />
                        </template>
                    </UFormGroup>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton class="flex-1" block color="black" @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton class="flex-1" type="submit" block :disabled="!canSubmit" :loading="!canSubmit">
                            {{ $t('common.save') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>
