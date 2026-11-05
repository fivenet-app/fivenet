import { PropertyStringAllowedValuesMode, type PropertyDefinition } from '~~/gen/ts/resources/properties/properties';

export type StringPropertyInputMode = 'text' | 'suggestions' | 'strict';

// Kept separate from the drawer so the UI behavior remains easy to exercise
// without mounting its API and overlay dependencies.
export function stringPropertyInputMode(definition: PropertyDefinition | undefined): StringPropertyInputMode {
    if ((definition?.valueConstraints?.stringAllowedValues.length ?? 0) === 0) return 'text';
    return definition?.valueConstraints?.stringAllowedValuesMode === PropertyStringAllowedValuesMode.SUGGEST
        ? 'suggestions'
        : 'strict';
}
