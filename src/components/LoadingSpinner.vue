<script lang="ts" setup>
import { ref, computed, watch } from 'vue'
import { useLoading, ActiveLoader } from 'vue-loading-overlay'
import { useStore } from '../store/store';

const store = useStore();

const loadingState = computed(() => store.state.loader?.loading);

const loading = useLoading({
    isFullPage: true,
    canCancel: false,
    color: '#0c4a8c',
    loader: 'spinner',
    backgroundColor: '#3E3C3E',
});

const fullPage = ref(false);

let loader: undefined | ActiveLoader = undefined;

function showLoader(): ActiveLoader {
    return loading.show({
        isFullPage: fullPage.value,
    });
}

watch(loadingState, (newState: number | undefined) => {
    if (newState !== undefined) {
        if (newState > 0 && loader === undefined) {
            loader = showLoader();
        } else if (newState <= 0 && loader !== undefined) {
            loader.hide();
            loader = undefined;
        }
    }
});
</script>

<template></template>
