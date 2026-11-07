import { mount } from '@vue/test-utils';
import { describe, expect, it, vi } from 'vitest';
import { defineComponent } from 'vue';
import ProfilePropertyBadges from '~/components/partials/properties/ProfilePropertyBadges.vue';
import {
    PropertyDefinition,
    PropertyDisplayMode,
    PropertyEmptyValueDisplay,
    PropertyScope,
    PropertyValueType,
    type Property,
} from '~~/gen/ts/resources/properties/properties';
import { UserProps } from '~~/gen/ts/resources/users/props/props';
import { User } from '~~/gen/ts/resources/users/user';

vi.mock('vue-i18n', () => ({
    useI18n: () => ({ t: (key: string) => key, te: () => false }),
}));

const UBadgeStub = defineComponent({ name: 'UBadge', template: '<span><slot /></span>' });
const UIconStub = defineComponent({ name: 'UIcon', template: '<i />' });
const WantedBadgeStub = defineComponent({
    name: 'WantedBadge',
    props: { wanted: { type: Boolean, required: true } },
    template: '<span class="wanted-badge" :data-wanted="wanted" />',
});

function property(key: string, value?: string): Property {
    return {
        key,
        definition: PropertyDefinition.create({
            scope: PropertyScope.USER,
            key,
            valueType: key === 'user.wanted' ? PropertyValueType.BOOL : PropertyValueType.STRING,
            enabled: true,
            displayMode: PropertyDisplayMode.PROFILE_BADGE,
            emptyValueDisplay: PropertyEmptyValueDisplay.NA,
        }),
        value: value === undefined ? undefined : { value: { oneofKind: 'stringValue', stringValue: value } },
    };
}

describe('ProfilePropertyBadges', () => {
    const stubs = { UBadge: UBadgeStub, UIcon: UIconStub, WantedBadge: WantedBadgeStub };

    it('shows an empty wanted property as the not-wanted state', () => {
        const wrapper = mount(ProfilePropertyBadges, {
            props: { properties: [property('user.wanted')] },
            global: { stubs },
        });

        expect(wrapper.find('.wanted-badge').attributes('data-wanted')).toBe('false');
    });

    it('uses the property job and grade labels and marks a property override', () => {
        const wrapper = mount(ProfilePropertyBadges, {
            props: {
                properties: [property('user.job', 'police')],
                user: User.create({
                    jobLabel: 'Unemployed',
                    jobGradeLabel: 'Unemployed',
                    props: UserProps.create({
                        job: { name: 'police', label: 'Police', grades: [] },
                        jobGrade: { grade: 2, label: 'Sergeant' },
                    }),
                }),
            },
            global: { stubs },
        });

        expect(wrapper.text()).toContain('Police (Sergeant)');
        expect(wrapper.text()).toContain('*');
    });

    it('does not show a grade for the configured unemployed job', () => {
        const wrapper = mount(ProfilePropertyBadges, {
            props: {
                properties: [property('user.job')],
                user: User.create({
                    job: 'unemployed',
                    jobLabel: 'Unemployed',
                    jobGrade: 0,
                    jobGradeLabel: 'Unemployed',
                }),
            },
            global: { stubs },
        });

        expect(wrapper.text()).toContain('Unemployed');
        expect(wrapper.text()).not.toContain('(Unemployed)');
    });
});
