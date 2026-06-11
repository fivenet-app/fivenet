<script lang="ts" setup>
import { tileLayers, type TileLayerKeys } from '~/types/livemap';

const props = withDefaults(
    defineProps<{
        modelValue: TileLayerKeys;
        disabled?: boolean;
        position?: 'topleft' | 'topright' | 'bottomleft' | 'bottomright';
    }>(),
    {
        disabled: false,
        position: 'topright',
    },
);

const emit = defineEmits<{
    (e: 'update:modelValue', value: TileLayerKeys): void;
}>();

const tileLayerItems = computed(() => [...tileLayers]);
</script>

<template>
    <LControl :position="position">
        <UPopover :ui="{ content: 'w-full' }">
            <UTooltip :text="$t('common.layer', 2)">
                <UButton
                    class="border border-black/20 bg-clip-padding p-1.5"
                    size="xl"
                    icon="i-mdi-layers-triple"
                    :disabled="disabled"
                    :ui="{ leadingIcon: 'size-5!' }"
                />
            </UTooltip>

            <template #content>
                <div class="w-full max-w-xl divide-y divide-default py-1">
                    <div class="px-1 pb-0.5">
                        <p class="truncate text-base font-bold text-highlighted">
                            {{ $t('common.layer', 2) }}
                        </p>

                        <URadioGroup
                            class="overflow-y-hidden"
                            :model-value="props.modelValue"
                            :items="tileLayerItems"
                            value-key="key"
                            :ui-radio="{ inner: 'ms-1' }"
                            :ui="{ fieldset: 'grid auto-cols-auto grid-flow-col gap-1' }"
                            @update:model-value="(value) => emit('update:modelValue', value)"
                        >
                            <template #label="{ item }">
                                {{ $t(item.label) }}
                            </template>
                        </URadioGroup>
                    </div>
                </div>
            </template>
        </UPopover>
    </LControl>
</template>
