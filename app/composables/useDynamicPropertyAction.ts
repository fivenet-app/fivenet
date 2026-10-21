import { PropertyDefinitionSource, type Property } from '~~/gen/ts/resources/properties/properties';

export function canActOnDynamicProperty(property: Property): boolean {
    if (property.definition?.source !== PropertyDefinitionSource.CUSTOM) return false;

    const value = property.value?.value;
    if (!value || value.oneofKind === undefined) return false;

    return true;
}

export function canSetDynamicProperty(property: Property): boolean {
    return property.definition?.source === PropertyDefinitionSource.CUSTOM && property.definition.enabled;
}

export function dynamicPropertiesForAction(
    properties: readonly Property[] | undefined,
    hasSetPermission: (property: Property) => boolean,
): Property[] {
    return (properties ?? []).filter((property) => canSetDynamicProperty(property) && hasSetPermission(property));
}
