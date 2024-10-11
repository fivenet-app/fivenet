<template>
    <UAlert
        :title="$t('system.streamer_mode.title')"
        :description="$t('system.streamer_mode.description')"
        icon="i-mdi-broadcast"
        color="amber"
        variant="subtle"
        :actions="[
            { label: $t('common.settings'), to: '/settings', variant: 'outline', color: 'black', icon: 'i-mdi-cog-outline' },
        ]"
    />
</template>
