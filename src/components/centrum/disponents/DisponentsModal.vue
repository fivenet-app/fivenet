<script lang="ts" setup>
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
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
                        <UBadge color="gray">
                            {{ $t('common.mode') }}: {{ $t(`enums.centrum.CentrumMode.${CentrumMode[getCurrentMode ?? 0]}`) }}
                        </UBadge>
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="$emit('close')" />
                </div>
            </template>

            <DataNoDataBlock
                v-if="disponents && disponents.length === 0"
                icon="i-mdi-monitor"
                :type="$t('common.disponents', 2)"
                class="mt-5"
            />
            <UPageGrid v-else>
                <UPageCard v-for="disponent in disponents" :key="disponent.userId" :title="disponent.firstname">
                    <PhoneNumberBlock :number="disponent.phoneNumber" />
                </UPageCard>
            </UPageGrid>
        </UCard>
    </UModal>
</template>
