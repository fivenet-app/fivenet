import { toDate } from '~/utils/time';
import { PropertyRenderer, PropertyValueFormat, type Property } from '~~/gen/ts/resources/properties/properties';

export function pinnedProperties(properties: readonly Property[] | undefined): Property[] {
    return (properties ?? []).filter(
        (property) => property.definition?.pinned && property.value?.value.oneofKind !== undefined,
    );
}

export function usePropertyRenderer() {
    const { t, te } = useI18n();
    const numberFormatter = useDisplayNumberFormatWithOptions({
        style: 'decimal',
        maximumFractionDigits: 2,
    });
    const currencyFormatter = useDisplayNumberFormat();
    const dateFormatter = useDateFormatterWithOptions(undefined, undefined, { dateStyle: 'medium' });
    const dateTimeFormatter = useDateFormatterWithOptions(undefined, undefined, { dateStyle: 'medium', timeStyle: 'short' });

    function propertyLabel(property: Property): string {
        if (property.definition?.label) return property.definition.label;
        const key = property.definition?.displayKey || property.key;
        return te(key) ? t(key) : property.key;
    }

    function propertyScalarValue(property: Property): string | number | boolean | undefined {
        const value = property.value?.value;
        if (!value || value.oneofKind === undefined) return undefined;
        if (value.oneofKind === 'stringValue') return value.stringValue;
        if (value.oneofKind === 'integerValue') return Number(value.integerValue);
        if (value.oneofKind === 'numberValue') return value.numberValue;
        if (value.oneofKind === 'boolValue') return value.boolValue;
        if (value.oneofKind === 'dateValue') {
            return `${value.dateValue.year}-${String(value.dateValue.month).padStart(2, '0')}-${String(value.dateValue.day).padStart(2, '0')}`;
        }
        return toDate(value.timestampValue)?.toLocaleString();
    }

    function renderProperty(property: Property): string {
        const value = propertyScalarValue(property);
        const rawValue = property.value?.value;
        if (
            rawValue?.oneofKind === 'dateValue' &&
            (property.definition?.renderer === PropertyRenderer.DATE ||
                property.definition?.valueFormat === PropertyValueFormat.DATE)
        ) {
            return dateFormatter.format(
                new Date(rawValue.dateValue.year, rawValue.dateValue.month - 1, rawValue.dateValue.day),
            );
        }
        if (
            rawValue?.oneofKind === 'timestampValue' &&
            (property.definition?.renderer === PropertyRenderer.DATETIME ||
                property.definition?.valueFormat === PropertyValueFormat.DATETIME)
        ) {
            return dateTimeFormatter.format(toDate(rawValue.timestampValue));
        }
        const rendererConfig = property.definition?.rendererConfig?.config;
        if (value === '') return t('common.na');
        if (
            typeof value === 'string' &&
            property.definition?.renderer === PropertyRenderer.STRING &&
            rendererConfig?.oneofKind === 'string'
        ) {
            const { prefix, suffix, truncate } = rendererConfig.string;
            const rendered = truncate && value.length > 64 ? `${value.slice(0, 61)}...` : value;
            return `${prefix}${rendered}${suffix}`;
        }
        if (property.definition?.renderer === PropertyRenderer.TRAFFIC_POINTS && typeof value === 'number') {
            return t('common.point', value);
        }
        if (property.key === 'user.sex' && typeof value === 'string') {
            const sex = value.toLowerCase();
            return t(`common.sex_mapping.${['m', 'f', 'n'].includes(sex) ? sex : 'n'}`);
        }
        if (property.definition?.renderer === PropertyRenderer.HEIGHT && typeof value === 'number') {
            return `${numberFormatter.format(value)}cm`;
        }
        if (typeof value === 'boolean') return value ? t('common.yes') : t('common.no');
        if (typeof value === 'number') {
            const numberConfig =
                rendererConfig?.oneofKind === 'integer'
                    ? rendererConfig.integer
                    : rendererConfig?.oneofKind === 'number'
                      ? rendererConfig.number
                      : undefined;
            if (numberConfig) {
                const formatter = new Intl.NumberFormat(numberFormatter.resolvedOptions().locale, {
                    style: 'decimal',
                    maximumFractionDigits: numberConfig.decimalPlaces ?? 2,
                });
                return `${numberConfig.prefix}${formatter.format(value)}${numberConfig.unit ? ` ${numberConfig.unit}` : ''}${numberConfig.suffix}`;
            }
            if (property.definition?.valueFormat === PropertyValueFormat.MONEY) {
                return currencyFormatter.format(value);
            }
            const decimalPlaces = property.definition?.decimalPlaces;
            const formatNumber = (number: number): string =>
                decimalPlaces === undefined
                    ? numberFormatter.format(number)
                    : new Intl.NumberFormat(numberFormatter.resolvedOptions().locale, {
                          style: 'decimal',
                          maximumFractionDigits: decimalPlaces,
                      }).format(number);
            if (property.definition?.valueFormat === PropertyValueFormat.PERCENTAGE) {
                return `${formatNumber(value)}%`;
            }
            return `${formatNumber(value)}${property.definition?.unit ? ` ${property.definition.unit}` : ''}`;
        }

        return value ?? t('common.na');
    }

    return { propertyLabel, propertyScalarValue, renderProperty };
}
