<script lang="ts" setup>
import ThemePickerButton from '~/components/partials/ThemePickerButton.vue';

const settingsStore = useSettingsStore();
const { design, streamerMode } = storeToRefs(settingsStore);
</script>

<template>
    <UFormField
        class="grid grid-cols-2 items-center gap-2"
        name="darkMode"
        :label="$t('components.auth.user_settings.color_mode')"
    >
        <UColorModeSelect class="w-full" color="neutral" variant="outline" />
    </UFormField>

    <UFormField class="grid grid-cols-2 items-center gap-2" name="primaryColor" :label="$t('common.color', 1)">
        <div class="grid grid-cols-2 gap-1 sm:grid-cols-3 lg:grid-cols-4">
            <ThemePickerButton
                v-for="color in [...primaryColors, ...backgroundColors]"
                :key="color.label"
                :label="color.label"
                :chip="color.chip.color"
                :selected="design.ui.primary === color.label"
                @click="design.ui.primary = color.label"
            />
        </div>
    </UFormField>

    <UFormField
        class="grid grid-cols-2 items-center gap-2"
        name="grayColor"
        :label="$t('components.auth.user_settings.background_color')"
    >
        <div class="grid grid-cols-2 gap-1 sm:grid-cols-3 lg:grid-cols-4">
            <ThemePickerButton
                v-for="color in [...primaryColors, ...backgroundColors]"
                :key="color.label"
                :label="color.label"
                :chip="color.chip.color"
                :selected="design.ui.gray === color.label"
                @click="design.ui.gray = color.label"
            />
        </div>
    </UFormField>

    <UFormField
        class="flex items-center gap-4"
        name="streamerMode"
        :label="$t('components.auth.user_settings.streamer_mode.title')"
        :description="$t('components.auth.user_settings.streamer_mode.description')"
        :ui="{ wrapper: 'min-w-0 flex-1', container: 'flex shrink-0 items-center justify-end' }"
    >
        <USwitch v-model="streamerMode" name="streamerMode" />
    </UFormField>
</template>
