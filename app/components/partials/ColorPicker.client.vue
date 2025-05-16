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
            color="white"
            :disabled="disabled"
            :block="block"
            :icon="!hideIcon ? 'i-mdi-palette' : ''"
            :label="!hideIcon ? '' : '&nbsp;'"
            :style="{ backgroundColor: color }"
            v-bind="$attrs"
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

                        <UButton class="-my-1" color="gray" variant="ghost" icon="i-mdi-window-close" @click="open = false" />
                    </div>
                </template>

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

                <template #footer>
                    <UButton class="flex-1" color="black" block @click="open = false">
                        {{ $t('common.close', 1) }}
                    </UButton>
                </template>
            </UCard>
        </UModal>
    </template>

    <UPopover v-else v-model:open="open" :popper="{ placement: 'bottom-start' }" v-bind="$attrs">
        <UButton
            :class="$attrs.class"
            variant="outline"
            color="white"
            :disabled="disabled"
            :block="block"
            :icon="!hideIcon ? 'i-mdi-palette' : ''"
            :label="!hideIcon ? '' : '&nbsp;'"
            :style="{ backgroundColor: color }"
            v-bind="$attrs"
            @touchstart="open = true"
        />

        <template #panel>
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
