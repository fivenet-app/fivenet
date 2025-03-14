import type { SoundsPaths } from '#build/sounds/index.d';
import type { ComposableOptions } from '@vueuse/sound';
import { useSettingsStore } from '~/stores/settings';

export function useSounds(name: SoundsPaths, options?: ComposableOptions) {
    const settingsStore = useSettingsStore();
    const { audio } = storeToRefs(settingsStore);

    if (!options) {
        options = {};
    }
    options.volume = audio.value.notificationsVolume;

    const { play, stop } = useSound(name, options);
    return {
        play,
        stop,
    };
}
