<script lang="ts" setup>
import { useInternetStore } from '~/store/internet';
import type { GetPageResponse } from '~~/gen/ts/services/internet/internet';
import HomePage from './HomePage.vue';
import NotFound from './NotFound.vue';

const props = defineProps<{
    tabId: number;
}>();

const internetStore = useInternetStore();
const { tabs } = storeToRefs(internetStore);

const tab = computed(() => tabs.value.find((t) => t.id === props.tabId));

async function getPage(): Promise<GetPageResponse | undefined> {
    if (!tab.value) {
        return;
    }

    try {
        const call = getGRPCInternetClient().getPage({
            url: tab.value.url,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        const err = e as RpcError;
        if (err.message.includes('ErrDomainNotFound')) {
            return undefined;
        }

        handleGRPCError(err);
        throw err;
    }
}

const { data: page } = useLazyAsyncData(`internet-page-${tab.value?.url}`, () => getPage());

// TODO load page from server
</script>

<template>
    <UDashboardPanelContent class="h-full overflow-x-auto p-0">
        <HomePage v-if="(!tab && !page) || (!page && tab?.label === '')" />
        <NotFound v-else-if="!page" />
        <template v-else>
            {{ tab }}
        </template>
    </UDashboardPanelContent>
</template>
