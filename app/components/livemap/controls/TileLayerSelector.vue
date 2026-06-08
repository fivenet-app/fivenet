<script lang="ts" setup>
import { tileLayers } from '~/types/livemap';

const props = withDefaults(
    defineProps<{
        modelValue: string;
        disabled?: boolean;
        position?: 'topleft' | 'topright' | 'bottomleft' | 'bottomright';
    }>(),
    {
        disabled: false,
        position: 'topright',
    },
);

const emit = defineEmits<{
    (e: 'update:modelValue', value: string): void;
}>();

const layer = computed({
    get: () => props.modelValue,
    set: (value: string) => emit('update:modelValue', value),
});
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
                            v-model="layer"
                            class="overflow-y-hidden"
                            :items="tileLayers"
                            value-key="key"
                            :ui-radio="{ inner: 'ms-1' }"
                            :ui="{ fieldset: 'grid auto-cols-auto grid-flow-col gap-1' }"
                        >
                            <template #label="{ item }">
                                {{ $t(item.label ?? item.id) }}
                            </template>
                        </URadioGroup>
                    </div>
                </div>
            </template>
        </UPopover>
    </LControl>
</template>
