<script lang="ts" setup>
import { MonitorIcon } from 'mdi-vue3';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import { useCentrumStore } from '~/store/centrum';
import { CentrumMode } from '~~/gen/ts/resources/dispatch/settings';

const centrumStore = useCentrumStore();
const { disponents, getCurrentMode } = storeToRefs(centrumStore);
</script>

<template>
    <div class="py-2">
        <div class="px-2 sm:px-6 lg:px-8">
            <div class="relative">
                <h3 class="text-2xl font-semibold leading-6">
                    {{ $t('common.disponents', 2) }}
                </h3>
            </div>
            <div class="sm:flex sm:items-center pt-4">
                <div class="sm:flex-auto">
                    <p class="text-neutral text-sm" :class="disponents && disponents.length > 0 ? 'mb-4' : ''">
                        <span
                            class="inline-flex items-center rounded-md bg-gray-400/10 px-2 py-1 text-xs font-medium text-gray-400 ring-1 ring-inset ring-gray-400/20"
                        >
                            {{ $t('common.mode') }}: {{ $t(`enums.centrum.CentrumMode.${CentrumMode[getCurrentMode ?? 0]}`) }}
                        </span>
                    </p>

                    <DataNoDataBlock
                        v-if="disponents && disponents.length === 0"
                        :icon="MonitorIcon"
                        :type="$t('common.disponents', 2)"
                        class="mt-5"
                    />
                    <div v-else class="grid grid-cols-3 gap-4">
                        <CitizenInfoPopover
                            v-for="disponent in disponents"
                            :key="disponent.userId"
                            type="button"
                            text-class="text-neutral bg-primary-500 hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-0.5"
                            :user="disponent"
                        />
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
