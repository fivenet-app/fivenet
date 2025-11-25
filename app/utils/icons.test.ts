import { describe, expect, it } from 'vitest';
import { convertComponentIconNameToDynamic, convertDynamicIconNameToComponent } from './icons';

describe('convertDynamicIconNameToComponent', () => {
    it('should convert dynamic icon name to component name', () => {
        const dynamicName = 'i-mdi-home-outline';
        const result = convertDynamicIconNameToComponent(dynamicName);
        expect(result).toBe('HomeOutlineIcon');
    });

    it('should return the input if it does not start with i-mdi-', () => {
        const dynamicName = 'custom-icon';
        const result = convertDynamicIconNameToComponent(dynamicName);
        expect(result).toBe('custom-icon');
    });
});

describe('convertComponentIconNameToDynamic', () => {
    it('should convert component name to dynamic icon name', () => {
        const componentName = 'HomeOutlineIcon';
        const result = convertComponentIconNameToDynamic(componentName);
        expect(result).toBe('i-mdi-home-outline');
    });

    it('should return the input if it already starts with i-mdi-', () => {
        const dynamicName = 'i-mdi-home-outline';
        const result = convertComponentIconNameToDynamic(dynamicName);
        expect(result).toBe('i-mdi-home-outline');
    });
});
