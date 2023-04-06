import { defineStore } from 'pinia';

export const useLoaderStore = defineStore('loaderStore', () => {
    const loading = ref(0);

    function show(): void {
        loading.value++;
    }

    function hide(): void {
        if (loading.value > 0) {
            loading.value--;
        }
    }

    return {
        loading,
        show,
        hide,
    };
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useLoaderStore, import.meta.hot));
}
