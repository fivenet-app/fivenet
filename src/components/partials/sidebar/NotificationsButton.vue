<script setup lang="ts">
import { BellOutlineIcon } from 'mdi-vue3';
</script>

<template>
    <div class="relative flex-shrink-0">
        <NuxtLink
            :to="{ name: 'notifications' }"
            class="flex rounded-full bg-base-800 text-sm ring-2 ring-base-600 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
        >
            <span class="sr-only">{{ $t('components.partials.sidebar_notifications') }}</span>
            <BellOutlineIcon
                class="h-10 w-auto rounded-full bg-base-800 p-1 text-base-300 hover:text-base-100 hover:transition-colors"
            />
        </NuxtLink>
    </div>
</template>
