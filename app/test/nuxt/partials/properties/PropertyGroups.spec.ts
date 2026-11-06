import { mount } from '@vue/test-utils';
import { describe, expect, it, vi } from 'vitest';
import { defineComponent } from 'vue';
import PropertyGroups from '~/components/partials/properties/PropertyGroups.vue';
import {
    PropertyDefinition,
    PropertyDisplayMode,
    PropertyEmptyValueDisplay,
    PropertyScope,
    PropertyValueType,
    type Property,
} from '~~/gen/ts/resources/properties/properties';

vi.mock('vue-i18n', () => ({
    useI18n: () => ({ t: (key: string) => key, te: () => false }),
}));

const UCardStub = defineComponent({ name: 'UCard', template: '<section><slot name="header" /><slot /></section>' });
const PropertyRenderStub = defineComponent({
    name: 'PropertyRender',
    props: { property: { type: Object, required: true } },
    template: '<span>{{ property.key }}</span>',
});

function property(key: string, emptyValueDisplay: PropertyEmptyValueDisplay): Property {
    return {
        key,
        definition: PropertyDefinition.create({
            scope: PropertyScope.USER,
            key,
            valueType: PropertyValueType.STRING,
            enabled: true,
            displayMode: PropertyDisplayMode.DEFAULT,
            emptyValueDisplay,
        }),
    };
}

describe('PropertyGroups', () => {
    it('shows unset properties configured to display N/A and hides the others', () => {
        const wrapper = mount(PropertyGroups, {
            props: {
                properties: [
                    property('user.show_empty', PropertyEmptyValueDisplay.NA),
                    property('user.hide_empty', PropertyEmptyValueDisplay.HIDE),
                ],
            },
            global: {
                stubs: { UCard: UCardStub, PropertyRender: PropertyRenderStub },
            },
        });

        expect(wrapper.text()).toContain('user.show_empty');
        expect(wrapper.text()).not.toContain('user.hide_empty');
    });
});
