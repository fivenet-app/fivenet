<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import { getSettingsLawsClient } from '~~/gen/ts/clients';
import type { Law } from '~~/gen/ts/resources/laws/laws';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = defineProps<{
    law: Law;
}>();

const emit = defineEmits<{
    (e: 'update:law', update: { id: number; law: Law }): void;
    (e: 'close'): void;
}>();

const notifications = useNotificationsStore();

const settingsLawsClient = await getSettingsLawsClient();

const schema = z.object({
    name: z.string().min(3).max(128),
    description: z.union([z.string().min(3).max(1024), z.string().length(0).optional()]),
    hint: z.union([z.string().min(3).max(512), z.string().length(0).optional()]),
    fine: z.number({ coerce: true }).nonnegative().max(999_999_999),
    detentionTime: z.number({ coerce: true }).nonnegative().max(999_999_999),
    stvoPoints: z.number({ coerce: true }).nonnegative().max(999_999_999),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    name: props.law.name,
    description: props.law.description,
    hint: props.law.hint,
    fine: props.law.fine ?? 0,
    detentionTime: props.law.detentionTime ?? 0,
    stvoPoints: props.law.stvoPoints ?? 0,
});

async function saveLaw(lawBookId: number, id: number, values: Schema): Promise<void> {
    try {
        const call = settingsLawsClient.createOrUpdateLaw({
            law: {
                id: id < 0 ? 0 : id,
                lawbookId: lawBookId,
                name: values.name,
                description: values.description,
                hint: values.hint,
                fine: values.fine,
                detentionTime: values.detentionTime,
                stvoPoints: values.stvoPoints,
            },
        });
        const { response } = await call;

        emit('update:law', { id: id, law: response.law! });

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await saveLaw(props.law.lawbookId, props.law.id, event.data).finally(() =>
        useTimeoutFn(() => (canSubmit.value = true), 400),
    );
}, 1000);
</script>

<template>
    <UForm class="my-2 flex flex-1 flex-col gap-2" :schema="schema" :state="state" @submit="onSubmitThrottle">
        <div class="flex flex-1 flex-row gap-2">
            <UFormField class="text-sm font-medium">
                <UButtonGroup class="inline-flex w-full" orientation="vertical">
                    <UTooltip :text="$t('common.save')">
                        <UButton type="submit" variant="link" icon="i-mdi-content-save" />
                    </UTooltip>

                    <UTooltip :text="$t('common.cancel')">
                        <UButton variant="link" icon="i-mdi-cancel" @click="$emit('close')" />
                    </UTooltip>
                </UButtonGroup>
            </UFormField>

            <UFormField class="flex-1 text-sm font-medium" :label="$t('common.law')" name="name">
                <UInput v-model="state.name" name="name" type="text" class="w-full" :placeholder="$t('common.law')" />
            </UFormField>
        </div>

        <div class="flex flex-1 justify-between gap-2">
            <UFormField :label="$t('common.fine')" name="fine">
                <UInputNumber
                    v-model="state.fine"
                    name="fine"
                    :min="0"
                    :step="1000"
                    :format-options="{
                        style: 'currency',
                        currency: 'USD',
                        currencyDisplay: 'code',
                        currencySign: 'accounting',
                    }"
                    :placeholder="$t('common.fine')"
                />
            </UFormField>

            <UFormField :label="$t('common.detention_time')" name="detentionTime">
                <UInputNumber
                    v-model="state.detentionTime"
                    name="detentionTime"
                    :min="0"
                    :step="1"
                    :placeholder="$t('common.detention_time')"
                />
            </UFormField>

            <UFormField :label="$t('common.traffic_infraction_points')" name="stvoPoints">
                <UInputNumber
                    v-model="state.stvoPoints"
                    name="stvoPoints"
                    :min="0"
                    :step="1"
                    :placeholder="$t('common.traffic_infraction_points')"
                />
            </UFormField>
        </div>

        <UFormField :label="$t('common.description')" name="description">
            <UTextarea
                v-model="state.description"
                name="description"
                type="text"
                class="w-full"
                :placeholder="$t('common.description')"
            />
        </UFormField>

        <UFormField :label="$t('common.hint')" name="hint">
            <UTextarea v-model="state.hint" name="hint" type="text" class="w-full" :placeholder="$t('common.hint')" />
        </UFormField>
    </UForm>
</template>
