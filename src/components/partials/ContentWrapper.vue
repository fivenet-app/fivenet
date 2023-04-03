<script lang="ts" setup>
defineProps({
    classes: {
        type: String,
        required: false,
        default: '',
    }
});
</script>

<template>
    <main :class="classes">
        <div class="py-4 mx-auto max-w-7xl">
            <slot />
        </div>
    </main>
</template>
