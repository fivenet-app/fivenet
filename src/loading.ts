import { useLoading, ActiveLoader, Props } from 'vue-loading-overlay';

const $loading = useLoading({
    isFullPage: true,
    canCancel: false,
    color: '#0c4a8c',
    loader: 'spinner',
    backgroundColor: '#ffffff',
});

let loader: ActiveLoader | null = null;

export function hideLoader() {
    // destroy previous
    if (loader) {
        loader.hide();
        loader = null;
    }
}

export function showLoader() {
    hideLoader();
    loader = $loading.show();
}
