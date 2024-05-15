<script setup lang="ts">
import { format, isToday } from 'date-fns';
import type { Mail } from '~~/gen/ts/resources/mailer/mail';
import ProfilePictureImg from '../partials/citizens/ProfilePictureImg.vue';

withDefaults(
    defineProps<{
        mail: Mail;
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
                    :src="mail.from?.avatar?.url"
                    :name="`${mail.from?.firstname} ${mail.from?.lastname}`"
                    size="lg"
                />

                <div class="min-w-0">
                    <p class="font-semibold text-gray-900 dark:text-white">
                        {{ mail.from?.firstname }} {{ mail.from?.lastname }}
                    </p>
                    <p class="font-medium text-gray-500 dark:text-gray-400">
                        {{ mail.subject }}
                    </p>
                </div>
            </div>

            <p class="font-medium text-gray-900 dark:text-white">
                {{
                    isToday(toDate(mail.createdAt))
                        ? format(toDate(mail.createdAt), 'HH:mm')
                        : format(toDate(mail.createdAt), 'dd MMM')
                }}
            </p>
        </div>

        <UDivider class="my-5" />

        <div class="flex-1">
            <p class="text-lg">
                {{ mail.body }}
            </p>
        </div>

        <UDivider class="my-5" />

        <form @submit.prevent>
            <UTextarea
                color="gray"
                required
                size="xl"
                :rows="5"
                :placeholder="$t('components.inbox.reply_to', { name: `${mail.from?.firstname} ${mail.from?.lastname}` })"
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
