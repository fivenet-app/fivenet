<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import { getJobsColleaguesClient } from '~~/gen/ts/clients';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { SetColleaguePropsResponse } from '~~/gen/ts/services/jobs/colleagues';

const props = defineProps<{
    userId: number;
}>();

const emit = defineEmits<{
    (e: 'refresh'): void;
}>();

const modelValue = defineModel<string | undefined>({ default: undefined });

const notifications = useNotificationsStore();

const jobsColleaguesClient = await getJobsColleaguesClient();

const schema = z.object({
    reason: z.coerce.string().min(3).max(255),
    note: z.coerce.string().min(0).max(512),
});

type Schema = z.output<typeof schema>;

const editing = ref(false);

const state = reactive<Schema>({
    reason: '',
    note: modelValue.value ?? '',
});

const { snapshotDirty: changed, syncSnapshot } = useSnapshotChanges(state, {
    dirty: editing,
    serializer: (value) =>
        JSON.stringify({
            note: value.note ?? '',
        }),
});

function setFromProps(): void {
    state.note = modelValue.value ?? '';
    syncSnapshot();
}

watch(modelValue, () => {
    setFromProps();
});

async function setJobsUserNote(values: Schema): Promise<undefined | SetColleaguePropsResponse> {
    try {
        const call = jobsColleaguesClient.setColleagueProps({
            reason: values.reason,
            props: {
                userId: props.userId,
                job: '',
                note: values.note,
                groups: [],
            },
        });
        const { response } = await call;

        editing.value = false;
        state.reason = '';
        emit('refresh');

        state.note = response.props?.note ?? '';
        syncSnapshot();

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

const { submit, isSubmitting, canSubmit } = useSubmitGuard(async (event: FormSubmitEvent<Schema>) => {
    await setJobsUserNote(event.data);
});
</script>

<template>
    <UForm class="flex flex-1 flex-col gap-2" :schema="schema" :state="state" @submit="submit">
        <div>
            <UTooltip v-if="!editing" :text="$t('common.edit')">
                <UButton icon="i-mdi-pencil" @click="editing = true" />
            </UTooltip>
            <UTooltip v-else :text="$t('common.cancel')">
                <UButton
                    icon="i-mdi-cancel"
                    color="error"
                    @click="
                        state.reason = '';
                        editing = false;
                        setFromProps();
                    "
                />
            </UTooltip>
        </div>

        <div class="flex flex-1 flex-col gap-2 sm:flex-row">
            <UFormField class="flex-1" name="note" :label="$t('common.note')">
                <UTextarea v-if="editing" v-model="state.note" class="w-full" block :rows="6" :maxrows="10" name="note" />
                <p v-else class="prose whitespace-pre-wrap text-toned dark:prose-invert">
                    {{ modelValue ?? $t('common.na') }}
                </p>
            </UFormField>
        </div>

        <template v-if="editing">
            <UFormField name="reason" :label="$t('common.reason')" required>
                <UInput v-model="state.reason" class="w-full" type="text" :disabled="!changed" />
            </UFormField>

            <UButton
                type="submit"
                block
                icon="i-mdi-content-save"
                :disabled="!canSubmit || !changed"
                :loading="isSubmitting"
                :label="$t('common.save')"
            />
        </template>
    </UForm>
</template>
