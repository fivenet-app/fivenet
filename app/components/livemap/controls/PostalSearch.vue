<script lang="ts" setup>
import { useLivemapStore } from '~/stores/livemap';
import type { Postal } from '~/types/livemap';
import PostalSearchSelect from './PostalSearchSelect.vue';

defineOptions({
    inheritAttrs: false,
});

const livemapStore = useLivemapStore();
const { location } = storeToRefs(livemapStore);

const selectedPostal = ref<Postal | undefined>();

watch(selectedPostal, () => {
    if (!selectedPostal.value) return;

    location.value = selectedPostal.value;
});
</script>

<template>
    <PostalSearchSelect v-model="selectedPostal" class="w-full max-w-40" size="xs" v-bind="$attrs" />
</template>
