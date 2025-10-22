<script setup lang="ts">
import { FabricHtmlInput } from '~/composables/fabric/FabricHtmlInput';
import { useFabricEditor } from '~/composables/useFabricEditor';

// Get methods from composable to manipulate canvas
const { addText, addPlaceholder, addRectangle, addCircle, canvas } = useFabricEditor();

function addInput(): void {
    const input = new FabricHtmlInput({ left: 100, top: 100, text: 'Hello' });
    canvas.value?.add(input);
}
</script>

<template>
    <UCard>
        <template #header>
            <h3 class="text-xs font-bold tracking-wider uppercase">{{ $t('common.shape', 2) }}</h3>
        </template>

        <div class="flex flex-row flex-wrap gap-2">
            <!-- Add Text -->
            <UButton size="xs" icon="i-mdi-plus" @click="addText">Text</UButton>
            <!-- Add Placeholder -->
            <UButton size="xs" icon="i-mdi-plus" @click="addPlaceholder">Placeholder</UButton>
            <!-- Add Rectangle -->
            <UButton size="xs" icon="i-mdi-plus" @click="addRectangle">Rectangle</UButton>
            <!-- Add Circle -->
            <UButton size="xs" icon="i-mdi-plus" @click="addCircle">Circle</UButton>
            <!-- Add Input -->
            <UButton size="xs" icon="i-mdi-plus" @click="addInput">Input</UButton>
        </div>
    </UCard>
</template>
