import { durationToSeconds } from '~/utils/duration';
import type { Duration } from '~~/gen/ts/google/protobuf/duration';

type DurationUnit = 'second' | 'minute' | 'hour' | 'day';

function _useIntlNumberFormat(opts?: Intl.NumberFormatOptions): Intl.NumberFormat {
    const { display } = useAppConfig();

    return new Intl.NumberFormat(display.intlLocale, {
        style: 'currency',
        currency: display.currencyName,
        trailingZeroDisplay: 'stripIfInteger',
        maximumFractionDigits: 2,
        ...opts,
    });
}

export const useIntlNumberFormat = createSharedComposable(() => _useIntlNumberFormat());
export const useIntlNumberFormatWithOptions = _useIntlNumberFormat;

function _useDateFormatter(
    dateStyle?: Intl.DateTimeFormatOptions['dateStyle'],
    timeStyle?: Intl.DateTimeFormatOptions['timeStyle'],
    opts?: Intl.DateTimeFormatOptions,
): Intl.DateTimeFormat {
    const { display } = useAppConfig();
    const { locale } = useI18n();

    return new Intl.DateTimeFormat(display.intlLocale ?? locale.value, {
        dateStyle: dateStyle,
        timeStyle: timeStyle,
        ...opts,
    });
}

export const useDateFormatter = createSharedComposable(_useDateFormatter);
export const useDateFormatterWithOptions = _useDateFormatter;

export const useDetentionTimeFormatter = createSharedComposable(() => {
    const { quickButtons } = useAppConfig();
    const { t } = useI18n();
    const { format } = useIntlNumberFormatWithOptions({
        style: 'decimal',
    });

    return (months: number) => {
        if (months > 1 || months === 0) {
            return `${format(months)} ${quickButtons.penaltyCalculator?.detentionTimeUnit?.plural ?? t('common.month', 2)}`;
        }
        return `${format(months)} ${quickButtons.penaltyCalculator?.detentionTimeUnit?.singular ?? t('common.month', 1)}`;
    };
});

export const useDurationFormatter = createSharedComposable(() => {
    const { t } = useI18n();
    const { format } = useIntlNumberFormatWithOptions({
        style: 'decimal',
        maximumFractionDigits: 2,
    });

    const secondsPerUnit: Record<DurationUnit, number> = {
        second: 1,
        minute: 60,
        hour: 60 * 60,
        day: 24 * 60 * 60,
    };

    const orderedUnits: DurationUnit[] = ['day', 'hour', 'minute', 'second'];

    return (duration?: Duration, unit?: DurationUnit): string => {
        if (!duration) {
            return format(0);
        }

        const totalSeconds = Math.max(0, durationToSeconds(duration));
        const resolvedUnit = unit ?? orderedUnits.find((current) => totalSeconds >= secondsPerUnit[current]) ?? 'second';

        const value = totalSeconds / secondsPerUnit[resolvedUnit];
        const pluralization = value === 1 ? 1 : 2;

        return `${format(value)} ${t(`common.time_ago.${resolvedUnit}`, pluralization)}`;
    };
});
