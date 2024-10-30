<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import type { Law } from '~~/gen/ts/resources/laws/laws';

const props = defineProps<{
    law: Law;
    startInEdit?: boolean;
}>();

const emits = defineEmits<{
    (e: 'update:law', update: { id: string; law: Law }): void;
    (e: 'close'): void;
}>();

const schema = z.object({
    name: z.string().min(3).max(128),
    description: z.union([z.string().min(3).max(512), z.string().length(0).optional()]),
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

async function saveLaw(lawBookId: string, id: string, values: Schema): Promise<void> {
    try {
        const call = getGRPCRectorLawsClient().createOrUpdateLaw({
            law: {
                id: parseInt(id) < 0 ? '0' : id,
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

        emits('update:law', { id, law: response.law! });

        editing.value = false;
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

const editing = ref(props.startInEdit);
</script>

<template>
    <UForm :schema="schema" :state="state" class="my-2 flex flex-1 flex-col gap-2" @submit="onSubmitThrottle">
        <div class="flex flex-1 flex-row gap-2">
            <UFormGroup class="text-sm font-medium">
                <UButtonGroup class="inline-flex w-full" orientation="vertical">
                    <UButton type="submit" variant="link" icon="i-mdi-content-save" :title="$t('common.save')" />
                    <UButton
                        variant="link"
                        icon="i-mdi-cancel"
                        :title="$t('common.cancel')"
                        @click="
                            editing = false;
                            $emit('close');
                        "
                    />
                </UButtonGroup>
            </UFormGroup>

            <UFormGroup :label="$t('common.law')" class="flex-1 text-sm font-medium">
                <UInput v-model="state.name" name="name" type="text" :placeholder="$t('common.law')" />
            </UFormGroup>
        </div>

        <div class="flex flex-1 gap-2">
            <UFormGroup :label="$t('common.fine')" class="whitespace-nowrap text-left">
                <UInput
                    v-model="state.fine"
                    name="fine"
                    type="number"
                    :placeholder="$t('common.fine')"
                    leading-icon="i-mdi-dollar"
                />
            </UFormGroup>

            <UFormGroup :label="$t('common.detention_time')" class="whitespace-nowrap text-left">
                <UInput
                    v-model="state.detentionTime"
                    name="detentionTime"
                    type="number"
                    :placeholder="$t('common.detention_time')"
                />
            </UFormGroup>

            <UFormGroup :label="$t('common.traffic_infraction_points')" class="whitespace-nowrap text-left">
                <UInput
                    v-model="state.stvoPoints"
                    name="stvoPoints"
                    type="number"
                    :placeholder="$t('common.traffic_infraction_points')"
                />
            </UFormGroup>
        </div>

        <UFormGroup :label="$t('common.description')" class="text-left">
            <UInput v-model="state.description" name="description" type="text" :placeholder="$t('common.description')" />
        </UFormGroup>

        <UFormGroup :label="$t('common.hint')" class="text-left">
            <UInput v-model="state.hint" name="hint" type="text" :placeholder="$t('common.hint')" />
        </UFormGroup>
    </UForm>
</template>
