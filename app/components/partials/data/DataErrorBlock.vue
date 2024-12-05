<script lang="ts" setup>
import type { Error as CommonError } from '~~/gen/ts/resources/common/error';

const props = defineProps<{
    title?: string;
    message?: string;
    error?: Error;
    retry?: () => Promise<unknown>;
    retryMessage?: string;
}>();

const err = ref<CommonError | undefined>();

function setFromProps(): void {
    if (props.error) {
        err.value = parseError(props.error);
    } else {
        err.value = undefined;
    }
}

setFromProps();
watch(props, setFromProps);

const disabled = ref(true);

const { start } = useTimeoutFn(() => (disabled.value = false), 1250);
</script>

<template>
    <UAlert
        color="red"
        icon="i-mdi-close-circle"
        class="relative my-2 block w-full min-w-60"
        :title="
            err?.title
                ? $t(err.title.key, err.title.parameters)
                : (title ?? $t('components.partials.data_error_block.default_title'))
        "
        :description="
            err?.content
                ? $t(err.content.key, err.content.parameters)
                : (message ?? $t('components.partials.data_error_block.default_message'))
        "
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
