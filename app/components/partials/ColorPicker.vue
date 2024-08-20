<script lang="ts" setup>
import { breakpointsTailwind, useBreakpoints } from '@vueuse/core';
import { ColorPicker } from 'vue3-colorpicker';
import 'vue3-colorpicker/style.css';

const props = defineProps<{
    modelValue: string;
    disabled?: boolean;
    hideIcon?: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:model-value', value: string): void;
    (e: 'close'): void;
}>();

defineOptions({
    inheritAttrs: false,
});

const color = computed({
    get: () => props.modelValue,
    set: (value) => {
        emit('update:model-value', value);
        emit('close');
    },
});

const breakpoints = useBreakpoints(breakpointsTailwind);

const smallerBreakpoint = breakpoints.smaller('sm');

const open = ref(false);
</script>

<template>
    <template v-if="smallerBreakpoint">
        <UButton
            variant="outline"
            color="white"
            :disabled="disabled"
            block
            :icon="!hideIcon ? 'i-mdi-palette' : ''"
            :label="!hideIcon ? '' : '&nbsp;'"
            :style="{ backgroundColor: color }"
            :class="$attrs.class"
            @click="open = true"
            @touchstart="open = true"
        />

        <UModal v-model="open">
            <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl font-semibold leading-6">
                            {{ $t('common.color') }}
                        </h3>

                        <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="open = false" />
                    </div>
                </template>

                <div class="flex flex-1 items-center">
                    <div class="mx-auto">
                        <ColorPicker
                            v-model:pureColor="color"
                            is-widget
                            format="hex"
                            picker-type="chrome"
                            disable-alpha
                            disable-history
                        />
                    </div>
                </div>

                <template #footer>
                    <UButton color="black" block class="flex-1" @click="open = false">
                        {{ $t('common.close', 1) }}
                    </UButton>
                </template>
            </UCard>
        </UModal>
    </template>

    <UPopover v-else v-model:open="open" :popper="{ placement: 'bottom-start' }">
        <UButton
            variant="outline"
            color="white"
            :disabled="disabled"
            block
            :icon="!hideIcon ? 'i-mdi-palette' : ''"
            :label="!hideIcon ? '' : '&nbsp;'"
            :style="{ backgroundColor: color }"
            :class="$attrs.class"
            @touchstart="open = true"
        />

        <template #panel>
            <ColorPicker v-model:pureColor="color" is-widget format="hex" picker-type="chrome" disable-alpha disable-history />
        </template>
    </UPopover>
</template>

<style>
.vc-input-toggle {
    display: none !important;
}
</style>
