<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { useNotificatorStore } from '~/store/notificator';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { SetJobsUserPropsResponse } from '~~/gen/ts/services/jobs/jobs';

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

const notifications = useNotificatorStore();

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

async function setJobsUserNote(values: Schema): Promise<undefined | SetJobsUserPropsResponse> {
    try {
        const call = $grpc.jobs.jobs.setJobsUserProps({
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

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
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
    <UForm :schema="schema" :state="state" class="flex flex-1 flex-col gap-2" @submit="onSubmitThrottle">
        <div>
            <UButton v-if="!editing" icon="i-mdi-pencil" @click="editing = true" />
            <UButton
                v-else
                icon="i-mdi-cancel"
                color="red"
                @click="
                    state.note = modelValue ?? '';
                    editing = false;
                "
            />
        </div>

        <div class="flex flex-1 flex-col gap-2 sm:flex-row">
            <UFormGroup name="note" class="flex-1" :label="$t('common.note')">
                <UTextarea v-if="editing" v-model="state.note" block :rows="6" :maxrows="10" name="note" />
                <p v-else class="prose dark:prose-invert whitespace-pre-wrap text-base-800 dark:text-base-300">
                    {{ modelValue ?? $t('common.na') }}
                </p>
            </UFormGroup>
        </div>

        <template v-if="changed">
            <UFormGroup name="reason" :label="$t('common.reason')" required>
                <UInput v-model="state.reason" type="text" />
            </UFormGroup>

            <UButton type="submit" block icon="i-mdi-content-save" :disabled="!canSubmit" :loading="!canSubmit">
                {{ $t('common.save') }}
            </UButton>
        </template>
    </UForm>
</template>
