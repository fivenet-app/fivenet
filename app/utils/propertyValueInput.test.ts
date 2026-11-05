import { describe, expect, it } from 'vitest';
import {
    PropertyDefinition,
    PropertyDefinitionSource,
    PropertyScope,
    PropertyStringAllowedValuesMode,
    PropertyValueType,
} from '~~/gen/ts/resources/properties/properties';
import { stringPropertyInputMode } from './propertyValueInput';

function definition(mode?: PropertyStringAllowedValuesMode, values: string[] = []): PropertyDefinition {
    return PropertyDefinition.create({
        scope: PropertyScope.USER,
        key: 'user.callsign',
        source: PropertyDefinitionSource.CUSTOM,
        valueType: PropertyValueType.STRING,
        valueConstraints: {
            stringAllowedValues: values,
            stringAllowedValuesMode: mode ?? PropertyStringAllowedValuesMode.UNSPECIFIED,
        },
    });
}

describe('stringPropertyInputMode', () => {
    it('uses a free-text input when no values are defined', () => {
        expect(stringPropertyInputMode(definition())).toBe('text');
    });

    it('uses a locked dropdown for strict lists, including legacy definitions', () => {
        expect(stringPropertyInputMode(definition(undefined, ['Alpha']))).toBe('strict');
        expect(stringPropertyInputMode(definition(PropertyStringAllowedValuesMode.STRICT, ['Alpha']))).toBe('strict');
    });

    it('uses an input with suggestions when the list is suggestive', () => {
        expect(stringPropertyInputMode(definition(PropertyStringAllowedValuesMode.SUGGEST, ['Alpha']))).toBe('suggestions');
    });
});
