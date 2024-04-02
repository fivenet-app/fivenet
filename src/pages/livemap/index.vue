<script lang="ts" setup>
import { MonitorIcon, RobotIcon } from 'mdi-vue3';
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
                        class="group mt-0.5 flex w-full flex-row items-center justify-center rounded-md p-1 text-xs font-medium hover:bg-primary-100/10 hover:transition-all"
                        :class="
                            getCurrentMode === CentrumMode.AUTO_ROUND_ROBIN
                                ? 'bg-info-400/10 text-info-500 ring-info-400/20'
                                : disponents.length === 0
                                  ? 'bg-warn-400/10 text-warn-500 ring-warn-400/20'
                                  : 'bg-success-500/10 text-success-400 ring-success-500/20'
                        "
                        @click="openDisponents = true"
                    >
                        <template v-if="getCurrentMode !== CentrumMode.AUTO_ROUND_ROBIN">
                            <MonitorIcon class="mr-1 size-5" />
                            <span class="truncate">
                                {{ $t('common.disponent', disponents.length) }}
                            </span>
                        </template>
                        <template v-else>
                            <RobotIcon class="mr-1 size-5" />
                            <span class="truncate">
                                {{ $t('enums.centrum.CentrumMode.AUTO_ROUND_ROBIN') }}
                            </span>
                        </template>
                    </UButton>
                </template>
            </UDashboardNavbar>

            <LivemapHolder class="size-full" />
        </UDashboardPanel>
    </UDashboardPage>
</template>
