<script setup lang="ts">
import type { Email } from '~~/gen/ts/resources/mailer/emails/email';

const props = withDefaults(
    defineProps<{
        modelValue?: Email;
        emails: Email[];
        loaded: boolean;
    }>(),
    {
        modelValue: undefined,
    },
);

const emit = defineEmits<{
    (e: 'update:modelValue', value: Email | undefined): void;
}>();

const emailRefs = ref<Map<number, Element>>(new Map<number, Element>());

const selectedEmail = computed({
    get() {
        return props.modelValue;
    },
    set(value: Email | undefined) {
        emit('update:modelValue', value);
    },
});

watch(selectedEmail, async () => {
    if (!selectedEmail.value) return;

    const ref = emailRefs.value.get(selectedEmail.value?.id);
    if (ref) {
        ref.scrollIntoView({ block: 'nearest' });
    }
});

defineShortcuts({
    arrowdown: () => {
        const index = props.emails.findIndex((thread) => thread.id === selectedEmail.value?.id);

        if (index === -1) {
            selectedEmail.value = props.emails[0];
        } else if (index < props.emails.length - 1) {
            selectedEmail.value = props.emails[index + 1];
        }
    },
    arrowup: () => {
        const index = props.emails.findIndex((mail) => mail.id === selectedEmail.value?.id);

        if (index === -1) {
            selectedEmail.value = props.emails[props.emails.length - 1];
        } else if (index > 0) {
            selectedEmail.value = props.emails[index - 1];
        }
    },
});
</script>

<template>
    <UDashboardPanel :ui="{ root: 'min-h-0 pb-(--page-content-bottom-offset)', body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <div v-if="!loaded" class="space-y-2">
            <USkeleton class="h-[73px] w-full" />
            <USkeleton class="h-[73px] w-full" />
            <USkeleton class="h-[73px] w-full" />
            <USkeleton class="h-[73px] w-full" />
        </div>

        <template v-else>
            <div v-for="email in emails" :key="email.id" :ref="(el) => emailRefs.set(email.id, el as Element)">
                <div
                    class="cursor-pointer border-l-2 p-4 text-sm focus-visible:ring-2 focus-visible:ring-primary focus-visible:outline-none focus-visible:ring-inset"
                    role="button"
                    tabindex="0"
                    :aria-current="selectedEmail?.id === email.id ? 'true' : undefined"
                    :class="[
                        selectedEmail && selectedEmail.id === email.id
                            ? 'border-primary-500 bg-primary-100 dark:border-primary-400 dark:bg-primary-900/25'
                            : email.deactivated
                              ? 'border-red-500 bg-red-100 dark:border-red-400 dark:bg-red-900/25'
                              : 'border-default hover:border-primary-500/25 hover:bg-primary-100/50 dark:hover:border-primary-400/25 dark:hover:bg-primary-900/10',
                        email.deactivated && selectedEmail?.id === email.id && 'ring-2 ring-red-500/50',
                    ]"
                    @click="selectedEmail = email"
                    @keydown.enter="selectedEmail = email"
                    @keydown.space.prevent="selectedEmail = email"
                >
                    <div class="flex min-w-0 items-start justify-between gap-3">
                        <div class="min-w-0">
                            <p
                                class="truncate font-semibold"
                                :class="[selectedEmail && selectedEmail.id === email.id && 'text-highlighted']"
                            >
                                {{ email.label || (email.userId ? $t('common.personal_email') : email.email) }}
                            </p>
                            <p class="truncate text-xs text-muted">
                                {{ email.email }}
                            </p>
                        </div>

                        <div class="flex shrink-0 items-center gap-1">
                            <UBadge
                                v-if="email.userId"
                                color="neutral"
                                variant="subtle"
                                size="xs"
                                :label="$t('common.personal_email')"
                            />
                            <UBadge v-if="email.deactivated" color="error" size="xs" :label="$t('common.disabled')" />
                        </div>
                    </div>
                </div>

                <USeparator />
            </div>

            <slot />
        </template>
    </UDashboardPanel>
</template>
