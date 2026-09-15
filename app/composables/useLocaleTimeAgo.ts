import { useI18n } from '#i18n';
import {
    formatTimeAgo,
    useNow,
    type FormatTimeAgoOptions,
    type UseTimeAgoMessages,
    type UseTimeAgoOptions,
    type UseTimeAgoUnitNamesDefault,
} from '@vueuse/core';
import { computed, type ComputedRef } from 'vue';

// Based on https://github.com/vueuse/vueuse/issues/1592#issuecomment-1341786344
// https://github.com/vueuse/vueuse/issues/1592#issuecomment-1381020982

function useLocaleTimeAgoMessages(): ComputedRef<UseTimeAgoMessages<UseTimeAgoUnitNamesDefault>> {
    const { t } = useI18n();

    return computed(() => ({
        justNow: t('common.time_ago.just-now'),
        past: (n) => (n.match(/\d/) ? t('common.time_ago.ago', [n]) : n),
        future: (n) => (n.match(/\d/) ? t('common.time_ago.in', [n]) : n),
        month: (n, past) =>
            n === 1
                ? past
                    ? t('common.time_ago.last-month')
                    : t('common.time_ago.next-month')
                : `${n} ${t('common.time_ago.month', n)}`,
        year: (n, past) =>
            n === 1
                ? past
                    ? t('common.time_ago.last-year')
                    : t('common.time_ago.next-year')
                : `${n} ${t('common.time_ago.year', n)}`,
        day: (n, past) =>
            n === 1
                ? past
                    ? t('common.time_ago.yesterday')
                    : t('common.time_ago.tomorrow')
                : `${n} ${t('common.time_ago.day', n)}`,
        week: (n, past) =>
            n === 1
                ? past
                    ? t('common.time_ago.last-week')
                    : t('common.time_ago.next-week')
                : `${n} ${t('common.time_ago.week', n)}`,
        hour: (n) => `${n} ${t('common.time_ago.hour', n)}`,
        minute: (n) => `${n} ${t('common.time_ago.minute', n)}`,
        second: (n) => `${n} ${t('common.time_ago.second', n)}`,
        invalid: t('common.unknown'),
    }));
}

export function useLocaleTimeAgo(date: Date, options?: UseTimeAgoOptions<false>): ComputedRef<string> {
    const { scheduler, ...formatOptions } = options ?? {};
    const now = scheduler ? useNow({ scheduler }) : useMinuteClock();
    const messages = useLocaleTimeAgoMessages();

    return computed(() =>
        formatTimeAgo(
            date,
            {
                ...formatOptions,
                messages: messages.value,
                fullDateFormatter: (value: Date) => value.toLocaleDateString(),
            },
            now.value,
        ),
    );
}

export function useLocaleTimeAgoFormatter(): (date: Date, options?: FormatTimeAgoOptions, now?: Date | number) => string {
    const messages = useLocaleTimeAgoMessages();

    return (date, options, now) =>
        formatTimeAgo(
            date,
            {
                ...options,
                messages: messages.value,
                fullDateFormatter: (value: Date) => value.toLocaleDateString(),
            },
            now,
        );
}
