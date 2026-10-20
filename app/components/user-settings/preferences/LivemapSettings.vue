<script lang="ts" setup>
import { z } from 'zod';
import CentrumSidebarPlacementSettings from '~/components/user-settings/preferences/CentrumSidebarPlacementSettings.vue';

const props = withDefaults(
    defineProps<{
        showSidebarPlacement?: boolean;
    }>(),
    {
        showSidebarPlacement: true,
    },
);

const { can } = useAuth();

const settingsStore = useSettingsStore();
const { livemap } = storeToRefs(settingsStore);

const livemapStore = useLivemapStore();
const { markerDragEnabled } = storeToRefs(livemapStore);

const schema = z.object({
    markerSize: z.coerce.number().min(14).max(32),
    centerSelectedMarker: z.coerce.boolean(),
    showUnitNames: z.coerce.boolean(),
    showUnitStatus: z.coerce.boolean(),
    showAllDispatches: z.coerce.boolean(),
    showGrid: z.coerce.boolean(),
    useUnitColor: z.coerce.boolean(),
});
</script>

<template>
    <UForm class="space-y-4" :schema="schema" :state="livemap">
        <UFormField
            class="flex items-center gap-4"
            name="centerSelectedMarker"
            :label="$t('components.livemap.center_selected_marker')"
            :description="$t('components.livemap.center_selected_marker_description')"
            :ui="{ wrapper: 'min-w-0 flex-1', container: 'flex shrink-0 items-center justify-end' }"
        >
            <USwitch v-model="livemap.centerSelectedMarker" />
        </UFormField>

        <UFormField
            class="flex items-center gap-4"
            name="markerSize"
            :label="$t('components.livemap.settings.marker_size')"
            :description="$t('components.livemap.settings.marker_size_description')"
            :ui="{ wrapper: 'min-w-0 flex-1', container: 'flex shrink-0 items-center justify-end' }"
        >
            <UInputNumber v-model="livemap.markerSize" :min="14" :max="32" :step="2" />
        </UFormField>

        <UFormField
            class="flex items-center gap-4"
            name="showUnitNames"
            :label="$t('components.livemap.show_unit_names')"
            :description="$t('components.livemap.show_unit_names_description')"
            :ui="{ wrapper: 'min-w-0 flex-1', container: 'flex shrink-0 items-center justify-end' }"
        >
            <USwitch v-model="livemap.showUnitNames" />
        </UFormField>

        <UFormField
            class="flex items-center gap-4"
            name="showUnitStatus"
            :label="$t('components.livemap.show_unit_status')"
            :description="$t('components.livemap.show_unit_status_description')"
            :ui="{ wrapper: 'min-w-0 flex-1', container: 'flex shrink-0 items-center justify-end' }"
        >
            <USwitch v-model="livemap.showUnitStatus" />
        </UFormField>

        <UFormField
            class="flex items-center gap-4"
            name="showAllDispatches"
            :label="$t('components.livemap.show_all_dispatches')"
            :description="$t('components.livemap.show_all_dispatches_description')"
            :ui="{ wrapper: 'min-w-0 flex-1', container: 'flex shrink-0 items-center justify-end' }"
        >
            <USwitch v-model="livemap.showAllDispatches" />
        </UFormField>

        <UFormField
            class="flex items-center gap-4"
            name="showGrid"
            :label="$t('components.livemap.show_grid')"
            :description="$t('components.livemap.show_grid_description')"
            :ui="{ wrapper: 'min-w-0 flex-1', container: 'flex shrink-0 items-center justify-end' }"
        >
            <USwitch v-model="livemap.showGrid" />
        </UFormField>

        <UFormField
            class="flex items-center gap-4"
            name="useUnitColor"
            :label="$t('components.livemap.use_unit_color')"
            :description="$t('components.livemap.use_unit_color_description')"
            :ui="{ wrapper: 'min-w-0 flex-1', container: 'flex shrink-0 items-center justify-end' }"
        >
            <USwitch v-model="livemap.useUnitColor" />
        </UFormField>

        <UFormField
            v-if="can('livemap.LivemapService/CreateOrUpdateMarker').value"
            class="flex items-center gap-4"
            name="markerDragEnabled"
            :label="$t('components.livemap.enable_marker_dragging')"
            :description="$t('components.livemap.enable_marker_dragging_description')"
            :ui="{ wrapper: 'min-w-0 flex-1', container: 'flex shrink-0 items-center justify-end' }"
        >
            <USwitch v-model="markerDragEnabled" />
        </UFormField>
    </UForm>

    <template v-if="props.showSidebarPlacement">
        <USeparator />

        <CentrumSidebarPlacementSettings />
    </template>
</template>
