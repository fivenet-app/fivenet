<script lang="ts" setup>
import type { AlertAction } from '#ui/types';

const props = withDefaults(
    defineProps<{
        title?: string;
        message?: string;
        icon?: string;
        type?: string;
        padding?: string;
        actions?: AlertAction[];
        focus?: () => void | Promise<void>;
        retry?: () => Promise<unknown>;
    }>(),
    {
        title: undefined,
        message: undefined,
        icon: 'i-mdi-magnify',
        type: undefined,
        padding: 'p-4',
        actions: () => [],
        focus: undefined,
        retry: undefined,
    },
);

const { t } = useI18n();

const actions = computed(() =>
    props.actions.length > 0
        ? props.actions
        : [
              props.focus ? { label: t('common.search'), color: 'gray' } : undefined,
              props.retry ? { label: t('common.refresh'), icon: 'i-mdi-refresh', click: () => props.retry!() } : undefined,
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
    <UAlert
        variant="outline"
        :icon="icon"
        class="block w-full"
        :class="padding"
        :title="title"
        :description="message ?? $t('common.not_found', [type ?? $t('common.data')])"
        :actions="actions"
        @click="click()"
    />
</template>
