<script lang="ts" setup>
// eslint-disable-next-line camelcase
import { integer, max, max_value, min, min_value, required } from '@vee-validate/rules';
import { useConfirmDialog, useThrottleFn, useTimeoutFn } from '@vueuse/core';
import { CancelIcon, ContentSaveIcon, PencilIcon, TrashCanIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import ConfirmDialog from '~/components/partials/ConfirmDialog.vue';
import { Law } from '~~/gen/ts/resources/laws/laws';

const props = defineProps<{
    law: Law;
    startInEdit?: boolean;
}>();

const emits = defineEmits<{
    (e: 'deleted', id: string): void;
    (e: 'update:law', update: { id: string; law: Law }): void;
}>();

const { $grpc } = useNuxtApp();

async function deleteLaw(id: string): Promise<void> {
    const i = parseInt(id);
    if (i < 0) {
        emits('deleted', id);
        return;
    }

    try {
        const call = $grpc.getRectorLawsClient().deleteLaw({
            id,
        });
        await call;

        emits('deleted', id);
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

interface FormData {
    name: string;
    description?: string;
    fine: number;
    detentionTime: number;
    stvoPoints: number;
}

async function saveLaw(lawBookId: string, id: string, values: FormData): Promise<void> {
    try {
        const call = $grpc.getRectorLawsClient().createOrUpdateLaw({
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
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

defineRule('required', required);
defineRule('max', max);
defineRule('max_value', max_value);
defineRule('min', min);
defineRule('min_value', min_value);
defineRule('integer', integer);

const { handleSubmit, setValues } = useForm<FormData>({
    validationSchema: {
        name: { required: true, min: 3, max: 128 },
        description: { required: true, min: 6, max: 500 },
        fine: { required: false, integer: true, min_value: 0, max_value: 999_999_999 },
        detentionTime: { required: false, integer: true, min_value: 0, max_value: 999_999_999 },
        stvoPoints: { required: false, integer: true, min_value: 0, max_value: 999_999_999 },
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
        await saveLaw(props.law.lawbookId, props.law.id, values).finally(() =>
            useTimeoutFn(() => (canSubmit.value = true), 400),
        ),
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
        <td class="flex flex-row py-2 pl-4 pr-3 text-sm font-medium sm:pl-1">
            <UButton class="pl-2" :title="$t('common.edit')" @click="editing = true">
                <PencilIcon class="size-5" />
            </UButton>
            <UButton class="pl-2" :title="$t('common.delete')" @click="reveal()">
                <TrashCanIcon class="size-5" />
            </UButton>
        </td>
        <td class="py-2 pl-4 pr-3 text-sm font-medium sm:pl-1">
            {{ law.name }}
        </td>
        <td class="whitespace-nowrap p-1 text-left text-accent-200">${{ law.fine }}</td>
        <td class="whitespace-nowrap p-1 text-left text-accent-200">
            {{ law.detentionTime }}
        </td>
        <td class="whitespace-nowrap p-1 text-left text-accent-200">
            {{ law.stvoPoints }}
        </td>
        <td class="p-1 text-left text-sm font-medium text-accent-200">
            {{ law.description }}
        </td>
    </tr>
    <tr v-else>
        <td class="py-2 pl-4 pr-3 text-sm font-medium sm:pl-1">
            <UButton :title="$t('common.save')" @click="onSubmitThrottle">
                <ContentSaveIcon class="size-5" />
            </UButton>
            <UButton
                :title="$t('common.cancel')"
                @click="
                    editing = false;
                    parseInt(law.id) < 0 && $emit('deleted', law.id);
                "
            >
                <CancelIcon class="size-5" />
            </UButton>
        </td>
        <td class="py-2 pl-4 pr-3 text-sm font-medium sm:pl-1">
            <VeeField
                name="name"
                type="text"
                :placeholder="$t('common.crime')"
                :label="$t('common.crime')"
                class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                @focusin="focusTablet(true)"
                @focusout="focusTablet(false)"
            />
            <VeeErrorMessage name="name" as="p" class="mt-2 text-sm text-error-400" />
        </td>
        <td class="whitespace-nowrap p-1 text-left text-accent-200">
            <VeeField
                name="fine"
                type="text"
                :placeholder="$t('common.fine')"
                :label="$t('common.fine')"
                class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
            />
            <VeeErrorMessage name="fine" as="p" class="mt-2 text-sm text-error-400" />
        </td>
        <td class="whitespace-nowrap p-1 text-left text-accent-200">
            <VeeField
                name="detentionTime"
                type="text"
                :placeholder="$t('common.detention_time')"
                :label="$t('common.detention_time')"
                class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                @focusin="focusTablet(true)"
                @focusout="focusTablet(false)"
            />
            <VeeErrorMessage name="detentionTime" as="p" class="mt-2 text-sm text-error-400" />
        </td>
        <td class="whitespace-nowrap p-1 text-left text-accent-200">
            <VeeField
                name="stvoPoints"
                type="text"
                :placeholder="$t('common.traffic_infraction_points')"
                :label="$t('common.traffic_infraction_points')"
                class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                @focusin="focusTablet(true)"
                @focusout="focusTablet(false)"
            />
            <VeeErrorMessage name="stvoPoints" as="p" class="mt-2 text-sm text-error-400" />
        </td>
        <td class="p-1 text-left text-accent-200">
            <VeeField
                name="description"
                type="text"
                :placeholder="$t('common.description')"
                :label="$t('common.description')"
                class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                @focusin="focusTablet(true)"
                @focusout="focusTablet(false)"
            />
            <VeeErrorMessage name="description" as="p" class="mt-2 text-sm text-error-400" />
        </td>
    </tr>
</template>
