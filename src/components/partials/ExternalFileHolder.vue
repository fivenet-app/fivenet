<script lang="ts" setup>
const props = defineProps<{
    path: string;
}>();

const { data: license, pending, error } = await useFetch(props.path);
</script>

<template>
    <div class="max-w-full text-black dark:text-white">
        <p v-if="pending || !license" class="text-lg">{{ $t('common.loading', $t('common.licenses', 1)) }}</p>
        <p v-else-if="error" class="text-lg">Error loading: {{ error?.message }}</p>
        <code v-else-if="license" class="block whitespace-pre-line px-4" v-text="license"></code>
        <p v-else class="text-lg text-black">Unknown Error while loading license. Please check your internet connection.</p>
    </div>
</template>
