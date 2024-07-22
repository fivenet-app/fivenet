<script lang="ts" setup>
import { type RoutesNamedLocations } from '@typed-router';

type Hint = { id: string; keyboard?: boolean; to?: RoutesNamedLocations };

const hints = shuffle([
    {
        id: 'commandpalette',
        keyboard: true,
    },
    {
        id: 'startpage',
        to: { name: 'settings' },
    },
    {
        id: 'documenteditor',
        to: { name: 'settings' },
    },
] as Hint[]);
</script>

<template>
    <UCard
        :ui="{
            body: { padding: 'px-2 py-3 sm:p-4' },
            header: { padding: 'px-2 py-3 sm:p-4' },
            footer: { padding: 'px-2 py-2 sm:p-4' },
        }"
    >
        <template #header>
            <div class="inline-flex items-center">
                <UIcon name="i-mdi-information-slab-circle" class="size-6" />
                <span class="ml-1 shrink-0 font-semibold">{{ $t('components.hints.start_text') }}</span>
            </div>
        </template>

        <UCarousel :items="hints" :ui="{ item: 'basis-full' }" arrows>
            <template #default="{ item: hint }">
                <div class="mx-auto mb-2 flex items-center gap-1 text-base">
                    <span class="grow">{{ $t(`components.hints.${hint.id}.content`) }}</span>

                    <div v-if="hint.keyboard || hint.to" class="flex-initial">
                        <UKbd v-if="hint.keyboard" size="md">
                            {{ $t(`components.hints.${hint.id}.keyboard`) }}
                        </UKbd>
                        <UButton v-else-if="hint.to" variant="soft" :to="hint.to">
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
