<script lang="ts" setup>
import MapPositionPicker from '~/components/livemap/MapPositionPicker.vue';
import type { TileLayerKeys } from '~/types/livemap';

const props = withDefaults(
    defineProps<{
        open: boolean;
        title: string;
        summary: string;
        x: number;
        y: number;
        zoom: number;
        layer?: TileLayerKeys;
        disabled?: boolean;
        showControls?: boolean;
    }>(),
    {
        layer: undefined,
        disabled: false,
        showControls: false,
    },
);

const emit = defineEmits<{
    (e: 'update:open', value: boolean): void;
    (e: 'update:x', value: number): void;
    (e: 'update:y', value: number): void;
    (e: 'update:zoom', value: number): void;
    (e: 'update:layer', value: TileLayerKeys): void;
}>();

const openModel = computed({
    get: () => props.open,
    set: (value: boolean) => emit('update:open', value),
});
</script>

<template>
    <UModal v-model:open="openModel" fullscreen :title="title" :ui="{ body: 'flex flex-col gap-4' }">
        <template #body>
            <div class="flex items-start justify-between gap-4">
                <div>
                    <div class="font-medium">{{ title }}</div>
                    <div class="text-xs text-neutral-500 dark:text-neutral-400">
                        {{ summary }}
                    </div>
                </div>
            </div>

            <ClientOnly>
                <MapPositionPicker
                    v-if="open"
                    :x="x"
                    :y="y"
                    :zoom="zoom"
                    :layer="layer"
                    :disabled="disabled"
                    :show-controls="showControls"
                    container-class="relative h-[calc(100vh-12rem)] w-full overflow-hidden rounded-lg border border-neutral-300 dark:border-neutral-700"
                    @update:x="(value) => emit('update:x', value)"
                    @update:y="(value) => emit('update:y', value)"
                    @update:zoom="(value) => emit('update:zoom', value)"
                    @update:layer="(value) => emit('update:layer', value)"
                />
            </ClientOnly>
        </template>

        <template #footer>
            <UButton class="flex-1" block color="neutral" :label="$t('common.close', 1)" @click="openModel = false" />
        </template>
    </UModal>
</template>
