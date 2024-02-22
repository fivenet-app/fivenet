<script lang="ts" setup></script>

<template>
    <div class="mt-3 overflow-hidden bg-base-800 text-neutral shadow sm:rounded-lg">
        <div id="settings" class="px-4 py-5 sm:px-6">
            <h3 class="text-base font-semibold leading-6">
                <slot name="title" />
            </h3>
            <p class="mt-1 max-w-2xl text-sm">
                <slot name="description" />
            </p>
        </div>
        <div class="border-t border-base-400 px-4 py-5 sm:p-0">
            <dl class="sm:divide-y sm:divide-base-400">
                <slot />
            </dl>
        </div>
    </div>
</template>
