<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import ColorPicker from '~/components/partials/ColorPicker.vue';
import IconSelectMenu from '~/components/partials/IconSelectMenu.vue';
import { getJobsColleaguesClient } from '~~/gen/ts/clients';
import type { Label } from '~~/gen/ts/resources/jobs/labels/labels';
import type { CreateOrUpdateLabelResponse } from '~~/gen/ts/services/jobs/colleagues';

const props = defineProps<{
    label?: Label;
}>();

const emits = defineEmits<{
    (e: 'close', v: boolean): void;
    (e: 'refresh'): void;
}>();

const jobsColleaguesClient = await getJobsColleaguesClient();

const schema = z.object({
    id: z.coerce.number(),
    job: z.coerce.string().max(20).optional(),
    name: z.coerce.string().min(1).max(48),
    color: z.coerce.string().length(7),
    icon: z.coerce.string().max(128).optional(),
    sortOrder: z.coerce.number().nonnegative().default(0),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    id: 0,
    job: undefined,
    name: '',
    color: '#ffffff',
    icon: undefined,
    sortOrder: 0,
});

const { hasUnsavedChanges, confirmLeave, syncSnapshot } = useSnapshotChanges(state, {
    serializer: (value) =>
        JSON.stringify({
            id: value.id,
            job: value.job ?? '',
            name: value.name,
            color: value.color,
            icon: value.icon ?? '',
            sortOrder: value.sortOrder,
        }),
});

function setFromData(label: Label | undefined): void {
    if (!label) {
        state.id = 0;
        state.job = undefined;
        state.name = '';
        state.color = '#ffffff';
        state.icon = undefined;
        state.sortOrder = 0;
        syncSnapshot();
        return;
    }

    state.id = label.id;
    state.job = label.job;
    state.name = label.name;
    state.color = label.color;
    state.icon = label.icon;
    state.sortOrder = label.sortOrder;
    syncSnapshot();
}

watch(
    () => props.label,
    () => setFromData(props.label),
    { immediate: true },
);

async function createOrUpdateLabel(values: Schema): Promise<CreateOrUpdateLabelResponse> {
    try {
        const { response } = await jobsColleaguesClient.createOrUpdateLabel({
            label: {
                id: values.id ?? 0,
                job: values.job,
                name: values.name ?? '',
                color: values.color ?? '#ffffff',
                icon: values.icon,
                sortOrder: values.sortOrder ?? 0,
            },
        });

        if (!response?.label) return response;

        const label = response.label;

        state.id = label.id;
        state.job = label.job;
        state.name = label.name;
        state.color = label.color;
        state.icon = label.icon;
        state.sortOrder = label.sortOrder;
        syncSnapshot();

        emits('refresh');
        emits('close', true);

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref<boolean>(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await createOrUpdateLabel(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const formRef = useTemplateRef('formRef');

async function closeModal(): Promise<void> {
    if (!canSubmit.value) return;

    if (hasUnsavedChanges.value && !(await confirmLeave())) return;

    emits('close', false);
}
</script>

<template>
    <UDrawer :title="$t('pages.jobs.colleagues.labels.title')" :close="false" :dismissible="!hasUnsavedChanges && canSubmit">
        <template #header>
            <div class="flex w-full items-center justify-between gap-2">
                <h3 class="font-semibold text-highlighted">
                    {{ $t('pages.jobs.colleagues.labels.title') }}
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
            <div class="mx-auto w-full max-w-(--breakpoint-xl)">
                <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                    <UFormField class="flex-1" name="name" :label="$t('common.name')">
                        <UInput v-model="state.name" class="w-full" name="name" type="text" :placeholder="$t('common.name')" />
                    </UFormField>

                    <UFormField name="color" :label="$t('common.color')">
                        <ColorPicker v-model="state.color" class="w-full" name="color" />
                    </UFormField>

                    <UFormField name="icon" :label="$t('common.icon')">
                        <IconSelectMenu v-model="state.icon" class="w-full" name="icon" :hex-color="state.color" clear />
                    </UFormField>
                </UForm>
            </div>
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
                    color="primary"
                    icon="i-mdi-content-save"
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    :label="$t('common.save')"
                    @click="() => formRef?.submit()"
                />
            </UFieldGroup>
        </template>
    </UDrawer>
</template>
