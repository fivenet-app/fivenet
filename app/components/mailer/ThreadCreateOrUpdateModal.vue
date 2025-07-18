<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import TiptapEditor from '~/components/partials/editor/TiptapEditor.vue';
import { useMailerStore } from '~/stores/mailer';
import type { MessageAttachment } from '~~/gen/ts/resources/mailer/message';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import { defaultEmptyContent } from './helpers';
import TemplateSelector from './TemplateSelector.vue';
import ThreadAttachmentsForm from './ThreadAttachmentsForm.vue';

const { isOpen } = useModal();

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

    isOpen.value = false;
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
</script>

<template>
    <UModal fullscreen>
        <UForm class="flex flex-1 flex-col" :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard
                :ui="{
                    ring: '',
                    divide: 'divide-y divide-gray-100 dark:divide-gray-800',
                    base: 'flex flex-1 flex-col',
                    body: { base: 'flex flex-1 flex-col' },
                }"
            >
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl font-semibold leading-6">
                            {{ $t('components.mailer.create_thread') }}
                        </h3>

                        <UButton class="-my-1" color="gray" variant="ghost" icon="i-mdi-window-close" @click="isOpen = false" />
                    </div>
                </template>

                <div class="mx-auto">
                    <div class="flex w-full max-w-screen-xl flex-1 flex-col">
                        <div class="flex w-full flex-col items-center justify-between gap-1">
                            <UFormGroup class="w-full flex-1" name="sender" :label="$t('common.sender')">
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
                                        :options="emails"
                                        :placeholder="$t('common.mail')"
                                        searchable
                                        :searchable-placeholder="$t('common.search_field')"
                                        :search-attributes="['label', 'email']"
                                        trailing
                                        by="id"
                                    >
                                        <template #label>
                                            <span class="overflow-hidden truncate">
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

                                        <template #option="{ option }">
                                            <span class="truncate">
                                                {{
                                                    (option?.label && option?.label !== ''
                                                        ? option?.label + ' (' + option.email + ')'
                                                        : undefined) ??
                                                    (option?.userId
                                                        ? $t('common.personal_email') +
                                                          (isSuperuser ? ' (' + option.email + ')' : '')
                                                        : undefined) ??
                                                    option?.email ??
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

                                        <template #option-empty="{ query: search }">
                                            <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                        </template>

                                        <template #empty> {{ $t('common.not_found', [$t('common.mail', 2)]) }} </template>
                                    </USelectMenu>
                                </ClientOnly>
                            </UFormGroup>

                            <UFormGroup class="w-full flex-1" name="title" :label="$t('common.title')">
                                <UInput
                                    v-model="state.title"
                                    class="font-semibold"
                                    type="text"
                                    size="lg"
                                    :placeholder="$t('common.title')"
                                    :disabled="!canSubmit"
                                    :ui="{ icon: { trailing: { pointer: '' } } }"
                                >
                                    <template #trailing>
                                        <UButton
                                            v-show="state.title !== ''"
                                            color="gray"
                                            variant="link"
                                            icon="i-mdi-close"
                                            :padded="false"
                                            @click="state.title = ''"
                                        />
                                    </template>
                                </UInput>
                            </UFormGroup>

                            <UFormGroup class="w-full flex-1" name="recipients" :label="$t('common.recipient', 2)">
                                <ClientOnly>
                                    <USelectMenu
                                        v-model="state.recipients"
                                        :placeholder="$t('common.recipient')"
                                        block
                                        multiple
                                        trailing
                                        searchable
                                        :options="[
                                            ...state.recipients,
                                            ...addressBook.filter((r) => !state.recipients.includes(r)),
                                        ]"
                                        :searchable-placeholder="$t('common.mail', 1)"
                                        creatable
                                        :disabled="!canSubmit"
                                    >
                                        <template #label>&nbsp;</template>

                                        <template #option-create="{ option }">
                                            <span class="shrink-0">{{ $t('common.recipient') }}: {{ option.label }}</span>
                                        </template>

                                        <template #option-empty="{ query: search }">
                                            <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                        </template>

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
                                        <UButton variant="solid" color="gray" :label="recipient.label" />

                                        <UButton
                                            variant="outline"
                                            icon="i-mdi-close"
                                            color="error"
                                            @click="state.recipients.splice(idx, 1)"
                                        />
                                    </UButtonGroup>
                                </div>
                            </UFormGroup>
                        </div>

                        <UFormGroup
                            class="flex flex-1 flex-col"
                            name="content"
                            :label="$t('common.message')"
                            :ui="{
                                container: 'flex flex-1 flex-col',
                                label: { base: 'flex flex-1' },
                            }"
                        >
                            <template #label>
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
                        </UFormGroup>

                        <ThreadAttachmentsForm
                            v-if="can('documents.DocumentsService/ListDocuments').value"
                            v-model="state.attachments"
                            :can-submit="canSubmit"
                        />
                    </div>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton class="flex-1" block color="black" @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton
                            class="flex-1"
                            type="submit"
                            :disabled="!canSubmit"
                            block
                            :label="$t('components.mailer.send')"
                            trailing-icon="i-mdi-paper-airplane"
                        />
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>
