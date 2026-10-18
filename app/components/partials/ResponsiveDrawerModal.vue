<script lang="ts" setup>
import { useMediaQuery } from '@vueuse/core';

defineOptions({ inheritAttrs: false });

withDefaults(
    defineProps<{
        title?: string;
        close?: boolean;
        dismissible?: boolean;
        modal?: boolean;
        overlay?: boolean;
        ui?: Record<string, string>;
    }>(),
    {
        title: undefined,
        close: true,
        dismissible: true,
        modal: true,
        overlay: true,
        ui: undefined,
    },
);

const isDesktop = useMediaQuery('(min-width: 768px)');
</script>

<template>
    <UModal
        v-if="isDesktop"
        v-bind="$attrs"
        :title="title"
        :close="close"
        :dismissible="dismissible"
        :modal="modal"
        :overlay="overlay"
        :ui="ui"
    >
        <template v-for="(_, slot) in $slots" #[slot]="slotProps">
            <slot :name="slot" v-bind="slotProps ?? {}" />
        </template>
    </UModal>

    <UDrawer
        v-else
        v-bind="$attrs"
        :title="title"
        :close="close"
        :dismissible="dismissible"
        :modal="modal"
        :overlay="overlay"
        :ui="ui"
    >
        <template v-for="(_, slot) in $slots" #[slot]="slotProps">
            <slot :name="slot" v-bind="slotProps ?? {}" />
        </template>
    </UDrawer>
</template>
