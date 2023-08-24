<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { max, max_value, min, min_value, required } from '@vee-validate/rules';
import { CancelIcon, ContentSaveIcon, PencilIcon, TrashCanIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import { Law } from '~~/gen/ts/resources/laws/laws';

const props = defineProps<{
    law: Law;
    startInEdit?: boolean;
}>();

const emits = defineEmits<{
    (e: 'deleted', id: bigint): void;
}>();

const { $grpc } = useNuxtApp();

async function deleteLaw(id: bigint): Promise<void> {
    return new Promise(async (res, rej) => {
        if (id < 0) {
            emits('deleted', id);
            return;
        }

        try {
            const call = $grpc.getRectorClient().deleteLaw({
                id: id,
            });
            await call;

            emits('deleted', id);

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function saveLaw(lawBookId: bigint, id: bigint, values: FormData): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getRectorClient().createOrUpdateLaw({
                id: BigInt(id < 0 ? 0 : id),
                lawbookId: lawBookId,
                name: values.name,
                description: values.description,
                fine: BigInt(values.fine),
                detentionTime: BigInt(values.detentionTime),
                stvoPoints: BigInt(values.stvoPoints),
            });
            const { response } = await call;

            props.law.id = response.id;
            props.law.createdAt = response.createdAt;
            props.law.updatedAt = response.updatedAt;

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

defineRule('required', required);
defineRule('max', max);
defineRule('max_value', max_value);
defineRule('min', min);
defineRule('min_value', min_value);

interface FormData {
    name: string;
    description?: string;
    fine: bigint;
    detentionTime: bigint;
    stvoPoints: bigint;
}

const { handleSubmit } = useForm<FormData>({
    validationSchema: {
        name: { required: true, min: 3, max: 128 },
        description: { required: true, min: 6, max: 500 },
        fine: { required: false, min_value: 0, max_value: 999_999_999 },
        detentionTime: { required: false, min_value: 0, max: 999_999_999 },
        stvoPoints: { required: false, min_value: 0, max: 999_999_999 },
    },
    initialValues: {
        name: props.law.name,
        description: props.law.description,
        fine: props.law.fine,
        detentionTime: props.law.detentionTime,
        stvoPoints: props.law.stvoPoints,
    },
    validateOnMount: true,
});

const onSubmit = handleSubmit(async (values): Promise<void> => await saveLaw(props.law.lawbookId, props.law.id, values));

const editing = ref(props.startInEdit);
</script>

<template>
    <tr v-if="!editing">
        <td class="py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            <button type="button" class="pl-2" @click="editing = true">
                <PencilIcon class="w-6 h-6" />
            </button>
            <button type="button" class="pl-2" @click="deleteLaw(law.id)">
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
            <button type="button" @click="onSubmit">
                <ContentSaveIcon class="w-6 h-6" />
            </button>
            <button
                type="button"
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
            />
            <VeeErrorMessage name="stvoPoints" as="p" class="mt-2 text-sm text-error-400" />
        </td>
        <td class="px-1 py-1 text-left text-base-200">
            <VeeField
                name="description"
                type="text"
                :placeholder="$t('common.detention_time')"
                :label="$t('common.detention_time')"
                class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
            />
            <VeeErrorMessage name="description" as="p" class="mt-2 text-sm text-error-400" />
        </td>
    </tr>
</template>
