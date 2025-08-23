<script setup lang="ts">
import type { Email } from '~~/gen/ts/resources/mailer/email';

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

const emailRefs = ref(new Map<number, Element>());

const selectedEmail = computed({
    get() {
        return props.modelValue;
    },
    set(value: Email | undefined) {
        emit('update:modelValue', value);
    },
});

watch(selectedEmail, async () => {
    if (!selectedEmail.value) {
        return;
    }

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
    <UDashboardPanelContent class="p-0 sm:pb-0">
        <div v-if="!loaded" class="space-y-2">
            <USkeleton class="h-[73px] w-full" />
            <USkeleton class="h-[73px] w-full" />
            <USkeleton class="h-[73px] w-full" />
            <USkeleton class="h-[73px] w-full" />
        </div>

        <template v-else>
            <div v-for="(email, index) in emails" :key="index" :ref="(el) => emailRefs.set(email.id, el as Element)">
                <div
                    class="cursor-pointer border-l-2 p-4 text-sm"
                    :class="[
                        selectedEmail && selectedEmail.id === email.id
                            ? 'border-primary-500 dark:border-primary-400 bg-primary-100 dark:bg-primary-900/25'
                            : 'hover:border-primary-500/25 dark:hover:border-primary-400/25 hover:bg-primary-100/50 dark:hover:bg-primary-900/10 border-white dark:border-gray-900',
                        email.deactivated ? 'dark:bg-red-900/25 border-red-500 bg-red-100 dark:border-red-400' : '',
                    ]"
                    @click="selectedEmail = email"
                >
                    <div
                        class="flex items-center justify-between gap-3"
                        :class="[selectedEmail && selectedEmail.id === email.id && 'font-semibold']"
                    >
                        <span class="truncate font-semibold">
                            {{ email.email }}
                        </span>

                        <UBadge v-if="email.deactivated" color="error" size="xs" :label="$t('common.disabled')" />
                    </div>
                    <div class="flex items-center justify-between">
                        <p>{{ $t('common.label', 1) }}: {{ email.label ?? $t('common.na') }}</p>
                    </div>
                </div>

                <USeparator />
            </div>

            <slot />
        </template>
    </UDashboardPanelContent>
</template>
