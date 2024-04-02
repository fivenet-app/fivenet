<script lang="ts" setup>
import { MonitorIcon } from 'mdi-vue3';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import { useCentrumStore } from '~/store/centrum';
import { CentrumMode } from '~~/gen/ts/resources/centrum/settings';

defineProps<{
    open: boolean;
}>();

defineEmits<{
    (e: 'close'): void;
}>();

const centrumStore = useCentrumStore();
const { disponents, getCurrentMode } = storeToRefs(centrumStore);
</script>

<template>
    <UModal :model-value="open" :ui="{ width: '!max-w-2xl' }" @update:model-value="$emit('close')">
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('common.disponents', 2) }}
                    </h3>

                    <UButton
                        color="gray"
                        variant="ghost"
                        icon="i-heroicons-x-mark-20-solid"
                        class="-my-1"
                        @click="$emit('close')"
                    />
                </div>
            </template>

            <UBadge color="gray">
                {{ $t('common.mode') }}: {{ $t(`enums.centrum.CentrumMode.${CentrumMode[getCurrentMode ?? 0]}`) }}
            </UBadge>

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
                    text-class="text-neutral bg-primary-500 hover:bg-primary-100/10  font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-0.5"
                    :user="disponent"
                />
            </div>
        </UCard>
    </UModal>
</template>
