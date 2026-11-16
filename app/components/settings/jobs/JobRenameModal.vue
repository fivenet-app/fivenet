<script lang="ts" setup>
import type { Form, FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import ResponsiveDrawerModal from '~/components/partials/ResponsiveDrawerModal.vue';
import type { Job } from '~~/gen/ts/resources/jobs/jobs';
import { getSettingsSystemClient } from '~~/gen/ts/clients';

const props = defineProps<{
    jobs: Job[];
    onRenamed: () => Promise<void> | void;
}>();

const emit = defineEmits<{
    close: [renamed: boolean];
}>();

const { t } = useI18n();

const systemClient = await getSettingsSystemClient();

const schema = z
    .object({
        sourceJob: z.string().min(1).max(50),
        targetJob: z.string().min(1).max(50),
    })
    .refine((value) => value.sourceJob !== value.targetJob, {
        path: ['targetJob'],
        message: t('zod.custom.job_rename.target_different'),
    });

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    sourceJob: '',
    targetJob: '',
});

const formRef = useTemplateRef<Form<typeof schema>>('formRef');
const { hasUnsavedChanges, confirmLeave, syncSnapshot } = useSnapshotChanges(state);

async function renameJob(values: Schema): Promise<void> {
    try {
        await systemClient.renameJob({
            sourceJob: values.sourceJob,
            targetJob: values.targetJob,
        }).response;
        await props.onRenamed();
        syncSnapshot();
        emit('close', true);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const { submit, isSubmitting, canSubmit } = useSubmitGuard(async (event: FormSubmitEvent<Schema>) => {
    await renameJob(event.data);
});

async function closeModal(): Promise<void> {
    if (!canSubmit.value) return;
    if (hasUnsavedChanges.value && !(await confirmLeave())) return;
    emit('close', false);
}
</script>

<template>
    <ResponsiveDrawerModal
        :title="$t('pages.settings.jobs.rename.title')"
        :close="false"
        :dismissible="!hasUnsavedChanges && canSubmit"
        :ui="{ body: 'mx-auto w-full max-w-xl' }"
    >
        <template #header>
            <div class="flex w-full items-center justify-between gap-2">
                <h3 class="font-semibold text-highlighted">
                    {{ $t('pages.settings.jobs.rename.title') }}
                </h3>
                <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-close"
                    :disabled="!canSubmit"
                    :aria-label="$t('common.close', 1)"
                    @click="closeModal"
                />
            </div>
        </template>

        <template #body>
            <p class="mb-4 text-sm text-muted">
                {{ $t('pages.settings.jobs.rename.description') }}
            </p>

            <UForm ref="formRef" :schema="schema" :state="state" class="space-y-4" @submit="submit">
                <UFormField name="sourceJob" :label="$t('common.source')" required>
                    <USelectMenu
                        v-model="state.sourceJob"
                        class="w-full"
                        :items="props.jobs"
                        value-key="name"
                        label-key="label"
                        :filter-fields="['label', 'name']"
                        :search-input="{ placeholder: $t('common.search_field') }"
                    >
                        <template v-if="state.sourceJob" #default>
                            {{ props.jobs.find((job) => job.name === state.sourceJob)?.label ?? state.sourceJob }}
                            ({{ state.sourceJob }})
                        </template>
                        <template #item-label="{ item }"> {{ item.label }} ({{ item.name }}) </template>
                    </USelectMenu>
                </UFormField>

                <UFormField name="targetJob" :label="$t('common.target')" required>
                    <UInput
                        v-model="state.targetJob"
                        class="w-full"
                        name="targetJob"
                        type="text"
                        autocomplete="off"
                        maxlength="50"
                    />
                </UFormField>
            </UForm>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton
                    class="flex-1"
                    color="neutral"
                    block
                    :disabled="!canSubmit"
                    :label="$t('common.close', 1)"
                    @click="closeModal"
                />
                <UButton
                    class="flex-1"
                    color="warning"
                    block
                    icon="i-mdi-pencil"
                    :loading="isSubmitting"
                    :disabled="!canSubmit"
                    :label="$t('pages.settings.jobs.rename.action')"
                    @click="formRef?.submit()"
                />
            </UFieldGroup>
        </template>
    </ResponsiveDrawerModal>
</template>
