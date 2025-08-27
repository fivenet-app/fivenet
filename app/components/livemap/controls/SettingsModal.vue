<script lang="ts" setup>
import { z } from 'zod';
import { useSettingsStore } from '~/stores/settings';

defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const settingsStore = useSettingsStore();
const { livemap } = storeToRefs(settingsStore);

const schema = z.object({
    markerSize: z.coerce.number().min(14).max(32),
    centerSelectedMarker: z.coerce.boolean(),
    showUnitNames: z.coerce.boolean(),
    showUnitStatus: z.coerce.boolean(),
    showAllDispatches: z.coerce.boolean(),
    showGrid: z.coerce.boolean(),
});
</script>

<template>
    <UModal :title="$t('common.setting', 2)">
        <template #body>
            <UForm :schema="schema" :state="livemap">
                <UFormField name="centerSelectedMarker" :label="$t('components.livemap.center_selected_marker')">
                    <USwitch v-model="livemap.centerSelectedMarker" />
                </UFormField>

                <UFormField name="markerSize" :label="$t('components.livemap.settings.marker_size')">
                    <USlider
                        v-model="livemap.markerSize"
                        class="my-auto h-1.5 w-full cursor-grab rounded-full"
                        :min="14"
                        :max="32"
                        :step="2"
                    />
                    <span class="text-sm">{{ livemap.markerSize }}</span>
                </UFormField>

                <UFormField name="showUnitNames" :label="$t('components.livemap.show_unit_names')">
                    <USwitch v-model="livemap.showUnitNames" />
                </UFormField>

                <UFormField name="showUnitStatus" :label="$t('components.livemap.show_unit_status')">
                    <USwitch v-model="livemap.showUnitStatus" />
                </UFormField>

                <UFormField name="showAllDispatches" :label="$t('components.livemap.show_all_dispatches')">
                    <USwitch v-model="livemap.showAllDispatches" />
                </UFormField>

                <UFormField name="showGrid" :label="$t('components.livemap.show_grid')">
                    <USwitch v-model="livemap.showGrid" />
                </UFormField>
            </UForm>
        </template>

        <template #footer>
            <UButton class="flex-1" block color="neutral" @click="$emit('close', false)">
                {{ $t('common.close', 1) }}
            </UButton>
        </template>
    </UModal>
</template>
