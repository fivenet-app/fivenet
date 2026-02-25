<script lang="ts" setup>
import type { Error as CommonError } from '~~/gen/ts/resources/common/error';

const props = withDefaults(
    defineProps<{
        title?: string;
        message?: string;
        error?: Error;
        retry?: () => Promise<unknown>;
        retryMessage?: string;
        close?: () => void;
        padded?: boolean;
    }>(),
    {
        title: undefined,
        message: undefined,
        error: undefined,
        retry: undefined,
        retryMessage: undefined,
        close: undefined,
        padded: false,
    },
);

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
    <div :class="padded ? 'm-2' : ''">
        <UAlert
            class="relative w-full min-w-50"
            color="error"
            icon="i-mdi-close-circle"
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
                              variant: 'solid',
                              color: 'neutral',
                              label: retryMessage ?? $t('common.retry'),
                              disabled: disabled,
                              onClick: async () => {
                                  start();
                                  retry && retry();
                              },
                          },
                      ]
                    : []
            "
            :close-button="
                close !== undefined
                    ? {
                          icon: 'i-mdi-window-close',
                          color: 'neutral',
                          variant: 'link',
                      }
                    : undefined
            "
            @close="close && close()"
        />
    </div>
</template>
