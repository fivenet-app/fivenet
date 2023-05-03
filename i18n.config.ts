import type { I18nOptions } from 'vue-i18n';
import type { NuxtApp } from 'nuxt/dist/app/index';

export default function (nuxt: NuxtApp) {
    return {
        legacy: false,
        locale: 'de',
        fallbackLocale: 'en',
        escapeParameterHtml: true,
        datetimeFormats: {
            en: {
                short: {
                    year: 'numeric',
                    month: 'short',
                    day: '2-digit',
                },
                long: {
                    year: 'numeric',
                    month: 'long',
                    day: '2-digit',
                    weekday: 'long',
                    hour: 'numeric',
                    minute: 'numeric',
                },
            },
            de: {
                short: {
                    year: 'numeric',
                    month: 'short',
                    day: '2-digit',
                },
                long: {
                    year: 'numeric',
                    month: '2-digit',
                    day: '2-digit',
                    weekday: 'long',
                    hour: '2-digit',
                    minute: '2-digit',
                },
            },
        },
    } as I18nOptions;
}
