<script lang="ts" setup>
import { z } from 'zod';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import Pagination from '~/components/partials/Pagination.vue';
import StreamerModeAlert from '~/components/partials/StreamerModeAlert.vue';
import type { ListAccountsResponse } from '~~/gen/ts/services/settings/accounts';
import AccountEditModal from './AccountEditModal.vue';

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const modal = useModal();

const schema = z.object({
    license: z.string().max(64).optional(),
    enabled: z.boolean(),
});

type Schema = z.output<typeof schema>;

const query = reactive<Schema>({
    enabled: true,
});

const page = useRouteQuery('page', '1', { transform: Number });
const offset = computed(() =>
    accounts.value?.pagination?.pageSize ? accounts.value?.pagination?.pageSize * (page.value - 1) : 0,
);

const sort = useRouteQueryObject<TableSortable>('sort', {
    column: 'username',
    direction: 'asc',
});

const {
    data: accounts,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`settings-accounts-${sort.value.column}:${sort.value.direction}-${page.value}`, () => listAccounts());

async function listAccounts(): Promise<ListAccountsResponse> {
    try {
        const call = $grpc.settings.accounts.listAccounts({
            pagination: {
                offset: offset.value,
            },
            sort: sort.value,
            enabled: query.enabled,
            license: query.license,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watch(offset, async () => refresh());
watchDebounced(query, async () => refresh(), { debounce: 200, maxWait: 1250 });

async function deleteAccount(id: number): Promise<void> {
    try {
        const call = $grpc.settings.accounts.deleteAccount({
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

const settingsStore = useSettingsStore();
const { streamerMode } = storeToRefs(settingsStore);

const columns = [
    {
        key: 'actions',
        label: t('common.action', 2),
        sortable: false,
    },
    {
        key: 'username',
        label: t('common.username'),
        sortable: true,
    },
    {
        key: 'enabled',
        label: t('common.enabled'),
    },
    {
        key: 'createdAt',
        label: t('common.created_at'),
    },
    {
        key: 'updatedAt',
        label: t('common.updated_at'),
    },
    {
        key: 'license',
        label: t('common.license'),
        sortable: true,
    },
];
</script>

<template>
    <template v-if="streamerMode">
        <UDashboardNavbar :title="$t('pages.settings.accounts.title')">
            <template #right>
                <PartialsBackButton fallback-to="/settings" />
            </template>
        </UDashboardNavbar>

        <UDashboardPanelContent>
            <StreamerModeAlert />
        </UDashboardPanelContent>
    </template>
    <template v-else>
        <UDashboardNavbar :title="$t('pages.settings.accounts.title')">
            <template #right>
                <PartialsBackButton fallback-to="/settings" />
            </template>
        </UDashboardNavbar>

        <UDashboardToolbar>
            <UForm class="w-full" :schema="schema" :state="query" @submit="refresh()">
                <div class="flex w-full flex-row gap-2">
                    <UFormGroup class="flex-1" :label="$t('common.search')" name="license">
                        <UInput
                            ref="input"
                            v-model="query.license"
                            type="text"
                            name="license"
                            :placeholder="$t('common.license')"
                            block
                            leading-icon="i-mdi-search"
                            @keydown.esc="$event.target.blur()"
                        >
                            <template #trailing>
                                <UKbd value="/" />
                            </template>
                        </UInput>
                    </UFormGroup>

                    <UFormGroup
                        class="flex flex-initial flex-col"
                        name="enabled"
                        :label="$t('common.enabled')"
                        :ui="{ container: 'flex-1 flex' }"
                    >
                        <div class="flex flex-1 items-center">
                            <UToggle v-model="query.enabled">
                                <span class="sr-only">
                                    {{ $t('common.enabled') }}
                                </span>
                            </UToggle>
                        </div>
                    </UFormGroup>
                </div>
            </UForm>
        </UDashboardToolbar>

        <DataErrorBlock
            v-if="error"
            :title="$t('common.unable_to_load', [$t('common.account', 2)])"
            :error="error"
            :retry="refresh"
        />

        <UTable
            v-else
            class="flex-1"
            :loading="loading"
            :columns="columns"
            :rows="accounts?.accounts"
            :empty-state="{ icon: 'i-mdi-account-multiple', label: $t('common.not_found', [$t('common.account', 2)]) }"
        >
            <template #actions-data="{ row: account }">
                <UTooltip :text="$t('common.update')">
                    <UButton
                        variant="link"
                        icon="i-mdi-pencil"
                        @click="
                            modal.open(AccountEditModal, {
                                account: account,
                                'onUpdate:account': () => refresh(),
                            })
                        "
                    />
                </UTooltip>

                <UTooltip :text="$t('common.delete')">
                    <UButton
                        variant="link"
                        icon="i-mdi-delete"
                        color="error"
                        @click="
                            modal.open(ConfirmModal, {
                                confirm: async () => deleteAccount(account.id),
                            })
                        "
                    />
                </UTooltip>
            </template>

            <template #username-data="{ row: account }">
                <UTooltip :text="`${$t('common.id')}: ${account.id}`">
                    <span class="text-gray-900 dark:text-white">
                        {{ account.username === '' ? $t('common.na') : account.username }}
                    </span>
                </UTooltip>
            </template>

            <template #enabled-data="{ row: account }">
                <UBadge
                    :color="account.enabled ? 'success' : 'error'"
                    :label="account.enabled ? t('common.yes') : t('common.no')"
                />
            </template>

            <template #createdAt-data="{ row: account }">
                <GenericTime :value="toDate(account.createdAt)" />
            </template>

            <template #updatedAt-data="{ row: account }">
                <GenericTime :value="toDate(account.updatedAt)" />
            </template>

            <template #license-data="{ row: account }">
                <pre class="text-gray-900 dark:text-white" v-text="account.license" />
            </template>
        </UTable>

        <Pagination v-model="page" :pagination="accounts?.pagination" :loading="loading" :refresh="refresh" />
    </template>
</template>
