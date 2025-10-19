<!-- eslint-disable @typescript-eslint/no-explicit-any -->
<script setup lang="ts">
const props = defineProps<{
    api?: any | null;
}>();
</script>

<template>
    <UCard class="h-full">
        <template #header>Properties</template>
        <div class="space-y-3">
            <UFormField label="Fill">
                <UInput type="color" @change="(e) => props.api?.setColor?.('fill', (e.target as HTMLInputElement).value)" />
            </UFormField>
            <UFormField label="Stroke">
                <UInput type="color" @change="(e) => props.api?.setColor?.('stroke', (e.target as HTMLInputElement).value)" />
            </UFormField>
            <UFormField label="Stroke width">
                <UInput
                    type="number"
                    placeholder="1"
                    @change="(e) => props.api?.setStrokeWidth?.(Number((e.target as HTMLInputElement).value || 1))"
                />
            </UFormField>
            <UFormField label="Opacity">
                <USlider
                    :model-value="props.api?.getOpacity() ?? 0"
                    :min="0"
                    :max="1"
                    :step="0.05"
                    @update:model-value="(e) => props.api?.setOpacity?.(e as number)"
                />
            </UFormField>
        </div>
    </UCard>
</template>
