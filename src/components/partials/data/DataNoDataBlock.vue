<script lang="ts" setup>
const props = withDefaults(
    defineProps<{
        message?: string;
        icon?: string;
        type?: string;
        focus?: Function;
        padding?: string;
    }>(),
    {
        message: undefined,
        icon: 'i-mdi-magnify',
        type: undefined,
        focus: undefined,
        padding: 'p-6',
    },
);

function click() {
    if (props.focus) {
        props.focus();
    }
}
</script>

<template>
    <UButton
        variant="ghost"
        :disabled="!focus"
        :icon="icon"
        block
        class="relative block w-full rounded-lg border border-dashed border-base-300 text-center"
        :class="padding"
        @click="click()"
    >
        <span class="mt-2 block text-sm font-semibold">
            <span v-if="message">
                {{ message }}
            </span>
            <span v-else>
                {{ $t('common.not_found', [type ?? $t('common.data')]) }}
            </span>
        </span>
    </UButton>
</template>
