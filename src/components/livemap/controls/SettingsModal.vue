<script lang="ts" setup>
import { useSettingsStore } from '~/store/settings';

const { isOpen } = useModal();

const settingsStore = useSettingsStore();
const { livemap } = storeToRefs(settingsStore);
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-xl' }">
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('common.setting', 2) }}
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                </div>
            </template>

            <UForm :schema="undefined" :state="{}">
                <UFormGroup name="centerSelectedMarker" :label="$t('components.livemap.center_selected_marker')">
                    <UToggle v-model="livemap.centerSelectedMarker">
                        <span class="sr-only">{{ $t('components.livemap.center_selected_marker') }}</span>
                    </UToggle>
                </UFormGroup>

                <UFormGroup name="livemapMarkerSize" :label="$t('components.livemap.settings.marker_size')">
                    <URange
                        v-model="livemap.markerSize"
                        name="livemapMarkerSize"
                        class="my-auto h-1.5 w-full cursor-grab rounded-full"
                        :min="16"
                        :max="30"
                        :step="2"
                    />
                    <span class="text-sm">{{ livemap.markerSize }}</span>
                </UFormGroup>

                <UFormGroup name="showUnitNames" :label="$t('components.livemap.show_unit_names')">
                    <UToggle v-model="livemap.showUnitNames">
                        <span class="sr-only">{{ $t('components.livemap.show_unit_names') }}</span>
                    </UToggle>
                </UFormGroup>

                <UFormGroup name="showUnitStatus" :label="$t('components.livemap.show_unit_status')">
                    <UToggle v-model="livemap.showUnitStatus">
                        <span class="sr-only">{{ $t('components.livemap.show_unit_status') }}</span>
                    </UToggle>
                </UFormGroup>

                <UFormGroup name="showAllDispatches" :label="$t('components.livemap.show_all_dispatches')">
                    <UToggle v-model="livemap.showAllDispatches">
                        <span class="sr-only">{{ $t('components.livemap.show_all_dispatches') }}</span>
                    </UToggle>
                </UFormGroup>
            </UForm>

            <template #footer>
                <UButton block class="flex-1" color="black" @click="isOpen = false">
                    {{ $t('common.close', 1) }}
                </UButton>
            </template>
        </UCard>
    </UModal>
</template>
