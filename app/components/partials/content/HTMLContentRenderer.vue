<script lang="ts" setup>
import type { JSONNode } from '~~/gen/ts/resources/common/content/content';

defineProps<{
    value: JSONNode;
}>();
</script>

<template>
    <component
        :is="value.tag === 'body' ? 'div' : value.tag"
        :id="value.id !== '' ? value.id : undefined"
        v-bind="value.attributes"
    >
        <template v-if="value.text !== ''">{{ value.text }}</template>
        <HTMLContentRenderer v-for="(child, idx) in value.children" :key="idx" :value="child" />
    </component>
</template>
