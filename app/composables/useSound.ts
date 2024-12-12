import { useState } from '#imports';
import { useSettingsStore } from '~/store/settings';

export type SoundBites = 'centrum/attention' | 'centrum/message-incoming' | 'centrum/morse-sos' | 'notification';

export type Sound = {
    id?: string;
    name: SoundBites;
    volume?: number;
    rate?: number;
};

export function useSound() {
    const sounds = useState<Sound[]>('sounds', () => []);

    const settingsStore = useSettingsStore();
    const { audio } = storeToRefs(settingsStore);

    function play(sound: Partial<Sound>) {
        const body = {
            id: new Date().getTime().toString(),
            ...sound,
        };
        if (!sound.volume) {
            sound.volume = audio.value.notificationsVolume;
        }

        const index = sounds.value.findIndex((s: Sound) => s.id === body.id);
        if (index === -1) {
            sounds.value.push(body as Sound);
        }

        return body;
    }

    function stop(id: string) {
        sounds.value = sounds.value.filter((s: Sound) => s.id !== id);
    }

    return {
        play,
        stop,
    };
}
