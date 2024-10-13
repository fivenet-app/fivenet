<script lang="ts" setup>
import CitizenActivityFeedEntry from '~/components/citizens/info/CitizenActivityFeedEntry.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { useAuthStore } from '~/store/auth';
import type { ListUserActivityResponse } from '~~/gen/ts/services/citizenstore/citizenstore';

const props = defineProps<{
    userId: number;
}>();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`citizeninfo-activity-${props.userId}-${page.value}`, () => listUserActivity());

async function listUserActivity(): Promise<ListUserActivityResponse> {
    try {
        const call = getGRPCCitizenStoreClient().listUserActivity({
            pagination: {
                offset: offset.value,
            },
            userId: props.userId,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watch(offset, async () => refresh());
</script>

<template>
    <UAlert
        v-if="userId === activeChar?.userId && !attr('CitizenStoreService.ListUserActivity', 'Fields', 'Own').value"
        variant="subtle"
        color="red"
        icon="i-mdi-denied"
        :title="$t('components.citizens.CitizenInfoActivityFeed.own.title')"
        :description="$t('components.citizens.CitizenInfoActivityFeed.own.message')"
    />

    <div v-else>
        <DataPendingBlock
            v-if="loading"
            :message="$t('common.loading', [`${$t('common.citizen', 1)} ${$t('common.activity')}`])"
        />
        <DataErrorBlock
            v-else-if="error"
            :title="$t('common.not_found', [`${$t('common.citizen', 1)} ${$t('common.activity')}`])"
            :retry="refresh"
        />
        <DataNoDataBlock
            v-else-if="!data || data?.activity.length === 0"
            :type="`${$t('common.citizen', 1)} ${$t('common.activity')}`"
            icon="i-mdi-pulse"
        />

        <ul v-else role="list" class="divide-y divide-gray-100 dark:divide-gray-800">
            <li
                v-for="activity in data?.activity"
                :key="activity.id"
                class="hover:border-primary-500/25 dark:hover:border-primary-400/25 hover:bg-primary-100/50 dark:hover:bg-primary-900/10 border-white py-2 dark:border-gray-900"
            >
                <CitizenActivityFeedEntry :activity="activity" />
            </li>
        </ul>

        <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />
    </div>
</template>
