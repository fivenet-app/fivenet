<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { max, min, required } from '@vee-validate/rules';
import { useConfirmDialog, useThrottleFn } from '@vueuse/core';
import { CancelIcon, ContentSaveIcon, PencilIcon, TrashCanIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import ConfirmDialog from '~/components/partials/ConfirmDialog.vue';
import { Law, LawBook } from '~~/gen/ts/resources/laws/laws';
import LawEntry from '~/components/rector/laws/LawEntry.vue';

const props = defineProps<{
    modelValue: LawBook;
    laws: Law[];
    startInEdit?: boolean;
}>();

const emit = defineEmits<{
    (e: 'deleted', id: string): void;
    (e: 'update:modelValue', book: LawBook): void;
    (e: 'update:laws', laws: Law[]): void;
    (e: 'update:law', law: Law): void;
}>();

const { $grpc } = useNuxtApp();

async function deleteLawBook(id: string): Promise<void> {
    const i = parseInt(id, 10);
    if (i < 0) {
        emit('deleted', id);
        return;
    }

    try {
        const call = $grpc.getRectorClient().deleteLawBook({ id });
        await call;

        emit('deleted', id);
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

interface FormData {
    id: string;
    name: string;
    description?: string;
}

async function saveLawBook(id: string, values: FormData): Promise<LawBook> {
    const i = parseInt(id, 10);

    try {
        const call = $grpc.getRectorClient().createOrUpdateLawBook({
            lawBook: {
                id: i < 0 ? '0' : id,
                name: values.name,
                description: values.description,
                laws: [],
            },
        });
        const { response } = await call;

        editing.value = false;

        emit('update:modelValue', response.lawBook!);

        return response.lawBook!;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

defineRule('required', required);
defineRule('max', max);
defineRule('min', min);

const { handleSubmit, setValues } = useForm<FormData>({
    validationSchema: {
        name: { required: true, min: 3, max: 128 },
        description: { required: false, min: 3, max: 255 },
    },
    validateOnMount: true,
});
setValues({
    name: props.modelValue.name,
    description: props.modelValue.description,
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<LawBook> =>
        await saveLawBook(props.modelValue.id, values).finally(() => setTimeout(() => (canSubmit.value = true), 400)),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);

function deletedLaw(id: string): void {
    emit(
        'update:laws',
        props.laws.filter((b) => b.id !== id),
    );
}

const lastNewId = ref(-1);

function addLaw(): void {
    emit('update:laws', [
        ...props.laws,
        {
            lawbookId: props.modelValue.id,
            id: lastNewId.value.toString(),
            name: '',
            fine: 0,
            detentionTime: 0,
            stvoPoints: 0,
        },
    ]);
    lastNewId.value--;
}

const { isRevealed, reveal, confirm, cancel, onConfirm } = useConfirmDialog();

onConfirm(async (id) => deleteLawBook(id));

const editing = ref(props.startInEdit);
</script>

<template>
    <ConfirmDialog :open="isRevealed" :cancel="cancel" :confirm="() => confirm(modelValue.id)" />

    <div class="my-2">
        <div v-if="!editing" class="flex items-center gap-x-2 text-neutral">
            <button type="button" :title="$t('common.edit')" @click="editing = true">
                <PencilIcon class="h-5 w-5" />
            </button>
            <button type="button" :title="$t('common.delete')" @click="reveal()">
                <TrashCanIcon class="h-5 w-5" />
            </button>
            <h2 class="text-xl">{{ modelValue.name }}</h2>
            <p v-if="modelValue.description" class="pl-2">- {{ $t('common.description') }}: {{ modelValue.description }}</p>
            <button
                type="button"
                class="ml-auto rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                @click="addLaw"
            >
                {{ $t('pages.rector.laws.add_new_law') }}
            </button>
        </div>
        <form v-else class="flex w-full flex-row items-start gap-x-4 text-neutral" @submit.prevent="onSubmitThrottle">
            <button type="submit" :title="$t('common.save')">
                <ContentSaveIcon class="h-5 w-5" />
            </button>
            <button
                type="button"
                :title="$t('common.cancel')"
                @click="
                    editing = false;
                    parseInt(modelValue.id, 10) < 0 && $emit('deleted', modelValue.id);
                "
            >
                <CancelIcon class="h-5 w-5" />
            </button>

            <div class="flex-initial">
                <label for="name">
                    {{ $t('common.law_book') }}
                </label>
                <VeeField
                    name="name"
                    type="text"
                    :placeholder="$t('common.law_book')"
                    :label="$t('common.law_book')"
                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                    @focusin="focusTablet(true)"
                    @focusout="focusTablet(false)"
                />
                <VeeErrorMessage name="name" as="p" class="mt-2 text-sm text-error-400" />
            </div>
            <div class="flex-auto">
                <label for="description">
                    {{ $t('common.description') }}
                </label>
                <VeeField
                    name="description"
                    type="text"
                    :placeholder="$t('common.description')"
                    :label="$t('common.description')"
                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                    @focusin="focusTablet(true)"
                    @focusout="focusTablet(false)"
                />
                <VeeErrorMessage name="description" as="p" class="mt-2 text-sm text-error-400" />
            </div>
        </form>
        <table class="min-w-full divide-y divide-base-600">
            <thead>
                <tr>
                    <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-neutral sm:pl-1">
                        {{ $t('common.action', 2) }}
                    </th>
                    <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-neutral sm:pl-1">
                        {{ $t('common.crime') }}
                    </th>
                    <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-neutral">
                        {{ $t('common.fine') }}
                    </th>
                    <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-neutral">
                        {{ $t('common.detention_time') }}
                    </th>
                    <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-neutral">
                        {{ $t('common.traffic_infraction_points', 2) }}
                    </th>
                    <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-neutral">
                        {{ $t('common.description') }}
                    </th>
                </tr>
            </thead>
            <tbody class="divide-y divide-base-800">
                <LawEntry
                    v-for="law in modelValue.laws"
                    :key="law.id"
                    :law="law"
                    :start-in-edit="parseInt(law.id, 10) < 0"
                    @update:law="$emit('update:law', $event)"
                    @deleted="deletedLaw($event)"
                />
            </tbody>
        </table>
    </div>
</template>
