<script lang="ts" setup>
import { ref, computed, watch } from 'vue'
import { useLoading, ActiveLoader } from 'vue-loading-overlay'
import { useLoaderStore } from '../store/loader';
import 'vue-loading-overlay/dist/css/index.css';

const store = useLoaderStore();

const loadingState = computed(() => store.loading);

const loading = useLoading({
    isFullPage: true,
    canCancel: false,
    color: '#0c4a8c',
    loader: 'spinner',
    backgroundColor: '#3E3C3E',
});

const fullPage = ref(false);

let loader: undefined | ActiveLoader = undefined;

async function showLoader(): Promise<ActiveLoader> {
    return loading.show({
        isFullPage: fullPage.value,
    });
}

watch(loadingState, async (newState: number | undefined) => {
    if (newState !== undefined) {
        if (newState > 0 && loader === undefined) {
            loader = await showLoader();
        } else if (newState <= 0 && loader !== undefined) {
            loader.hide();
            loader = undefined;
        }
    }
});
</script>

<template></template>
