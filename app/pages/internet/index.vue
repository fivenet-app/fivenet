<script lang="ts" setup>
import { urlHomePage } from '~/components/internet/helper';
import InternetTab from '~/components/internet/InternetTab.vue';
import { useInternetStore } from '~/store/internet';

useHead({
    title: 'common.internet',
});
definePageMeta({
    title: 'common.internet',
    requiresAuth: true,
    permission: 'TODOService.TODOMethod',
});

const internetStore = useInternetStore();
const { selectedTab, tabs } = storeToRefs(internetStore);

// Set thread as query param for persistence between reloads
const router = useRouter();

function updateQuery(): void {
    if (!selectedTab.value) {
        router.replace({ query: {} });
    } else {
        // Hash is specified here to prevent the page from scrolling to the top
        router.replace({ query: { tab: selectedTab.value }, hash: '#' });
    }
}

onMounted(async () => {
    if (tabs.value.length === 0) {
        await internetStore.newTab(true);
    }

    if (!selectedTab.value && tabs.value[0]) {
        await internetStore.selectTab(tabs.value[0].id);
    }

    updateQuery();
});

watch(selectedTab, updateQuery, { deep: true });

const tab = computed(() => tabs.value.find((t) => t.id === selectedTab.value));
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
                                link.id === selectedTab &&
                                    'after:bg-primary-500 dark:after:bg-primary-400 text-gray-900 after:rounded-full dark:text-white',
                            ]"
                        >
                            <span v-if="link.id === selectedTab" class="sr-only"> Current page: </span>
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

            <UDashboardToolbar
                v-if="tab"
                :ui="{ wrapper: 'bg-gray-100 dark:bg-gray-800 p-0 gap-x-0', container: 'gap-x-0 gap-y-0' }"
            >
                <div class="flex flex-1 items-center gap-1">
                    <UButtonGroup>
                        <UButton variant="ghost" color="white" icon="i-mdi-chevron-left" />
                        <UButton variant="ghost" color="white" icon="i-mdi-chevron-right" />
                    </UButtonGroup>

                    <UButton
                        :disabled="tab.url.startsWith(urlHomePage)"
                        variant="ghost"
                        color="white"
                        icon="i-mdi-home"
                        @click="tab && (tab.url = urlHomePage)"
                    />

                    <UButton variant="ghost" color="white" icon="i-mdi-refresh" />

                    <ClientOnly>
                        <UInputMenu v-model="tab.url" class="mx-1 flex-1" />
                    </ClientOnly>
                </div>
            </UDashboardToolbar>

            <UDashboardPanelContent class="p-0">
                <UTabs
                    v-model="selectedTab"
                    :items="tabs"
                    :ui="{ wrapper: 'space-y-0', list: { base: 'hidden' }, padding: 'p-0' }"
                >
                    <template #item="{ item }">
                        <Suspense>
                            <InternetTab :tab-id="item.id" />
                        </Suspense>
                    </template>
                </UTabs>
            </UDashboardPanelContent>
        </UDashboardPanel>
    </UDashboardPage>
</template>
