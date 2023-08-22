<script lang="ts" setup>
const props = defineProps<{
    path: string;
}>();

const { data: license, pending, error } = await useFetch(props.path);
</script>

<template>
    <p v-if="pending || !license" class="text-black text-lg">Loading...</p>
    <p v-else-if="error" class="text-black text-lg">Error loading: {{ error?.message }}</p>
    <code v-else-if="license" v-text="license" class="max-w-full mt-2 p-4 block bg-white text-black whitespace-pre-line"></code>
    <p v-else class="text-black text-lg">Unknown Error while loading license.</p>
</template>
