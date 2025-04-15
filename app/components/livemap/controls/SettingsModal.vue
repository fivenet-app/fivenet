<script lang="ts" setup>
import { z } from 'zod';
import { useSettingsStore } from '~/stores/settings';

const { isOpen } = useModal();

const settingsStore = useSettingsStore();
const { livemap } = storeToRefs(settingsStore);

const schema = z.object({
    markerSize: z.number().min(14).max(32),
    centerSelectedMarker: z.boolean(),
    showUnitNames: z.boolean(),
    showUnitStatus: z.boolean(),
    showAllDispatches: z.boolean(),
});
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-xl' }">
        <UCard :ui="{ ring: '', divide: 'divide-y divide-neutral-100 dark:divide-neutral-800' }">
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('common.setting', 2) }}
                    </h3>

                    <UButton color="neutral" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                </div>
            </template>

            <UForm :schema="schema" :state="livemap">
                <UFormField name="centerSelectedMarker" :label="$t('components.livemap.center_selected_marker')">
                    <USwitch v-model="livemap.centerSelectedMarker">
                        <span class="sr-only">{{ $t('components.livemap.center_selected_marker') }}</span>
                    </USwitch>
                </UFormField>

                <UFormField name="markerSize" :label="$t('components.livemap.settings.marker_size')">
                    <URange
                        v-model="livemap.markerSize"
                        class="my-auto h-1.5 w-full cursor-grab rounded-full"
                        :min="14"
                        :max="32"
                        :step="2"
                    />
                    <span class="text-sm">{{ livemap.markerSize }}</span>
                </UFormField>

                <UFormField name="showUnitNames" :label="$t('components.livemap.show_unit_names')">
                    <USwitch v-model="livemap.showUnitNames">
                        <span class="sr-only">{{ $t('components.livemap.show_unit_names') }}</span>
                    </USwitch>
                </UFormField>

                <UFormField name="showUnitStatus" :label="$t('components.livemap.show_unit_status')">
                    <USwitch v-model="livemap.showUnitStatus">
                        <span class="sr-only">{{ $t('components.livemap.show_unit_status') }}</span>
                    </USwitch>
                </UFormField>

                <UFormField name="showAllDispatches" :label="$t('components.livemap.show_all_dispatches')">
                    <USwitch v-model="livemap.showAllDispatches">
                        <span class="sr-only">{{ $t('components.livemap.show_all_dispatches') }}</span>
                    </USwitch>
                </UFormField>
            </UForm>

            <template #footer>
                <UButton block class="flex-1" color="neutral" @click="isOpen = false">
                    {{ $t('common.close', 1) }}
                </UButton>
            </template>
        </UCard>
    </UModal>
</template>
