<script lang="ts" setup>
import { z } from 'zod';
import type { FormSubmitEvent } from '#ui/types';
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

async function deleteLawBook(id: string): Promise<void> {
    const i = parseInt(id);
    if (i < 0) {
        emit('deleted', id);
        return;
    }

    try {
        const call = getGRPCRectorLawsClient().deleteLawBook({ id });
        await call;

        emit('deleted', id);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const schema = z.object({
    name: z.string().min(3).max(128),
    description: z.union([z.string().min(3).max(255), z.string().length(0).optional()]),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    name: props.modelValue.name,
    description: props.modelValue.description,
});

async function saveLawBook(id: string, values: Schema): Promise<LawBook> {
    const i = parseInt(id);

    try {
        const call = getGRPCRectorLawsClient().createOrUpdateLawBook({
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
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await saveLawBook(props.modelValue.id, event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
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
            <UForm
                v-else
                :schema="schema"
                :state="state"
                class="flex w-full flex-row items-start gap-x-2"
                @submit="onSubmitThrottle"
            >
                <UButton :title="$t('common.save')" variant="link" icon="i-mdi-content-save" />
                <UButton
                    :title="$t('common.cancel')"
                    variant="link"
                    icon="i-mdi-cancel"
                    @click="
                        editing = false;
                        parseInt(modelValue.id) < 0 && $emit('deleted', modelValue.id);
                    "
                />

                <UFormGroup name="name" :label="$t('common.law_book')" class="flex-initial">
                    <UInput
                        v-model="state.name"
                        name="name"
                        type="text"
                        :placeholder="$t('common.law_book')"
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
                    />
                </UFormGroup>

                <UFormGroup name="description" :label="$t('common.description')" class="flex-auto">
                    <UInput
                        v-model="state.description"
                        name="description"
                        type="text"
                        :placeholder="$t('common.description')"
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
                    />
                </UFormGroup>
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
