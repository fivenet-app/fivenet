<script setup lang="ts">
import { FabricHtmlInput } from '~/composables/fabric/FabricHtmlInput';
import { useFabricEditor } from '~/composables/useFabricEditor';

// Get methods from composable to manipulate canvas
const { addText, addCurvedText, addPlaceholder, addRectangle, addCircle, addImage, canvas } = useFabricEditor();

const { jobProps } = useAuth();

function addInput(): void {
    const input = new FabricHtmlInput({ left: 100, top: 100, value: 'Hello' });
    canvas.value?.add(input);
}

const text = ref<string>('Curved Text');
</script>

<template>
    <UCard>
        <template #header>
            <h3 class="text-xs font-bold tracking-wider uppercase">{{ $t('common.shape', 2) }}</h3>
        </template>

        <div class="flex flex-row flex-wrap gap-2">
            <!-- Add Text -->
            <UButton size="xs" icon="i-mdi-plus" :label="$t('components.fabric_editor.text')" @click="addText" />

            <UPopover mode="hover">
                <UButton size="xs" icon="i-mdi-plus" :label="$t('components.fabric_editor.curved_text.title')" />

                <template #content>
                    <div class="flex flex-col gap-2 p-2">
                        <UFormField :label="$t('common.content')">
                            <UInput v-model="text" type="text" class="w-full" />
                        </UFormField>

                        <UButton size="xs" :label="$t('components.fabric_editor.curved_text.half_circle_up')" @click="addCurvedText(text, 100, { arcAngleDeg: 180, clockwise: true })" />
                        <UButton size="xs" :label="$t('components.fabric_editor.curved_text.half_circle_down')" @click="addCurvedText(text, 100, { arcAngleDeg: 180, clockwise: false })" />
                        <UButton size="xs" :label="$t('components.fabric_editor.curved_text.full_circle')" @click="addCurvedText(text, 100, { arcAngleDeg: 360, clockwise: true })" />
                    </div>
                </template>
            </UPopover>

            <UButton size="xs" icon="i-mdi-plus" :label="$t('components.fabric_editor.placeholder')" @click="addPlaceholder" />
            <UButton size="xs" icon="i-mdi-plus" :label="$t('components.fabric_editor.rectangle')" @click="addRectangle" />
            <UButton size="xs" icon="i-mdi-plus" :label="$t('components.fabric_editor.circle')" @click="addCircle" />
            <UButton size="xs" icon="i-mdi-plus" :label="$t('components.fabric_editor.input')" @click="addInput" />

            <UButton v-if="jobProps?.logoFile" size="xs" icon="i-mdi-plus" :label="$t('common.logo')" @click="addImage('/' + jobProps.logoFile.filePath)" />
        </div>
    </UCard>
</template>
