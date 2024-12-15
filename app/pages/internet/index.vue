<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { joinURL, splitURL, urlHomePage } from '~/components/internet/helper';
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
const { activeTab, selectedTab, tabs } = storeToRefs(internetStore);

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
        internetStore.newTab(true);
    }

    if (!selectedTab.value && tabs.value[0]) {
        internetStore.selectTab(tabs.value[0].id);
    }

    if (activeTab.value) {
        state.url = joinURL(activeTab.value.domain, activeTab.value.path);
    }

    updateQuery();
});

watch(selectedTab, updateQuery, { deep: true });

function goToPage(domain: string, path?: string): void {
    state.url = domain + (path && path !== '' ? path : '');

    internetStore.goTo(domain, path);
}

const schema = z.object({
    url: z.string().max(128),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    url: '',
});

const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    if (!activeTab.value) {
        return;
    }

    const split = splitURL(event.data.url);
    if (!split) {
        return;
    }

    goToPage(split.domain, split.path);
}, 500);
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
                v-if="activeTab"
                :ui="{ wrapper: 'bg-gray-100 dark:bg-gray-800 p-0 gap-x-0', container: 'gap-x-0 gap-y-0 mx-1' }"
            >
                <div class="flex flex-1 items-center gap-1">
                    <UButton
                        variant="ghost"
                        color="white"
                        icon="i-mdi-chevron-left"
                        :disabled="activeTab.history.length === 0"
                        @click="internetStore.back()"
                    />

                    <UButton
                        :disabled="activeTab.domain === urlHomePage"
                        variant="ghost"
                        color="white"
                        icon="i-mdi-home"
                        @click="goToPage(urlHomePage)"
                    />

                    <UButton variant="ghost" color="white" icon="i-mdi-refresh" />

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
                >
                    <template #item="{ item }">
                        <Suspense>
                            <InternetTab :tab-id="item.id" @url-change="goToPage($event.domain, $event.path ?? '')" />
                        </Suspense>
                    </template>
                </UTabs>
            </UDashboardPanelContent>
        </UDashboardPanel>
    </UDashboardPage>
</template>
