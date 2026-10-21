import { describe, expect, it } from 'vitest';
import {
    PropertyDefinition,
    PropertyDefinitionSource,
    PropertyScope,
    PropertyValueType,
    type Property,
} from '~~/gen/ts/resources/properties/properties';
import { pinnedProperties } from './usePropertyRenderer';

function property(key: string, pinned: boolean, hasValue: boolean): Property {
    return {
        key,
        definition: PropertyDefinition.create({
            scope: PropertyScope.USER,
            key,
            valueType: PropertyValueType.STRING,
            source: PropertyDefinitionSource.CUSTOM,
            enabled: true,
            pinned,
        }),
        value: hasValue ? { value: { oneofKind: 'stringValue', stringValue: key } } : undefined,
    };
}

describe('pinnedProperties', () => {
    it('returns only pinned properties with values', () => {
        const result = pinnedProperties([
            property('user.pinned', true, true),
            property('user.empty', true, false),
            property('user.unpinned', false, true),
        ]);

        expect(result.map(({ key }) => key)).toEqual(['user.pinned']);
    });
});
