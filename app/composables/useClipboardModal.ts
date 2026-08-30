import ClipboardModal from '~/components/clipboard/modal/ClipboardModal.vue';

export function useClipboardModal() {
    const overlay = useOverlay();
    const clipboardModal = overlay.create(ClipboardModal);

    return {
        open: () => clipboardModal.open(),
    };
}
