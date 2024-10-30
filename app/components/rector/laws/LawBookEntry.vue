<script lang="ts" setup>
import type Table from '#ui/components/data/Table.vue';
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import LawEntry from '~/components/rector/laws/LawEntry.vue';
import type { Law, LawBook } from '~~/gen/ts/resources/laws/laws';

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

const { t } = useI18n();

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

async function deleteLaw(id: string): Promise<void> {
    const i = parseInt(id);
    if (i < 0) {
        deletedLaw(id);
        return;
    }

    try {
        const call = getGRPCRectorLawsClient().deleteLaw({
            id,
        });
        await call;

        deletedLaw(id);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const columns = [
    {
        key: 'actions',
        label: '',
        sortable: false,
    },
    {
        key: 'crime',
        label: t('common.crime'),
    },
    {
        key: 'fine',
        label: t('common.fine'),
    },
    {
        key: 'detentionTime',
        label: t('common.detention_time'),
    },
    {
        key: 'service',
        label: t('common.traffic_infraction_points', 2),
    },
    {
        key: 'description',
        label: t('common.description'),
    },
];

const table = useTemplateRef<typeof Table>('table');

watch(table, () => {
    console.log('table', table.value);
    table.value?.toggleOpened(0);
});

const editing = ref(props.startInEdit);
</script>

<template>
    <UCard v-if="lawBook" class="overflow-y-auto">
        <template #header>
            <div v-if="!editing" class="flex items-center gap-x-2">
                <UButtonGroup class="inline-flex">
                    <UButton variant="link" icon="i-mdi-pencil" :title="$t('common.edit')" @click="editing = true" />

                    <UButton
                        variant="link"
                        icon="i-mdi-trash-can"
                        color="red"
                        :title="$t('common.delete')"
                        @click="
                            modal.open(ConfirmModal, {
                                confirm: async () => deleteLawBook(lawBook!.id),
                            })
                        "
                    />
                </UButtonGroup>

                <div class="inline-flex w-full flex-col">
                    <h2 class="text-xl">{{ lawBook.name }}</h2>

                    <p v-if="lawBook.description">{{ $t('common.description') }}: {{ lawBook.description }}</p>
                </div>

                <UButton color="gray" trailing-icon="i-mdi-plus" @click="addLaw">
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

        <UTable
            ref="table"
            :columns="columns"
            :rows="laws"
            :expand-button="{ icon: 'i-mdi-pencil', color: 'primary' }"
            :ui="{ wrapper: '' }"
        >
            <template #expand="{ row: law, index }">
                <LawEntry
                    :law="law"
                    :start-in-edit="parseInt(law.id) < 0"
                    @update:law="
                        $emit('update:law', $event);
                        table?.toggleOpened(index);
                    "
                    @close="
                        table?.toggleOpened(index);
                        if (parseInt(law.id) < 0) {
                            deleteLaw(law.id);
                        }
                    "
                />
            </template>

            <template #actions-data="{ row: law }">
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
            </template>

            <template #crime-data="{ row: law }">
                <span class="truncate text-gray-900 dark:text-white">
                    {{ law.name }}
                </span>
            </template>

            <template #fine-data="{ row: law }">{{ $n(law.fine, 'currency') }}</template>

            <template #detentionTime-data="{ row: law }">
                {{ law.detentionTime }}
            </template>

            <template #stvoPoints-data="{ row: law }">
                {{ law.stvoPoints }}
            </template>

            <template #description-data="{ row: law }">
                <span class="line-clamp-2 truncate hover:line-clamp-4">
                    {{ law.description }}
                </span>

                <span class="line-clamp-2 truncate hover:line-clamp-4">
                    <span class="font-semibold">{{ $t('common.hint') }}:</span> {{ law.hint }}
                </span>
            </template>
        </UTable>
    </UCard>
</template>
