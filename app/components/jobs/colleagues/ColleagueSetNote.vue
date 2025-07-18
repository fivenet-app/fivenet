<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { SetColleaguePropsResponse } from '~~/gen/ts/services/jobs/jobs';

const props = defineProps<{
    modelValue: string | undefined;
    userId: number;
}>();

const emit = defineEmits<{
    (e: 'update:namePrefix', value?: string): void;
    (e: 'update:nameSuffix', value?: string): void;
    (e: 'refresh'): void;
}>();

const { modelValue } = useVModels(props, emit);

const { $grpc } = useNuxtApp();

const notifications = useNotificationsStore();

const schema = z.object({
    reason: z.string().min(3).max(255),
    note: z.string().min(0).max(512),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    reason: '',
    note: props.modelValue ?? '',
});

watch(props, () => {
    state.note = props.modelValue ?? '';
});

const changed = ref(false);

async function setJobsUserNote(values: Schema): Promise<undefined | SetColleaguePropsResponse> {
    try {
        const call = $grpc.jobs.jobs.setColleagueProps({
            reason: values.reason,
            props: {
                userId: props.userId,
                job: '',
                note: values.note,
            },
        });
        const { response } = await call;

        editing.value = false;
        changed.value = false;
        state.reason = '';
        emit('refresh');

        state.note = response.props?.note ?? '';

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await setJobsUserNote(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

watch(state, () => {
    if (state.note === props.modelValue) {
        changed.value = false;
    } else {
        changed.value = true;
    }
});

const editing = ref(false);
</script>

<template>
    <UForm class="flex flex-1 flex-col gap-2" :schema="schema" :state="state" @submit="onSubmitThrottle">
        <div>
            <UTooltip v-if="!editing" :text="$t('common.edit')">
                <UButton icon="i-mdi-pencil" @click="editing = true" />
            </UTooltip>
            <UTooltip v-else :text="$t('common.cancel')">
                <UButton
                    icon="i-mdi-cancel"
                    color="error"
                    @click="
                        state.note = modelValue ?? '';
                        editing = false;
                    "
                />
            </UTooltip>
        </div>

        <div class="flex flex-1 flex-col gap-2 sm:flex-row">
            <UFormGroup class="flex-1" name="note" :label="$t('common.note')">
                <UTextarea v-if="editing" v-model="state.note" block :rows="6" :maxrows="10" name="note" />
                <p v-else class="prose dark:prose-invert whitespace-pre-wrap text-base-800 dark:text-base-300">
                    {{ modelValue ?? $t('common.na') }}
                </p>
            </UFormGroup>
        </div>

        <template v-if="editing">
            <UFormGroup name="reason" :label="$t('common.reason')" required>
                <UInput v-model="state.reason" type="text" :disabled="!changed" />
            </UFormGroup>

            <UButton type="submit" block icon="i-mdi-content-save" :disabled="!canSubmit || !changed" :loading="!canSubmit">
                {{ $t('common.save') }}
            </UButton>
        </template>
    </UForm>
</template>
