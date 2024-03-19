<script lang="ts" setup>
import { InformationSlabCircleIcon } from 'mdi-vue3';
import type { RoutesNamedLocations } from '@typed-router';
import GenericBadge from '~/components/partials/elements/GenericBadge.vue';

withDefaults(
    defineProps<{
        id: string;
        keyboard?: boolean;
        link?: RoutesNamedLocations;
        layout?: 'horizontal' | 'vertical';
    }>(),
    {
        keyboard: false,
        link: undefined,
        layout: 'vertical',
    },
);
</script>

<template>
    <div class="pointer-events-none inset-x-0 min-w-full max-w-full sm:flex sm:justify-center sm:pb-2 lg:px-1">
        <div
            class="pointer-events-auto flex items-center justify-between gap-x-6 bg-primary-900 px-6 py-2.5 sm:rounded-xl sm:border-2 sm:border-neutral/20 sm:py-3 sm:pl-4 sm:pr-3.5"
        >
            <p
                class="inline-flex items-center gap-1 text-sm leading-6 text-white"
                :class="layout === 'horizontal' ? 'flex-row' : 'flex-col'"
            >
                <span class="inline-flex items-center">
                    <InformationSlabCircleIcon class="h-7 w-7" aria-hidden="true" />
                    <strong class="mx-1 shrink-0 font-semibold">{{ $t('components.hints.start_text') }}</strong>
                </span>
                <span class="grow">{{ $t(`components.hints.${id}.content`) }} </span>
                <GenericBadge v-if="keyboard" class="ml-1 text-black" color="gray">
                    <kbd class="font-sans">{{ $t(`components.hints.${id}.keyboard`) }}</kbd>
                </GenericBadge>
                <NuxtLink v-else-if="link" :to="link" class="ml-1 text-accent-200 underline">
                    {{ $t('components.hints.click_me') }}
                </NuxtLink>
            </p>
        </div>
    </div>
</template>
