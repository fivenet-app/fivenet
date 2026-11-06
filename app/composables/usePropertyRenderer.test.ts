import { describe, expect, it } from 'vitest';
import {
    PropertyDefinition,
    PropertyDefinitionSource,
    PropertyEmptyValueDisplay,
    PropertyScope,
    PropertyValueType,
    type Property,
} from '~~/gen/ts/resources/properties/properties';
import { pinnedProperties } from './usePropertyRenderer';

function property(
    key: string,
    pinned: boolean,
    hasValue: boolean,
    emptyValueDisplay = PropertyEmptyValueDisplay.HIDE,
): Property {
    return {
        key,
        definition: PropertyDefinition.create({
            scope: PropertyScope.USER,
            key,
            valueType: PropertyValueType.STRING,
            source: PropertyDefinitionSource.CUSTOM,
            enabled: true,
            pinned,
            emptyValueDisplay,
        }),
        value: hasValue ? { value: { oneofKind: 'stringValue', stringValue: key } } : undefined,
    };
}

describe('pinnedProperties', () => {
    it('returns pinned properties with a value or configured empty-value display', () => {
        const result = pinnedProperties([
            property('user.pinned', true, true),
            property('user.empty', true, false),
            property('user.empty_na', true, false, PropertyEmptyValueDisplay.NA),
            property('user.unpinned', false, true),
        ]);

        expect(result.map(({ key }) => key)).toEqual(['user.pinned', 'user.empty_na']);
    });
});
