<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import type { Law } from '~~/gen/ts/resources/laws/laws';

const props = defineProps<{
    law: Law;
}>();

const emit = defineEmits<{
    (e: 'update:law', update: { id: number; law: Law }): void;
    (e: 'close'): void;
}>();

const { $grpc } = useNuxtApp();

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
        const call = $grpc.settings.laws.createOrUpdateLaw({
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
            <UFormGroup class="text-sm font-medium">
                <UButtonGroup class="inline-flex w-full" orientation="vertical">
                    <UTooltip :text="$t('common.save')">
                        <UButton type="submit" variant="link" icon="i-mdi-content-save" />
                    </UTooltip>

                    <UTooltip :text="$t('common.cancel')">
                        <UButton variant="link" icon="i-mdi-cancel" @click="$emit('close')" />
                    </UTooltip>
                </UButtonGroup>
            </UFormGroup>

            <UFormGroup class="flex-1 text-sm font-medium" :label="$t('common.law')" name="name">
                <UInput v-model="state.name" name="name" type="text" :placeholder="$t('common.law')" />
            </UFormGroup>
        </div>

        <div class="flex flex-1 gap-2">
            <UFormGroup class="whitespace-nowrap text-left" :label="$t('common.fine')" name="fine">
                <UInput
                    v-model="state.fine"
                    name="fine"
                    type="number"
                    :min="0"
                    :placeholder="$t('common.fine')"
                    leading-icon="i-mdi-dollar"
                />
            </UFormGroup>

            <UFormGroup class="whitespace-nowrap text-left" :label="$t('common.detention_time')" name="detentionTime">
                <UInput
                    v-model="state.detentionTime"
                    name="detentionTime"
                    type="number"
                    :min="0"
                    :placeholder="$t('common.detention_time')"
                />
            </UFormGroup>

            <UFormGroup class="whitespace-nowrap text-left" :label="$t('common.traffic_infraction_points')" name="stvoPoints">
                <UInput
                    v-model="state.stvoPoints"
                    name="stvoPoints"
                    type="number"
                    :min="0"
                    :placeholder="$t('common.traffic_infraction_points')"
                />
            </UFormGroup>
        </div>

        <UFormGroup class="text-left" :label="$t('common.description')" name="description">
            <UInput v-model="state.description" name="description" type="text" :placeholder="$t('common.description')" />
        </UFormGroup>

        <UFormGroup class="text-left" :label="$t('common.hint')" name="hint">
            <UInput v-model="state.hint" name="hint" type="text" :placeholder="$t('common.hint')" />
        </UFormGroup>
    </UForm>
</template>
