<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import NotSupportedTabletBlock from '~/components/partials/NotSupportedTabletBlock.vue';
import { useSettingsStore } from '~/stores/settings';
import { getFilestoreFilestoreClient } from '~~/gen/ts/clients';
import type { File } from '~~/gen/ts/resources/file/file';
import type { UploadFileResponse } from '~~/gen/ts/resources/file/filestore';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
    (e: 'uploaded', file: File): void;
}>();

const { fileUpload } = useAppConfig();

const notifications = useNotificationsStore();

const settingsStore = useSettingsStore();
const { nuiEnabled } = storeToRefs(settingsStore);

const filestoreFilestoreClient = await getFilestoreFilestoreClient();

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

async function upload(values: Schema): Promise<UploadFileResponse | undefined> {
    if (!values.file[0]) {
        return;
    }

    try {
        const call = filestoreFilestoreClient.upload({
            prefix: values.category,
            name: values.name,
        });
        const { response } = await call;

        if (response.file) {
            emit('uploaded', response.file);

            notifications.add({
                title: { key: 'notifications.action_successful.title', parameters: {} },
                description: { key: 'notifications.action_successful.content', parameters: {} },
                type: NotificationType.SUCCESS,
            });
        }

        emit('close', false);

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await upload(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UModal :title="$t('common.upload')">
        <template #body>
            <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
                <UFormField class="flex-1" name="category" :label="$t('common.category')" required>
                    <ClientOnly>
                        <USelectMenu
                            v-model="state.category"
                            :items="categories"
                            :search-input="{ placeholder: $t('common.search_field') }"
                        />
                    </ClientOnly>
                </UFormField>

                <UFormField class="flex-1" name="name" :label="$t('common.name')" required>
                    <UInput v-model="state.name" type="text" name="name" :placeholder="$t('common.name')" />
                </UFormField>

                <UFormField class="flex-1" name="file" :label="$t('common.file')">
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
                </UFormField>
            </UForm>
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton class="flex-1" block color="neutral" @click="$emit('close', false)">
                    {{ $t('common.close', 1) }}
                </UButton>

                <UButton class="flex-1" type="submit" block :disabled="!canSubmit" :loading="!canSubmit">
                    {{ $t('common.save') }}
                </UButton>
            </UButtonGroup>
        </template>
    </UModal>
</template>
