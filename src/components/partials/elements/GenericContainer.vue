<script lang="ts" setup>
defineOptions({
    inheritAttrs: true,
});
</script>
<template>
    <div class="overflow-hidden rounded-lg bg-base-700 shadow">
        <div class="px-4 py-5 sm:p-6">
            <slot />
        </div>
    </div>
</template>
