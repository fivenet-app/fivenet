<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import LawEntry from '~/components/rector/laws/LawEntry.vue';
import { Law, LawBook } from '~~/gen/ts/resources/laws/laws';

const props = defineProps<{
    modelValue: LawBook | undefined;
    laws: Law[];
    startInEdit?: boolean;
}>();

const emits = defineEmits<{
    (e: 'deleted', id: string): void;
    (e: 'update:modelValue', book?: LawBook): void;
    (e: 'update:laws', laws: Law[]): void;
    (e: 'update:law', update: { id: string; law: Law }): void;
}>();

const lawBook = useVModel(props, 'modelValue', emits);

const modal = useModal();

const schema = z.object({
    name: z.string().min(3).max(128),
    description: z.union([z.string().min(3).max(255), z.string().length(0).optional()]),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    name: '',
    description: '',
});

async function deleteLawBook(id: string): Promise<void> {
    const i = parseInt(id);
    if (i < 0) {
        emits('deleted', id);
        return;
    }

    try {
        const call = getGRPCRectorLawsClient().deleteLawBook({ id });
        await call;

        emits('deleted', id);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

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

        lawBook.value = response.lawBook;

        return response.lawBook!;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    if (!lawBook.value) {
        return;
    }

    canSubmit.value = false;
    await saveLawBook(lawBook.value.id, event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

function deletedLaw(id: string): void {
    emits(
        'update:laws',
        props.laws.filter((b) => b.id !== id),
    );
}

const lastNewId = ref(-1);

function addLaw(): void {
    if (!lawBook.value) {
        return;
    }

    emits('update:laws', [
        ...props.laws,
        {
            lawbookId: lawBook.value.id,
            id: lastNewId.value.toString(),
            name: '',
            fine: 0,
            detentionTime: 0,
            stvoPoints: 0,
        },
    ]);
    lastNewId.value--;
}

function resetForm(): void {
    state.name = lawBook.value?.name ?? '';
    state.description = lawBook.value?.description;
}

onMounted(() => resetForm());
watch(props, () => resetForm());

const editing = ref(props.startInEdit);
</script>

<template>
    <UCard v-if="lawBook">
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
                                confirm: async () => deleteLawBook(lawBook!.id),
                            })
                        "
                    />
                </UButtonGroup>

                <div class="inline-flex flex-col">
                    <h2 class="text-xl">{{ lawBook.name }}</h2>

                    <p v-if="lawBook.description">{{ $t('common.description') }}: {{ lawBook.description }}</p>
                </div>

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
                <UButton type="submit" :title="$t('common.save')" variant="link" icon="i-mdi-content-save" />
                <UButton
                    :title="$t('common.cancel')"
                    variant="link"
                    icon="i-mdi-cancel"
                    @click="
                        editing = false;
                        parseInt(lawBook.id) < 0 && $emit('deleted', lawBook.id);
                    "
                />

                <UFormGroup name="name" :label="$t('common.law_book')" class="flex-initial">
                    <UInput v-model="state.name" name="name" type="text" :placeholder="$t('common.law_book')" />
                </UFormGroup>

                <UFormGroup name="description" :label="$t('common.description')" class="flex-auto">
                    <UInput
                        v-model="state.description"
                        name="description"
                        type="text"
                        :placeholder="$t('common.description')"
                    />
                </UFormGroup>
            </UForm>
        </template>

        <div class="table min-w-full divide-y divide-base-600">
            <div class="table-header-group">
                <div class="table-row">
                    <div class="table-cell py-3.5 pl-4 pr-3 text-left text-sm font-semibold sm:pl-1">
                        {{ $t('common.action', 2) }}
                    </div>
                    <div class="table-cell py-3.5 pl-4 pr-3 text-left text-sm font-semibold sm:pl-1">
                        {{ $t('common.crime') }}
                    </div>
                    <div class="table-cell px-2 py-3.5 text-left text-sm font-semibold">
                        {{ $t('common.fine') }}
                    </div>
                    <div class="table-cell px-2 py-3.5 text-left text-sm font-semibold">
                        {{ $t('common.detention_time') }}
                    </div>
                    <div class="table-cell px-2 py-3.5 text-left text-sm font-semibold">
                        {{ $t('common.traffic_infraction_points', 2) }}
                    </div>
                    <div class="table-cell px-2 py-3.5 text-left text-sm font-semibold">
                        {{ $t('common.description') }}
                    </div>
                </div>
            </div>
            <div class="table-row-group divide-y divide-base-800">
                <LawEntry
                    v-for="law in laws"
                    :key="law.id"
                    :law="law"
                    :start-in-edit="parseInt(law.id) < 0"
                    @update:law="$emit('update:law', $event)"
                    @deleted="deletedLaw($event)"
                />
            </div>
        </div>
    </UCard>
</template>
