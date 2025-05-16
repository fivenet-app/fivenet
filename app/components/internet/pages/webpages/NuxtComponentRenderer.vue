<script lang="ts" setup>
import { UCard, UContainer, ULandingCard, ULandingHero, ULandingSection, UPage, UPageGrid } from '#components';
import type { Component } from 'vue';
import type { ContentNode } from '~~/gen/ts/resources/internet/page';

const props = defineProps<{
    value: ContentNode;
}>();

const availableComponents: Record<string, Component> = {
    UPage: UPage,
    UPageGrid: UPageGrid,
    UContainer: UContainer,
    UCard: UCard,
    ULandingHero: ULandingHero,
    ULandingSection: ULandingSection,
    ULandingCard: ULandingCard,
};

const component = availableComponents[Object.keys(availableComponents).find((c) => c === props.value.tag) ?? ''];

defineOptions({
    inheritAttrs: false,
});
</script>

<template>
    <template v-if="value.text">
        {{ value.text }}
    </template>
    <component
        :is="value.tag === 'body' ? 'div' : (component ?? value.tag)"
        v-else
        :id="!!value.id ? value.id : undefined"
        :disabled="value.tag === 'input' ? true : undefined"
        v-bind="value.attrs"
    >
        <template v-for="(slot, slotIdx) in value.slots" :key="slotIdx" #[slot.tag]>
            <NuxtComponentRenderer :value="slot" />
        </template>
        <NuxtComponentRenderer v-for="(child, idx) in value.content" :key="idx" :value="child" />
    </component>
</template>
