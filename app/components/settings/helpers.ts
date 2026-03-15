import type { BadgeProps } from '@nuxt/ui';
import { EventAction, EventResult } from '~~/gen/ts/resources/audit/audit';

export function eventActionToBadgeColor(et: EventAction): BadgeProps['color'] {
    switch (et) {
        case EventAction.VIEWED:
            return 'info';
        case EventAction.CREATED:
            return 'success';
        case EventAction.UPDATED:
            return 'warning';
        default:
            return 'error';
    }
}

export function eventResultToBadgeColor(er: EventResult): BadgeProps['color'] {
    switch (er) {
        case EventResult.SUCCEEDED:
            return 'success';
        case EventResult.FAILED:
            return 'warning';
        case EventResult.ERRORED:
            return 'error';
        default:
            return 'info';
    }
}

type IntlLocaleOption = {
    code: string;
    name: string;
    icon?: string;
};

type CurrencyOption = {
    code: string;
    name: string;
    flag?: string;
};

export const intlLocales: IntlLocaleOption[] = [
    { code: 'en-US', name: 'English (United States)', icon: 'i-mdi-translate' },
    { code: 'en-GB', name: 'English (United Kingdom)', icon: 'i-mdi-translate' },
    { code: 'de-DE', name: 'Deutsch (Deutschland)', icon: 'i-mdi-translate' },
    { code: 'fr-FR', name: 'Francais (France)', icon: 'i-mdi-translate' },
    { code: 'es-ES', name: 'Espanol (Espana)', icon: 'i-mdi-translate' },
    { code: 'it-IT', name: 'Italiano (Italia)', icon: 'i-mdi-translate' },
    { code: 'nl-NL', name: 'Nederlands (Nederland)', icon: 'i-mdi-translate' },
    { code: 'pt-BR', name: 'Portugues (Brasil)', icon: 'i-mdi-translate' },
    { code: 'tr-TR', name: 'Turkce (Turkiye)', icon: 'i-mdi-translate' },
    { code: 'pl-PL', name: 'Polski (Polska)', icon: 'i-mdi-translate' },
    { code: 'cs-CZ', name: 'Cestina (Cesko)', icon: 'i-mdi-translate' },
    { code: 'da-DK', name: 'Dansk (Danmark)', icon: 'i-mdi-translate' },
    { code: 'sv-SE', name: 'Svenska (Sverige)', icon: 'i-mdi-translate' },
    { code: 'nb-NO', name: 'Norsk bokmal (Norge)', icon: 'i-mdi-translate' },
    { code: 'fi-FI', name: 'Suomi (Suomi)', icon: 'i-mdi-translate' },
    { code: 'ru-RU', name: 'Russkiy (Rossiya)', icon: 'i-mdi-translate' },
    { code: 'uk-UA', name: 'Ukrayinska (Ukrayina)', icon: 'i-mdi-translate' },
    { code: 'ar-SA', name: 'al-Arabiyya (as-Saudiyya)', icon: 'i-mdi-translate' },
    { code: 'hi-IN', name: 'Hindi (Bharat)', icon: 'i-mdi-translate' },
    { code: 'ja-JP', name: 'Nihongo (Nihon)', icon: 'i-mdi-translate' },
    { code: 'ko-KR', name: 'Hangugeo (Daehan Minguk)', icon: 'i-mdi-translate' },
    { code: 'zh-CN', name: 'Zhongwen (Zhongguo)', icon: 'i-mdi-translate' },
    { code: 'zh-TW', name: 'Zhongwen (Taiwan)', icon: 'i-mdi-translate' },
];

export const currencies: CurrencyOption[] = [
    { code: 'USD', name: 'US Dollar (USD)', flag: 'i-mdi-cash' },
    { code: 'EUR', name: 'Euro (EUR)', flag: 'i-mdi-cash' },
    { code: 'GBP', name: 'British Pound (GBP)', flag: 'i-mdi-cash' },
    { code: 'JPY', name: 'Japanese Yen (JPY)', flag: 'i-mdi-cash' },
    { code: 'CNY', name: 'Chinese Yuan (CNY)', flag: 'i-mdi-cash' },
    { code: 'AUD', name: 'Australian Dollar (AUD)', flag: 'i-mdi-cash' },
    { code: 'CAD', name: 'Canadian Dollar (CAD)', flag: 'i-mdi-cash' },
    { code: 'CHF', name: 'Swiss Franc (CHF)', flag: 'i-mdi-cash' },
    { code: 'SEK', name: 'Swedish Krona (SEK)', flag: 'i-mdi-cash' },
    { code: 'NOK', name: 'Norwegian Krone (NOK)', flag: 'i-mdi-cash' },
    { code: 'DKK', name: 'Danish Krone (DKK)', flag: 'i-mdi-cash' },
    { code: 'PLN', name: 'Polish Zloty (PLN)', flag: 'i-mdi-cash' },
    { code: 'CZK', name: 'Czech Koruna (CZK)', flag: 'i-mdi-cash' },
    { code: 'HUF', name: 'Hungarian Forint (HUF)', flag: 'i-mdi-cash' },
    { code: 'RON', name: 'Romanian Leu (RON)', flag: 'i-mdi-cash' },
    { code: 'TRY', name: 'Turkish Lira (TRY)', flag: 'i-mdi-cash' },
    { code: 'BRL', name: 'Brazilian Real (BRL)', flag: 'i-mdi-cash' },
    { code: 'MXN', name: 'Mexican Peso (MXN)', flag: 'i-mdi-cash' },
    { code: 'INR', name: 'Indian Rupee (INR)', flag: 'i-mdi-cash' },
    { code: 'RUB', name: 'Russian Ruble (RUB)', flag: 'i-mdi-cash' },
    { code: 'UAH', name: 'Ukrainian Hryvnia (UAH)', flag: 'i-mdi-cash' },
    { code: 'AED', name: 'UAE Dirham (AED)', flag: 'i-mdi-cash' },
    { code: 'SAR', name: 'Saudi Riyal (SAR)', flag: 'i-mdi-cash' },
    { code: 'ZAR', name: 'South African Rand (ZAR)', flag: 'i-mdi-cash' },
];
