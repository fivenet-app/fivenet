<script lang="ts" setup>
import { VueSignaturePad, type Signature } from '@selemondev/vue3-signature-pad';

const settingsStore = useSettingsStore();
const { signature: signatureSettings } = storeToRefs(settingsStore);

const modelValue = defineModel<string>({ required: true });

const colors = [
    {
        color: 'rgb(0, 0, 0)',
    },
    {
        color: 'rgb(51, 133, 255)',
    },
    {
        color: 'rgb(255, 85, 51)',
    },
];
const signature = ref<Signature>();

function handleUndo() {
    signature.value?.undo();
}

function handleClearCanvas() {
    signature.value?.clearCanvas();
}

function handleSaveSignature() {
    if (!signature.value) return;

    modelValue.value = signature.value.saveSignature('image/svg');
}
</script>

<template>
    <UCard :ui="{ root: 'grow-0 max-w-[950px]', body: 'p-0 sm:p-0', footer: 'p-2 sm:px-2' }">
        <div class="relative">
            <VueSignaturePad
                ref="signature"
                height="400px"
                width="950px"
                :min-width="signatureSettings.minStrokeWidth"
                :max-width="signatureSettings.maxStrokeWidth"
                :options="{
                    penColor: signatureSettings.penColor,
                    backgroundColor: 'rgb(255, 255, 255)',
                }"
                @end-stroke="handleSaveSignature"
            />

            <UButtonGroup class="absolute bottom-0 left-0 flex flex-row">
                <UBadge icon="i-mdi-color" class="!cursor-default rounded-l-none" size="lg" />

                <UButton
                    v-for="(color, idx) in colors"
                    :key="color.color"
                    :style="{ background: color.color }"
                    class="inline-flex w-8 items-center"
                    :class="idx === colors.length - 1 ? 'rounded-br-none' : ''"
                    @click="signatureSettings.penColor = color.color"
                >
                    <UIcon
                        v-if="signatureSettings.penColor === color.color"
                        class="size-full text-white"
                        name="i-mdi-check-bold"
                    />
                </UButton>
            </UButtonGroup>

            <UButtonGroup class="absolute right-0 bottom-0 flex flex-row">
                <UTooltip :text="$t('common.undo')">
                    <UButton icon="i-mdi-undo" class="rounded-bl-none" @click="handleUndo" />
                </UTooltip>

                <UTooltip :text="$t('common.clear')">
                    <UButton icon="i-mdi-clear" class="rounded-r-none" @click="handleClearCanvas" />
                </UTooltip>
            </UButtonGroup>
        </div>

        <template #footer>
            <UCollapsible>
                <UButton
                    :label="$t('common.setting', 2)"
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-signature-text"
                    block
                    trailing-icon="i-mdi-chevron-down"
                    :ui="{
                        trailingIcon: 'group-data-[state=open]:rotate-180 transition-transform duration-200',
                    }"
                />

                <template #content>
                    <div class="my-2 flex w-full flex-col space-y-4">
                        <UFormField
                            class="grid grid-cols-2 items-center gap-2"
                            :label="$t('components.partials.SignaturePad.min_stroke_width')"
                        >
                            <USlider v-model="signatureSettings.minStrokeWidth" :min="1" :max="10" />
                            <p>{{ signatureSettings.minStrokeWidth }}</p>
                        </UFormField>

                        <UFormField
                            class="grid grid-cols-2 items-center gap-2"
                            :label="$t('components.partials.SignaturePad.max_stroke_width')"
                        >
                            <USlider v-model="signatureSettings.maxStrokeWidth" :min="1" :max="10" />
                            <p>{{ signatureSettings.maxStrokeWidth }}</p>
                        </UFormField>
                    </div>
                </template>
            </UCollapsible>
        </template>
    </UCard>
</template>
