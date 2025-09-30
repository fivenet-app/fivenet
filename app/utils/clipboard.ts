import { useClipboard } from '@vueuse/core';
import { copyToClipboard } from '~/composables/nui';
import { useSettingsStore } from '~/stores/settings';

export async function copyToClipboardWrapper(text: string): Promise<void> {
    const settingsStore = useSettingsStore();
    const { nuiEnabled } = storeToRefs(settingsStore);

    if (nuiEnabled.value) {
        return await copyToClipboard(text);
    } else {
        const clipboard = useClipboard({
            legacy: true,
            read: false,
        });

        return await clipboard.copy(text);
    }
}
