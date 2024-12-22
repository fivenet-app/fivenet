<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { useMailerStore } from '~/store/mailer';
import { useNotificatorStore } from '~/store/notificator';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import TiptapEditor from '../partials/editor/TiptapEditor.vue';
import TemplateSelector from './TemplateSelector.vue';
import { defaultEmptyContent } from './helpers';

const { isOpen } = useModal();

const { activeChar, isSuperuser } = useAuth();

const notifications = useNotificatorStore();

const mailerStore = useMailerStore();
const { draft: state, addressBook, emails, selectedEmail } = storeToRefs(mailerStore);

const schema = z.object({
    title: z.string().min(3).max(255),
    content: z.string().min(1).max(2048),
    recipients: z
        .object({ label: z.string().min(6).max(80) })
        .array()
        .min(1)
        .max(20),
});

type Schema = z.output<typeof schema>;

async function createThread(values: Schema): Promise<void> {
    if (!selectedEmail.value?.id) {
        return;
    }

    await mailerStore.createThread({
        thread: {
            id: '0',
            recipients: [],
            creatorEmailId: selectedEmail.value.id,
            creatorId: activeChar.value!.userId,
            title: values.title,
        },

        message: {
            id: '0',
            threadId: '0',
            senderId: selectedEmail.value?.id,
            title: values.title,
            content: {
                rawContent: values.content,
            },
            creatorId: activeChar.value!.userId,
            creatorJob: activeChar.value!.job,
        },

        recipients: [...new Set(values.recipients.map((r) => r.label.trim()))],
    });

    notifications.add({
        title: { key: 'notifications.action_successfull.title', parameters: {} },
        description: { key: 'notifications.action_successfull.content', parameters: {} },
        type: NotificationType.SUCCESS,
    });

    // Clear draft data
    state.value.title = '';
    state.value.content = '';
    state.value.recipients = [];

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
        <UForm :schema="schema" :state="state" class="flex flex-1 flex-col" @submit="onSubmitThrottle">
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

                        <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                    </div>
                </template>

                <div class="mx-auto flex w-full max-w-screen-xl flex-1 overflow-y-hidden">
                    <div class="flex w-full flex-1 flex-col gap-2 overflow-y-hidden">
                        <div class="flex w-full flex-col items-center justify-between gap-1">
                            <UFormGroup name="sender" :label="$t('common.sender')" class="w-full flex-1">
                                <ClientOnly>
                                    <UInput
                                        v-if="emails.length === 1"
                                        type="text"
                                        disabled
                                        :model-value="
                                            (selectedEmail?.label && selectedEmail?.label !== ''
                                                ? selectedEmail?.label + ' (' + selectedEmail.email + ')'
                                                : undefined) ??
                                            selectedEmail?.email ??
                                            $t('common.none')
                                        "
                                        class="pt-1"
                                    />
                                    <USelectMenu
                                        v-else
                                        v-model="selectedEmail"
                                        :options="emails"
                                        :placeholder="$t('common.mail')"
                                        searchable
                                        :searchable-placeholder="$t('common.search_field')"
                                        :search-attributes="['label', 'email']"
                                        trailing
                                        class="pt-1"
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
                                                    color="red"
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
                                                color="red"
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

                            <UFormGroup name="title" :label="$t('common.title')" class="w-full flex-1">
                                <UInput
                                    v-model="state.title"
                                    type="text"
                                    size="lg"
                                    class="font-semibold"
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

                            <UFormGroup name="recipients" class="w-full flex-1" :label="$t('common.recipient', 2)">
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
                                            color="red"
                                            @click="state.recipients.splice(idx, 1)"
                                        />
                                    </UButtonGroup>
                                </div>
                            </UFormGroup>
                        </div>

                        <UFormGroup
                            name="content"
                            :label="$t('common.message')"
                            class="flex flex-1 flex-col"
                            :ui="{
                                container: 'flex flex-1',
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
                    </div>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton block class="flex-1" color="black" @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton
                            type="submit"
                            :disabled="!canSubmit"
                            block
                            class="flex-1"
                            :label="$t('components.mailer.send')"
                            trailing-icon="i-mdi-paper-airplane"
                        />
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>
