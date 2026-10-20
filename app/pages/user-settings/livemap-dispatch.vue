<script lang="ts" setup>
import DispatchLayoutSettings from '~/components/user-settings/preferences/DispatchLayoutSettings.vue';
import LivemapSettings from '~/components/user-settings/preferences/LivemapSettings.vue';
</script>

<template>
    <div class="space-y-4">
        <UPageCard
            id="user-settings-livemap"
            class="scroll-mt-24 transition-shadow target:ring-2 target:ring-primary"
            :title="$t('components.auth.user_settings.livemap_dispatch.livemap_title')"
            :description="$t('components.auth.user_settings.livemap_dispatch.livemap_description')"
        >
            <LivemapSettings />
        </UPageCard>

        <UPageCard
            id="user-settings-dispatch"
            class="scroll-mt-24 transition-shadow target:ring-2 target:ring-primary"
            :title="$t('components.auth.user_settings.livemap_dispatch.dispatch_title')"
            :description="$t('components.auth.user_settings.livemap_dispatch.dispatch_description')"
        >
            <DispatchLayoutSettings auto-save />
        </UPageCard>
    </div>
</template>
