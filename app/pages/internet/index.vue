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
const { tab } = storeToRefs(internetStore);

// Set thread as query param for persistence between reloads
const router = useRouter();

function updateQuery(): void {
    // Hash is specified here to prevent the page from scrolling to the top
    router.replace({
        query: { url: joinURL(tab.value.domain, tab.value.path) },
        hash: '#',
    });
}

onMounted(async () => {
    updateQuery();
});

watch(tab, () => updateQuery(), { deep: true });
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardToolbar
                :ui="{
                    wrapper: 'p-0 gap-x-0',
                    container:
                        'gap-x-0 gap-y-0 justify-stretch items-stretch h-full flex flex-row p-0 px bg-gray-100 dark:bg-gray-800 overflow-x-hidden',
                }"
            >
                <UDashboardNavbarToggle class="px-4 py-2" />

                <UHorizontalNavigation
                    :links="[tab]"
                    :ui="{
                        container: 'flex-1 divide-x divide-gray-200 dark:divide-gray-600',
                        inner: 'flex-1',
                        base: 'justify-center',
                        wrapper: 'overflow-x-auto h-[60px]',
                    }"
                >
                    <template #default="{ link }">
                        <span
                            class="group-hover:text-primary relative truncate text-left"
                            :class="[
                                'after:bg-primary-500 dark:after:bg-primary-400 text-gray-900 after:rounded-full dark:text-white',
                            ]"
                        >
                            <span class="sr-only"> Current page: </span>
                            {{ link.label === '' ? $t('common.home') : link.label }}
                        </span>
                    </template>
                </UHorizontalNavigation>
            </UDashboardToolbar>

            <UDashboardPanelContent class="p-0">
                <InternetTab v-model="tab" />
            </UDashboardPanelContent>
        </UDashboardPanel>
    </UDashboardPage>
</template>
