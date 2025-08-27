<script lang="ts" setup>
import { breakpointsTailwind, useBreakpoints } from '@vueuse/core';
import { ColorPicker } from 'vue3-colorpicker';
import 'vue3-colorpicker/style.css';

const props = defineProps<{
    modelValue: string | undefined;
    disabled?: boolean;
    hideIcon?: boolean;
    block?: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: string): void;
    (e: 'close'): void;
}>();

defineOptions({
    inheritAttrs: false,
});

const color = computed({
    get: () => props.modelValue ?? '',
    set: (value) => {
        emit('update:modelValue', value);
        emit('close');
    },
});

const breakpoints = useBreakpoints(breakpointsTailwind);

const smallerBreakpoint = breakpoints.smaller('sm');

const colorMode = useColorMode();
const theme = computed(() => (colorMode.value === 'dark' ? 'black' : 'white'));

const open = ref(false);
</script>

<template>
    <template v-if="smallerBreakpoint">
        <UButton
            :class="$attrs.class"
            variant="outline"
            color="neutral"
            :disabled="disabled"
            :block="block"
            :icon="!hideIcon ? 'i-mdi-palette' : ''"
            :label="!hideIcon ? '' : '&nbsp;'"
            :style="{ backgroundColor: color }"
            v-bind="$attrs"
            @click="open = true"
            @touchstart="open = true"
        />

        <UModal v-model:open="open" :title="$t('common.color')">
            <template #body>
                <div class="flex flex-1 items-center">
                    <ColorPicker
                        v-model:pure-color="color"
                        is-widget
                        format="hex"
                        picker-type="chrome"
                        disable-alpha
                        disable-history
                        :theme="theme"
                    />
                </div>
            </template>

            <template #footer>
                <UButton class="flex-1" color="neutral" block @click="open = false">
                    {{ $t('common.close', 1) }}
                </UButton>
            </template>
        </UModal>
    </template>

    <UPopover v-else v-model:open="open" v-bind="$attrs">
        <UButton
            :class="$attrs.class"
            variant="outline"
            color="neutral"
            :disabled="disabled"
            :block="block"
            :icon="!hideIcon ? 'i-mdi-palette' : ''"
            :label="!hideIcon ? '' : '&nbsp;'"
            :style="{ backgroundColor: color }"
            v-bind="$attrs"
            @touchstart="open = true"
        />

        <template #content>
            <ColorPicker
                v-model:pure-color="color"
                is-widget
                format="hex"
                picker-type="chrome"
                disable-alpha
                disable-history
                :theme="theme"
            />
        </template>
    </UPopover>
</template>

<style>
.vc-input-toggle {
    display: none !important;
}

@media not all and screen(sm) {
    .vc-colorpicker {
        box-shadow: none !important;
        flex: 1;
    }
}
</style>
