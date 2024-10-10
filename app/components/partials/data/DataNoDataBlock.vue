<script lang="ts" setup>
const props = withDefaults(
    defineProps<{
        message?: string;
        icon?: string;
        type?: string;
        focus?: () => void | Promise<void>;
        padding?: string;
    }>(),
    {
        message: undefined,
        icon: 'i-mdi-magnify',
        type: undefined,
        focus: undefined,
        padding: 'p-4',
    },
);

async function click() {
    if (props.focus) {
        props.focus();
    }
}
</script>

<template>
    <UButton
        variant="outline"
        :disabled="!focus"
        :icon="icon"
        block
        class="block w-full text-center"
        :class="padding"
        @click="click()"
    >
        <span class="mt-2 block text-sm font-semibold">
            <template v-if="message">
                {{ message }}
            </template>
            <template v-else>
                {{ $t('common.not_found', [type ?? $t('common.data')]) }}
            </template>
        </span>
    </UButton>
</template>
