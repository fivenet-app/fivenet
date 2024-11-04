import { defineI18nConfig } from '#i18n';

export default defineI18nConfig(() => ({
    legacy: false,
    locale: 'de',
    fallbackLocale: 'en',
    inheritLocale: 'en',
    escapeParameterHtml: true,
    warnHtmlInMessage: 'off',
    numberFormats: {
        en: {
            currency: {
                style: 'currency',
                currency: 'USD',
                notation: 'standard',
            },
            decimal: {
                style: 'decimal',
                minimumFractionDigits: 2,
                maximumFractionDigits: 2,
            },
            percent: {
                style: 'percent',
                useGrouping: false,
            },
        },
        de: {
            currency: {
                style: 'currency',
                currency: 'USD',
                notation: 'standard',
            },
            decimal: {
                style: 'decimal',
                minimumFractionDigits: 2,
                maximumFractionDigits: 2,
            },
            percent: {
                style: 'percent',
                useGrouping: false,
            },
        },
    },
    datetimeFormats: {
        en: {
            date: {
                year: 'numeric',
                month: '2-digit',
                day: '2-digit',
                weekday: 'short',
            },
            short: {
                year: 'numeric',
                month: '2-digit',
                day: '2-digit',
                hour: 'numeric',
                minute: 'numeric',
            },
            long: {
                year: 'numeric',
                month: 'short',
                day: '2-digit',
                weekday: 'short',
                hour: 'numeric',
                minute: 'numeric',
            },
            compact: {
                hour12: false,
                hour: 'numeric',
                minute: 'numeric',
                second: 'numeric',
            },
            time: {
                hour12: false,
                hour: 'numeric',
                minute: 'numeric',
            },
        },
        de: {
            date: {
                year: 'numeric',
                month: 'short',
                day: '2-digit',
            },
            short: {
                year: 'numeric',
                month: '2-digit',
                day: '2-digit',
                hour: 'numeric',
                minute: 'numeric',
            },
            long: {
                year: 'numeric',
                month: '2-digit',
                day: '2-digit',
                weekday: 'long',
                hour: '2-digit',
                minute: '2-digit',
            },
            compact: {
                hour12: false,
                hour: 'numeric',
                minute: 'numeric',
                second: 'numeric',
            },
            time: {
                hour12: false,
                hour: 'numeric',
                minute: 'numeric',
            },
        },
    },
}));
