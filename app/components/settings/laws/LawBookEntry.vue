<script lang="ts" setup>
import { UButton, UTooltip } from '#components';
import type { FormSubmitEvent, TableColumn } from '@nuxt/ui';
import type { ExpandedState } from '@tanstack/vue-table';
import { z } from 'zod';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import LawEntry from '~/components/settings/laws/LawEntry.vue';
import { getSettingsLawsClient } from '~~/gen/ts/clients';
import type { Law, LawBook } from '~~/gen/ts/resources/laws/laws';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

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

const { display } = useAppConfig();

const { can } = useAuth();

const notifications = useNotificationsStore();

const lawBook = useVModel(props, 'modelValue', emit);

const laws = useVModel(props, 'laws', emit);

const overlay = useOverlay();

const formatter = new Intl.NumberFormat(display.intlLocale, {
    style: 'currency',
    currency: display.currencyName,
});

const settingsLawsClient = await getSettingsLawsClient();

const schema = z.object({
    name: z.coerce.string().min(3).max(128),
    description: z.union([z.coerce.string().min(3).max(255), z.coerce.string().length(0).optional()]),
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

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        return response.lawBook!;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    if (!lawBook.value) return;

    canSubmit.value = false;
    await saveLawBook(lawBook.value.id, event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

function deletedLaw(id: number): void {
    laws.value = laws.value.filter((b) => b.id !== id);
}

const lastNewId = ref(-1);

const lawEntriesRefs = ref(new Map<number, Element>());

function addLaw(): void {
    if (!lawBook.value) return;

    const law = {
        lawbookId: lawBook.value.id,
        id: lastNewId.value,
        name: '',
        fine: 0,
        detentionTime: 0,
        stvoPoints: 0,
    };
    laws.value.push(law);

    if (typeof expanded.value === 'object') {
        expanded.value = { ...expanded.value, [laws.value.length - 1]: true };
    } else {
        expanded.value = { [laws.value.length - 1]: true };
    }

    useTimeoutFn(() => {
        const ref = lawEntriesRefs.value.get(law.id);
        if (ref) {
            ref.scrollIntoView({ block: 'nearest' });
        }
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

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const editing = ref(props.startInEdit);

const expanded = ref<ExpandedState>({});

const columns = computed(
    () =>
        [
            {
                id: 'expand',
                cell: ({ row }) =>
                    h(UTooltip, { text: t('common.expand_collapse') }, () =>
                        h(UButton, {
                            color: row.getIsExpanded() ? 'neutral' : 'primary',
                            variant: 'ghost',
                            icon: row.getIsExpanded() ? 'i-mdi-chevron-up' : 'i-mdi-pencil',
                            square: true,
                            'aria-label': 'Expand',
                            ui: {
                                leadingIcon: ['transition-transform', row.getIsExpanded() ? 'duration-200 rotate-180' : ''],
                            },
                            onClick: () => row.toggleExpanded(),
                        }),
                    ),
            },
            {
                id: 'actions',
            },
            {
                accessorKey: 'crime',
                header: t('common.crime'),
                cell: ({ row }) => row.original.name,
            },
            {
                accessorKey: 'fine',
                header: t('common.fine'),
                cell: ({ row }) => formatter.format(row.original.fine ?? 0),
            },
            {
                accessorKey: 'detentionTime',
                header: t('common.detention_time'),
                cell: ({ row }) => row.original.detentionTime,
            },
            {
                accessorKey: 'stvoPoints',
                header: t('common.traffic_infraction_points'),
                cell: ({ row }) => row.original.stvoPoints,
            },
            {
                accessorKey: 'description',
                header: t('common.description'),
                cell: ({ row }) =>
                    h('span', { class: 'line-clamp-2 truncate hover:line-clamp-4' }, [
                        row.original.description,
                        row.original.hint !== undefined && row.original.hint !== ''
                            ? h('span', { class: 'font-semibold' }, `${t('common.hint')}: ${row.original.hint}`)
                            : null,
                    ]),
            },
        ] as TableColumn<Law>[],
);

const confirmModal = overlay.create(ConfirmModal);
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
                                confirmModal.open({
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
                    <UInput v-model="state.name" name="name" type="text" class="w-full" :placeholder="$t('common.law_book')" />
                </UFormField>

                <UFormField class="flex-auto" name="description" :label="$t('common.description')">
                    <UInput
                        v-model="state.description"
                        name="description"
                        type="text"
                        class="w-full"
                        :placeholder="$t('common.description')"
                    />
                </UFormField>
            </UForm>
        </template>

        <UTable
            v-model:expanded="expanded"
            :columns="columns"
            :data="laws"
            :expand-button="{ icon: 'i-mdi-pencil', color: 'primary' }"
            :pagination-options="{ manualPagination: true }"
            :sorting-options="{ manualSorting: true }"
            :empty="$t('common.not_found', [$t('common.law', 2)])"
            :ui="{ tr: 'data-[expanded=true]:bg-elevated/50' }"
            sticky
        >
            <template #expanded="{ row }">
                <LawEntry :law="row.original" @update:law="$emit('update:law', $event)" @close="row.toggleExpanded()" />
            </template>

            <template #actions-cell="{ row: law }">
                <UTooltip v-if="can('settings.LawsService/DeleteLawBook').value" :text="$t('common.delete')">
                    <UButton
                        variant="link"
                        icon="i-mdi-delete"
                        color="error"
                        @click="
                            () => {
                                confirmModal.open({
                                    confirm: async () => deleteLaw(law.original.id),
                                });
                            }
                        "
                    />
                </UTooltip>
            </template>
        </UTable>
    </UCard>
</template>
