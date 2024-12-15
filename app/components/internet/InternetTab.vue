<script lang="ts" setup>
import { useInternetStore } from '~/store/internet';
import type { GetPageResponse } from '~~/gen/ts/services/internet/internet';
import { urlHomePage } from './helper';
import HomePage from './pages/HomePage.vue';
import NotFound from './pages/NotFound.vue';

const props = defineProps<{
    tabId: number;
}>();

defineEmits<{
    (e: 'urlChange', url: { domain: string; path?: string }): void;
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
            domain: tab.value.domain,
            path: tab.value.path,
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

const isDev = import.meta.dev;

const { data: page } = useLazyAsyncData(`internet-page-${tab.value?.domain}:${tab.value?.path}`, () => getPage());
</script>

<template>
    <UDashboardPanelContent class="p-0">
        <HomePage v-if="tab?.domain === urlHomePage || tab?.domain === ''" v-model="tab" />
        <NotFound v-else-if="!page" v-model="tab" />
        <template v-else>
            <template v-if="isDev">
                <span>Tab: {{ tab }}</span>
                <span>Page: {{ page }}</span>
            </template>
        </template>
    </UDashboardPanelContent>
</template>
