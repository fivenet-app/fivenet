<script lang="ts" setup>
import { max, min, required } from '@vee-validate/rules';
import { defineRule } from 'vee-validate';
import { Law, LawBook } from '~~/gen/ts/resources/laws/laws';
import LawEntry from '~/components/rector/laws/LawEntry.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';

const props = defineProps<{
    modelValue: LawBook;
    laws: Law[];
    startInEdit?: boolean;
}>();

const emit = defineEmits<{
    (e: 'deleted', id: string): void;
    (e: 'update:modelValue', book: LawBook): void;
    (e: 'update:laws', laws: Law[]): void;
    (e: 'update:law', update: { id: string; law: Law }): void;
}>();

const { $grpc } = useNuxtApp();

async function deleteLawBook(id: string): Promise<void> {
    const i = parseInt(id);
    if (i < 0) {
        emit('deleted', id);
        return;
    }

    try {
        const call = $grpc.getRectorLawsClient().deleteLawBook({ id });
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
    const i = parseInt(id);

    try {
        const call = $grpc.getRectorLawsClient().createOrUpdateLawBook({
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
        await saveLawBook(props.modelValue.id, values).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400)),
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

const modal = useModal();

const editing = ref(props.startInEdit);
</script>

<template>
    <UCard>
        <template #header>
            <div v-if="!editing" class="flex items-center gap-x-2">
                <UButtonGroup class="inline-flex w-full">
                    <UButton variant="link" icon="i-mdi-pencil" :title="$t('common.edit')" @click="editing = true" />
                    <UButton
                        variant="link"
                        icon="i-mdi-trash-can"
                        :title="$t('common.delete')"
                        @click="
                            modal.open(ConfirmModal, {
                                confirm: async () => deleteLawBook(modelValue.id),
                            })
                        "
                    />
                </UButtonGroup>

                <h2 class="text-xl">{{ modelValue.name }}</h2>

                <p v-if="modelValue.description" class="pl-2">- {{ $t('common.description') }}: {{ modelValue.description }}</p>

                <UButton @click="addLaw">
                    {{ $t('pages.rector.laws.add_new_law') }}
                </UButton>
            </div>
            <UForm v-else :state="{}" class="flex w-full flex-row items-start gap-x-4">
                <UButton :title="$t('common.save')" icon="i-mdi-content-save-icon" @click="onSubmitThrottle" />
                <UButton
                    :title="$t('common.cancel')"
                    icon="i-mdi-cancel"
                    @click="
                        editing = false;
                        parseInt(modelValue.id) < 0 && $emit('deleted', modelValue.id);
                    "
                />

                <div class="flex-initial">
                    <label for="name">
                        {{ $t('common.law_book') }}
                    </label>
                    <VeeField
                        name="name"
                        type="text"
                        :placeholder="$t('common.law_book')"
                        :label="$t('common.law_book')"
                        class="placeholder:text-accent-200 block w-full rounded-md border-0 bg-base-700 py-1.5 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
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
                        class="placeholder:text-accent-200 block w-full rounded-md border-0 bg-base-700 py-1.5 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
                    />
                    <VeeErrorMessage name="description" as="p" class="mt-2 text-sm text-error-400" />
                </div>
            </UForm>
        </template>

        <table class="min-w-full divide-y divide-base-600">
            <thead>
                <tr>
                    <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold sm:pl-1">
                        {{ $t('common.action', 2) }}
                    </th>
                    <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold sm:pl-1">
                        {{ $t('common.crime') }}
                    </th>
                    <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold">
                        {{ $t('common.fine') }}
                    </th>
                    <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold">
                        {{ $t('common.detention_time') }}
                    </th>
                    <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold">
                        {{ $t('common.traffic_infraction_points', 2) }}
                    </th>
                    <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold">
                        {{ $t('common.description') }}
                    </th>
                </tr>
            </thead>
            <tbody class="divide-y divide-base-800">
                <LawEntry
                    v-for="law in modelValue.laws"
                    :key="law.id"
                    :law="law"
                    :start-in-edit="parseInt(law.id) < 0"
                    @update:law="$emit('update:law', $event)"
                    @deleted="deletedLaw($event)"
                />
            </tbody>
        </table>
    </UCard>
</template>
