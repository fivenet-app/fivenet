<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
// eslint-disable-next-line camelcase
import { max, max_value, min, min_value, required } from '@vee-validate/rules';
import { useConfirmDialog, useThrottleFn } from '@vueuse/core';
import { CancelIcon, ContentSaveIcon, PencilIcon, TrashCanIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import ConfirmDialog from '~/components/partials/ConfirmDialog.vue';
import { Law } from '~~/gen/ts/resources/laws/laws';

const props = defineProps<{
    law: Law;
    startInEdit?: boolean;
}>();

const emit = defineEmits<{
    (e: 'deleted', id: bigint): void;
}>();

const { $grpc } = useNuxtApp();

async function deleteLaw(id: bigint): Promise<void> {
    if (id < 0) {
        emit('deleted', id);
        return;
    }

    try {
        const call = $grpc.getRectorClient().deleteLaw({
            id,
        });
        await call;

        emit('deleted', id);
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

interface FormData {
    name: string;
    description?: string;
    fine: bigint;
    detentionTime: bigint;
    stvoPoints: bigint;
}

async function saveLaw(lawBookId: bigint, id: bigint, values: FormData): Promise<void> {
    try {
        const call = $grpc.getRectorClient().createOrUpdateLaw({
            law: {
                id: BigInt(id < 0 ? 0 : id),
                lawbookId: lawBookId,
                name: values.name,
                description: values.description,
                fine: BigInt(values.fine),
                detentionTime: BigInt(values.detentionTime),
                stvoPoints: BigInt(values.stvoPoints),
            },
        });
        const { response } = await call;
        const law = response.law;
        if (law === undefined) {
            throw new Error('failed to get law from server response');
        }

        props.law.id = law.id;
        props.law.createdAt = law.createdAt;
        props.law.updatedAt = law.updatedAt;
        props.law.name = law.name;
        props.law.description = law.description;
        props.law.fine = law.fine;
        props.law.detentionTime = law.detentionTime;
        props.law.stvoPoints = law.stvoPoints;

        editing.value = false;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

defineRule('required', required);
defineRule('max', max);
defineRule('max_value', max_value);
defineRule('min', min);
defineRule('min_value', min_value);

const { handleSubmit, setValues } = useForm<FormData>({
    validationSchema: {
        name: { required: true, min: 3, max: 128 },
        description: { required: true, min: 6, max: 500 },
        fine: { required: false, min_value: 0, max_value: 999_999_999 },
        detentionTime: { required: false, min_value: 0, max_value: 999_999_999 },
        stvoPoints: { required: false, min_value: 0, max_value: 999_999_999 },
    },
    validateOnMount: true,
});

setValues({
    name: props.law.name,
    description: props.law.description,
    fine: props.law.fine,
    detentionTime: props.law.detentionTime,
    stvoPoints: props.law.stvoPoints,
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await saveLaw(props.law.lawbookId, props.law.id, values).finally(() => setTimeout(() => (canSubmit.value = true), 400)),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);

const { isRevealed, reveal, confirm, cancel, onConfirm } = useConfirmDialog();

onConfirm(async (id) => deleteLaw(id));

const editing = ref(props.startInEdit);
</script>

<template>
    <ConfirmDialog
        :open="isRevealed"
        :title="$t('components.partials.confirm_dialog.title')"
        :description="$t('components.partials.confirm_dialog.description')"
        :cancel="cancel"
        :confirm="() => confirm(law.id)"
    />

    <tr v-if="!editing">
        <td class="py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0 flex flex-row">
            <button type="button" class="pl-2" :title="$t('common.edit')" @click="editing = true">
                <PencilIcon class="w-6 h-6" />
            </button>
            <button type="button" class="pl-2" :title="$t('common.delete')" @click="reveal()">
                <TrashCanIcon class="w-6 h-6" />
            </button>
        </td>
        <td class="py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ law.name }}
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-base-200">${{ law.fine }}</td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-base-200">
            {{ law.detentionTime }}
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-base-200">
            {{ law.stvoPoints }}
        </td>
        <td class="px-1 py-1 text-sm font-medium text-left text-base-200">
            {{ law.description }}
        </td>
    </tr>
    <tr v-else>
        <td class="py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            <button type="button" :title="$t('common.save')" @click="onSubmitThrottle">
                <ContentSaveIcon class="w-6 h-6" />
            </button>
            <button
                type="button"
                :title="$t('common.cancel')"
                @click="
                    editing = false;
                    law.id < BigInt(0) && $emit('deleted', law.id);
                "
            >
                <CancelIcon class="w-6 h-6" />
            </button>
        </td>
        <td class="py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            <VeeField
                name="name"
                type="text"
                :placeholder="$t('common.crime')"
                :label="$t('common.crime')"
                class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                @focusin="focusTablet(true)"
                @focusout="focusTablet(false)"
            />
            <VeeErrorMessage name="name" as="p" class="mt-2 text-sm text-error-400" />
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-base-200">
            <VeeField
                name="fine"
                type="number"
                :placeholder="$t('common.fine')"
                :label="$t('common.fine')"
                class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
            />
            <VeeErrorMessage name="fine" as="p" class="mt-2 text-sm text-error-400" />
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-base-200">
            <VeeField
                name="detentionTime"
                type="number"
                :placeholder="$t('common.detention_time')"
                :label="$t('common.detention_time')"
                class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                @focusin="focusTablet(true)"
                @focusout="focusTablet(false)"
            />
            <VeeErrorMessage name="detentionTime" as="p" class="mt-2 text-sm text-error-400" />
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-base-200">
            <VeeField
                name="stvoPoints"
                type="text"
                :placeholder="$t('common.traffic_infraction_points')"
                :label="$t('common.traffic_infraction_points')"
                class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                @focusin="focusTablet(true)"
                @focusout="focusTablet(false)"
            />
            <VeeErrorMessage name="stvoPoints" as="p" class="mt-2 text-sm text-error-400" />
        </td>
        <td class="px-1 py-1 text-left text-base-200">
            <VeeField
                name="description"
                type="text"
                :placeholder="$t('common.description')"
                :label="$t('common.description')"
                class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                @focusin="focusTablet(true)"
                @focusout="focusTablet(false)"
            />
            <VeeErrorMessage name="description" as="p" class="mt-2 text-sm text-error-400" />
        </td>
    </tr>
</template>
