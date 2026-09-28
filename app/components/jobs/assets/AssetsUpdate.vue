<script lang="ts" setup>
import type { Form, FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import NotSupportedTabletBlock from '~/components/partials/NotSupportedTabletBlock.vue';
import { useSettingsStore } from '~/stores/settings';
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
const settingsStore = useSettingsStore();
const { nuiEnabled } = storeToRefs(settingsStore);
const jobsClient = await getSettingsJobsassetsClient();
const state = reactive<Schema>({ displayName: '' });
const replacementFile = ref<File>();
const replacementAssetId = ref(0);

const { resizeAndUpload: replaceAsset } = useFileUploader(
    () => jobsClient.replaceJobAsset({}),
    'jobassets',
    replacementAssetId,
);

const formSnapshot = computed(() => ({
    displayName: state.displayName,
}));

const { hasUnsavedChanges, confirmLeave, syncSnapshot } = useSnapshotChanges(formSnapshot, {
    dirty: computed(() => replacementFile.value !== undefined),
});

watch(
    () => props.asset,
    (asset) => {
        state.displayName = asset ? getJobAssetDisplayName(asset) : '';
        replacementFile.value = undefined;
        syncSnapshot();
    },
    { immediate: true },
);

async function closeEditor(): Promise<void> {
    if (hasUnsavedChanges.value && !(await confirmLeave())) return;

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
        replacementFile.value = undefined;
        syncSnapshot();
        emit('updated');
        closeEditor();
    } catch (e) {
        handleGRPCError(e as RpcError);
    }
}

const { submit: submitAsset, isSubmitting } = useSubmitGuard(saveAsset);
</script>

<template>
    <UModal
        :open="asset !== undefined"
        :title="$t('common.edit')"
        :close="false"
        :dismissible="!hasUnsavedChanges"
        @update:open="(open) => !open && closeEditor()"
    >
        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" class="grid gap-4" @submit="submitAsset">
                <UFormField name="displayName" :label="$t('common.display_name')" required>
                    <UInput v-model="state.displayName" class="w-full" />
                </UFormField>

                <UFormField :label="$t('common.file')">
                    <NotSupportedTabletBlock v-if="nuiEnabled" />
                    <div v-else class="grid gap-2">
                        <UFileUpload
                            v-model="replacementFile"
                            position="inside"
                            :accept="fileUpload.types.images.join(',')"
                            :placeholder="$t('common.image')"
                            :label="$t('common.replace')"
                            :description="
                                $t('components.jobs.job_assets.replace_file_description', {
                                    allowedFileTypes: $t('common.allowed_file_types'),
                                })
                            "
                        />
                        <PartialsContentGuidelinesAlert />
                    </div>
                </UFormField>
            </UForm>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" :label="$t('common.cancel')" @click="closeEditor" />

                <UButton
                    class="flex-1"
                    :disabled="!state.displayName.trim() || isSubmitting"
                    :loading="isSubmitting"
                    :label="$t('common.save')"
                    @click="formRef?.submit()"
                />
            </UFieldGroup>
        </template>
    </UModal>
</template>
