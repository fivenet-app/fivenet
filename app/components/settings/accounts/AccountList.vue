<script lang="ts" setup>
import type { Row, TableMeta } from '@tanstack/vue-table';
import { UBadge, UButton, UTooltip } from '#components';
import type { Form, TableColumn } from '@nuxt/ui';
import { h } from 'vue';
import { z } from 'zod';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import Pagination from '~/components/partials/Pagination.vue';
import StreamerModeAlert from '~/components/partials/StreamerModeAlert.vue';
import TableSortButton from '~/components/partials/TableSortButton.vue';
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
    license: z.string().max(64).optional(),
    onlyDisabled: z.coerce.boolean().default(false),
    username: z.string().max(64).optional(),
    externalId: z.string().max(64).optional(),
    group: z.string().max(64).optional(),

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

type Schema = z.output<typeof schema>;

const query = useSearchForm('settings_accounts', schema);

const formRef = useTemplateRef<Form<typeof schema>>('formRef');
const { validatedQuery, commitValidatedQuery } = useFormSearchValidation<typeof schema>(query, formRef);

const accountsKey = computed(
    () =>
        `settings-accounts-${validatedQuery.value.license}-${validatedQuery.value.onlyDisabled}-${validatedQuery.value.username}-${validatedQuery.value.externalId}-${JSON.stringify(validatedQuery.value.sorting)}-${validatedQuery.value.page}`,
);

const { data: accounts, status, refresh, error } = useLazyAsyncData(accountsKey, () => listAccounts(validatedQuery.value));

async function listAccounts(values: Schema): Promise<ListAccountsResponse> {
    try {
        const call = settingsAccountsClient.listAccounts({
            pagination: {
                offset: calculateOffset(values.page, accounts.value?.pagination),
            },
            sort: values.sorting,
            onlyDisabled: values.onlyDisabled,
            license: values.license,
            username: values.username,
            externalId: values.externalId,
            group: values.group,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function deleteAccount(id: number): Promise<void> {
    try {
        const call = settingsAccountsClient.deleteAccount({
            id,
        });
        const { response } = await call;

        const idx = accounts.value?.accounts.findIndex((f) => f.id === id);
        if (idx !== undefined && idx > -1 && accounts.value && accounts.value.accounts[idx]) {
            accounts.value.accounts[idx]!.deletedAt = response.deletedAt;
        }
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const confirmModal = overlay.create(ConfirmModal);
const accountEditModal = overlay.create(AccountEditModal);

const meta = computed(
    () =>
        ({
            class: {
                tr: (row: Row<Account>) => {
                    return row.original.deletedAt
                        ? 'bg-warning-100/10 hover:bg-warning-200/10 dark:bg-warning-800/10 dark:hover:bg-warning-700/10'
                        : '';
                },
            },
        }) as TableMeta<Account>,
);

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
                            ? h(UTooltip, { text: !row.original.deletedAt ? t('common.delete') : t('common.restore') }, [
                                  h(UButton, {
                                      color: !row.original.deletedAt ? 'error' : 'success',
                                      icon: !row.original.deletedAt ? 'i-mdi-delete' : 'i-mdi-restore',
                                      variant: 'link',
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
                accessorKey: 'id',
                header: t('common.id'),
                cell: ({ row }) => h(IDCopyBadge, { id: row.original.id, icon: 'i-mdi-content-copy', variant: 'link' }),
            },
            {
                accessorKey: 'username',
                header: ({ column }) => {
                    return h(TableSortButton, {
                        column,
                        label: t('common.username'),
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
                        ? h(CitizenInfoPopover, {
                              userId: row.original.lastChar,
                              key: row.original.lastChar,
                          })
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
                    return h(TableSortButton, {
                        column,
                        label: t('common.license'),
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
                    ref="formRef"
                    class="my-2 flex w-full flex-1 flex-col gap-2"
                    :schema="schema"
                    :state="query"
                    @submit="commitValidatedQuery"
                >
                    <div class="flex flex-1 flex-row gap-2">
                        <UFormField class="flex-1" :label="$t('common.search')" name="license">
                            <UInput
                                ref="input"
                                v-model="query.license"
                                class="w-full"
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
                                        class="w-full"
                                        type="text"
                                        name="username"
                                        :placeholder="$t('common.username')"
                                        block
                                    />
                                </UFormField>

                                <UFormField
                                    class="flex-1"
                                    name="externalId"
                                    :label="$t('components.auth.SocialLogins.external_id')"
                                >
                                    <UInput
                                        v-model="query.externalId"
                                        class="w-full"
                                        type="text"
                                        name="externalId"
                                        :placeholder="$t('components.auth.SocialLogins.external_id')"
                                        block
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
                    :meta="meta"
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
