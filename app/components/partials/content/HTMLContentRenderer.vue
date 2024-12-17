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
    <UCheckbox
        v-else-if="value.tag === 'input' && value.attrs.type === 'checkbox'"
        disabled
        :model-value="!!value.attrs.checked"
        :ui="{
            wrapper: '',
            container: '',
            base: 'h-5 w-5',
        }"
    />
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
