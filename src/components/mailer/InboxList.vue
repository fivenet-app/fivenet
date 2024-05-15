<script setup lang="ts">
import { format, isToday } from 'date-fns';
import type { Mail } from '~~/gen/ts/resources/mailer/mail';

const props = withDefaults(
    defineProps<{
        modelValue?: Mail;
        mails: Mail[];
    }>(),
    {
        mails: () => [],
    },
);

const emit = defineEmits<{
    (e: 'update:model-value', value: Mail | undefined): void;
}>();

const mailsRefs = ref<Map<string, Element>>(new Map());

const selectedMail = computed({
    get() {
        return props.modelValue;
    },
    set(value: Mail | undefined) {
        emit('update:model-value', value);
    },
});

watch(selectedMail, () => {
    if (!selectedMail.value) {
        return;
    }

    const ref = mailsRefs.value.get(selectedMail.value?.id);
    if (ref) {
        ref.scrollIntoView({ block: 'nearest' });
    }
});

defineShortcuts({
    arrowdown: () => {
        const index = props.mails.findIndex((mail) => mail.id === selectedMail.value?.id);

        if (index === -1) {
            selectedMail.value = props.mails[0];
        } else if (index < props.mails.length - 1) {
            selectedMail.value = props.mails[index + 1];
        }
    },
    arrowup: () => {
        const index = props.mails.findIndex((mail) => mail.id === selectedMail.value?.id);

        if (index === -1) {
            selectedMail.value = props.mails[props.mails.length - 1];
        } else if (index > 0) {
            selectedMail.value = props.mails[index - 1];
        }
    },
});
</script>

<template>
    <UDashboardPanelContent class="p-0">
        <div
            v-for="(mail, index) in mails"
            :key="index"
            :ref="
                (el) => {
                    mailsRefs.set(mail.id, el as Element);
                }
            "
        >
            <div
                class="cursor-pointer border-l-2 p-4 text-sm"
                :class="[
                    mail.unread ? 'text-gray-900 dark:text-white' : 'text-gray-600 dark:text-gray-300',
                    selectedMail && selectedMail.id === mail.id
                        ? 'border-primary-500 dark:border-primary-400 bg-primary-100 dark:bg-primary-900/25'
                        : 'hover:border-primary-500/25 dark:hover:border-primary-400/25 hover:bg-primary-100/50 dark:hover:bg-primary-900/10 border-white dark:border-gray-900',
                ]"
                @click="selectedMail = mail"
            >
                <div class="flex items-center justify-between" :class="[mail.unread && 'font-semibold']">
                    <div class="flex items-center gap-3">
                        {{ mail.from?.firstname }} {{ mail.from?.lastname }}

                        <UChip v-if="mail.unread" />
                    </div>

                    <span>{{
                        isToday(toDate(mail.createdAt))
                            ? format(toDate(mail.createdAt), 'HH:mm')
                            : format(toDate(mail.createdAt), 'dd MMM')
                    }}</span>
                </div>
                <p :class="[mail.unread && 'font-semibold']">
                    {{ mail.subject }}
                </p>
                <p class="line-clamp-1 text-gray-400 dark:text-gray-500">
                    {{ mail.body }}
                </p>
            </div>

            <UDivider />
        </div>
    </UDashboardPanelContent>
</template>
