<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import TiptapEditor from '~/components/partials/editor/TiptapEditor.vue';
import { useMailerStore } from '~/stores/mailer';
import type { MessageAttachment } from '~~/gen/ts/resources/mailer/message';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import { defaultEmptyContent } from './helpers';
import TemplateSelector from './TemplateSelector.vue';
import ThreadAttachmentsForm from './ThreadAttachmentsForm.vue';

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const { can, activeChar, isSuperuser } = useAuth();

const notifications = useNotificationsStore();

const mailerStore = useMailerStore();
const { draft: state, addressBook, emails, selectedEmail } = storeToRefs(mailerStore);

const schema = z.object({
    title: z.string().min(3).max(255),
    content: z.string().min(1).max(2048),
    recipients: z
        .object({ label: z.string().min(6).max(80) })
        .array()
        .min(1)
        .max(20)
        .default([]),
    attachments: z.custom<MessageAttachment>().array().max(3).default([]),
});

type Schema = z.output<typeof schema>;

async function createThread(values: Schema): Promise<void> {
    if (!selectedEmail.value?.id) {
        return;
    }

    await mailerStore.createThread({
        thread: {
            id: 0,
            recipients: [],
            creatorEmailId: selectedEmail.value.id,
            creatorId: activeChar.value!.userId,
            title: values.title,
        },

        message: {
            id: 0,
            threadId: 0,
            senderId: selectedEmail.value?.id,
            title: values.title,
            content: {
                rawContent: values.content,
            },
            creatorId: activeChar.value!.userId,
            creatorJob: activeChar.value!.job,
            data: {
                attachments: values.attachments.filter((a) => {
                    if (a.data.oneofKind === 'document') {
                        return a.data.document.id > 0;
                    }

                    return false;
                }),
            },
        },

        recipients: [...new Set(values.recipients.map((r) => r.label.trim()))],
    });

    notifications.add({
        title: { key: 'notifications.action_successful.title', parameters: {} },
        description: { key: 'notifications.action_successful.content', parameters: {} },
        type: NotificationType.SUCCESS,
    });

    // Clear draft data
    state.value.title = '';
    state.value.content = '';
    state.value.recipients = [];
    state.value.attachments = [];

    emit('close', false);
}

function onCreate(item: string): void {
    if (state.value.recipients.findIndex((r) => r.label.toLowerCase() === item.toLowerCase()) !== -1) {
        return;
    }

    state.value.recipients.push({ label: item });
}

onBeforeMount(() => {
    if (
        (!state.value.content || state.value.content === '' || state.value.content === '<p><br></p>') &&
        !!selectedEmail.value?.settings?.signature
    ) {
        state.value.content = defaultEmptyContent + selectedEmail.value?.settings?.signature;
    }
});

watch(
    state.value.recipients,
    () =>
        (state.value.recipients = state.value.recipients.filter(
            (item, idx) => state.value.recipients.findIndex((r) => r.label.toLowerCase() === item.label.toLowerCase()) === idx,
        )),
);

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await createThread(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 1000));
}, 1000);

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UModal :title="$t('components.mailer.create_thread')" fullscreen>
        <template #body>
            <UForm ref="formRef" class="flex flex-1 flex-col" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <div class="mx-auto">
                    <div class="flex w-full max-w-(--breakpoint-xl) flex-1 flex-col">
                        <div class="flex w-full flex-col items-center justify-between gap-1">
                            <UFormField class="w-full flex-1" name="sender" :label="$t('common.sender')">
                                <ClientOnly>
                                    <UInput
                                        v-if="emails.length === 1"
                                        class="pt-1"
                                        type="text"
                                        disabled
                                        :model-value="
                                            (selectedEmail?.label && selectedEmail?.label !== ''
                                                ? selectedEmail?.label + ' (' + selectedEmail.email + ')'
                                                : undefined) ??
                                            selectedEmail?.email ??
                                            $t('common.none')
                                        "
                                    />
                                    <USelectMenu
                                        v-else
                                        v-model="selectedEmail"
                                        class="pt-1"
                                        :items="emails"
                                        :placeholder="$t('common.mail')"
                                        :search-input="{ placeholder: $t('common.search_field') }"
                                        :filter-fields="['label', 'email']"
                                        trailing
                                    >
                                        <template #item-label>
                                            <span class="truncate overflow-hidden">
                                                {{
                                                    (selectedEmail?.label && selectedEmail?.label !== ''
                                                        ? selectedEmail?.label + ' (' + selectedEmail.email + ')'
                                                        : undefined) ??
                                                    selectedEmail?.email ??
                                                    $t('common.none')
                                                }}

                                                <UBadge
                                                    v-if="selectedEmail?.deactivated"
                                                    color="error"
                                                    size="xs"
                                                    :label="$t('common.disabled')"
                                                />
                                            </span>
                                        </template>

                                        <template #item="{ item }">
                                            <span class="truncate">
                                                {{
                                                    (item?.label && item?.label !== ''
                                                        ? item?.label + ' (' + item.email + ')'
                                                        : undefined) ??
                                                    (item?.userId
                                                        ? $t('common.personal_email') +
                                                          (isSuperuser ? ' (' + item.email + ')' : '')
                                                        : undefined) ??
                                                    item?.email ??
                                                    $t('common.none')
                                                }}
                                            </span>

                                            <UBadge
                                                v-if="selectedEmail?.deactivated"
                                                color="error"
                                                size="xs"
                                                :label="$t('common.disabled')"
                                            />
                                        </template>

                                        <template #empty> {{ $t('common.not_found', [$t('common.mail', 2)]) }} </template>
                                    </USelectMenu>
                                </ClientOnly>
                            </UFormField>

                            <UFormField class="w-full flex-1" name="title" :label="$t('common.title')">
                                <UInput
                                    v-model="state.title"
                                    class="font-semibold"
                                    type="text"
                                    size="lg"
                                    :placeholder="$t('common.title')"
                                    :disabled="!canSubmit"
                                    :ui="{ trailing: 'pe-1' }"
                                >
                                    <template #trailing>
                                        <UButton
                                            v-if="state.title !== ''"
                                            color="neutral"
                                            variant="link"
                                            icon="i-mdi-close"
                                            aria-controls="search"
                                            @click="state.title = ''"
                                        />
                                    </template>
                                </UInput>
                            </UFormField>

                            <UFormField class="w-full flex-1" name="recipients" :label="$t('common.recipient', 2)">
                                <ClientOnly>
                                    <USelectMenu
                                        v-model="state.recipients"
                                        :placeholder="$t('common.recipient')"
                                        block
                                        multiple
                                        trailing
                                        :items="[
                                            ...state.recipients,
                                            ...addressBook.filter((r) => !state.recipients.includes(r)),
                                        ]"
                                        :search-input="{ placeholder: $t('common.mail', 1) }"
                                        :disabled="!canSubmit"
                                        create-item
                                        @create="(item) => onCreate(item)"
                                    >
                                        <template #empty>
                                            {{ $t('common.not_found', [$t('common.recipient', 2)]) }}
                                        </template>
                                    </USelectMenu>
                                </ClientOnly>

                                <div class="mt-2 flex snap-x flex-row flex-wrap gap-2 overflow-x-auto">
                                    <UButtonGroup
                                        v-for="(recipient, idx) in state.recipients"
                                        :key="idx"
                                        size="sm"
                                        orientation="horizontal"
                                    >
                                        <UButton variant="solid" color="neutral" :label="recipient.label" />

                                        <UButton
                                            variant="outline"
                                            icon="i-mdi-close"
                                            color="error"
                                            @click="state.recipients.splice(idx, 1)"
                                        />
                                    </UButtonGroup>
                                </div>
                            </UFormField>
                        </div>

                        <UFormField
                            class="flex flex-1 flex-col"
                            name="content"
                            :label="$t('common.message')"
                            :ui="{
                                container: 'flex flex-1 flex-col',
                            }"
                        >
                            <template #item-label>
                                <div class="flex flex-1 flex-col items-center sm:flex-row">
                                    <span class="flex-1">{{ $t('common.message', 1) }}</span>

                                    <TemplateSelector v-model="state.content" class="ml-auto" />
                                </div>
                            </template>

                            <ClientOnly>
                                <TiptapEditor
                                    v-model="state.content"
                                    class="flex-1 overflow-y-hidden"
                                    :disabled="!canSubmit"
                                    wrapper-class="min-h-96"
                                />
                            </ClientOnly>
                        </UFormField>

                        <ThreadAttachmentsForm
                            v-if="can('documents.DocumentsService/ListDocuments').value"
                            v-model="state.attachments"
                            :can-submit="canSubmit"
                        />
                    </div>
                </div>
            </UForm>
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton class="flex-1" block color="neutral" :label="$t('common.close', 1)" @click="$emit('close', false)" />

                <UButton
                    class="flex-1"
                    :disabled="!canSubmit"
                    block
                    :label="$t('components.mailer.send')"
                    trailing-icon="i-mdi-paper-airplane"
                    @click="() => formRef?.submit()"
                />
            </UButtonGroup>
        </template>
    </UModal>
</template>
