<template>
    <UAlert
        :title="$t('system.streamer_mode.title')"
        :description="$t('system.streamer_mode.description')"
        icon="i-mdi-broadcast"
        color="warning"
        variant="subtle"
        :actions="[
            {
                label: $t('components.auth.user_settings.title'),
                to: '/user-settings',
                variant: 'outline',
                color: 'neutral',
                icon: 'i-mdi-account-cog-outline',
            },
        ]"
    />
</template>
