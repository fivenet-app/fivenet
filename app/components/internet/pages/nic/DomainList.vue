<script lang="ts" setup>
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { useInternetStore } from '~/stores/internet';
import type { ListDomainsResponse } from '~~/gen/ts/services/internet/domain';

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const internetStore = useInternetStore();

const page = useRouteQuery('page', '1', { transform: Number });

const { data, pending: loading, refresh } = useLazyAsyncData(`internet-domain-list-${page.value}`, () => listDomains());

async function listDomains(): Promise<ListDomainsResponse> {
    try {
        const call = $grpc.internet.domain.listDomains({
            pagination: {
                offset: calculateOffset(page.value, data.value?.pagination),
            },
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const columns = [
    {
        label: t('common.name'),
        key: 'name',
    },
    {
        label: t('common.approve', 2),
        key: 'active',
    },
    {
        label: t('common.expires_at'),
        key: 'expiresAt',
    },
    {
        label: t('common.action', 2),
        key: 'actions',
    },
];

// TODO
</script>

<template>
    <div>
        <UTable
            :rows="data?.domains"
            :columns="columns"
            :empty-state="{ icon: 'i-mdi-domain', label: $t('common.not_found', [$t('common.domain', 2)]) }"
        >
            <template #name-data="{ row }">
                <p class="text-gray-900 dark:text-white">{{ row.name }}.{{ row.tld.name }}</p>
            </template>

            <template #active-data="{ row }">
                <UBadge
                    :color="row.active ? 'success' : 'orange'"
                    :label="row.active ? $t('common.approve', 2) : $t('common.decline', 2)"
                />
            </template>

            <template #expiresAt-data="{ row }">
                <GenericTime v-if="row.expiresAt" :value="row.expiresAt" />
                <span v-else>∞</span>
            </template>

            <template #actions-data="{ row }">
                <UButton variant="link" color="emerald" icon="i-mdi-payment" />

                <UButton variant="link" icon="i-mdi-link" @click="internetStore.goTo(`${row.name}.${row.tld.name}`)" />

                <UButton variant="link" color="amber" icon="i-mdi-cog-transfer" />
                <!-- TODO handle actions -->
            </template>
        </UTable>

        <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />
    </div>
</template>
