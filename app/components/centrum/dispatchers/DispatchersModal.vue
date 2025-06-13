<script lang="ts" setup>
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import { useCentrumStore } from '~/stores/centrum';
import { CentrumMode } from '~~/gen/ts/resources/centrum/settings';

const { isOpen } = useModal();

const centrumStore = useCentrumStore();
const { dispatchers, anyDispatchersActive, getCurrentMode } = storeToRefs(centrumStore);
</script>

<template>
    <UModal :ui="{ width: '!max-w-2xl' }">
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('common.dispatchers', 2) }}
                        <UBadge color="gray">
                            {{ $t('common.mode') }}: {{ $t(`enums.centrum.CentrumMode.${CentrumMode[getCurrentMode ?? 0]}`) }}
                        </UBadge>
                    </h3>

                    <UButton class="-my-1" color="gray" variant="ghost" icon="i-mdi-window-close" @click="isOpen = false" />
                </div>
            </template>

            <DataNoDataBlock
                v-if="!dispatchers || !anyDispatchersActive"
                icon="i-mdi-monitor"
                :type="$t('common.dispatcher')"
            />
            <template v-else>
                <UPageGrid v-for="dispas in dispatchers" :key="dispas.job" class="gap-4 p-4" :cols="2">
                    <UPageCard
                        v-for="dispatcher in dispas.dispatchers"
                        :key="dispatcher.userId"
                        :title="`${dispatcher.firstname} ${dispatcher.lastname}`"
                        :ui="{
                            title: 'text-gray-900 dark:text-white text-base font-semibold flex items-center gap-1.5 line-clamp-2 whitespace-break-spaces',
                        }"
                    >
                        <PhoneNumberBlock :number="dispatcher.phoneNumber" />
                    </UPageCard>
                </UPageGrid>
            </template>

            <template #footer>
                <UButton class="flex-1" color="black" block @click="isOpen = false">
                    {{ $t('common.close', 1) }}
                </UButton>
            </template>
        </UCard>
    </UModal>
</template>
