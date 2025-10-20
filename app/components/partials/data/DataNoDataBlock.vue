<script lang="ts" setup>
import type { ButtonProps } from '@nuxt/ui';

const props = withDefaults(
    defineProps<{
        title?: string;
        message?: string;
        icon?: string;
        type?: string;
        actions?: ButtonProps[];
        focus?: () => void | Promise<void>;
        retry?: () => Promise<unknown>;
        padded?: boolean;
    }>(),
    {
        title: undefined,
        message: undefined,
        icon: 'i-mdi-magnify',
        type: undefined,
        actions: () => [],
        focus: undefined,
        retry: undefined,
        padded: true,
    },
);

const { t } = useI18n();

const actions = computed<ButtonProps[]>(() =>
    props.actions.length > 0
        ? props.actions
        : [
              props.focus ? { label: t('common.search'), icon: 'i-mdi-search', onClick: () => props.focus!() } : undefined,
              props.retry ? { label: t('common.refresh'), icon: 'i-mdi-refresh', onClick: () => props.retry!() } : undefined,
          ].flatMap((item) => (item !== undefined ? [item] : [])),
);

async function click() {
    if (props.retry) {
        props.retry();
    } else if (props.focus) {
        props.focus();
    }
}
</script>

<template>
    <div :class="padded ? 'm-2' : ''">
        <UAlert
            :icon="icon"
            variant="outline"
            :title="title"
            :description="message ?? $t('common.not_found', [type ?? $t('common.data')])"
            :actions="actions"
            @click="click()"
        />
    </div>
</template>
