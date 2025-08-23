<script lang="ts" setup>
defineProps<{
    hideButton?: boolean;
}>();

const w = window;
</script>

<template>
    <UAlert
        :title="$t('system.not_supported_on_tablet.title')"
        :description="$t('system.not_supported_on_tablet.description')"
        icon="i-mdi-info-circle"
        variant="subtle"
        color="warning"
        :ui="{ inner: 'w-auto m-w-0' }"
        :actions="
            !hideButton
                ? [
                      {
                          label: $t('system.not_supported_on_tablet.button'),
                          icon: 'i-mdi-link-variant',
                          color: 'warning',
                          variant: 'solid',
                          onClick: () => openURLInWindow(w.location.href),
                      },
                  ]
                : undefined
        "
    />
</template>
