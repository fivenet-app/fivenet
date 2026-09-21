<script lang="ts" setup>
withDefaults(
    defineProps<{
        icon: string;
        iconClass?: string;
        inline?: boolean;
    }>(),
    {
        iconClass: undefined,
        inline: false,
    },
);
</script>

<template>
    <li
        class="flex-initial border-default p-2 hover:border-primary-500/25 hover:bg-primary-100/50 dark:hover:border-primary-400/25 dark:hover:bg-primary-900/10"
    >
        <div class="flex gap-3">
            <div class="my-auto flex size-7 shrink-0 items-center justify-center rounded-full">
                <UIcon :class="[iconClass, 'size-full']" :name="icon" :inline="inline" />
            </div>

            <div class="min-w-0 flex-1 space-y-1">
                <div class="flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between">
                    <h3 class="min-w-0 text-sm font-medium break-words">
                        <slot name="title" />
                    </h3>

                    <p class="shrink-0 text-sm text-dimmed">
                        <slot name="timestamp" />
                    </p>
                </div>

                <slot />
            </div>
        </div>
    </li>
</template>
