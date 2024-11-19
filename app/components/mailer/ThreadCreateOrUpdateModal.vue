<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { useMailerStore } from '~/store/mailer';
import DocEditor from '../partials/DocEditor.vue';
import TemplateSelector from './TemplateSelector.vue';
import { defaultEmptyContent } from './helpers';

const { isOpen } = useModal();

const { activeChar } = useAuth();

const mailerStore = useMailerStore();
const { draft: state, selectedEmail } = storeToRefs(mailerStore);

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
            content: values.content,
            creatorId: activeChar.value!.userId,
            creatorJob: activeChar.value!.job,
        },

        recipients: [...new Set(values.recipients.map((r) => r.label))],
    });

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

                <div class="mx-auto flex w-full max-w-screen-xl">
                    <div class="flex w-full flex-col gap-2">
                        <div class="flex flex-1 flex-col items-center justify-between gap-1">
                            <UFormGroup name="sender" :label="$t('common.sender')" class="w-full flex-1">
                                <UInput type="text" disabled :model-value="selectedEmail?.email ?? 'N/A'" />
                            </UFormGroup>

                            <UFormGroup name="title" :label="$t('common.title')" class="w-full flex-1">
                                <UInput
                                    v-model="state.title"
                                    type="text"
                                    size="lg"
                                    class="font-semibold"
                                    :placeholder="$t('common.title')"
                                    :disabled="!canSubmit"
                                />
                            </UFormGroup>

                            <UFormGroup name="recipients" class="w-full flex-1" :label="$t('common.recipient', 2)">
                                <ClientOnly>
                                    <USelectMenu
                                        v-model="state.recipients"
                                        :placeholder="$t('common.recipient')"
                                        block
                                        multiple
                                        trailing
                                        value-attribute="label"
                                        searchable
                                        :options="state.recipients"
                                        :searchable-placeholder="$t('common.mail', 1)"
                                        creatable
                                        :disabled="!canSubmit"
                                    >
                                        <template #label>
                                            {{
                                                state.recipients.length > 0
                                                    ? state.recipients.map((r) => r.label).join(', ')
                                                    : $t('common.none_selected', [$t('common.recipient', 2)])
                                            }}
                                        </template>

                                        <template #option-create="{ option }">
                                            <span class="flex-shrink-0">{{ $t('common.recipient') }}: {{ option.label }}</span>
                                        </template>

                                        <template #option-empty="{ query: search }">
                                            <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                        </template>

                                        <template #empty>
                                            {{ $t('common.not_found', [$t('common.recipient', 2)]) }}
                                        </template>
                                    </USelectMenu>
                                </ClientOnly>
                            </UFormGroup>
                        </div>

                        <UDivider class="my-2" />

                        <UFormGroup
                            :label="$t('common.message', 1)"
                            name="content"
                            :ui="{
                                wrapper: 'flex flex-1 flex-col',
                                container: 'flex flex-1 flex-col',
                                label: { base: 'flex flex-1' },
                            }"
                        >
                            <template #label>
                                <div class="flex flex-1 flex-col items-center sm:flex-row">
                                    <span>{{ $t('common.message', 2) }}</span>

                                    <TemplateSelector v-model="state.content" class="ml-auto" />
                                </div>
                            </template>

                            <ClientOnly>
                                <DocEditor v-model="state.content" class="h-full w-full flex-1" :disabled="!canSubmit" />
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
