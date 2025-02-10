import { useClipboard } from '@vueuse/core';
import { copyToClipboard } from '~/composables/nui';
import { useSettingsStore } from '~/store/settings';

export function copyToClipboardWrapper(text: string): Promise<void> {
    const settingsStore = useSettingsStore();
    const { nuiEnabled } = storeToRefs(settingsStore);

    if (nuiEnabled.value) {
        return copyToClipboard(text);
    } else {
        const clipboard = useClipboard();
        return clipboard.copy(text);
    }
}
