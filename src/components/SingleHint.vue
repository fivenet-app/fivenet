<script lang="ts" setup generic="T extends RoutesNamesList, P extends string, E extends boolean = false">
import { InformationSlabCircleIcon } from 'mdi-vue3';
import type { NuxtRoute, RoutesNamesList } from '@typed-router';

const props = defineProps<{
    hintId: string;
    layout?: 'horizontal' | 'vertical';
    keyboard?: boolean;
    to?: NuxtRoute<T, P, E>;
    external?: E;
    linkTarget?: '_blank' | null;
}>();
</script>

<template>
    <div class="pointer-events-none inset-x-0 max-w-full sm:flex sm:justify-center sm:pb-2 lg:px-1">
        <div
            class="bg-primary-900 sm:border-neutral/20 pointer-events-auto flex items-center justify-between gap-x-6 px-6 py-2.5 sm:rounded-xl sm:border-2 sm:py-3 sm:pl-4 sm:pr-3.5"
        >
            <p
                class="inline-flex items-center gap-1 text-sm leading-6 text-white"
                :class="layout === 'horizontal' ? 'flex-row' : 'flex-col'"
            >
                <span class="inline-flex items-center">
                    <InformationSlabCircleIcon class="size-7" />
                    <strong class="mx-1 shrink-0 font-semibold">{{ $t('components.hints.start_text') }}</strong>
                </span>
                <span class="grow">{{ $t(`components.hints.${hintId}.content`) }} </span>
                <UKbd v-if="keyboard" class="ml-1">
                    {{ $t(`components.hints.${hintId}.keyboard`) }}
                </UKbd>
                <UButton
                    v-else-if="props.to"
                    variant="link"
                    :to="props.to"
                    :external="props.external"
                    :target="linkTarget ?? null"
                    class="underline"
                >
                    {{ $t('components.hints.click_me') }}
                </UButton>
            </p>
        </div>
    </div>
</template>
