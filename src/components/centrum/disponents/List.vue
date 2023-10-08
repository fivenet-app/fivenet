<script lang="ts" setup>
import { MonitorIcon } from 'mdi-vue3';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import { useCentrumStore } from '~/store/centrum';

const centrumStore = useCentrumStore();
const { disponents } = storeToRefs(centrumStore);
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
                    <div class="grid grid-cols-3 gap-4">
                        <DataNoDataBlock
                            v-if="disponents && disponents.length === 0"
                            :icon="MonitorIcon"
                            :type="$t('common.disponents', 2)"
                            class="mt-5"
                        />
                        <template v-else>
                            <CitizenInfoPopover
                                v-for="disponent in disponents"
                                :key="disponent.userId"
                                type="button"
                                textClass="text-white bg-primary-500 hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-0.5"
                                :user="disponent"
                            />
                        </template>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
