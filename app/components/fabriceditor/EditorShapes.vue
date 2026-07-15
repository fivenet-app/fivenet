<script setup lang="ts">
import { FabricHtmlInput } from '~/composables/fabric/FabricHtmlInput';
import { useFabricEditor } from '~/composables/useFabricEditor';

defineProps<{
    disableShapeInput?: boolean;
}>();

const { t } = useI18n();

// Get methods from composable to manipulate canvas
const {
    addText,
    addCurvedText,
    addPolyline,
    addPolygon,
    addTriangle,
    addEllipse,
    addPlaceholder,
    addRectangle,
    addCircle,
    addImage,
    canvas,
} = useFabricEditor();

const { jobProps } = useAuth();

function addInput(): void {
    const input = new FabricHtmlInput({ left: 100, top: 100, value: t('common.text') });
    canvas.value?.add(input);
}

const text = ref<string>(t('components.fabric_editor.curved_text.title').valueOf());
</script>

<template>
    <UCard>
        <template #header>
            <h3 class="text-xs font-bold tracking-wider uppercase">{{ $t('common.shape', 2) }}</h3>
        </template>

        <div class="flex flex-col gap-2">
            <div class="grid grid-cols-2 items-center gap-2">
                <UButton size="xs" icon="i-mdi-format-textbox" :label="$t('components.fabric_editor.text')" @click="addText" />

                <UPopover mode="hover">
                    <UButton size="xs" icon="i-mdi-circle-outline" :label="$t('components.fabric_editor.curved_text.title')" />

                    <template #content>
                        <div class="flex flex-col gap-2 p-2">
                            <UFormField :label="$t('common.content')">
                                <UInput v-model="text" class="w-full" type="text" />
                            </UFormField>

                            <UButton
                                icon="i-mdi-circle-half"
                                size="xs"
                                :label="$t('components.fabric_editor.curved_text.half_circle_up')"
                                :ui="{ leadingIcon: 'rotate-90' }"
                                @click="addCurvedText(text, 100, { arcAngleDeg: 180, clockwise: true })"
                            />
                            <UButton
                                icon="i-mdi-circle-half"
                                size="xs"
                                :label="$t('components.fabric_editor.curved_text.half_circle_down')"
                                :ui="{ leadingIcon: 'rotate-270' }"
                                @click="addCurvedText(text, 100, { arcAngleDeg: 180, clockwise: false })"
                            />
                            <UButton
                                icon="i-mdi-circle-outline"
                                size="xs"
                                :label="$t('components.fabric_editor.curved_text.full_circle')"
                                @click="addCurvedText(text, 100, { arcAngleDeg: 360, clockwise: true })"
                            />
                        </div>
                    </template>
                </UPopover>

                <template v-if="disableShapeInput">
                    <UButton
                        size="xs"
                        icon="i-mdi-note-text"
                        :label="$t('components.fabric_editor.placeholder')"
                        @click="addPlaceholder"
                    />
                    <UButton
                        size="xs"
                        icon="i-mdi-form-textbox"
                        :label="$t('components.fabric_editor.input')"
                        @click="addInput"
                    />
                </template>
            </div>

            <USeparator />

            <div class="grid grid-cols-2 items-center gap-2">
                <UButton
                    size="xs"
                    icon="i-mdi-vector-polyline"
                    :label="$t('components.fabric_editor.polyline')"
                    @click="addPolyline"
                />
                <UButton
                    size="xs"
                    icon="i-mdi-vector-polygon"
                    :label="$t('components.fabric_editor.polygon')"
                    @click="addPolygon"
                />
                <UButton
                    size="xs"
                    icon="i-mdi-vector-triangle"
                    :label="$t('components.fabric_editor.triangle')"
                    @click="addTriangle"
                />
                <UButton
                    size="xs"
                    icon="i-mdi-vector-rectangle"
                    :label="$t('components.fabric_editor.rectangle')"
                    @click="addRectangle"
                />
                <UButton
                    size="xs"
                    icon="i-mdi-vector-ellipse"
                    :label="$t('components.fabric_editor.ellipse')"
                    @click="addEllipse"
                />
                <UButton
                    size="xs"
                    icon="i-mdi-vector-circle"
                    :label="$t('components.fabric_editor.circle')"
                    @click="addCircle"
                />
            </div>

            <UButton
                v-if="jobProps?.logoFile"
                size="xs"
                icon="i-mdi-image"
                :label="$t('common.logo')"
                @click="addImage('/api/filestore/' + jobProps.logoFile.filePath)"
            />
        </div>
    </UCard>
</template>
