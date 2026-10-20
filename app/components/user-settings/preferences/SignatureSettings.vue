<script lang="ts" setup>
import { signatureColors } from '~/utils/signature';

const settingsStore = useSettingsStore();
const { signature: signatureSettings } = storeToRefs(settingsStore);
</script>

<template>
    <div class="space-y-4">
        <UFormField class="grid grid-cols-2 items-center gap-2" name="signature.penColor" :label="$t('common.color', 1)">
            <div class="flex justify-end">
                <UFieldGroup>
                    <UButton
                        v-for="color in signatureColors"
                        :key="color.color"
                        class="inline-flex w-8 items-center"
                        :style="{ background: color.color }"
                        @click="signatureSettings.penColor = color.color"
                    >
                        <UIcon
                            v-if="signatureSettings.penColor === color.color"
                            class="size-full text-white"
                            name="i-mdi-check-bold"
                        />
                    </UButton>
                </UFieldGroup>
            </div>
        </UFormField>

        <UFormField
            class="grid grid-cols-2 items-center gap-2"
            name="signature.minStrokeWidth"
            :label="$t('components.partials.signature_pad.min_stroke_width')"
        >
            <div class="flex items-center justify-end gap-2">
                <span>{{ signatureSettings.minStrokeWidth }}</span>
                <USlider v-model="signatureSettings.minStrokeWidth" class="w-full max-w-48" :min="1" :max="10" />
            </div>
        </UFormField>

        <UFormField
            class="grid grid-cols-2 items-center gap-2"
            name="signature.maxStrokeWidth"
            :label="$t('components.partials.signature_pad.max_stroke_width')"
        >
            <div class="flex items-center justify-end gap-2">
                <span>{{ signatureSettings.maxStrokeWidth }}</span>
                <USlider v-model="signatureSettings.maxStrokeWidth" class="w-full max-w-48" :min="1" :max="10" />
            </div>
        </UFormField>
    </div>
</template>
