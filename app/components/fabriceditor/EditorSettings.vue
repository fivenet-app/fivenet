<script lang="ts" setup>
const { documentSize, snapThreshold } = useFabricEditor();
</script>

<template>
    <!-- Canvas Settings Section -->
    <UCard>
        <template #header>
            <h3 class="text-xs font-bold tracking-wider uppercase">Canvas Settings</h3>
        </template>

        <div class="flex items-center gap-2">
            <UFormField help="Width">
                <UInput
                    v-model.number="documentSize.width"
                    type="number"
                    size="xs"
                    class="w-20"
                    placeholder="Width"
                    :disabled="documentSize.disabled"
                />
                {{ documentSize.disabled }}
            </UFormField>

            <UFormField>
                <span class="text-xs text-muted">x</span>
            </UFormField>

            <UFormField help="Height">
                <UInput
                    v-model.number="documentSize.height"
                    type="number"
                    size="xs"
                    class="w-20"
                    placeholder="Height"
                    :disabled="documentSize.disabled"
                />
            </UFormField>

            <UFormField>
                <span class="text-xs text-muted">px</span>
            </UFormField>
        </div>

        <UFormField label="Snap Threshold">
            <UInputNumber v-model="snapThreshold" size="xs" class="w-20" />
        </UFormField>
    </UCard>
</template>
