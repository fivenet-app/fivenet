<script lang="ts" setup>
import type Table from '#ui/components/data/Table.vue';
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import LawEntry from '~/components/settings/laws/LawEntry.vue';
import { getSettingsLawsClient } from '~~/gen/ts/clients';
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

const { t } = useI18n();

const { can } = useAuth();

const lawBook = useVModel(props, 'modelValue', emit);

const laws = useVModel(props, 'laws', emit);

const modal = useOverlay();

const settingsLawsClient = await getSettingsLawsClient();

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
        const call = settingsLawsClient.deleteLawBook({
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
        const call = settingsLawsClient.createOrUpdateLawBook({
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

const tableRef = useTemplateRef<typeof Table>('tableRef');

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    if (!lawBook.value) {
        return;
    }

    canSubmit.value = false;
    await saveLawBook(lawBook.value.id, event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

function deletedLaw(id: number): void {
    laws.value = laws.value.filter((b) => b.id !== id);
}

const lastNewId = ref(-1);

const lawEntriesRefs = ref(new Map<number, Element>());

function addLaw(): void {
    if (!lawBook.value) {
        return;
    }

    const law = {
        lawbookId: lawBook.value.id,
        id: lastNewId.value,
        name: '',
        fine: 0,
        detentionTime: 0,
        stvoPoints: 0,
    };
    laws.value.push(law);

    useTimeoutFn(() => {
        const ref = lawEntriesRefs.value.get(law.id);
        if (ref) {
            ref.scrollIntoView({ block: 'nearest' });
        }

        tableRef.value?.toggleOpened(law);
    }, 100);

    lastNewId.value--;
}

function resetForm(): void {
    state.name = lawBook.value?.name ?? '';
    state.description = lawBook.value?.description;
}

onBeforeMount(() => resetForm());
watch(props, () => resetForm());

async function deleteLaw(id: number): Promise<void> {
    if (id < 0) {
        deletedLaw(id);
        return;
    }

    try {
        const call = settingsLawsClient.deleteLaw({
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
        accessorKey: 'actions',
        label: '',
        sortable: false,
    },
    {
        accessorKey: 'crime',
        label: t('common.crime'),
    },
    {
        accessorKey: 'fine',
        label: t('common.fine'),
    },
    {
        accessorKey: 'detentionTime',
        label: t('common.detention_time'),
    },
    {
        accessorKey: 'service',
        label: t('common.traffic_infraction_points', 2),
    },
    {
        accessorKey: 'description',
        label: t('common.description'),
    },
];

const expand = ref({
    openedRows: [],
    row: {},
});

const editing = ref(props.startInEdit);
</script>

<template>
    <UCard v-if="lawBook" class="overflow-y-auto">
        <template #header>
            <div v-if="!editing" class="inline-flex w-full items-center gap-x-2">
                <UButtonGroup class="inline-flex">
                    <UTooltip :text="$t('common.edit')">
                        <UButton variant="link" icon="i-mdi-pencil" @click="editing = true" />
                    </UTooltip>

                    <UTooltip v-if="can('settings.LawsService/DeleteLawBook').value" :text="$t('common.delete')">
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

                <UTooltip class="shrink-0" :text="$t('pages.settings.laws.add_new_law')">
                    <UButton color="neutral" trailing-icon="i-mdi-plus" @click="addLaw">
                        {{ $t('pages.settings.laws.add_new_law') }}
                    </UButton>
                </UTooltip>
            </div>
            <UForm
                v-else
                class="flex w-full flex-row items-start gap-x-2"
                :schema="schema"
                :state="state"
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

                <UFormField class="flex-initial" name="name" :label="$t('common.law_book')">
                    <UInput v-model="state.name" name="name" type="text" :placeholder="$t('common.law_book')" />
                </UFormField>

                <UFormField class="flex-auto" name="description" :label="$t('common.description')">
                    <UInput
                        v-model="state.description"
                        name="description"
                        type="text"
                        :placeholder="$t('common.description')"
                    />
                </UFormField>
            </UForm>
        </template>

        <UTable
            ref="tableRef"
            v-model:expand="expand"
            :columns="columns"
            :data="laws"
            :expand-button="{ icon: 'i-mdi-pencil', color: 'primary' }"
            :ui="{ wrapper: '' }"
            :empty-state="{
                icon: 'i-mdi-gavel',
                label: $t('common.not_found', [$t('common.law', 2)]),
            }"
        >
            <template #expand="{ row: law }">
                <LawEntry
                    :law="law"
                    @update:law="
                        $emit('update:law', $event);
                        tableRef?.toggleOpened(law);
                    "
                    @close="
                        tableRef?.toggleOpened(law);
                        if (law.id < 0) {
                            deleteLaw(law.id);
                        }
                    "
                />
            </template>

            <template #actions-cell="{ row: law }">
                <UTooltip v-if="can('settings.LawsService/DeleteLawBook').value" :text="$t('common.delete')">
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

            <template #crime-cell="{ row: law }">
                <span :ref="(ref) => lawEntriesRefs.set(law.id, ref as Element)" class="truncate text-highlighted">
                    {{ law.name }}
                </span>
            </template>

            <template #fine-cell="{ row: law }">{{ $n(law.fine, 'currency') }}</template>

            <template #detentionTime-cell="{ row: law }">
                {{ law.detentionTime }}
            </template>

            <template #stvoPoints-cell="{ row: law }">
                {{ law.stvoPoints }}
            </template>

            <template #description-cell="{ row: law }">
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
