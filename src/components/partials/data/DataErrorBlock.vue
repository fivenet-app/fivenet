<script lang="ts" setup>
defineProps<{
    title?: string;
    message?: string;
    retry?: () => Promise<any>;
    retryMessage?: string;
}>();

const disabled = ref(true);

const { start } = useTimeoutFn(() => (disabled.value = false), 1250);
</script>

<template>
    <UAlert
        color="red"
        icon="i-mdi-close-circle"
        class="relative block w-full min-w-60"
        :title="title ?? $t('components.partials.data_error_block.default_title')"
        :description="message ?? $t('components.partials.data_error_block.default_message')"
        :actions="
            retry !== undefined
                ? [
                      {
                          variant: 'soft',
                          color: 'red',
                          label: retryMessage ?? $t('common.retry'),
                          disabled: disabled,
                          click: async () => {
                              start();
                              retry && retry();
                          },
                      },
                  ]
                : []
        "
    />
</template>
