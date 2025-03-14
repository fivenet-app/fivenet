<script lang="ts" setup>
import type { ButtonColor, ButtonVariant } from '#ui/types';
import EmailBlock from '~/components/partials/citizens/EmailBlock.vue';
import { useNotificatorStore } from '~/stores/notificator';
import type { ClassProp } from '~/typings';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = withDefaults(
    defineProps<{
        email: string | undefined;
        textClass?: ClassProp;
        trailing?: boolean;
        variant?: ButtonVariant;
        color?: ButtonColor;
        hideNaText?: boolean;
    }>(),
    {
        textClass: '',
        trailing: true,
        variant: 'solid',
        color: 'gray',
        hideNaText: false,
    },
);

const notifications = useNotificatorStore();

function copyEmail(): void {
    if (!props.email) {
        return;
    }

    copyToClipboardWrapper(props.email);

    notifications.add({
        title: { key: 'notifications.clipboard.email_address_copied.title', parameters: {} },
        description: { key: 'notifications.clipboard.email_address_copied.content', parameters: {} },
        timeout: 3250,
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
    <UPopover v-else :ui="{ trigger: 'inline-flex w-auto', wrapper: 'inline-block' }">
        <UButton v-bind="$attrs" :variant="variant" :color="color" :trailing-icon="trailing ? 'i-mdi-chevron-down' : undefined">
            <span class="truncate" :class="textClass"> {{ email ?? $t('common.na') }} </span>
        </UButton>

        <template #panel>
            <div class="flex flex-col gap-2 p-4">
                <div class="grid w-full grid-cols-2 gap-2">
                    <EmailBlock :email="email" hide-email />

                    <UButton icon="i-mdi-content-copy" :label="$t('common.copy')" variant="link" @click="copyEmail" />
                </div>

                <div class="flex flex-col gap-2 text-gray-900 dark:text-white">
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
