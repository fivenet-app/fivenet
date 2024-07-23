import { useI18n } from '#i18n';
import type { ComputedRef } from '@vue/reactivity';
import { useTimeAgo, type UseTimeAgoMessages, type UseTimeAgoOptions, type UseTimeAgoUnitNamesDefault } from '@vueuse/core';

// Based on https://github.com/vueuse/vueuse/issues/1592#issuecomment-1341786344
// https://github.com/vueuse/vueuse/issues/1592#issuecomment-1381020982

export function useLocaleTimeAgo(date: Date, options?: UseTimeAgoOptions<false, any>): ComputedRef<string> {
    const { t } = useI18n();

    const I18N_MESSAGES: UseTimeAgoMessages<UseTimeAgoUnitNamesDefault> = {
        justNow: t('common.time_ago.just-now'),
        past: (n) => (n.match(/\d/) ? t('common.time_ago.ago', [n]) : n),
        future: (n) => (n.match(/\d/) ? t('common.time_ago.in', [n]) : n),
        month: (n, past) =>
            n === 1
                ? past
                    ? t('common.time_ago.last-month')
                    : t('common.time_ago.next-month')
                : `${n} ${t(`common.time_ago.month`, n)}`,
        year: (n, past) =>
            n === 1
                ? past
                    ? t('common.time_ago.last-year')
                    : t('common.time_ago.next-year')
                : `${n} ${t(`common.time_ago.year`, n)}`,
        day: (n, past) =>
            n === 1
                ? past
                    ? t('common.time_ago.yesterday')
                    : t('common.time_ago.tomorrow')
                : `${n} ${t(`common.time_ago.day`, n)}`,
        week: (n, past) =>
            n === 1
                ? past
                    ? t('common.time_ago.last-week')
                    : t('common.time_ago.next-week')
                : `${n} ${t(`common.time_ago.week`, n)}`,
        hour: (n) => `${n} ${t('common.time_ago.hour', n)}`,
        minute: (n) => `${n} ${t('common.time_ago.minute', n)}`,
        second: (n) => `${n} ${t(`common.time_ago.second`, n)}`,
        invalid: t('common.unknown'),
    };

    if (options === undefined) {
        options = {
            updateInterval: 30_000,
            messages: I18N_MESSAGES,
            fullDateFormatter: (date: Date) => date.toLocaleDateString(),
        };
    } else {
        options.messages = I18N_MESSAGES;
        options.fullDateFormatter = (date: Date) => date.toLocaleDateString();
    }

    return useTimeAgo(date, options);
}
