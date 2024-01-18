<script lang="ts" setup>
const props = defineProps<{
    path: string;
}>();

const { data: license, pending, error } = await useFetch(props.path);
</script>

<template>
    <p v-if="pending || !license" class="text-lg text-black">Loading...</p>
    <p v-else-if="error" class="text-lg text-black">Error loading: {{ error?.message }}</p>
    <code
        v-else-if="license"
        v-text="license"
        class="mt-2 block max-w-full whitespace-pre-line bg-neutral p-4 text-black"
    ></code>
    <p v-else class="text-lg text-black">Unknown Error while loading license.</p>
</template>
