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

const emit = defineEmits<{
    (e: 'deleted', id: number): void;
    (e: 'update:modelValue', book?: LawBook): void;
    (e: 'update:laws', laws: Law[]): void;
    (e: 'update:law', update: { id: number; law: Law }): void;
}>();

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const { can } = useAuth();

const lawBook = useVModel(props, 'modelValue', emit);

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

async function deleteLawBook(id: number): Promise<void> {
    if (id < 0) {
        emit('deleted', id);
        return;
    }

    try {
        const call = $grpc.rector.rectorLaws.deleteLawBook({
            id: id,
        });
        await call;

        emit('deleted', id);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function saveLawBook(id: number, values: Schema): Promise<LawBook> {
    try {
        const call = $grpc.rector.rectorLaws.createOrUpdateLawBook({
            lawBook: {
                id: id < 0 ? 0 : id,
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

function deletedLaw(id: number): void {
    emit(
        'update:laws',
        props.laws.filter((b) => b.id !== id),
    );
}

const lastNewId = ref(-1);

function addLaw(): void {
    if (!lawBook.value) {
        return;
    }

    emit('update:laws', [
        ...props.laws,
        {
            lawbookId: lawBook.value.id,
            id: lastNewId.value,
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

async function deleteLaw(id: number): Promise<void> {
    if (id < 0) {
        deletedLaw(id);
        return;
    }

    try {
        const call = $grpc.rector.rectorLaws.deleteLaw({
            id: id,
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

const expand = ref({
    openedRows: [],
    row: {},
});

const editing = ref(props.startInEdit);
</script>

<template>
    <UCard v-if="lawBook" class="overflow-y-auto">
        <template #header>
            <div v-if="!editing" class="flex items-center gap-x-2">
                <UButtonGroup class="inline-flex">
                    <UTooltip :text="$t('common.edit')">
                        <UButton variant="link" icon="i-mdi-pencil" @click="editing = true" />
                    </UTooltip>

                    <UTooltip v-if="can('RectorLawsService.DeleteLawBook').value" :text="$t('common.delete')">
                        <UButton
                            variant="link"
                            icon="i-mdi-delete"
                            color="error"
                            @click="
                                modal.open(ConfirmModal, {
                                    confirm: async () => deleteLawBook(lawBook!.id),
                                })
                            "
                        />
                    </UTooltip>
                </UButtonGroup>

                <div class="inline-flex w-full flex-col">
                    <h2 class="text-xl">{{ lawBook.name }}</h2>

                    <p v-if="lawBook.description">{{ $t('common.description') }}: {{ lawBook.description }}</p>
                </div>

                <UTooltip :text="$t('pages.rector.laws.add_new_law')">
                    <UButton color="gray" trailing-icon="i-mdi-plus" @click="addLaw">
                        {{ $t('pages.rector.laws.add_new_law') }}
                    </UButton>
                </UTooltip>
            </div>
            <UForm
                v-else
                :schema="schema"
                :state="state"
                class="flex w-full flex-row items-start gap-x-2"
                @submit="onSubmitThrottle"
            >
                <UTooltip :text="$t('common.save')">
                    <UButton type="submit" variant="link" icon="i-mdi-content-save" />
                </UTooltip>

                <UTooltip :text="$t('common.cancel')">
                    <UButton
                        variant="link"
                        icon="i-mdi-cancel"
                        @click="
                            editing = false;
                            lawBook.id < 0 && $emit('deleted', lawBook.id);
                        "
                    />
                </UTooltip>

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
            v-model:expand="expand"
            :columns="columns"
            :rows="laws"
            :expand-button="{ icon: 'i-mdi-pencil', color: 'primary' }"
            :ui="{ wrapper: '' }"
            :empty-state="{
                icon: 'i-mdi-gavel',
                label: $t('common.not_found', [$t('common.law', 2)]),
            }"
        >
            <template #expand="{ row: law, index }">
                <LawEntry
                    :law="law"
                    :start-in-edit="law.id < 0"
                    @update:law="
                        $emit('update:law', $event);
                        table?.toggleOpened(index);
                    "
                    @close="
                        table?.toggleOpened(index);
                        if (law.id < 0) {
                            deleteLaw(law.id);
                        }
                    "
                />
            </template>

            <template #actions-data="{ row: law }">
                <UTooltip v-if="can('RectorLawsService.DeleteLawBook').value" :text="$t('common.delete')">
                    <UButton
                        variant="link"
                        icon="i-mdi-delete"
                        color="error"
                        @click="
                            modal.open(ConfirmModal, {
                                confirm: async () => deleteLaw(law.id),
                            })
                        "
                    />
                </UTooltip>
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

                <span v-if="law.hint !== undefined && law.hint !== ''" class="line-clamp-2 truncate hover:line-clamp-4">
                    <span class="font-semibold">{{ $t('common.hint') }}:</span> {{ law.hint }}
                </span>
            </template>
        </UTable>
    </UCard>
</template>
