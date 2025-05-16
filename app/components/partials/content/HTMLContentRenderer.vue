<script lang="ts" setup>
import { NuxtImg } from '#components';
import type { JSONNode } from '~~/gen/ts/resources/common/content/content';

defineProps<{
    value: JSONNode;
}>();

const tagRemapping: Record<string, Component> = {
    img: NuxtImg,
};
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
    <br v-else-if="value.tag === 'br'" v-bind="value.attrs" />
    <component
        :is="value.tag === 'body' ? 'div' : (tagRemapping[value.tag] ?? value.tag)"
        v-else
        :id="!!value.id ? value.id : undefined"
        :disabled="value.tag === 'input' ? true : undefined"
        v-bind="value.attrs"
    >
        <HTMLContentRenderer v-for="(child, idx) in value.content" :key="idx" :value="child" />
        <UIcon v-if="value.tag === 'a' && !value.attrs.href?.startsWith('/')" class="ml-0.5 size-4" name="i-mdi-open-in-new" />
    </component>
</template>

<style scoped>
@supports (height: 1lh) {
    p:empty {
        height: 1lh;
    }
}
@supports not (height: 1lh) {
    p:empty::after {
        content: '\00A0';
    }
}
</style>
