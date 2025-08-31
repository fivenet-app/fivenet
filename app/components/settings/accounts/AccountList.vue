<script lang="ts" setup>
import { UBadge, UButton, UTooltip } from '#components';
import type { TableColumn } from '@nuxt/ui';
import { h } from 'vue';
import { z } from 'zod';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import Pagination from '~/components/partials/Pagination.vue';
import StreamerModeAlert from '~/components/partials/StreamerModeAlert.vue';
import { getSettingsAccountsClient } from '~~/gen/ts/clients';
import type { Account } from '~~/gen/ts/resources/accounts/accounts';
import type { SortByColumn } from '~~/gen/ts/resources/common/database/database';
import type { ListAccountsResponse } from '~~/gen/ts/services/settings/accounts';
import AccountEditModal from './AccountEditModal.vue';

const { t } = useI18n();

const overlay = useOverlay();

const settingsStore = useSettingsStore();
const { streamerMode } = storeToRefs(settingsStore);

const settingsAccountsClient = await getSettingsAccountsClient();

const schema = z.object({
    license: z.string().max(64).optional(),
    enabled: z.coerce.boolean().default(true),
    username: z.string().max(64).optional(),
    externalId: z.string().max(64).optional(),

    sorting: z
        .object({
            columns: z.custom<SortByColumn>().array().max(3).default([]),
        })
        .default({
            columns: [
                {
                    id: 'username',
                    desc: false,
                },
            ],
        }),
    page: pageNumberSchema,
});

const query = useSearchForm('settings_accounts', schema);

const {
    data: accounts,
    status,
    refresh,
    error,
} = useLazyAsyncData(
    () =>
        `settings-accounts-${query.license}-${query.enabled}-${query.username}-${query.externalId}-${JSON.stringify(query.sorting)}-${query.page}`,
    () => listAccounts(),
);

async function listAccounts(): Promise<ListAccountsResponse> {
    try {
        const call = settingsAccountsClient.listAccounts({
            pagination: {
                offset: calculateOffset(query.page, accounts.value?.pagination),
            },
            sort: query.sorting,
            enabled: query.enabled,
            license: query.license,
            username: query.username,
            externalId: query.externalId,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watchDebounced(query, async () => refresh(), { debounce: 200, maxWait: 1250 });

async function deleteAccount(id: number): Promise<void> {
    try {
        const call = settingsAccountsClient.deleteAccount({
            id,
        });
        await call;

        const idx = accounts.value?.accounts.findIndex((f) => f.id === id);
        if (idx !== undefined && idx > -1 && accounts.value !== null) {
            accounts.value?.accounts.splice(idx, 1);
        }
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const confirmModal = overlay.create(ConfirmModal);
const accountEditModal = overlay.create(AccountEditModal);

const appConfig = useAppConfig();

const columns = computed(
    () =>
        [
            {
                id: 'actions',
                cell: ({ row }) =>
                    h('div', [
                        h(UTooltip, { text: t('common.update') }, [
                            h(UButton, {
                                variant: 'link',
                                icon: 'i-mdi-pencil',
                                onClick: () => {
                                    accountEditModal.open({
                                        account: row.original,
                                        'onUpdate:account': () => refresh(),
                                    });
                                },
                            }),
                        ]),
                        h(UTooltip, { text: t('common.delete') }, [
                            h(UButton, {
                                variant: 'link',
                                icon: 'i-mdi-delete',
                                color: 'error',
                                onClick: () => {
                                    confirmModal.open({
                                        confirm: async () => deleteAccount(row.original.id),
                                    });
                                },
                            }),
                        ]),
                    ]),
            },
            {
                accessorKey: 'username',
                header: ({ column }) => {
                    const isSorted = column.getIsSorted();

                    return h(UButton, {
                        color: 'neutral',
                        variant: 'ghost',
                        label: t('common.username'),
                        icon: isSorted
                            ? isSorted === 'asc'
                                ? appConfig.custom.icons.sortAsc
                                : appConfig.custom.icons.sortDesc
                            : appConfig.custom.icons.sort,
                        class: '-mx-2.5',
                        onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
                    });
                },
                sortable: true,
                cell: ({ row }) =>
                    h(UTooltip, { text: `${t('common.id')}: ${row.original.id}` }, [
                        h(
                            'span',
                            { class: 'text-highlighted' },
                            row.original.username === '' ? t('common.na') : row.original.username,
                        ),
                    ]),
            },
            {
                accessorKey: 'enabled',
                header: t('common.enabled'),
                cell: ({ row }) =>
                    h(UBadge, {
                        color: row.original.enabled ? 'success' : 'error',
                        label: row.original.enabled ? t('common.yes') : t('common.no'),
                    }),
            },
            {
                accessorKey: 'createdAt',
                header: t('common.created_at'),
                cell: ({ row }) => h(GenericTime, { value: toDate(row.original.createdAt) }),
            },
            {
                accessorKey: 'updatedAt',
                header: t('common.updated_at'),
                cell: ({ row }) => h(GenericTime, { value: toDate(row.original.updatedAt) }),
            },
            {
                accessorKey: 'license',
                header: ({ column }) => {
                    const isSorted = column.getIsSorted();

                    return h(UButton, {
                        color: 'neutral',
                        variant: 'ghost',
                        label: t('common.license'),
                        icon: isSorted
                            ? isSorted === 'asc'
                                ? appConfig.custom.icons.sortAsc
                                : appConfig.custom.icons.sortDesc
                            : appConfig.custom.icons.sort,
                        class: '-mx-2.5',
                        onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
                    });
                },
                sortable: true,
                cell: ({ row }) => h('pre', { class: 'text-highlighted' }, row.original.license),
            },
        ] as TableColumn<Account>[],
);
</script>

<template>
    <UDashboardPanel>
        <template #header>
            <UDashboardNavbar :title="$t('pages.settings.accounts.title')">
                <template #right>
                    <PartialsBackButton fallback-to="/settings" />
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar>
                <UForm v-if="!streamerMode" class="w-full" :schema="schema" :state="query" @submit="refresh()">
                    <div class="flex w-full flex-row gap-2">
                        <UFormField class="flex-1" :label="$t('common.search')" name="license">
                            <UInput
                                ref="input"
                                v-model="query.license"
                                type="text"
                                name="license"
                                :placeholder="$t('common.license')"
                                block
                                leading-icon="i-mdi-search"
                            >
                                <template #trailing>
                                    <UKbd value="/" />
                                </template>
                            </UInput>
                        </UFormField>

                        <UFormField
                            class="flex flex-initial flex-col"
                            name="enabled"
                            :label="$t('common.enabled')"
                            :ui="{ container: 'flex-1 flex' }"
                        >
                            <div class="flex flex-1 items-center">
                                <USwitch v-model="query.enabled" />
                            </div>
                        </UFormField>
                    </div>

                    <UAccordion
                        class="mt-2"
                        color="neutral"
                        variant="soft"
                        size="sm"
                        :items="[{ label: $t('common.advanced_search'), slot: 'search' as const }]"
                    >
                        <template #search>
                            <div class="flex flex-row flex-wrap gap-1">
                                <UFormField class="flex-1" name="username" :label="$t('common.username')">
                                    <UInput
                                        v-model="query.username"
                                        type="text"
                                        name="username"
                                        :placeholder="$t('common.username')"
                                        block
                                    />
                                </UFormField>

                                <UFormField
                                    class="flex-1"
                                    name="externalId"
                                    :label="$t('components.auth.OAuth2Connections.external_id')"
                                >
                                    <UInput
                                        v-model="query.externalId"
                                        type="text"
                                        name="externalId"
                                        :placeholder="$t('components.auth.OAuth2Connections.external_id')"
                                        block
                                    />
                                </UFormField>
                            </div>
                        </template>
                    </UAccordion>
                </UForm>
            </UDashboardToolbar>
        </template>

        <template #body>
            <StreamerModeAlert v-if="streamerMode" />
            <template v-else>
                <DataErrorBlock
                    v-if="error"
                    :title="$t('common.unable_to_load', [$t('common.account', 2)])"
                    :error="error"
                    :retry="refresh"
                />

                <UTable
                    v-model:sorting="query.sorting.columns"
                    class="flex-1"
                    :loading="isRequestPending(status)"
                    :columns="columns"
                    :data="accounts?.accounts"
                    :empty="$t('common.not_found', [$t('common.account', 2)])"
                    :sorting-options="{ manualSorting: true }"
                    :pagination-options="{ manualPagination: true }"
                    sticky
                />

                <Pagination v-model="query.page" :pagination="accounts?.pagination" :status="status" :refresh="refresh" />
            </template>
        </template>
    </UDashboardPanel>
</template>
