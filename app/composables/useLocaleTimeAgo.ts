import { useI18n } from '#i18n';
import {
    formatTimeAgo,
    useTimeAgo,
    type FormatTimeAgoOptions,
    type UseTimeAgoMessages,
    type UseTimeAgoOptions,
    type UseTimeAgoUnitNamesDefault,
} from '@vueuse/core';
import type { ComputedRef } from 'vue';

// Based on https://github.com/vueuse/vueuse/issues/1592#issuecomment-1341786344
// https://github.com/vueuse/vueuse/issues/1592#issuecomment-1381020982

export function useLocaleTimeAgo(date: Date, options?: UseTimeAgoOptions<false>): ComputedRef<string> {
    const { t } = useI18n();

    const i18nMessages: UseTimeAgoMessages<UseTimeAgoUnitNamesDefault> = {
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
        options = { updateInterval: 30_000 };
    } else {
        options = { ...options };
    }

    return useTimeAgo(date, {
        ...options,
        messages: i18nMessages,
        fullDateFormatter: (value: Date) => value.toLocaleDateString(),
    });
}

export function useLocaleTimeAgoFormatter(): (
    date: Date,
    options?: FormatTimeAgoOptions,
    now?: Date | number,
) => string {
    const { t } = useI18n();

    const i18nMessages: UseTimeAgoMessages<UseTimeAgoUnitNamesDefault> = {
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

    return (date, options, now) =>
        formatTimeAgo(date, {
            ...options,
            messages: i18nMessages,
            fullDateFormatter: (value: Date) => value.toLocaleDateString(),
        }, now);
}
