<template>
    <div class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-5 sm:py-4">
        <dt class="text-sm font-medium">
            <slot name="title" />
        </dt>
        <dd class="mt-1 text-sm sm:col-span-2 sm:mt-0">
            <slot />
        </dd>
    </div>
</template>
