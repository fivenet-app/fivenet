<template>
    <div class="absolute left-1/2 top-10 z-10 inline-flex -translate-x-1/2 -translate-y-1/2 flex-col md:top-8">
        <div
            class="text-primary-500 bg-background inset-0 inline-flex items-center justify-center rounded-md border border-black/20 bg-clip-padding p-1.5"
        >
            <div class="mx-0.5 flex items-center gap-1 text-xs">
                <UIcon class="size-4 animate-spin" name="i-mdi-refresh" />
                {{ $t('components.livemap.restarting_datastream') }}
            </div>
        </div>
    </div>
</template>
