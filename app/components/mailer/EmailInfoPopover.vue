<script lang="ts" setup>
import type { ButtonProps } from '@nuxt/ui';
import EmailBlock from '~/components/partials/citizens/EmailBlock.vue';
import type { ClassProp } from '~/utils/types';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = withDefaults(
    defineProps<{
        email: string | undefined;
        textClass?: ClassProp;
        trailing?: boolean;
        variant?: ButtonProps['variant'];
        color?: ButtonProps['color'];
        hideNaText?: boolean;
    }>(),
    {
        textClass: '',
        trailing: true,
        variant: 'link',
        color: 'neutral',
        hideNaText: false,
    },
);

const notifications = useNotificationsStore();

function copyEmail(): void {
    if (!props.email) return;

    copyToClipboardWrapper(props.email);

    notifications.add({
        title: { key: 'notifications.clipboard.email_address_copied.title', parameters: {} },
        description: { key: 'notifications.clipboard.email_address_copied.content', parameters: {} },
        duration: 3250,
        type: NotificationType.INFO,
    });
}
</script>

<template>
    <template v-if="!email">
        <span class="inline-flex items-center">
            <template v-if="!hideNaText">
                {{ $t('common.na') }}
            </template>
        </span>
    </template>
    <UPopover v-else>
        <UButton
            class="inline-flex w-auto"
            :variant="variant"
            :color="color"
            :trailing-icon="trailing ? 'i-mdi-chevron-down' : undefined"
            :ui="{ base: 'py-0 sm:py-0 px-0 sm:px-0' }"
            v-bind="$attrs"
        >
            <span class="cursor-pointer truncate" :class="textClass"> {{ email ?? $t('common.na') }} </span>
        </UButton>

        <template #content>
            <div class="flex flex-col gap-2 p-4">
                <div class="grid w-full grid-cols-2 gap-2">
                    <EmailBlock :email="email" hide-email />

                    <UButton icon="i-mdi-content-copy" :label="$t('common.copy')" variant="link" @click="copyEmail" />
                </div>

                <div class="flex flex-col gap-2 text-highlighted">
                    <div class="flex flex-col gap-1 text-sm font-normal">
                        <p>
                            <span class="font-semibold">{{ $t('common.mail') }}:</span>
                            {{ email }}
                        </p>
                    </div>
                </div>
            </div>
        </template>
    </UPopover>
</template>
