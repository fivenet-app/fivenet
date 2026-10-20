<script lang="ts" setup>
import AccessibilitySettings from '~/components/user-settings/preferences/AccessibilitySettings.vue';
import AppearanceSettings from '~/components/user-settings/preferences/AppearanceSettings.vue';
import DocumentsSettings from '~/components/user-settings/preferences/DocumentsSettings.vue';
import OverviewQuickAccessSettings from '~/components/user-settings/preferences/OverviewQuickAccessSettings.vue';
import StartupSettings from '~/components/user-settings/preferences/StartupSettings.vue';
</script>

<template>
    <div class="space-y-4">
        <UPageCard
            id="user-settings-appearance"
            class="scroll-mt-24 transition-shadow target:ring-2 target:ring-primary"
            :title="$t('components.auth.user_settings.appearance.title')"
            :description="$t('components.auth.user_settings.appearance.description')"
        >
            <AppearanceSettings />
        </UPageCard>

        <UPageCard
            id="user-settings-startpage"
            class="scroll-mt-24 transition-shadow target:ring-2 target:ring-primary"
            :title="$t('components.auth.user_settings.startup.title')"
            :description="$t('components.auth.user_settings.startup.description')"
        >
            <StartupSettings />
        </UPageCard>

        <UPageCard
            id="user-settings-accessibility"
            class="scroll-mt-24 transition-shadow target:ring-2 target:ring-primary"
            :title="$t('components.auth.user_settings.accessibility.title')"
            :description="$t('components.auth.user_settings.accessibility.description')"
        >
            <AccessibilitySettings />
        </UPageCard>

        <UPageCard
            id="user-settings-quick-access"
            class="scroll-mt-24 transition-shadow target:ring-2 target:ring-primary"
            :title="$t('components.auth.user_settings.quick_access.title')"
            :description="$t('components.auth.user_settings.quick_access.description')"
        >
            <OverviewQuickAccessSettings />
        </UPageCard>

        <UPageCard
            id="user-settings-documents"
            class="scroll-mt-24 transition-shadow target:ring-2 target:ring-primary"
            :title="$t('components.auth.user_settings.documents.title')"
            :description="$t('components.auth.user_settings.documents.description')"
        >
            <DocumentsSettings />
        </UPageCard>
    </div>
</template>
