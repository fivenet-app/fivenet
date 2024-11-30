<script lang="ts" setup>
import { useInternetStore } from '~/store/internet';
import HomePage from './HomePage.vue';

const props = defineProps<{
    tabId: number;
}>();

const internetStore = useInternetStore();
const { tabs } = storeToRefs(internetStore);

const tab = computed(() => tabs.value.find((t) => t.id === props.tabId));
// TODO
</script>

<template>
    <UDashboardPanelContent class="overflow-x-auto p-0">
        <Suspense>
            <HomePage v-if="!tab || tab.label === ''" />
            <template v-else>
                {{ tab }}
            </template>
        </Suspense>
    </UDashboardPanelContent>
</template>
