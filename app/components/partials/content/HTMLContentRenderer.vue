<script lang="ts" setup>
import type { JSONNode } from '~~/gen/ts/resources/common/content/content';

defineProps<{
    value: JSONNode;
}>();
</script>

<template>
    <template v-if="value.text">
        {{ value.text }}
    </template>
    <component
        :is="value.tag === 'body' ? 'div' : value.tag"
        v-else
        :id="value.id !== '' ? value.id : undefined"
        :disabled="value.tag === 'input' ? true : undefined"
        v-bind="value.attrs"
    >
        <HTMLContentRenderer v-for="(child, idx) in value.content" :key="idx" :value="child" />
    </component>
</template>
