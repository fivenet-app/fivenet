function _useIntlNumberFormat(opts?: Intl.NumberFormatOptions): Intl.NumberFormat {
    const { display } = useAppConfig();

    return new Intl.NumberFormat(display.intlLocale, {
        style: 'currency',
        currency: display.currencyName,
        trailingZeroDisplay: 'stripIfInteger',
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

    return (months: number) => {
        if (months > 1 || months === 0) {
            return `${months} ${quickButtons.penaltyCalculator?.detentionTimeUnit?.plural ?? t('common.month', 2)}`;
        }
        return `${months} ${quickButtons.penaltyCalculator?.detentionTimeUnit?.singular ?? t('common.month', 1)}`;
    };
});
