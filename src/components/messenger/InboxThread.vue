<script setup lang="ts">
import { format, isToday } from 'date-fns';
import type { Thread } from '~~/gen/ts/resources/messenger/thread';
import ProfilePictureImg from '../partials/citizens/ProfilePictureImg.vue';

withDefaults(
    defineProps<{
        thread: Thread;
        selected?: boolean;
    }>(),
    {
        selected: false,
    },
);
</script>

<template>
    <UDashboardPanelContent>
        <div class="flex justify-between">
            <div class="flex items-center gap-4">
                <ProfilePictureImg
                    :src="thread.from?.avatar?.url"
                    :name="`${thread.from?.firstname} ${thread.from?.lastname}`"
                    size="lg"
                />

                <div class="min-w-0">
                    <p class="font-semibold text-gray-900 dark:text-white">
                        {{ thread.from?.firstname }} {{ thread.from?.lastname }}
                    </p>
                    <p class="font-medium text-gray-500 dark:text-gray-400">
                        {{ thread.title }}
                    </p>
                </div>
            </div>

            <p class="font-medium text-gray-900 dark:text-white">
                {{
                    isToday(toDate(thread.createdAt))
                        ? format(toDate(thread.createdAt), 'HH:mm')
                        : format(toDate(thread.createdAt), 'dd MMM')
                }}
            </p>
        </div>

        <UDivider class="my-5" />

        <div class="flex-1">
            <p class="text-lg">
                {{ thread.body }}
            </p>
        </div>

        <UDivider class="my-5" />

        <form @submit.prevent>
            <UTextarea
                color="gray"
                required
                size="xl"
                :rows="5"
                :placeholder="$t('components.inbox.reply_to', { name: `${thread.from?.firstname} ${thread.from?.lastname}` })"
            >
                <UButton
                    type="submit"
                    color="black"
                    :label="$t('components.inbox.send')"
                    icon="i-mdi-paper-airplane"
                    class="absolute bottom-2.5 right-3.5"
                />
            </UTextarea>
        </form>
    </UDashboardPanelContent>
</template>
