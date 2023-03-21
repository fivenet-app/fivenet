<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useLoading, ActiveLoader } from 'vue-loading-overlay'
import { useStore } from '../store/store';

const store = useStore();

const loading = useLoading({
    isFullPage: true,
    canCancel: false,
    color: '#0c4a8c',
    loader: 'spinner',
    backgroundColor: '#3E3C3E',
});

const fullPage = ref(false);
const loadingState = computed(() => store.state.loader?.loading);

let loader: undefined | ActiveLoader = undefined;

function showLoader(): ActiveLoader {
    return loading.show({
        isFullPage: fullPage.value,
    });
}

watch(loadingState, (newState) => {
    console.log("LOADING STATE: " + newState);
    if (newState) {
        loader = showLoader();
    } else {
        loader?.hide();
    }
});
</script>

<template></template>
