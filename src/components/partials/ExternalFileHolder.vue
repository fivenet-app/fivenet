<script lang="ts" setup>
const props = defineProps<{
    path: string;
}>();

const { data: content, pending: loading, error, refresh } = useLazyFetch<string>(props.path);
</script>

<template>
    <div class="max-w-full px-4 text-black dark:text-white">
        <p v-if="loading || !content" class="text-lg">{{ $t('common.loading', $t('common.licenses', 1)) }}</p>
        <p v-else-if="error" class="text-lg">Error loading: {{ error?.message }}</p>
        <code v-else-if="content" class="block whitespace-pre-line" v-text="content"></code>
        <p v-else class="text-lg">
            Unknown Error while loading license. Please check your internet connection.
            <UButton :label="$t('common.retry')" @click="refresh" />
        </p>
    </div>
</template>
