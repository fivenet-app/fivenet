import type { HowlOptions, Howl as HowlType } from 'howler';
import { storeToRefs } from 'pinia';
import { unref, type MaybeRef } from 'vue';
import { useSettingsStore } from '~/stores/settings';

export type SoundKeys =
    | 'notification'
    // Centrum Sounds
    | 'centrum.attention'
    | 'centrum.dispatchSOS'
    | 'centrum.dispatchAssigned'
    | 'centrum.dispatchCompleted';

export type NotificationSound = {
    value: 'default' | 'custom' | 'none';
    custom?: string;
};

type StoredSound = {
    blob: Blob;
    filename?: string;
    ext?: string;
};

const SoundsRegister: Record<SoundKeys, string> = {
    notification: '/sounds/notification.aac',
    'centrum.attention': '/sounds/centrum/attention.aac',
    'centrum.dispatchSOS': '/sounds/centrum/dispatch_sos.aac',
    'centrum.dispatchAssigned': '/sounds/centrum/dispatch_assigned.aac',
    'centrum.dispatchCompleted': '/sounds/centrum/dispatch_completed.aac',
} as const;

const DB_NAME = 'fivenet-sounds';
const STORE_NAME = 'fivenet_sounds_v0';

const disabledSound: NotificationSound = { value: 'none' } as const;

type CachedSound = {
    howl: HowlType;
    src: string;
    objectUrl?: string;
    format?: string | string[];
};

const soundCache = new Map<SoundKeys, CachedSound>();
const loadInFlight = new Map<SoundKeys, Promise<CachedSound | undefined>>();

let dbPromise: Promise<IDBDatabase> | undefined;
let howlerPromise: Promise<typeof import('howler')> | undefined;

function getDB(): Promise<IDBDatabase> {
    if (!dbPromise) {
        dbPromise = new Promise((resolve, reject) => {
            const req = indexedDB.open(DB_NAME, 1);

            req.onupgradeneeded = () => {
                const db = req.result;
                if (!db.objectStoreNames.contains(STORE_NAME)) {
                    db.createObjectStore(STORE_NAME);
                }
            };

            req.onsuccess = () => resolve(req.result);
            req.onerror = () => reject(req.error ?? new Error('Failed to open IndexedDB'));
        });
    }

    return dbPromise;
}

export async function putSound(name: SoundKeys, blob: Blob, filename?: string): Promise<void> {
    try {
        const db = await getDB();
        const transaction = db.transaction([STORE_NAME], 'readwrite');
        const record: StoredSound = {
            blob,
            filename,
            ext: deriveExtension(filename, blob.type),
        };
        transaction.objectStore(STORE_NAME).put(record, name);

        // Evict cached Howl so the new blob is used next play
        const cached = soundCache.get(name);
        cached?.howl.unload();
        if (cached?.objectUrl) URL.revokeObjectURL(cached.objectUrl);
        soundCache.delete(name);
    } catch (e) {
        console.warn('Failed to store custom sound', e);
    }
}

export async function deleteSound(name: SoundKeys): Promise<void> {
    try {
        const db = await getDB();
        const transaction = db.transaction([STORE_NAME], 'readwrite');
        transaction.objectStore(STORE_NAME).delete(name);

        const cached = soundCache.get(name);
        cached?.howl.unload();
        if (cached?.objectUrl) URL.revokeObjectURL(cached.objectUrl);
        soundCache.delete(name);
    } catch (e) {
        console.warn('Failed to delete custom sound', e);
    }
}

export type UseSoundsOptions = Omit<HowlOptions, 'src' | 'volume' | 'onload'> & {
    volume?: MaybeRef<number>;
    /**
     * Alias for Howler `rate` for compatibility with @vueuse/sound API
     */
    playbackRate?: number;
    onload?: () => void;
};

export type UseSoundsReturn = {
    play: () => void;
    stop: () => void;
};

type ResolvedSource = { src: string; objectUrl?: string; format?: string | string[] };

async function resolveSoundSource(name: SoundKeys, sound: NotificationSound): Promise<ResolvedSource> {
    if (sound.value !== 'custom') {
        return { src: SoundsRegister[name] };
    }

    try {
        const db = await getDB();
        const tx = db.transaction([STORE_NAME], 'readonly');
        const store = tx.objectStore(STORE_NAME);

        const stored: Blob | StoredSound | undefined = await new Promise((resolve, reject) => {
            const req = store.get(name);
            req.onsuccess = () => resolve(req.result as Blob | StoredSound | undefined);
            req.onerror = () => reject(req.error ?? new Error('Failed to read sound blob'));
        });

        const blob = stored instanceof Blob ? stored : stored?.blob;
        const ext = stored instanceof Blob ? undefined : stored?.ext;

        if (blob) {
            const objectUrl = URL.createObjectURL(blob);
            return { src: objectUrl, objectUrl: objectUrl, format: ext ? [ext] : undefined };
        }
    } catch (e) {
        console.warn('Falling back to default sound after IndexedDB failure', e);
    }

    return { src: SoundsRegister[name] };
}

function mergeStructuralOptions(options?: UseSoundsOptions): Partial<HowlOptions> {
    if (!options) return {};

    const { sprite, format, html5, xhr } = options;
    return { sprite, format, html5, xhr };
}

async function getHowler(): Promise<typeof import('howler')> {
    if (!howlerPromise) howlerPromise = import('howler');

    return howlerPromise;
}

async function ensureHowl(
    name: SoundKeys,
    sound: NotificationSound,
    structuralOptions?: Partial<HowlOptions>,
): Promise<CachedSound | undefined> {
    if (sound.value === 'none') return undefined;

    const cached = soundCache.get(name);
    const inFlight = loadInFlight.get(name);
    if (inFlight) return inFlight;

    const loadPromise = (async () => {
        const { src, objectUrl, format } = await resolveSoundSource(name, sound);

        if (cached && cached.src === src) {
            loadInFlight.delete(name);
            return cached;
        }

        if (cached) {
            cached.howl.unload();
            if (cached.objectUrl) URL.revokeObjectURL(cached.objectUrl);
        }

        const { Howl } = await getHowler();

        const howl = new Howl({
            src: [src],
            preload: true,
            format: structuralOptions?.format ?? format,
            ...structuralOptions,
        });

        const newCached: CachedSound = { howl, src, objectUrl, format: structuralOptions?.format ?? format };
        soundCache.set(name, newCached);
        loadInFlight.delete(name);
        return newCached;
    })().catch((e) => {
        console.warn(`Failed to load sound ${name}`, e);
        loadInFlight.delete(name);
        return undefined;
    });

    loadInFlight.set(name, loadPromise);
    return loadPromise;
}

export function useSounds(name: SoundKeys, options?: UseSoundsOptions): UseSoundsReturn {
    const settingsStore = useSettingsStore();
    const { audio } = storeToRefs(settingsStore);

    const structuralOptions = mergeStructuralOptions(options);

    const play = () => {
        const soundSetting = audio.value.sounds[name] ?? disabledSound;
        ensureHowl(name, soundSetting, structuralOptions).then((cached) => {
            if (!cached) return;

            const volume = clamp(unref(options?.volume ?? audio.value.notificationsVolume) ?? 1, 0, 1);
            const rate = options?.rate ?? options?.playbackRate;

            if (options?.onload) {
                if (cached.howl.state() === 'loaded') {
                    options.onload();
                } else {
                    cached.howl.once('load', options.onload);
                }
            }

            const id = cached.howl.play();

            if (rate !== undefined) cached.howl.rate(rate, id);
            if (options?.loop !== undefined) cached.howl.loop(options.loop, id);
            cached.howl.volume(volume, id);
        });
    };

    const stop = () => {
        const soundSetting = audio.value.sounds[name] ?? disabledSound;
        ensureHowl(name, soundSetting, structuralOptions).then((cached) => {
            cached?.howl.stop();
        });
    };

    return { play, stop };
}

function clamp(value: number, min: number, max: number): number {
    return Math.min(Math.max(value, min), max);
}

function deriveExtension(filename?: string, mimeType?: string): string | undefined {
    if (filename) {
        const lastDot = filename.lastIndexOf('.');
        if (lastDot > -1 && lastDot < filename.length - 1) {
            return filename.slice(lastDot + 1).toLowerCase();
        }
    }

    if (!mimeType) return undefined;

    if (mimeType === 'audio/mpeg') return 'mp3';
    if (mimeType === 'audio/ogg') return 'ogg';
    if (mimeType === 'audio/wav' || mimeType === 'audio/x-wav') return 'wav';
    if (mimeType === 'audio/aac') return 'aac';

    return undefined;
}
