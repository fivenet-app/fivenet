<script lang="ts" setup>
import type { JSONNode } from '~~/gen/ts/resources/common/content/content';

const props = defineProps<{
    value: JSONNode;
}>();

// TODO load custom components
</script>

<template>
    <component
        :is="props.value.tag === 'body' ? 'div' : props.value.tag"
        :id="props.value.id !== '' ? props.value.id : undefined"
        v-bind="props.value.attributes"
    >
        {{ props.value.text }}
        <HTMLContentRenderer v-for="(child, idx) in props.value.children" :key="idx" :value="child" />
    </component>
</template>
