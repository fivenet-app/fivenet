<script lang="ts" setup>
import { joinURL } from '~/components/internet/helper';
import InternetTab from '~/components/internet/InternetTab.vue';
import { useInternetStore } from '~/store/internet';

useHead({
    title: 'common.internet',
});
definePageMeta({
    title: 'common.internet',
    requiresAuth: true,
});

const internetStore = useInternetStore();
const { activeTab, tabs } = storeToRefs(internetStore);

// Set thread as query param for persistence between reloads
const route = useRoute();
const router = useRouter();

function updateQuery(): void {
    if (!activeTab.value) {
        router.replace({ query: {} });
    } else {
        // Hash is specified here to prevent the page from scrolling to the top
        router.replace({
            query: { tab: activeTab.value.id, url: joinURL(activeTab.value.domain, activeTab.value.path) },
            hash: '#',
        });
    }
}

onMounted(async () => {
    if (tabs.value.length === 0) {
        internetStore.newTab(true);
    }

    if (!activeTab.value && tabs.value[0]) {
        internetStore.selectTab(tabs.value[0].id);
    }

    updateQuery();
});

const selectedTab = computed({
    get() {
        const tabId = parseInt(route.query.tab as string);
        const index = tabs.value.findIndex((item) => item.id === tabId);
        if (index === -1) {
            return 0;
        }

        return index;
    },
    set(value) {
        // Hash is specified here to prevent the page from scrolling to the top
        router.replace({ query: { tab: tabs.value[value]?.id }, hash: '#' });
    },
});

watch(activeTab, () => updateQuery(), { deep: true });
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardToolbar
                :ui="{
                    wrapper: 'p-0 gap-x-0',
                    container:
                        'gap-x-0 gap-y-0 justify-stretch items-stretch h-full flex flex-row bg-gray-100 p-0 px dark:bg-gray-800 overflow-x-hidden',
                }"
            >
                <UDashboardNavbarToggle class="px-4 py-2" />

                <UHorizontalNavigation
                    :links="
                        tabs.map((t) => ({
                            ...t,
                            click: () => internetStore.selectTab(t.id),
                        }))
                    "
                    :ui="{
                        container: 'divide-x divide-gray-200 dark:divide-gray-600',
                        inner: 'min-w-60 max-w-60',
                        wrapper: 'overflow-x-auto',
                    }"
                >
                    <template #default="{ link }">
                        <span
                            class="group-hover:text-primary relative flex-1 truncate text-left"
                            :class="[
                                link.id === activeTab?.id &&
                                    'after:bg-primary-500 dark:after:bg-primary-400 text-gray-900 after:rounded-full dark:text-white',
                            ]"
                        >
                            <span v-if="link.id === activeTab?.id" class="sr-only"> Current page: </span>
                            {{ link.label === '' ? $t('common.home') : link.label }}
                        </span>
                    </template>

                    <template #badge="{ link }">
                        <UButton icon="i-mdi-close" variant="ghost" color="black" @click="internetStore.closeTab(link.id)" />
                    </template>
                </UHorizontalNavigation>

                <ClientOnly>
                    <UTooltip :text="$t('components.internet.new_tab')">
                        <UButton icon="i-mdi-plus" variant="ghost" color="black" @click="internetStore.newTab()" />
                    </UTooltip>
                </ClientOnly>
            </UDashboardToolbar>

            <UDashboardPanelContent class="p-0">
                <UTabs
                    v-model="selectedTab"
                    :items="tabs"
                    :ui="{
                        wrapper: 'space-y-0 h-full',
                        list: { base: 'hidden' },
                        padding: 'p-0',
                        container: 'h-full',
                        base: 'h-full',
                    }"
                    :unmount="false"
                >
                    <template #item="{ item }">
                        <InternetTab :tab-id="item.id" />
                    </template>
                </UTabs>
            </UDashboardPanelContent>
        </UDashboardPanel>
    </UDashboardPage>
</template>
