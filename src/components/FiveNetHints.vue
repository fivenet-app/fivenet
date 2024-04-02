<script lang="ts" setup>
import { InformationSlabCircleIcon } from 'mdi-vue3';
import { type RoutesNamedLocations } from '@typed-router';

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
</script>

<template>
    <UCard>
        <template #header>
            <div class="inline-flex items-center">
                <InformationSlabCircleIcon class="size-7" />
                <strong class="mx-1 shrink-0 font-semibold">{{ $t('components.hints.start_text') }}</strong>
            </div>
        </template>

        <UCarousel :items="hints" :ui="{ item: 'basis-full' }" arrows class="mx-auto">
            <template #default="{ item: hint }">
                <div class="mx-2 mb-4 flex items-center gap-1 text-base">
                    <span class="grow">{{ $t(`components.hints.${hint.id}.content`) }}</span>
                    <div class="self-end">
                        <UKbd v-if="hint.keyboard" size="md">
                            {{ $t(`components.hints.${hint.id}.keyboard`) }}
                        </UKbd>
                        <UButton v-else-if="hint.link" variant="link" :to="hint.link" class="underline">
                            {{ $t('components.hints.click_me') }}
                        </UButton>
                    </div>
                </div>
            </template>

            <template #prev="{ onClick, disabled }">
                <UButton :disabled="disabled" @click="onClick">{{ $t('common.previous') }}</UButton>
            </template>

            <template #next="{ onClick, disabled }">
                <UButton :disabled="disabled" @click="onClick">{{ $t('common.next') }}</UButton>
            </template>
        </UCarousel>
    </UCard>
</template>
