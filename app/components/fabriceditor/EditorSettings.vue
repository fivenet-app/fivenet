<script lang="ts" setup>
const { documentSize, snapThreshold } = useFabricEditor();
</script>

<template>
    <!-- Canvas Settings Section -->
    <UCard>
        <template #header>
            <h3 class="text-xs font-bold tracking-wider uppercase">
                {{ $t('components.fabric_editor.canvas_settings') }}
            </h3>
        </template>

        <div class="flex items-center gap-2">
            <UFormField :label="$t('components.fabric_editor.width')">
                <UInputNumber
                    v-model="documentSize.width"
                    class="w-20"
                    size="xs"
                    :step="1"
                    :min="128"
                    :max="1280"
                    placeholder="Width"
                    :disabled="documentSize.disabled"
                />
            </UFormField>

            <UFormField>
                <span class="text-xs text-muted">x</span>
            </UFormField>

            <UFormField :label="$t('components.fabric_editor.height')">
                <UInputNumber
                    v-model="documentSize.height"
                    class="w-20"
                    size="xs"
                    :step="1"
                    :min="128"
                    :max="1280"
                    placeholder="Height"
                    :disabled="documentSize.disabled"
                />
            </UFormField>

            <UFormField>
                <span class="text-xs text-muted">px</span>
            </UFormField>
        </div>

        <UFormField :label="$t('components.fabric_editor.snap_threshold')">
            <div class="inline-flex items-center gap-1">
                <UInputNumber v-model="snapThreshold" class="w-20" size="xs" />
                <span class="text-xs text-muted">px</span>
            </div>
        </UFormField>
    </UCard>
</template>
