<template>
    <div class="m-2 overflow-hidden bg-base-800 shadow sm:rounded-lg">
        <div id="settings" class="px-4 py-5 sm:px-6">
            <h3 class="text-base font-semibold leading-6">
                <slot name="title" />
            </h3>
            <p class="mt-1 max-w-2xl text-sm">
                <slot name="description" />
            </p>
        </div>
        <div class="border-t border-base-400 px-4 py-5 sm:p-0">
            <dl class="divide-y divide-gray-200 dark:divide-gray-700">
                <slot />
            </dl>
        </div>
    </div>
</template>
