<script lang="ts" setup>
import type { Calendar } from '~~/gen/ts/resources/calendar/calendar';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';

defineProps<{
    calendar: Calendar;
}>();

const { isOpen } = useModal();
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template #header>
                <div class="flex flex-col gap-1">
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl font-semibold leading-6">{{ $t('common.calendar') }}: {{ calendar.name }}</h3>

                        <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                    </div>

                    <div class="flex flex-row items-center gap-2">
                        <span>{{ $t('common.creator') }}:</span>
                        <CitizenInfoPopover :user="calendar.creator" show-avatar-in-name />
                    </div>
                </div>
            </template>

            <div>
                <p class="flex-1">
                    {{ $t('common.description') }}:
                    {{ calendar.description ?? $t('common.na') }}
                </p>
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
