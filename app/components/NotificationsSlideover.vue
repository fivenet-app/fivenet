<script setup lang="ts">
import DNBToggle from './partials/notification/DNBToggle.vue';
import NotificationsList from './partials/notification/NotificationsList.vue';

const { isNotificationsSlideoverOpen } = useDashboard();
</script>

<template>
    <USlideover v-model="isNotificationsSlideoverOpen" :ui="{ width: 'w-screen sm:max-w-3xl' }">
        <UCard
            class="flex flex-1 flex-col"
            :ui="{
                body: {
                    base: 'flex-1 min-h-[calc(100dvh-(2*var(--header-height)))] max-h-[calc(100dvh-(2*var(--header-height)))] overflow-y-auto',
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
                            color="gray"
                            variant="ghost"
                            icon="i-mdi-window-close"
                            class="-my-1"
                            @click="isNotificationsSlideoverOpen = false"
                        />
                    </div>
                </div>
            </template>

            <div class="flex h-full flex-1 flex-col">
                <NotificationsList @clicked="isNotificationsSlideoverOpen = false" />
            </div>

            <template #footer>
                <UButtonGroup class="inline-flex w-full">
                    <UButton
                        block
                        class="flex-1"
                        @click="
                            navigateTo({ name: 'notifications' });
                            isNotificationsSlideoverOpen = false;
                        "
                    >
                        {{ $t('components.partials.sidebar_notifications') }}
                    </UButton>

                    <UButton color="black" block class="flex-1" @click="isNotificationsSlideoverOpen = false">
                        {{ $t('common.close', 1) }}
                    </UButton>
                </UButtonGroup>
            </template>
        </UCard>
    </USlideover>
</template>
