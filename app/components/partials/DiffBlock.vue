<script lang="ts" setup>
const props = defineProps<{
    diff: string;
}>();

const lines = computed<{ line: string; class?: string }[]>(() =>
    props.diff.split('\n').map((l) => {
        if (l.startsWith(' ')) {
            return { line: l };
        } else if (l.startsWith('@@') || l.startsWith('---') || l.startsWith('+++')) {
            return { line: l, class: 'text-primary-500' };
        } else if (l.startsWith('-')) {
            return { line: l, class: 'text-red-500  bg-red-200/20' };
        } else if (l.startsWith('+')) {
            return { line: l, class: 'text-green-500 bg-green-200/20' };
        } else {
            return { line: l };
        }
    }),
);
</script>

<template>
    <pre
        class="break-words whitespace-pre-line"
    ><span v-for="(line, idx) in lines" :key="idx" :class="line.class">{{ line.line }}<br></span></pre>
</template>
