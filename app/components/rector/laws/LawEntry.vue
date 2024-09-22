<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import type { Law } from '~~/gen/ts/resources/laws/laws';

const props = defineProps<{
    law: Law;
    startInEdit?: boolean;
}>();

const emits = defineEmits<{
    (e: 'deleted', id: string): void;
    (e: 'update:law', update: { id: string; law: Law }): void;
}>();

const schema = z.object({
    name: z.string().min(3).max(128),
    description: z.union([z.string().min(3).max(500), z.string().length(0).optional()]),
    fine: z.number({ coerce: true }).nonnegative().max(999_999_999),
    detentionTime: z.number({ coerce: true }).nonnegative().max(999_999_999),
    stvoPoints: z.number({ coerce: true }).nonnegative().max(999_999_999),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    name: props.law.name,
    description: props.law.description,
    fine: props.law.fine ?? 0,
    detentionTime: props.law.detentionTime ?? 0,
    stvoPoints: props.law.stvoPoints ?? 0,
});

async function deleteLaw(id: string): Promise<void> {
    const i = parseInt(id);
    if (i < 0) {
        emits('deleted', id);
        return;
    }

    try {
        const call = getGRPCRectorLawsClient().deleteLaw({
            id,
        });
        await call;

        emits('deleted', id);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function saveLaw(lawBookId: string, id: string, values: Schema): Promise<void> {
    try {
        const call = getGRPCRectorLawsClient().createOrUpdateLaw({
            law: {
                id: parseInt(id) < 0 ? '0' : id,
                lawbookId: lawBookId,
                name: values.name,
                description: values.description,
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

const modal = useModal();

const editing = ref(props.startInEdit);
</script>

<template>
    <div v-if="!editing" class="table-row">
        <div class="flex table-cell flex-row py-2 pl-4 pr-3 text-sm font-medium sm:pl-1">
            <UButtonGroup class="inline-flex w-full">
                <UButton variant="link" icon="i-mdi-pencil" :title="$t('common.edit')" @click="editing = true" />
                <UButton
                    variant="link"
                    icon="i-mdi-trash-can"
                    color="red"
                    :title="$t('common.delete')"
                    @click="
                        modal.open(ConfirmModal, {
                            confirm: async () => deleteLaw(law.id),
                        })
                    "
                />
            </UButtonGroup>
        </div>
        <div class="table-cell py-2 pl-4 pr-3 text-sm font-medium sm:pl-1">
            {{ law.name }}
        </div>
        <div class="table-cell whitespace-nowrap p-1 text-left">${{ law.fine }}</div>
        <div class="table-cell whitespace-nowrap p-1 text-left">
            {{ law.detentionTime }}
        </div>
        <div class="table-cell whitespace-nowrap p-1 text-left">
            {{ law.stvoPoints }}
        </div>
        <div class="table-cell p-1 text-left text-sm font-medium">
            {{ law.description }}
        </div>
    </div>

    <UForm v-else :schema="schema" :state="state" class="flex table-row flex-1" @submit="onSubmitThrottle">
        <UFormGroup class="table-cell py-2 pl-4 pr-3 text-sm font-medium sm:pl-1">
            <UButtonGroup class="inline-flex w-full">
                <UButton type="submit" variant="link" icon="i-mdi-content-save" :title="$t('common.save')" />
                <UButton
                    variant="link"
                    icon="i-mdi-cancel"
                    :title="$t('common.cancel')"
                    @click="
                        editing = false;
                        parseInt(law.id) < 0 && $emit('deleted', law.id);
                    "
                />
            </UButtonGroup>
        </UFormGroup>
        <UFormGroup class="table-cell py-2 pl-4 pr-3 text-sm font-medium sm:pl-1">
            <UInput v-model="state.name" name="name" type="text" :placeholder="$t('common.crime')" />
        </UFormGroup>
        <UFormGroup class="table-cell whitespace-nowrap p-1 text-left">
            <UInput
                v-model="state.fine"
                name="fine"
                type="number"
                :placeholder="$t('common.fine')"
                :label="$t('common.fine')"
            />
        </UFormGroup>
        <UFormGroup class="table-cell whitespace-nowrap p-1 text-left">
            <UInput
                v-model="state.detentionTime"
                name="detentionTime"
                type="number"
                :placeholder="$t('common.detention_time')"
            />
        </UFormGroup>
        <UFormGroup class="table-cell whitespace-nowrap p-1 text-left">
            <UInput
                v-model="state.stvoPoints"
                name="stvoPoints"
                type="number"
                :placeholder="$t('common.traffic_infraction_points')"
            />
        </UFormGroup>
        <UFormGroup class="table-cell p-1 text-left">
            <UInput v-model="state.description" name="description" type="text" :placeholder="$t('common.description')" />
        </UFormGroup>
    </UForm>
</template>
