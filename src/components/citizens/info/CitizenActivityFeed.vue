<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import { ListUserActivityResponse } from '~~/gen/ts/services/citizenstore/citizenstore';
import CitizenActivityFeedEntry from '~/components/citizens/info/CitizenActivityFeedEntry.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import { BulletinBoardIcon } from 'mdi-vue3';

const { $grpc } = useNuxtApp();

const props = defineProps<{
    userId: number;
}>();

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * page.value : 0));

const { data, pending, refresh, error } = useLazyAsyncData(`citizeninfo-activity-${props.userId}-${page.value}`, () =>
    listUserActivity(),
);

async function listUserActivity(): Promise<ListUserActivityResponse> {
    try {
        const call = $grpc.getCitizenStoreClient().listUserActivity({
            pagination: {
                offset: offset.value,
            },
            userId: props.userId,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

watch(offset, async () => refresh());
</script>

<template>
    <div>
        <DataPendingBlock
            v-if="pending"
            :message="$t('common.loading', [`${$t('common.citizen', 1)} ${$t('common.activity')}`])"
        />
        <DataErrorBlock
            v-else-if="error"
            :title="$t('common.not_found', [`${$t('common.citizen', 1)} ${$t('common.activity')}`])"
            :retry="refresh"
        />
        <DataNoDataBlock
            v-else-if="data === null || data?.activity.length === 0"
            :type="`${$t('common.document', 1)} ${$t('common.relation', 2)}`"
            :icon="BulletinBoardIcon"
        />
        <ul v-else role="list" class="divide-y divide-gray-200">
            <li v-for="activity in data?.activity" :key="activity.id" class="py-4">
                <CitizenActivityFeedEntry :activity="activity" />
            </li>
        </ul>

        <div class="flex justify-end px-3 py-3.5 border-t border-gray-200 dark:border-gray-700">
            <UPagination
                v-model="page"
                :page-count="data?.pagination?.pageSize ?? 0"
                :total="data?.pagination?.totalCount ?? 0"
            />
        </div>
    </div>
</template>
