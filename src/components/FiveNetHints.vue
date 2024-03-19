<script lang="ts" setup>
import { ArrowLeftBoldCircleIcon, ArrowRightBoldCircleIcon, InformationSlabCircleIcon } from 'mdi-vue3';
import { useCounter } from '@vueuse/core';
import { type RoutesNamedLocations } from '@typed-router';
import GenericBadge from '~/components/partials/elements/GenericBadge.vue';

type Hint = { id: string; keyboard?: boolean; link?: RoutesNamedLocations };

const hints = shuffle([
    {
        id: 'commandpalette',
        keyboard: true,
    },
    {
        id: 'startpage',
        link: { name: 'settings' },
    },
    {
        id: 'documenteditor',
        link: { name: 'settings' },
    },
] as Hint[]);

const hintsCount = hints.length;
const { count, inc, dec, reset, set } = useCounter(0, { min: 0, max: hintsCount - 1 });

const hint = ref<Hint>(hints[count.value]);

watch(count, () => {
    hint.value = hints[count.value];
});

function previousHint(): void {
    if (count.value <= 0) {
        set(hintsCount);
    } else {
        dec();
    }
}

function nextHint(): void {
    if (count.value >= hintsCount - 1) {
        reset();
    } else {
        inc();
    }
}
</script>

<template>
    <div class="pt-2">
        <div class="pointer-events-none inset-x-0 min-w-full max-w-full sm:flex sm:justify-center sm:px-6 sm:pb-5 lg:px-8">
            <div
                class="pointer-events-auto flex items-center justify-between gap-x-6 bg-primary-900 px-6 py-2.5 sm:rounded-xl sm:border-2 sm:border-neutral/20 sm:py-3 sm:pl-4 sm:pr-3.5"
            >
                <button type="button" class="text-white" @click="previousHint()">
                    <ArrowLeftBoldCircleIcon class="h-7 w-7" aria-hidden="true" />
                </button>
                <p class="inline-flex items-center gap-2 text-sm leading-6 text-white">
                    <InformationSlabCircleIcon class="h-7 w-7" aria-hidden="true" />
                    <strong class="mx-1 shrink-0 font-semibold">{{ $t('components.hints.start_text') }}</strong>
                    <span class="grow">{{ $t(`components.hints.${hint.id}.content`) }} </span>
                    <GenericBadge v-if="hint.keyboard" class="ml-1 text-black" color="gray">
                        <kbd class="font-sans">{{ $t(`components.hints.${hint.id}.keyboard`) }}</kbd>
                    </GenericBadge>
                    <NuxtLink v-else-if="hint.link" :to="hint.link" class="ml-1 text-accent-200 underline">
                        {{ $t('components.hints.click_me') }}
                    </NuxtLink>
                </p>
                <button type="button" class="text-white" @click="nextHint()">
                    <ArrowRightBoldCircleIcon class="h-7 w-7" aria-hidden="true" />
                </button>
            </div>
        </div>
    </div>
</template>
