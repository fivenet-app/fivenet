<script lang="ts" setup>
import type { CalendarEntry } from '~~/gen/ts/resources/calendar/calendar';
import CitizenInfoPopover from '../partials/citizens/CitizenInfoPopover.vue';

defineProps<{
    entry: CalendarEntry;
}>();

const { isOpen } = useModal();
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template #header>
                <div class="flex flex-col gap-1">
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl font-semibold leading-6">
                            {{ entry.title }}
                        </h3>

                        <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                    </div>

                    <p class="flex-1">
                        {{ $d(toDate(entry.startTime), 'long') }} -
                        {{ $d(toDate(entry.endTime), 'long') }}
                    </p>

                    <div class="flex flex-row items-center gap-2">
                        <span>{{ $t('common.creator') }}:</span>
                        <CitizenInfoPopover :user="entry.creator" show-avatar-in-name />
                    </div>
                </div>
            </template>

            <div>
                <p v-html="entry.content"></p>
            </div>

            <template #footer>
                <UButtonGroup class="inline-flex w-full">
                    <UButton color="black" block class="flex-1" @click="isOpen = false">
                        {{ $t('common.close', 1) }}
                    </UButton>
                </UButtonGroup>
            </template>
        </UCard>
    </UModal>
</template>
