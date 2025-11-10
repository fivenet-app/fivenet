<script setup lang="ts">
import { UButton, UFieldGroup, UTooltip } from '#components';
import type { TableColumn, TabsItem } from '@nuxt/ui';
import { computed, h } from 'vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import { useAuthStore } from '~/stores/auth';
import { getUser, useClipboardStore } from '~/stores/clipboard';
import { getCitizensCitizensClient } from '~~/gen/ts/clients';
import { type DocumentRelation, DocRelation } from '~~/gen/ts/resources/documents/documents';
import type { User } from '~~/gen/ts/resources/users/users';

const props = defineProps<{
    documentId?: number;
}>();

const modelValue = defineModel<DocumentRelation[]>({
    required: true,
});

const { t } = useI18n();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const clipboardStore = useClipboardStore();

const citizensCitizensClient = await getCitizensCitizensClient();

const items = ref<TabsItem[]>([
    {
        label: t('components.documents.document_managers.view_current'),
        icon: 'i-mdi-view-list-outline',
        slot: 'current' as const,
        value: 'current',
    },
    {
        label: t('common.clipboard'),
        icon: 'i-mdi-clipboard-list',
        slot: 'clipboard' as const,
        value: 'clipboard',
    },
    {
        label: t('components.documents.document_managers.add_new'),
        icon: 'i-mdi-account-search',
        slot: 'new' as const,
        value: 'new',
    },
]);

const queryCitizens = ref('');

const {
    data: citizens,
    status,
    refresh,
    error,
} = useLazyAsyncData(`document-${props.documentId?.toString()}-relations-citzens-${queryCitizens.value}`, () => listCitizens());

watchDebounced(queryCitizens, async () => await refresh(), {
    debounce: 200,
    maxWait: 1750,
});

async function listCitizens(): Promise<User[]> {
    try {
        const call = citizensCitizensClient.listCitizens({
            pagination: {
                offset: 0,
                pageSize: 8,
            },
            search: queryCitizens.value,
        });
        const { response } = await call;

        return response.users.filter((user) => !modelValue.value.find((r) => r.targetUserId === user.userId));
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

let lastId = 0;

async function addRelation(user: User, relation: DocRelation): Promise<void> {
    modelValue.value.push({
        id: lastId--,
        documentId: props.documentId ?? 0,
        sourceUserId: activeChar.value!.userId,
        sourceUser: activeChar.value!,
        targetUserId: user.userId,
        targetUser: user,
        relation: relation,
    });

    await refresh();
}

async function removeRelation(id: number): Promise<void> {
    const idx = modelValue.value.findIndex((r) => r.id === id);
    if (idx > -1) {
        modelValue.value.splice(idx, 1);
    }
    refresh();
}

const columnsCurrent = computed(
    () =>
        [
            {
                accessorKey: 'name',
                header: t('common.name'),
                cell: ({ row }) =>
                    h(CitizenInfoPopover, {
                        user: !row.original.targetUser?.userId ? undefined : row.original.targetUser,
                        userId: row.original.targetUserId,
                        showBirthdate: true,
                    }),
            },
            {
                accessorKey: 'creator',
                header: t('common.creator'),
                cell: ({ row }) =>
                    h(CitizenInfoPopover, {
                        user: !row.original.sourceUser?.userId ? undefined : row.original.sourceUser,
                        userId: row.original.sourceUserId,
                        trailing: false,
                    }),
            },
            {
                accessorKey: 'relation',
                header: t('common.relation', 1),
                cell: ({ row }) => t(`enums.documents.DocRelation.${DocRelation[row.original.relation]}`),
            },
            {
                id: 'actions',
                cell: ({ row }) =>
                    h(
                        UFieldGroup,
                        {},
                        {
                            default: () => [
                                h(
                                    UTooltip,
                                    { text: t('components.documents.document_managers.open_citizen') },
                                    {
                                        default: () =>
                                            h(UButton, {
                                                to: {
                                                    name: 'citizens-id',
                                                    params: { id: row.original.targetUserId },
                                                },
                                                target: '_blank',
                                                variant: 'link',
                                                icon: 'i-mdi-open-in-new',
                                            }),
                                    },
                                ),
                                h(
                                    UTooltip,
                                    { text: t('components.documents.document_managers.remove_relation') },
                                    {
                                        default: () =>
                                            h(UButton, {
                                                variant: 'link',
                                                icon: 'i-mdi-account-minus',
                                                color: 'error',
                                                onClick: () => removeRelation(row.original.id!),
                                            }),
                                    },
                                ),
                            ],
                        },
                    ),
            },
        ] as TableColumn<DocumentRelation>[],
);

const columnsClipboard = computed(
    () =>
        [
            {
                accessorKey: 'name',
                header: t('common.name'),
                cell: ({ row }) =>
                    h(CitizenInfoPopover, {
                        user: row.original,
                        showBirthdate: true,
                    }),
            },
            {
                accessorKey: 'job',
                header: t('common.job'),
                cell: ({ row }) => row.original.jobLabel,
            },
            {
                accessorKey: 'relations',
                header: t('components.documents.document_managers.add_relation'),
                cell: ({ row }) =>
                    h(
                        UFieldGroup,
                        {},
                        {
                            default: () => [
                                h(
                                    UTooltip,
                                    { text: t('components.documents.document_managers.mentioned') },
                                    {
                                        default: () =>
                                            h(UButton, {
                                                color: 'blue',
                                                icon: 'i-mdi-at',
                                                onClick: () => addRelation(getUser(row.original), DocRelation.MENTIONED),
                                            }),
                                    },
                                ),
                                h(
                                    UTooltip,
                                    { text: t('components.documents.document_managers.targets') },
                                    {
                                        default: () =>
                                            h(UButton, {
                                                color: 'warning',
                                                icon: 'i-mdi-target',
                                                onClick: () => addRelation(getUser(row.original), DocRelation.TARGETS),
                                            }),
                                    },
                                ),
                                h(
                                    UTooltip,
                                    { text: t('components.documents.document_managers.caused') },
                                    {
                                        default: () =>
                                            h(UButton, {
                                                color: 'error',
                                                icon: 'i-mdi-source-commit-start',
                                                onClick: () => addRelation(getUser(row.original), DocRelation.CAUSED),
                                            }),
                                    },
                                ),
                            ],
                        },
                    ),
            },
        ] as TableColumn<ClipboardUser>[],
);

const columnsNew = computed(
    () =>
        [
            {
                accessorKey: 'name',
                header: t('common.name'),
                cell: ({ row }) =>
                    h(CitizenInfoPopover, {
                        user: row.original,
                        showBirthdate: true,
                    }),
            },
            {
                accessorKey: 'job',
                header: t('common.job'),
                cell: ({ row }) => row.original.jobLabel,
            },
            {
                accessorKey: 'relations',
                header: t('components.documents.document_managers.add_relation'),
                cell: ({ row }) =>
                    h(
                        UFieldGroup,
                        {},
                        {
                            default: () => [
                                h(
                                    UTooltip,
                                    { text: t('components.documents.document_managers.mentioned') },
                                    {
                                        default: () =>
                                            h(UButton, {
                                                color: 'blue',
                                                icon: 'i-mdi-at',
                                                onClick: () => addRelation(row.original, DocRelation.MENTIONED),
                                            }),
                                    },
                                ),
                                h(
                                    UTooltip,
                                    { text: t('components.documents.document_managers.targets') },
                                    {
                                        default: () =>
                                            h(UButton, {
                                                color: 'warning',
                                                icon: 'i-mdi-target',
                                                onClick: () => addRelation(row.original, DocRelation.TARGETS),
                                            }),
                                    },
                                ),
                                h(
                                    UTooltip,
                                    { text: t('components.documents.document_managers.caused') },
                                    {
                                        default: () =>
                                            h(UButton, {
                                                color: 'error',
                                                icon: 'i-mdi-source-commit-start',
                                                onClick: () => addRelation(row.original, DocRelation.CAUSED),
                                            }),
                                    },
                                ),
                            ],
                        },
                    ),
            },
        ] as TableColumn<User>[],
);
</script>

<template>
    <UTabs default-value="current" :items="items" variant="link">
        <template #current>
            <UTable :columns="columnsCurrent" :data="modelValue" :empty="$t('common.not_found', [$t('common.relation', 2)])" />
        </template>

        <template #clipboard>
            <UTable
                :columns="columnsClipboard"
                :data="clipboardStore.$state.users"
                :empty="$t('common.not_found', [$t('common.citizen', 2)])"
            />
        </template>

        <template #new>
            <UFormField class="mb-2" name="name" :label="$t('common.search')">
                <UInput
                    v-model="queryCitizens"
                    type="text"
                    name="name"
                    :placeholder="`${$t('common.citizen', 1)} ${$t('common.name')}`"
                    leading-icon="i-mdi-search"
                />
            </UFormField>

            <div>
                <DataErrorBlock
                    v-if="error"
                    :title="$t('common.unable_to_load', [$t('common.citizen', 2)])"
                    :error="error"
                    :retry="refresh"
                />

                <UTable
                    v-else
                    :columns="columnsNew"
                    :loading="isRequestPending(status)"
                    :data="citizens"
                    :empty="$t('common.not_found', [$t('common.citizen', 2)])"
                />
            </div>
        </template>
    </UTabs>
</template>
