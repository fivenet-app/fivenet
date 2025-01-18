<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { splitURL, urlHomePage } from '~/components/internet/helper';
import { useInternetStore } from '~/store/internet';
import type { GetPageResponse } from '~~/gen/ts/services/internet/internet';
import HomePage from './pages/HomePage.vue';
import NotFound from './pages/NotFound.vue';
import WebPage from './pages/WebPage.vue';
import { localPages, localPagesDomains } from './pages/helpers';

const props = defineProps<{
    tabId: number;
}>();

const internetStore = useInternetStore();
const { activeTab, tabs } = storeToRefs(internetStore);

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

        handleGRPCError(err);
        throw err;
    }
}

const {
    data: page,
    status,
    refresh,
} = useLazyAsyncData(`internet-page-${tab.value?.domain}:${tab.value?.path}`, () => getPage(), {
    immediate: false,
});

watch(tab, async () => {
    if (!tab.value) {
        return;
    }

    state.url = tab.value.domain + (tab.value.path && tab.value.path !== '' ? tab.value.path : '/');

    // Skip local pages
    if (localPagesDomains.includes(tab.value?.domain)) {
        return;
    }

    if (tab.value.active) {
        refresh();
    }
});

function goToPage(domain: string, path?: string): void {
    state.url = domain + (path && path !== '' ? path : '/');
    internetStore.goTo(domain, path);
}

const schema = z.object({
    url: z.string().max(128),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    url: tab.value ? tab.value.domain + (tab.value.path && tab.value.path !== '' ? tab.value.path : '') : '',
});

const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    if (!activeTab.value) {
        return;
    }

    const split = splitURL(event.data.url);
    if (!split) {
        return;
    }

    goToPage(split.domain.toLowerCase(), split.path?.toLowerCase());
}, 500);
</script>

<template>
    <UDashboardToolbar
        v-if="activeTab"
        :ui="{ wrapper: 'min-h-[41px] bg-gray-100 dark:bg-gray-800 p-0 gap-x-0', container: 'gap-x-0 gap-y-0 mx-1' }"
    >
        <div class="flex flex-1 items-center gap-1">
            <UButton
                variant="ghost"
                color="white"
                icon="i-mdi-chevron-left"
                :disabled="activeTab.history.length === 0 || status === 'pending'"
                @click="internetStore.back()"
            />

            <UButton
                :disabled="activeTab.domain === urlHomePage || status === 'pending'"
                variant="ghost"
                color="white"
                icon="i-mdi-home"
                @click="goToPage(urlHomePage)"
            />

            <UButton
                variant="ghost"
                color="white"
                icon="i-mdi-refresh"
                :disabled="status === 'pending'"
                :loading="status === 'pending'"
                @click="refresh"
            />

            <UForm :schema="schema" :state="state" class="flex flex-1 items-center gap-1" @submit="onSubmitThrottle">
                <UInput v-model="state.url" type="text" class="flex-1" :ui="{ icon: { trailing: { pointer: '' } } }">
                    <template #trailing>
                        <UButton
                            v-show="state.url !== ''"
                            color="gray"
                            variant="link"
                            icon="i-mdi-close"
                            :padded="false"
                            @click="state.url = ''"
                        />
                    </template>
                </UInput>
            </UForm>
        </div>
    </UDashboardToolbar>

    <UDashboardPanelContent class="p-0">
        <HomePage v-if="tab?.domain === urlHomePage || tab?.domain === ''" v-model="tab" />
        <template v-else-if="tab?.domain && localPages[tab.domain]">
            <component :is="localPages[tab.domain]" v-model="tab" />
        </template>
        <NotFound v-else-if="!page?.page" v-model="tab" />
        <WebPage v-else v-model="tab" :page="page.page" />
    </UDashboardPanelContent>
</template>
