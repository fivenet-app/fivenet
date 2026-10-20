<script lang="ts" setup>
const settingsStore = useSettingsStore();
const { eventsDisabled } = storeToRefs(settingsStore);
</script>

<template>
    <UFormField
        class="flex items-center gap-4"
        :label="$t('components.auth.user_settings.events_disabled.label')"
        :description="$t('components.auth.user_settings.events_disabled.description')"
        name="eventsDisabled"
        :ui="{ wrapper: 'min-w-0 flex-1', container: 'flex shrink-0 items-center justify-end' }"
    >
        <USwitch v-model="eventsDisabled" />
    </UFormField>
</template>
