<script lang="ts" setup>
import type { Form, FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import type { JobAsset } from '~~/gen/ts/resources/jobs/job_asset';
import { getSettingsJobsassetsClient } from '~~/gen/ts/clients';

const props = defineProps<{
    asset?: JobAsset;
}>();

const emit = defineEmits<{
    close: [];
    updated: [];
}>();

const schema = z.object({
    displayName: z.string().trim().min(1).max(255),
});

type Schema = z.output<typeof schema>;

const { fileUpload } = useAppConfig();
const jobsClient = await getSettingsJobsassetsClient();
const state = reactive<Schema>({ displayName: '' });
const replacementFile = ref<File>();
const replacementAssetId = ref(0);

const { resizeAndUpload: replaceAsset } = useFileUploader(
    () => jobsClient.replaceJobAsset({}),
    'jobassets',
    replacementAssetId,
);

watch(
    () => props.asset,
    (asset) => {
        state.displayName = asset ? getJobAssetDisplayName(asset) : '';
        replacementFile.value = undefined;
    },
    { immediate: true },
);

function closeEditor(): void {
    emit('close');
}

const formRef = useTemplateRef<Form<typeof schema>>('formRef');

async function saveAsset(event: FormSubmitEvent<Schema>): Promise<void> {
    if (!props.asset) return;
    try {
        await jobsClient.updateJobAsset({
            id: props.asset.id,
            displayName: event.data.displayName,
        });
        if (replacementFile.value) {
            replacementAssetId.value = props.asset.id;
            await replaceAsset(replacementFile.value);
        }
        emit('updated');
        closeEditor();
    } catch (e) {
        handleGRPCError(e as RpcError);
    }
}
</script>

<template>
    <UModal :open="asset !== undefined" :title="$t('common.edit')" @update:open="(open) => !open && closeEditor()">
        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" class="grid gap-4" @submit="saveAsset">
                <UFormField name="displayName" :label="$t('common.display_name')" required>
                    <UInput v-model="state.displayName" class="w-full" />
                </UFormField>

                <UFormField :label="$t('common.file')">
                    <UFileUpload
                        v-model="replacementFile"
                        :accept="fileUpload.types.images.join(',')"
                        :label="$t('common.replace')"
                        :description="
                            $t('components.jobs.job_assets.replace_file_description', {
                                allowedFileTypes: $t('common.allowed_file_types'),
                            })
                        "
                    />
                </UFormField>
            </UForm>
        </template>

        <template #footer>
            <UFieldGroup class="flex w-full gap-2">
                <UButton class="flex-1" color="neutral" :label="$t('common.cancel')" @click="closeEditor" />
                <UButton
                    class="flex-1"
                    :disabled="!state.displayName.trim()"
                    :label="$t('common.save')"
                    @click="formRef?.submit()"
                />
            </UFieldGroup>
        </template>
    </UModal>
</template>
