<!-- eslint-disable @typescript-eslint/no-explicit-any -->
<script lang="ts" setup>
import { UCheckbox, UInput, UInputNumber, USelect } from '#components';
import { Canvas } from 'fabric';
import { onMounted, ref } from 'vue';
// eslint-disable-next-line @typescript-eslint/consistent-type-imports
import { FabricHtmlInput } from '~/composables/fabric/FabricHtmlInput';
import InputDatePicker from '../partials/InputDatePicker.vue';
import SignaturePad from '../partials/SignaturePad.vue';

const props = defineProps<{
    json: string;
}>();

type Field = {
    name: string;
    inputType: string;
    top: string;
    left: string;
    width: string;
    height: string;
    style: Record<string, any>;
    fieldProps: Record<string, any>;
};

const canvasRef = useTemplateRef('canvasRef');
const container = useTemplateRef('container');
const form = ref<Record<string, any>>({});
const fields = ref<Field[]>([]);

const getComponent = (field: Field) => {
    switch (field.inputType) {
        case 'number':
            return UInputNumber;
        case 'checkbox':
            return UCheckbox;
        case 'select':
            return USelect;
        case 'text':
        default:
            return UInput;
    }
};

onMounted(async () => {
    if (!canvasRef.value) return;

    const fabricCanvas = new Canvas(canvasRef.value, {
        selection: false,
        interactive: false,
        preserveObjectStacking: true,
    });

    fabricCanvas.loadFromJSON(props.json).then((canvas) => {
        fields.value = canvas
            .getObjects()
            .filter((obj) => obj.isType('html-input'))
            .map((obj) => {
                const o = obj as FabricHtmlInput;
                const { left, top, width, height } = obj.getBoundingRect();
                form.value[o.name] = o.value || '';

                canvas.remove(obj);
                return {
                    name: o.name,
                    inputType: o.inputType || 'text',
                    top: `${top}px`,
                    left: `${left}px`,
                    width: `${width}px`,
                    height: `${height}px`,
                    style: o.fieldProps?.style || {
                        fontSize: o.fieldProps?.fontSize ? `${o.fieldProps.fontSize}px` : '14px',
                        fontFamily: o.fieldProps?.fontFamily || 'Arial',
                        color: o.fieldProps?.textColor || '#000000',
                    },
                    fieldProps: o.fieldProps || {},
                };
            });

        canvas.setDimensions({
            width: container.value?.clientWidth || 800,
            height: container.value?.clientHeight || 600,
        });

        fabricCanvas.skipTargetFind = false;

        canvas.getObjects().forEach((obj) => {
            if (!obj.isType('textbox')) return;

            const text: string = obj.get('text');

            const replacements = [
                { search: '{{$citizen.Firstname}}', value: 'Prof. Dr. Philipp' },
                { search: '{{$citizen.Lastname}}', value: 'Scott' },
                { search: '{{$citizen.UserId}}', value: '26061' },
                { search: '{{.ActiveChar.UserId}}', value: '26061' },
                { search: '{{$citizen.Dateofbirth}}', value: '15.05.1982' },
                { search: '{{now | date "02.01.2006"}}', value: '23.10.2025' },
            ];
            replacements.forEach((replacement) => {
                if (text.includes(replacement.search)) {
                    obj.set('text', text.replace(replacement.search, replacement.value || ''));
                }
            });
        });

        canvas.requestRenderAll();
    });
});

watchDeep(form, () => {
    console.log('Form updated:', form.value);
});
</script>

<template>
    <div ref="container" class="relative h-full h-screen w-full w-screen">
        <canvas ref="canvasRef"></canvas>

        <div class="absolute inset-0 z-10">
            <!-- Render UI Inputs based on Fabric object metadata -->
            <div
                v-for="field in fields"
                :key="field.name"
                class="absolute"
                :style="{
                    top: field.top,
                    left: field.left,
                    width: field.width,
                    height: field.height,
                }"
            >
                <template v-if="field.inputType === 'date'">
                    <InputDatePicker
                        v-model="form[field.name]"
                        class="w-full"
                        date-format="date"
                        :button="{
                            variant: 'link',
                            style: field.style,
                        }"
                        hide-icon
                        v-bind="{ ...field.fieldProps }"
                    />
                </template>

                <template v-else-if="field.inputType === 'datetime'">
                    <InputDatePicker
                        v-model="form[field.name]"
                        class="w-full"
                        date-format="short"
                        v-bind="{ ...field.fieldProps }"
                    />
                </template>

                <template v-else-if="field.inputType === 'time'">
                    <InputDatePicker
                        v-model="form[field.name]"
                        class="w-full"
                        date-format="time"
                        v-bind="{ ...field.fieldProps }"
                    />
                </template>

                <template v-else-if="field.inputType === 'checkbox'">
                    <div
                        class="flex h-full w-full cursor-pointer items-center justify-center text-center text-lg select-none"
                        :class="[field.fieldProps?.border ? 'border' : '', field.fieldProps?.rounded ? 'rounded' : '']"
                        :style="{
                            ...field.style,
                            borderColor: field.fieldProps?.borderColor || '#000000',
                        }"
                        @click="form[field.name] = !form[field.name]"
                    >
                        {{ form[field.name] ? field.fieldProps?.mark || '✗' : '' }}
                    </div>
                </template>

                <template v-else-if="field.inputType === 'signature'">
                    <SignaturePad transparent />
                </template>

                <template v-else-if="field.inputType === 'textarea'">
                    <UTextarea
                        v-model="form[field.name]"
                        class="h-full w-full"
                        :style="field.style"
                        v-bind="{ ...field.fieldProps }"
                    />
                </template>

                <template v-else>
                    <component
                        :is="getComponent(field)"
                        v-model="form[field.name]"
                        :type="field.inputType"
                        class="w-full"
                        v-bind="{ ...field.fieldProps }"
                    />
                </template>
            </div>
        </div>
    </div>
</template>
