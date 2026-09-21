<template>
    <span class="sr-only" role="status">{{ $t('common.loading', [$t('common.activity')]) }}</span>

    <ul class="divide-y divide-default" role="list" aria-hidden="true">
        <li v-for="index in 10" :key="index" class="flex gap-3 p-2">
            <USkeleton class="mt-1 size-7 shrink-0 rounded-full" />

            <div class="min-w-0 flex-1 space-y-1">
                <div class="flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between">
                    <div class="flex flex-wrap items-center gap-2">
                        <USkeleton class="h-5 w-44" />
                        <USkeleton class="h-5 w-20" />
                        <USkeleton class="h-5 w-14" />
                    </div>
                    <USkeleton class="h-5 w-36 shrink-0" />
                </div>

                <div class="flex flex-col gap-2 text-sm sm:flex-row sm:items-center sm:justify-between">
                    <USkeleton class="h-5 w-full max-w-48" />
                    <USkeleton class="h-5 w-32 shrink-0" />
                </div>
            </div>
        </li>
    </ul>
</template>
