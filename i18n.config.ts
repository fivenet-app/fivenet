import type { I18nOptions } from 'vue-i18n';
import type { NuxtApp } from 'nuxt/dist/app/index';

export default function (nuxt: NuxtApp) {
    return {
        legacy: false,
        locale: 'de',
        fallbackLocale: 'en',
        escapeParameterHtml: true,
    } as I18nOptions;
}
