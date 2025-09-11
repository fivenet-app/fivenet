<script setup lang="ts">
import DNBToggle from '~/components/notifications/DNBToggle.vue';
import List from '~/components/notifications/List.vue';

const { isNotificationSlideoverOpen } = useDashboard();
</script>

<template>
    <USlideover
        v-model:open="isNotificationSlideoverOpen"
        :title="$t('common.notification', 2)"
        :ui="{ body: 'flex flex-col p-0 sm:p-0' }"
    >
        <template #actions>
            <DNBToggle />
        </template>

        <template #body>
            <List @clicked="isNotificationSlideoverOpen = false" />
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton
                    class="flex-1"
                    block
                    :label="$t('components.partials.sidebar_notifications')"
                    @click="
                        navigateTo({ name: 'notifications' });
                        isNotificationSlideoverOpen = false;
                    "
                />

                <UButton
                    class="flex-1"
                    color="neutral"
                    block
                    :label="$t('common.close', 1)"
                    @click="isNotificationSlideoverOpen = false"
                />
            </UButtonGroup>
        </template>
    </USlideover>
</template>
