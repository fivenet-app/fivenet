<script setup lang="ts">
import type { Textbox } from 'fabric';
import type { FabricCurvedText } from '~/composables/fabric/FabricCurvedText';
import type { FabricHtmlInput } from '~/composables/fabric/FabricHtmlInput';
import { svgPatterns, useFabricEditor } from '~/composables/useFabricEditor';
import { fonts } from '~/types/editor';
import ColorPicker from '../partials/ColorPicker.vue';
import EditorSettings from './EditorSettings.vue';
import EditorShapes from './EditorShapes.vue';

const { activeObject, canvas, applyPatternFill } = useFabricEditor();

// Update functions to modify the selected object's properties
const updateText = (val: string) => {
    if (!activeObject.value) return;
    activeObject.value.set('text', val);
    canvas.value?.renderAll();
};

const updateFontSize = (size: number) => {
    if (!isNaN(size) && activeObject.value) {
        activeObject.value.set('fontSize', size);
        canvas.value?.renderAll();
    }
};

const updateFontFamily = (fontFamily: string) => {
    if (activeObject.value) {
        activeObject.value.set('fontFamily', fontFamily);
        activeObject.value.set('text', (activeObject.value as Textbox).text + '-');
        activeObject.value.set('text', (activeObject.value as Textbox).text.slice(0, -1));
        canvas.value?.renderAll();
    }
};

const updateFillColor = (val: string) => {
    if (activeObject.value) {
        activeObject.value.set('fill', val);
        canvas.value?.renderAll();
    }
};

const updateStrokeColor = (val: string) => {
    if (activeObject.value) {
        activeObject.value.set('stroke', val);
        canvas.value?.renderAll();
    }
};

const updateStrokeWidth = (val: number) => {
    if (activeObject.value) {
        activeObject.value.set('strokeWidth', val);
        canvas.value?.renderAll();
    }
};

const updateStrokeDash = (dashArray: number[] | null) => {
    if (activeObject.value) {
        activeObject.value.set('strokeDashArray', dashArray);
        canvas.value?.renderAll();
    }
};

const updateOpacity = (opacity: number) => {
    if (!isNaN(opacity) && activeObject.value) {
        activeObject.value.set('opacity', opacity);
        canvas.value?.renderAll();
    }
};

const selectedPattern = ref<string | undefined>(undefined);
const selectedPatternColor = ref<string>('#333333');

watch([selectedPattern, selectedPatternColor], async () => {
    if (activeObject.value) {
        if (selectedPattern.value) {
            await applyPatternFill(selectedPattern.value, selectedPatternColor.value);
        } else {
            // If no pattern is selected, remove pattern fill
            activeObject.value.set('fill', selectedPatternColor.value);
            canvas.value?.renderAll();
        }
    }
});

const updateCurvedText = (val: string) => {
    if (!activeObject.value || !activeObject.value.isType('curved-text')) return;

    (activeObject.value as FabricCurvedText).set('text', val);
    (activeObject.value as FabricCurvedText).update();
    canvas.value?.renderAll();
};

const updateCurvedTextFontSize = (size: number) => {
    if (!isNaN(size) && activeObject.value && activeObject.value.isType('curved-text')) {
        const curvedText = activeObject.value as FabricCurvedText;
        curvedText.update(undefined, undefined, { ...curvedText.options, fontSize: size });
        canvas.value?.renderAll();
    }
};

const updateCurvedTextFontFamily = (fontFamily: string) => {
    if (activeObject.value && activeObject.value.isType('curved-text')) {
        const curvedText = activeObject.value as FabricCurvedText;
        curvedText.update(undefined, undefined, { ...curvedText.options, fontFamily });
        canvas.value?.renderAll();
    }
};

const updateCurvedTextFillColor = (val: string) => {
    if (activeObject.value && activeObject.value.isType('curved-text')) {
        const curvedText = activeObject.value as FabricCurvedText;
        curvedText.update(undefined, undefined, { ...curvedText.options, fill: val });
        canvas.value?.renderAll();
    }
};
</script>

<template>
    <div class="flex flex-col gap-4 overflow-y-auto p-2 text-sm">
        <EditorSettings />

        <UCard>
            <template #header>
                <h3 class="text-xs font-bold tracking-wider uppercase">{{ $t('common.propertie', 2) }}</h3>
            </template>

            <template v-if="activeObject">
                <!-- Text/Textbox properties -->
                <div v-if="activeObject.isType('textbox')" class="flex flex-col gap-2">
                    <UFormField :label="$t('common.content')">
                        <UInput
                            type="text"
                            :model-value="(activeObject as Textbox).text"
                            class="w-full"
                            @update:model-value="updateText($event)"
                        />
                    </UFormField>

                    <UFormField :label="$t('components.partials.tiptap_editor.font_size')">
                        <UInputNumber
                            :model-value="(activeObject as Textbox).fontSize"
                            :min="8"
                            :max="100"
                            :step="1"
                            class="w-full"
                            @update:model-value="($event) => updateFontSize($event ?? 0)"
                        />
                    </UFormField>

                    <UFormField :label="$t('components.partials.tiptap_editor.font_family')">
                        <UInputMenu
                            :model-value="(activeObject as Textbox).fontFamily"
                            class="w-full"
                            name="selectedFont"
                            :filter-fields="['label']"
                            :items="fonts"
                            :placeholder="$t('common.font', 1)"
                            value-key="value"
                            label-key="label"
                            @update:model-value="updateFontFamily"
                        >
                            <template #item-label="{ item }">
                                <span class="truncate" :style="{ fontFamily: item.value }">{{
                                    item.label.includes('.') ? $t(item.label) : item.label
                                }}</span>
                            </template>

                            <template #empty>
                                {{ $t('common.not_found', [$t('components.partials.tiptap_editor.font_family')]) }}
                            </template>
                        </UInputMenu>
                    </UFormField>

                    <UFormField label="Text Color">
                        <ColorPicker
                            :model-value="typeof activeObject.fill === 'string' ? activeObject.fill : '#000000'"
                            class="w-full"
                            @update:model-value="updateFillColor($event ?? '#000000')"
                        />
                    </UFormField>
                </div>

                <div v-else-if="activeObject.isType('curved-text')" class="flex flex-col gap-2">
                    <UFormField :label="$t('common.content')">
                        <UInput
                            type="text"
                            :model-value="(activeObject as FabricCurvedText).text"
                            class="w-full"
                            @update:model-value="updateCurvedText($event)"
                        />
                    </UFormField>

                    <UFormField :label="$t('components.partials.tiptap_editor.font_size')">
                        <UInputNumber
                            :model-value="(activeObject as FabricCurvedText).options?.fontSize ?? 16"
                            :min="8"
                            :max="100"
                            class="w-full"
                            @update:model-value="($event) => updateCurvedTextFontSize($event ?? 0)"
                        />
                    </UFormField>

                    <UFormField :label="$t('components.partials.tiptap_editor.font_family')">
                        <UInputMenu
                            :model-value="(activeObject as FabricCurvedText).options?.fontFamily"
                            class="w-full"
                            name="selectedFont"
                            :filter-fields="['label']"
                            :items="fonts"
                            :placeholder="$t('common.font', 1)"
                            value-key="value"
                            label-key="label"
                            @update:model-value="updateCurvedTextFontFamily"
                        >
                            <template #item-label="{ item }">
                                <span class="truncate" :style="{ fontFamily: item.value }">{{
                                    item.label.includes('.') ? $t(item.label) : item.label
                                }}</span>
                            </template>

                            <template #empty>
                                {{ $t('common.not_found', [$t('components.partials.tiptap_editor.font_family')]) }}
                            </template>
                        </UInputMenu>
                    </UFormField>

                    <UFormField label="Text Color">
                        {{ (activeObject as FabricCurvedText).fill }}

                        <ColorPicker
                            :model-value="
                                typeof (activeObject as FabricCurvedText).fill === 'string'
                                    ? ((activeObject as FabricCurvedText).fill as string)
                                    : '#000000'
                            "
                            class="w-full"
                            @update:model-value="updateCurvedTextFillColor($event ?? '#000000')"
                        />
                    </UFormField>
                </div>

                <!-- Shape (Rectangle/Circle) properties -->
                <div v-else-if="activeObject.isType('rect', 'circle')" class="flex flex-col gap-2">
                    <UFormField label="Fill Color">
                        <ColorPicker
                            :model-value="
                                typeof activeObject.fill === 'string' && !activeObject.fill.includes('url(')
                                    ? activeObject.fill
                                    : '#000000'
                            "
                            class="w-full"
                            @update:model-value="updateFillColor($event ?? '#000000')"
                        />
                    </UFormField>

                    <UFormField label="Stroke Color">
                        <ColorPicker
                            :model-value="
                                typeof activeObject.stroke === 'string' && !activeObject.stroke.includes('url(')
                                    ? activeObject.stroke
                                    : '#000000'
                            "
                            class="w-full"
                            @update:model-value="updateStrokeColor($event ?? '#000000')"
                        />
                    </UFormField>

                    <UFormField label="Stroke Width">
                        <USlider
                            :model-value="activeObject.strokeWidth"
                            :min="0"
                            :step="1"
                            :max="24"
                            class="w-full"
                            @update:model-value="updateStrokeWidth($event ?? 0)"
                        />
                    </UFormField>

                    <UFormField label="Stroke Pattern">
                        <USelectMenu
                            :model-value="activeObject.strokeDashArray"
                            :items="strokeDashes"
                            label-key="name"
                            value-key="value"
                            class="w-full"
                            @update:model-value="updateStrokeDash($event)"
                        />
                    </UFormField>

                    <UFormField label="Opacity">
                        <USlider
                            :min="0"
                            :max="1"
                            :step="0.1"
                            :model-value="activeObject.opacity ?? 1"
                            class="w-full"
                            @update:model-value="updateOpacity($event ?? 1)"
                        />
                    </UFormField>

                    <UFormField label="Pattern">
                        <USelectMenu
                            v-model="selectedPattern"
                            :items="svgPatterns"
                            label-key="name"
                            value-key="value"
                            class="w-full"
                        />
                    </UFormField>

                    <UFormField label="Pattern Color">
                        <ColorPicker v-model="selectedPatternColor" class="w-full" />
                    </UFormField>
                </div>

                <div v-else-if="activeObject.isType('image')" class="flex flex-col gap-2">
                    <UFormField label="Opacity">
                        <USlider
                            :min="0"
                            :max="1"
                            :step="0.1"
                            :model-value="activeObject.opacity ?? 1"
                            class="w-full"
                            @update:model-value="updateOpacity($event ?? 1)"
                        />
                    </UFormField>
                </div>

                <!-- HTML Input properties -->
                <div v-else-if="activeObject.isType('html-input')" class="flex flex-col gap-2">
                    <UFormField label="Name" required>
                        <UInput
                            type="text"
                            :model-value="(activeObject as FabricHtmlInput).name"
                            class="w-full"
                            @update:model-value="
                                (val) => {
                                    (activeObject as FabricHtmlInput)!.name = val;
                                    canvas?.renderAll();
                                }
                            "
                        />
                    </UFormField>

                    <UFormField label="Input Type">
                        <USelectMenu
                            :model-value="(activeObject as FabricHtmlInput).inputType"
                            :items="[
                                { name: 'Text', value: 'text' },
                                { name: 'Number', value: 'number' },
                                { name: 'Date', value: 'date' },
                                { name: 'Date & Time', value: 'datetime' },
                                { name: 'Time', value: 'time' },
                                { name: 'Checkbox', value: 'checkbox' },
                                { name: 'Select', value: 'select' },
                                { name: 'Signature', value: 'signature' },
                            ]"
                            label-key="name"
                            value-key="value"
                            class="w-full"
                            @update:model-value="
                                (val) => {
                                    (activeObject as FabricHtmlInput)!.inputType = val;
                                    (activeObject as FabricHtmlInput).label = toTitleCase(val);
                                    canvas?.renderAll();
                                }
                            "
                        />
                    </UFormField>

                    <UFormField v-if="(activeObject as FabricHtmlInput).inputType !== 'checkbox'" label="Placeholder">
                        <UInput
                            type="text"
                            :model-value="(activeObject as FabricHtmlInput).placeholder"
                            class="w-full"
                            @update:model-value="
                                (val) => {
                                    (activeObject as FabricHtmlInput)!.placeholder = val;
                                    canvas?.renderAll();
                                }
                            "
                        />
                    </UFormField>

                    <UFormField :label="$t('components.partials.tiptap_editor.font_size')">
                        <UInputNumber
                            :model-value="(activeObject as Textbox).fontSize"
                            :min="8"
                            :max="100"
                            class="w-full"
                            @update:model-value="($event) => updateFontSize($event ?? 0)"
                        />
                    </UFormField>

                    <UFormField :label="$t('components.partials.tiptap_editor.font_family')">
                        <UInputMenu
                            :model-value="(activeObject as Textbox).fontFamily"
                            class="w-full"
                            name="selectedFont"
                            :filter-fields="['label']"
                            :items="fonts"
                            :placeholder="$t('common.font', 1)"
                            value-key="value"
                            label-key="label"
                            @update:model-value="updateFontFamily"
                        >
                            <template #item-label="{ item }">
                                <span class="truncate" :style="{ fontFamily: item.value }">{{
                                    item.label.includes('.') ? $t(item.label) : item.label
                                }}</span>
                            </template>

                            <template #empty>
                                {{ $t('common.not_found', [$t('components.partials.tiptap_editor.font_family')]) }}
                            </template>
                        </UInputMenu>
                    </UFormField>

                    <UFormField label="Text Color">
                        <ColorPicker
                            :model-value="typeof activeObject.fill === 'string' ? activeObject.fill : '#000000'"
                            class="w-full"
                            @update:model-value="updateFillColor($event ?? '#000000')"
                        />
                    </UFormField>

                    <UFormField v-if="(activeObject as FabricHtmlInput).inputType === 'select'" label="Options">
                        <UInputTags
                            v-model="(activeObject as FabricHtmlInput).options"
                            class="w-full"
                            placeholder="Add options"
                            @change="() => canvas?.renderAll()"
                        />
                    </UFormField>

                    <UFormField v-if="(activeObject as FabricHtmlInput).inputType === 'checkbox'" label="Checkbox Style">
                        <USelectMenu
                            :model-value="(activeObject as FabricHtmlInput).fieldProps.mark || 'check'"
                            :items="[
                                { name: 'Checkmark', value: 'check' },
                                { name: 'Cross', value: 'x' },
                                { name: 'Boxed', value: 'box' },
                            ]"
                            label-key="name"
                            value-key="value"
                            class="w-full"
                            @update:model-value="
                                (val) => {
                                    (activeObject as FabricHtmlInput)!.fieldProps.mark = val;
                                    canvas?.renderAll();
                                }
                            "
                        />
                    </UFormField>
                </div>

                <!-- Placeholder for other object types if needed -->
            </template>

            <!-- No selection -->
            <div v-else class="text-sm text-muted">No object selected.</div>
        </UCard>

        <EditorShapes />
    </div>
</template>
