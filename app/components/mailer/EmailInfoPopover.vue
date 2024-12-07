<script lang="ts" setup>
import type { ButtonColor, ButtonVariant } from '#ui/types';
import type { ClassProp } from '~/typings';
import EmailBlock from '../partials/citizens/EmailBlock.vue';

withDefaults(
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
</script>

<template>
    <template v-if="!email">
        <span class="inline-flex items-center">
            <span v-if="!hideNaText">
                {{ $t('common.na') }}
            </span>
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

                    <UButton
                        icon="i-mdi-content-copy"
                        :label="$t('common.copy')"
                        variant="link"
                        @click="copyToClipboardWrapper(email)"
                    />
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
