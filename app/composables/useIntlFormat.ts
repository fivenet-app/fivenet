function _useIntlNumberFormat(): Intl.NumberFormat {
    const { display } = useAppConfig();

    return new Intl.NumberFormat(display.intlLocale, {
        style: 'currency',
        currency: display.currencyName,
        trailingZeroDisplay: 'stripIfInteger',
    });
}

export const useIntlNumberFormat = createSharedComposable(_useIntlNumberFormat);

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
