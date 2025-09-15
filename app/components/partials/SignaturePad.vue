<script lang="ts" setup>
import { VueSignaturePad, type Signature } from '@selemondev/vue3-signature-pad';

const options = ref({
    penColor: 'rgb(0,0,0)',
    backgroundColor: 'rgb(255, 255, 255)',
    maxWidth: 6,
    minWidth: 2,
});

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
    signature.value?.saveSignature() && alert(signature.value?.saveSignature());
}
</script>

<template>
    <UCard :ui="{ root: 'grow-0 max-w-[950px]', body: 'p-0 sm:p-0', footer: 'p-2 sm:px-2' }">
        <div class="relative">
            <VueSignaturePad
                ref="signature"
                height="400px"
                width="950px"
                :max-width="options.maxWidth"
                :min-width="options.minWidth"
                :options="{
                    penColor: options.penColor,
                    backgroundColor: options.backgroundColor,
                }"
            />

            <UButtonGroup class="absolute bottom-0 left-0 flex flex-row">
                <UTooltip :text="$t('common.color')">
                    <UButton icon="i-mdi-color" class="!cursor-default rounded-l-none" />
                </UTooltip>

                <UButton
                    v-for="(color, idx) in colors"
                    :key="color.color"
                    :style="{ background: color.color }"
                    class="inline-flex w-8 items-center"
                    :class="idx === colors.length - 1 ? 'rounded-br-none' : ''"
                    @click="options.penColor = color.color"
                >
                    <UIcon v-if="options.penColor === color.color" class="size-full text-white" name="i-mdi-check-bold" />
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
                        <UFormField class="grid grid-cols-2 items-center gap-2" label="Choose minimum pen line thickness">
                            <USlider v-model="options.minWidth" :min="1" :max="10" />
                            <p>{{ options.minWidth }}</p>
                        </UFormField>

                        <UFormField class="grid grid-cols-2 items-center gap-2" label="Choose maximum pen line thickness">
                            <USlider v-model="options.maxWidth" :min="1" :max="10" />
                            <p>{{ options.maxWidth }}</p>
                        </UFormField>
                    </div>
                </template>
            </UCollapsible>
        </template>
    </UCard>
</template>
