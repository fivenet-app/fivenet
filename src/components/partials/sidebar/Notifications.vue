<script setup lang="ts">
import { BellOutlineIcon } from 'mdi-vue3';
</script>

<template>
    <div class="relative flex-shrink-0">
        <NuxtLink
            :to="{ name: 'notifications' }"
            class="flex text-sm rounded-full bg-base-800 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
        >
            <span class="sr-only">{{ $t('components.partials.sidebar_notifications') }}</span>
            <BellOutlineIcon
                class="w-auto h-10 p-1 rounded-full hover:transition-colors text-base-300 bg-base-800 hover:text-base-100"
            />
        </NuxtLink>
    </div>
</template>
