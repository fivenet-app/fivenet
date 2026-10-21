import { describe, expect, it } from 'vitest';
import {
    PropertyDefinition,
    PropertyDefinitionSource,
    PropertyScope,
    PropertyValueType,
    type Property,
} from '~~/gen/ts/resources/properties/properties';
import { dynamicPropertiesForAction } from './useDynamicPropertyAction';

function property(key: string, source: PropertyDefinitionSource, enabled = true): Property {
    return {
        key,
        definition: PropertyDefinition.create({
            scope: PropertyScope.USER,
            key,
            valueType: PropertyValueType.STRING,
            source,
            enabled,
        }),
    };
}

describe('dynamicPropertiesForAction', () => {
    it('keeps only enabled custom properties with set permission', () => {
        const properties = [
            property('user.allowed', PropertyDefinitionSource.CUSTOM),
            property('user.denied', PropertyDefinitionSource.CUSTOM),
            property('user.disabled', PropertyDefinitionSource.CUSTOM, false),
            property('user.system', PropertyDefinitionSource.SYSTEM),
        ];

        const result = dynamicPropertiesForAction(properties, ({ key }) => key !== 'user.denied');

        expect(result.map(({ key }) => key)).toEqual(['user.allowed']);
    });
});
