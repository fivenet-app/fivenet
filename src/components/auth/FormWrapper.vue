<script lang="ts" setup>
import FiveNetLogo from '~/components/partials/logos/FiveNetLogo.vue';
</script>

<template>
    <div class="max-w-lg sm:min-w-[32rem] mx-auto">
        <div class="px-4 py-8 rounded-lg bg-base-800 sm:px-10">
            <FiveNetLogo class="h-auto mx-auto mb-2 w-36" />

            <slot />
        </div>
    </div>
</template>
