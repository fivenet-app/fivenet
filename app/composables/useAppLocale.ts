import type { LocaleObject } from '@nuxtjs/i18n';
import { createSharedComposable } from '@vueuse/core';
import type { Locale } from 'vue-i18n';
import { useSettingsStore } from '~/stores/settings';

type SetLocaleOptions = {
    persist: boolean;
};

export const useAppLocale = createSharedComposable(() => {
    const logger = useLogger('⚙️ Settings');
    const settingsStore = useSettingsStore();
    const appConfig = useAppConfig();
    const { locale, locales, setLocale } = useI18n();
    type AppLocale = Parameters<typeof setLocale>[0];
    const defaultLocale: AppLocale = 'en';

    const localeCodes = computed<string[]>(() =>
        locales.value.flatMap((entry) => (typeof entry === 'string' ? [entry] : [entry.code])),
    );

    function isAppLocale(candidate: string): candidate is AppLocale {
        return localeCodes.value.includes(candidate);
    }

    const fallbackLocale = computed<AppLocale>(() => {
        const appDefault = appConfig.defaultLocale || defaultLocale;
        if (isAppLocale(appDefault)) return appDefault;
        if (isAppLocale(defaultLocale)) return defaultLocale;

        const firstLocale = localeCodes.value[0];
        if (firstLocale && isAppLocale(firstLocale)) return firstLocale;

        return defaultLocale;
    });

    const currentLocale = computed<Locale>({
        get: () => locale.value,
        set: (next) => {
            void setUserLocale(next);
        },
    });

    function normalizeLocaleCode(candidate?: string | null): AppLocale | undefined {
        if (!candidate) return undefined;

        if (isAppLocale(candidate)) return candidate;

        const short = candidate.split('-')[0];
        if (short && isAppLocale(short)) return short;

        return undefined;
    }

    async function applyLocale(
        candidate: string | undefined,
        options: SetLocaleOptions = { persist: true },
    ): Promise<AppLocale> {
        const nextLocale = normalizeLocaleCode(candidate) ?? fallbackLocale.value;
        logger.info('Setting locale to', nextLocale);

        await setLocale(nextLocale);
        if (options.persist) settingsStore.locale = nextLocale as Locale;

        return nextLocale;
    }

    async function initLocale(): Promise<AppLocale> {
        const { $appConfigPromise } = useNuxtApp();
        await $appConfigPromise;

        const preferred = settingsStore.locale ?? appConfig.defaultLocale;
        return await applyLocale(preferred, { persist: settingsStore.locale !== undefined });
    }

    async function setUserLocale(next: string): Promise<AppLocale> {
        return await applyLocale(next, { persist: true });
    }

    const availableLocales = computed<LocaleObject[]>(() =>
        locales.value.flatMap((entry) => (typeof entry === 'string' ? [] : [entry])),
    );

    return {
        availableLocales,
        currentLocale,
        fallbackLocale,
        initLocale,
        setUserLocale,
    };
});
