<script lang="ts" setup>
import FiveNetLogo from '~/components/partials/logos/FiveNetLogo.vue';
</script>

<template>
    <div class="mx-auto max-w-lg sm:min-w-[32rem]">
        <div class="rounded-lg bg-base-800 px-4 py-8 sm:px-10">
            <FiveNetLogo class="mx-auto mb-2 h-auto w-36" />

            <slot />
        </div>
    </div>
</template>
