<script lang="ts" setup>
import DisponentsModal from '~/components/centrum/disponents/DisponentsModal.vue';
import LivemapHolder from '~/components/livemap/LivemapHolder.vue';
import { useCentrumStore } from '~/store/centrum';
import { CentrumMode } from '~~/gen/ts/resources/centrum/settings';

useHead({
    title: 'common.livemap',
});
definePageMeta({
    title: 'common.livemap',
    requiresAuth: true,
    permission: 'LivemapperService.Stream',
});

const centrumStore = useCentrumStore();
const { getCurrentMode, disponents } = storeToRefs(centrumStore);

const openDisponents = ref(false);
</script>

<template>
    <UDashboardPage class="size-full">
        <UDashboardPanel grow class="size-full">
            <UDashboardNavbar :title="$t('common.livemap')">
                <template #right>
                    <DisponentsModal :open="openDisponents" @close="openDisponents = false" />

                    <UButton
                        :icon="getCurrentMode !== CentrumMode.AUTO_ROUND_ROBIN ? 'i-mdi-monitor' : 'i-mdi-robot'"
                        :color="
                            getCurrentMode === CentrumMode.AUTO_ROUND_ROBIN
                                ? 'primary'
                                : disponents.length === 0
                                  ? 'amber'
                                  : 'green'
                        "
                        truncate
                        @click="openDisponents = true"
                    >
                        <template v-if="getCurrentMode !== CentrumMode.AUTO_ROUND_ROBIN">
                            {{ $t('common.disponent', disponents.length) }}
                        </template>
                        <template v-else>
                            {{ $t('enums.centrum.CentrumMode.AUTO_ROUND_ROBIN') }}
                        </template>
                    </UButton>
                </template>
            </UDashboardNavbar>

            <LivemapHolder class="size-full" />
        </UDashboardPanel>
    </UDashboardPage>
</template>
