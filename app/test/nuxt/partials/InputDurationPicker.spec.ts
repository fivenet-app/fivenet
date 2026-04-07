import { mount } from '@vue/test-utils';
import { describe, expect, it } from 'vitest';
import { defineComponent, nextTick } from 'vue';
import InputDurationPicker from '~/components/partials/InputDurationPicker.vue';
import type { Duration } from '~~/gen/ts/google/protobuf/duration';

const UInputNumberStub = defineComponent({
    name: 'UInputNumber',
    props: {
        modelValue: {
            type: Number,
            default: 0,
        },
    },
    emits: ['update:modelValue'],
    template: '<div class="u-input-number-stub" />',
});

const USelectMenuStub = defineComponent({
    name: 'USelectMenu',
    props: {
        modelValue: {
            type: String,
            default: undefined,
        },
    },
    emits: ['update:modelValue'],
    template: '<div class="u-select-menu-stub"><slot /></div>',
});

const UFormFieldStub = defineComponent({
    name: 'UFormField',
    template: '<div class="u-form-field-stub"><slot /></div>',
});

const UFieldGroupStub = defineComponent({
    name: 'UFieldGroup',
    template: '<div class="u-field-group-stub"><slot /></div>',
});

const UTooltipStub = defineComponent({
    name: 'UTooltip',
    template: '<div class="u-tooltip-stub"><slot /></div>',
});

const UButtonStub = defineComponent({
    name: 'UButton',
    props: {
        label: {
            type: String,
            default: '',
        },
    },
    emits: ['click'],
    template: '<button class="u-button-stub" @click="$emit(\'click\')">{{ label }}</button>',
});

const ClientOnlyStub = defineComponent({
    name: 'ClientOnly',
    template: '<div class="client-only-stub"><slot /></div>',
});

type PickerTestProps = { modelValue: Duration | undefined } & Record<string, unknown>;

function createWrapper(props: PickerTestProps) {
    return mount(InputDurationPicker, {
        props,
        global: {
            stubs: {
                UInputNumber: UInputNumberStub,
                USelectMenu: USelectMenuStub,
                UFormField: UFormFieldStub,
                UFieldGroup: UFieldGroupStub,
                UTooltip: UTooltipStub,
                UButton: UButtonStub,
                ClientOnly: ClientOnlyStub,
            },
            mocks: {
                $t: (key: string) => key,
            },
        },
    });
}

describe('InputDurationPicker', () => {
    it('emits protobuf duration in single mode', async () => {
        const wrapper = createWrapper({
            modelValue: undefined,
            mode: 'single',
        });

        wrapper.findComponent(UInputNumberStub).vm.$emit('update:modelValue', 90);
        await nextTick();

        expect(wrapper.emitted('update:modelValue')?.at(-1)?.[0]).toEqual({ seconds: 90, nanos: 0 });
    });

    it('auto-clamps to max bound', async () => {
        const wrapper = createWrapper({
            modelValue: undefined,
            mode: 'single',
            units: ['hour'],
            max: { seconds: 3600, nanos: 0 },
        });

        wrapper.findComponent(UInputNumberStub).vm.$emit('update:modelValue', 2);
        await nextTick();

        expect(wrapper.emitted('update:modelValue')?.at(-1)?.[0]).toEqual({ seconds: 3600, nanos: 0 });
    });

    it('emits protobuf duration in composite mode', async () => {
        const wrapper = createWrapper({
            modelValue: undefined,
            mode: 'composite',
            units: ['day', 'hour', 'minute'],
        });

        const inputs = wrapper.findAllComponents(UInputNumberStub);
        inputs[1]!.vm.$emit('update:modelValue', 2);
        inputs[2]!.vm.$emit('update:modelValue', 30);
        await nextTick();

        expect(wrapper.emitted('update:modelValue')?.at(-1)?.[0]).toEqual({ seconds: 9000, nanos: 0 });
    });

    it('respects allowed units', async () => {
        const wrapper = createWrapper({
            modelValue: undefined,
            mode: 'single',
            units: ['hour'],
        });

        wrapper.findComponent(UInputNumberStub).vm.$emit('update:modelValue', 1.5);
        await nextTick();

        expect(wrapper.emitted('update:modelValue')?.at(-1)?.[0]).toEqual({ seconds: 5400, nanos: 0 });
    });

    it('keeps clamped value when switching mode', async () => {
        const wrapper = createWrapper({
            modelValue: { seconds: 10800, nanos: 0 },
            mode: 'single',
            units: ['hour', 'minute'],
            max: { seconds: 7200, nanos: 0 },
        });
        await nextTick();

        expect(wrapper.emitted('update:modelValue')?.at(-1)?.[0]).toEqual({ seconds: 7200, nanos: 0 });

        await wrapper.setProps({
            modelValue: { seconds: 7200, nanos: 0 },
            mode: 'composite',
        });
        await nextTick();

        const inputs = wrapper.findAllComponents(UInputNumberStub);
        expect(inputs[0]!.props('modelValue')).toBe(2);
        expect(inputs[1]!.props('modelValue')).toBe(0);
    });

    it('emits undefined when clear button is clicked', async () => {
        const wrapper = createWrapper({
            modelValue: { seconds: 3600, nanos: 0 },
            clearable: true,
            mode: 'single',
        });

        await wrapper.findComponent(UButtonStub).trigger('click');

        expect(wrapper.emitted('update:modelValue')?.at(-1)?.[0]).toBeUndefined();
    });
});
