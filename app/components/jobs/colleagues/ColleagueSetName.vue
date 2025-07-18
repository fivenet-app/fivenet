<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { SetColleaguePropsResponse } from '~~/gen/ts/services/jobs/jobs';

const props = defineProps<{
    namePrefix: string | undefined;
    nameSuffix: string | undefined;
    userId: number;
}>();

const emit = defineEmits<{
    (e: 'update:namePrefix', value?: string): void;
    (e: 'update:nameSuffix', value?: string): void;
    (e: 'refresh'): void;
}>();

const { namePrefix, nameSuffix } = useVModels(props, emit);

const { $grpc } = useNuxtApp();

const notifications = useNotificationsStore();

const schema = z.object({
    reason: z.string().min(3).max(255),
    prefix: z.string().max(12).optional(),
    suffix: z.string().max(12).optional(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    reason: '',
    prefix: props.namePrefix ?? '',
    suffix: props.nameSuffix ?? '',
});

watch(props, () => {
    state.prefix = props.namePrefix ?? '';
    state.suffix = props.nameSuffix ?? '';
});

const changed = ref(false);

async function setJobsUserNote(values: Schema): Promise<undefined | SetColleaguePropsResponse> {
    try {
        const call = $grpc.jobs.jobs.setColleagueProps({
            reason: values.reason,
            props: {
                userId: props.userId,
                job: '',
                namePrefix: values.prefix,
                nameSuffix: values.suffix,
            },
        });
        const { response } = await call;

        changed.value = false;
        editing.value = false;
        state.reason = '';
        emit('refresh');

        state.prefix = response.props?.namePrefix;
        state.suffix = response.props?.nameSuffix;

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
    if (state.prefix === namePrefix.value && state.suffix === nameSuffix.value) {
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
                        state.prefix = namePrefix;
                        state.suffix = nameSuffix;
                        editing = false;
                    "
                />
            </UTooltip>
        </div>

        <div class="flex flex-col gap-2 sm:flex-row">
            <UFormGroup name="prefix" :label="$t('common.prefix')">
                <UInput v-if="editing" v-model="state.prefix" type="text" />
                <span v-else>{{ namePrefix ?? $t('common.na') }}</span>
            </UFormGroup>
            <UFormGroup name="suffix" :label="$t('common.suffix')">
                <UInput v-if="editing" v-model="state.suffix" type="text" />
                <span v-else>{{ nameSuffix ?? $t('common.na') }}</span>
            </UFormGroup>
        </div>

        <template v-if="editing">
            <UFormGroup name="reason" :label="$t('common.reason')" required>
                <UInput v-model="state.reason" type="text" :disabled="!changed" />
            </UFormGroup>

            <UButton type="submit" block icon="i-mdi-content-save" :disabled="!changed || !canSubmit" :loading="!canSubmit">
                {{ $t('common.save') }}
            </UButton>
        </template>
    </UForm>
</template>
