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
            <UButton size="xs" icon="i-mdi-plus" @click="addText">Text</UButton>

            <UPopover mode="hover">
                <UButton size="xs" icon="i-mdi-plus">Curved Text</UButton>

                <template #content>
                    <div class="flex flex-col gap-2 p-2">
                        <UFormField :label="$t('common.content')">
                            <UInput v-model="text" type="text" class="w-full" />
                        </UFormField>

                        <UButton size="xs" @click="addCurvedText(text, 100, { arcAngleDeg: 180, clockwise: true })">
                            {{ $t('fabricEditor.curvedText.halfCircleUp') }}
                        </UButton>
                        <UButton size="xs" @click="addCurvedText(text, 100, { arcAngleDeg: 180, clockwise: false })">
                            {{ $t('fabricEditor.curvedText.halfCircleDown') }}
                        </UButton>
                        <UButton size="xs" @click="addCurvedText(text, 100, { arcAngleDeg: 360, clockwise: true })">
                            {{ $t('fabricEditor.curvedText.fullCircle') }}
                        </UButton>
                    </div>
                </template>
            </UPopover>

            <UButton size="xs" icon="i-mdi-plus" @click="addPlaceholder">Placeholder</UButton>

            <UButton size="xs" icon="i-mdi-plus" @click="addRectangle">Rectangle</UButton>

            <UButton size="xs" icon="i-mdi-plus" @click="addCircle">Circle</UButton>

            <UButton size="xs" icon="i-mdi-plus" @click="addInput">Input</UButton>

            <UButton v-if="jobProps?.logoFile" size="xs" icon="i-mdi-plus" @click="addImage('/' + jobProps.logoFile.filePath)">
                {{ $t('common.logo') }}
            </UButton>
        </div>
    </UCard>
</template>
