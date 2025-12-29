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
    category: z.coerce.string().min(3).max(255),
    name: z.coerce.string().min(3).max(255),
    file: z.file().mime(fileUpload.types.images).max(fileUpload.fileSizes.fileStore).min(1).optional(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    category: '',
    name: '',
    file: undefined,
});

const categories = ['jobassets'];

async function upload(values: Schema): Promise<UploadFileResponse | undefined> {
    if (!values.file) return;

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

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UModal :title="$t('common.upload')">
        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <UFormField class="flex-1" name="category" :label="$t('common.category')" required>
                    <ClientOnly>
                        <USelectMenu
                            v-model="state.category"
                            :items="categories"
                            class="w-full"
                            :search-input="{ placeholder: $t('common.search_field') }"
                        />
                    </ClientOnly>
                </UFormField>

                <UFormField class="flex-1" name="name" :label="$t('common.name')" required>
                    <UInput v-model="state.name" type="text" name="name" class="w-full" :placeholder="$t('common.name')" />
                </UFormField>

                <UFormField class="flex-1" name="file" :label="$t('common.file')">
                    <NotSupportedTabletBlock v-if="nuiEnabled" />
                    <UFileUpload
                        v-else
                        v-model="state.file"
                        name="file"
                        class="mx-auto max-w-md flex-1"
                        :accept="fileUpload.types.images.join(',')"
                        :placeholder="$t('common.image')"
                        :label="$t('common.file_upload_label')"
                        :description="$t('common.allowed_file_types')"
                    />
                </UFormField>
            </UForm>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton class="flex-1" block color="neutral" :label="$t('common.close', 1)" @click="$emit('close', false)" />

                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    :label="$t('common.save')"
                    @click="() => formRef?.submit()"
                />
            </UFieldGroup>
        </template>
    </UModal>
</template>
