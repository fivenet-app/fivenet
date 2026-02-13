<script lang="ts" setup>
import { UBadge, UButton, UTooltip } from '#components';
import type { TableColumn } from '@nuxt/ui';
import { h } from 'vue';
import { z } from 'zod';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
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

const { accountId } = useAuth();

const settingsStore = useSettingsStore();
const { streamerMode } = storeToRefs(settingsStore);

const settingsAccountsClient = await getSettingsAccountsClient();

const schema = z.object({
    license: z.coerce.string().max(64).optional(),
    onlyDisabled: z.coerce.boolean().default(false),
    username: z.coerce.string().max(64).optional(),
    externalId: z.coerce.string().max(64).optional(),
    group: z.coerce.string().max(64).optional(),

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
        `settings-accounts-${query.license}-${query.onlyDisabled}-${query.username}-${query.externalId}-${JSON.stringify(query.sorting)}-${query.page}`,
    () => listAccounts(),
);

async function listAccounts(): Promise<ListAccountsResponse> {
    try {
        const call = settingsAccountsClient.listAccounts({
            pagination: {
                offset: calculateOffset(query.page, accounts.value?.pagination),
            },
            sort: query.sorting,
            onlyDisabled: query.onlyDisabled,
            license: query.license,
            username: query.username,
            externalId: query.externalId,
            group: query.group,
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
                        row.original.id !== accountId.value
                            ? h(UTooltip, { text: t('common.delete') }, [
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
                              ])
                            : undefined,
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
                cell: ({ row }) =>
                    h(UTooltip, { text: `${t('common.id', 1)}: ${row.original.id}` }, [
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
                accessorKey: 'group',
                header: t('common.group', 2),
                cell: ({ row }) =>
                    h('pre', { class: 'text-highlighted' }, row.original.groups?.groups.join(', ') || t('common.na')),
            },
            {
                accessorKey: 'lastChar',
                header: t('common.last_char'),
                cell: ({ row }) =>
                    row.original.lastChar
                        ? h(
                              CitizenInfoPopover,
                              {
                                  userId: row.original.lastChar,
                              },
                              row.original.lastChar,
                          )
                        : undefined,
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
                cell: ({ row }) => h('pre', { class: 'text-highlighted' }, row.original.license),
            },
        ] as TableColumn<Account>[],
);
</script>

<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardNavbar :title="$t('pages.settings.accounts.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <PartialsBackButton fallback-to="/settings" />
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar>
                <UForm
                    v-if="!streamerMode"
                    class="my-2 flex w-full flex-1 flex-col gap-2"
                    :schema="schema"
                    :state="query"
                    @submit="refresh()"
                >
                    <div class="flex flex-1 flex-row gap-2">
                        <UFormField class="flex-1" :label="$t('common.search')" name="license">
                            <UInput
                                ref="input"
                                v-model="query.license"
                                type="text"
                                name="license"
                                :placeholder="$t('common.license')"
                                block
                                leading-icon="i-mdi-search"
                                class="w-full"
                            >
                                <template #trailing>
                                    <UKbd value="/" />
                                </template>
                            </UInput>
                        </UFormField>

                        <UFormField
                            class="flex flex-initial flex-col"
                            name="onlyDisabled"
                            :label="$t('common.only_disabled')"
                            :ui="{ container: 'flex-1 flex' }"
                        >
                            <div class="flex flex-1 items-center">
                                <USwitch v-model="query.onlyDisabled" />
                            </div>
                        </UFormField>
                    </div>

                    <UCollapsible>
                        <UButton
                            class="group"
                            color="neutral"
                            variant="ghost"
                            trailing-icon="i-mdi-chevron-down"
                            :label="$t('common.advanced_search')"
                            :ui="{
                                trailingIcon: 'group-data-[state=open]:rotate-180 transition-transform duration-200',
                            }"
                            block
                        />

                        <template #content>
                            <div class="flex flex-row flex-wrap gap-1">
                                <UFormField class="flex-1" name="username" :label="$t('common.username')">
                                    <UInput
                                        v-model="query.username"
                                        type="text"
                                        name="username"
                                        :placeholder="$t('common.username')"
                                        block
                                        class="w-full"
                                    />
                                </UFormField>

                                <UFormField
                                    class="flex-1"
                                    name="externalId"
                                    :label="$t('components.auth.SocialLogins.external_id')"
                                >
                                    <UInput
                                        v-model="query.externalId"
                                        type="text"
                                        name="externalId"
                                        :placeholder="$t('components.auth.SocialLogins.external_id')"
                                        block
                                        class="w-full"
                                    />
                                </UFormField>
                            </div>
                        </template>
                    </UCollapsible>
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
                    v-else
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
