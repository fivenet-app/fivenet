import { copyToClipboard } from '~/composables/nui';
import { useSettingsStore } from '~/stores/settings';

export async function copyToClipboardWrapper(text: string): Promise<void> {
    const settingsStore = useSettingsStore();
    const { nuiEnabled } = storeToRefs(settingsStore);

    if (nuiEnabled.value) return await copyToClipboard(text);

    // Based on https://stackoverflow.com/a/79250919
    try {
        if (!navigator.userAgent.includes('Firefox')) {
            await navigator.permissions.query({
                // @ts-expect-error clipboard-write is not in the types (yet?)
                name: 'clipboard-write',
            });
        }

        return await navigator.clipboard.writeText(text);
    } catch {
        // Based on VueUse useClipboard fallback implementation, licensed under MIT License
        const textArea = document.createElement('textarea');
        textArea.value = text;
        textArea.style.position = 'absolute';
        textArea.style.opacity = '0';
        document.body.appendChild(textArea);
        textArea.focus();
        textArea.select();

        try {
            document.execCommand('copy');
        } finally {
            document.body.removeChild(textArea);
        }
    }
}
