<script setup lang="ts">
defineProps<{
    label: string;
    icon?: string;
    chip?: string;
    selected?: boolean;
}>();

const slots = defineSlots<{
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    leading: () => any;
}>();
</script>

<template>
    <UButton
        class="text-[11px] capitalize ring-default"
        :class="[selected ? 'bg-elevated' : 'hover:bg-elevated/50']"
        size="sm"
        color="neutral"
        variant="outline"
        :icon="icon"
        :label="label"
    >
        <template v-if="chip || !!slots.leading" #leading>
            <slot name="leading">
                <span
                    class="inline-block size-2 rounded-full"
                    :class="`bg-(--color-light) dark:bg-(--color-dark)`"
                    :style="{
                        '--color-light': `var(--color-${chip}-500)`,
                        '--color-dark': `var(--color-${chip}-400)`,
                    }"
                />
            </slot>
        </template>
    </UButton>
</template>
