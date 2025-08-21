<script setup lang="ts">
import DNBToggle from '~/components/notifications/DNBToggle.vue';
import NotificationList from '~/components/notifications/NotificationList.vue';

const { isNotificationSlideoverOpen } = useDashboard();
</script>

<template>
    <USlideover v-model="isNotificationSlideoverOpen" :ui="{ width: 'w-screen sm:max-w-3xl' }">
        <UCard
            class="flex flex-1 flex-col"
            :ui="{
                body: {
                    base: 'flex flex-1 min-h-[calc(100dvh-(2*var(--header-height)))] max-h-[calc(100dvh-(2*var(--header-height)))] overflow-y-auto',
                    padding: '',
                },
                ring: '',
                divide: 'divide-y divide-gray-100 dark:divide-gray-800',
            }"
        >
            <template #header>
                <div class="flex flex-col gap-1">
                    <div class="flex items-center justify-between">
                        <h3 class="inline-flex gap-2 text-2xl font-semibold leading-6">
                            {{ $t('common.notification', 2) }}
                        </h3>

                        <DNBToggle />

                        <UButton
                            class="-my-1"
                            color="gray"
                            variant="ghost"
                            icon="i-mdi-window-close"
                            @click="isNotificationSlideoverOpen = false"
                        />
                    </div>
                </div>
            </template>

            <div class="flex flex-1 flex-col">
                <NotificationList @clicked="isNotificationSlideoverOpen = false" />
            </div>

            <template #footer>
                <UButtonGroup class="inline-flex w-full">
                    <UButton
                        class="flex-1"
                        block
                        @click="
                            navigateTo({ name: 'notifications' });
                            isNotificationSlideoverOpen = false;
                        "
                    >
                        {{ $t('components.partials.sidebar_notifications') }}
                    </UButton>

                    <UButton class="flex-1" color="black" block @click="isNotificationSlideoverOpen = false">
                        {{ $t('common.close', 1) }}
                    </UButton>
                </UButtonGroup>
            </template>
        </UCard>
    </USlideover>
</template>
