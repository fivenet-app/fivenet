<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { useCompletorStore } from '~/store/completor';
import { useMailerStore } from '~/store/mailer';
import type { UserShort } from '~~/gen/ts/resources/users/users';
import DocEditor from '../partials/DocEditor.vue';

const { isOpen } = useModal();

const { activeChar } = useAuth();

const mailerStore = useMailerStore();
const { draft: state } = storeToRefs(mailerStore);

const completorStore = useCompletorStore();

const schema = z.object({
    title: z.string().min(3).max(255),
    content: z.string().min(1).max(2048),
    users: z.custom<UserShort>().array().min(1).max(25),
});

type Schema = z.output<typeof schema>;

const usersLoading = ref(false);

async function createThread(values: Schema): Promise<void> {
    await mailerStore.createThread({
        thread: {
            id: '0',
            recipients: [
                {
                    id: '0',
                    emailId: '1',
                    targetId: '2',
                },
            ],
            creatorEmailId: '1',
            creatorId: activeChar.value!.userId,
        },

        message: {
            id: '0',
            threadId: '0',
            title: values.title,
            content: values.content,
            creatorId: activeChar.value!.userId,
            creatorJob: activeChar.value!.job,
        },
    });

    isOpen.value = false;
}

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

                <div class="flex w-full">
                    <div class="flex w-full flex-col gap-2">
                        <div class="flex flex-1 items-center justify-between gap-1">
                            <UFormGroup name="title" :label="$t('common.title')" class="w-full flex-1">
                                <UInput
                                    v-model="state.title"
                                    type="text"
                                    size="xl"
                                    class="font-semibold text-gray-900 dark:text-white"
                                    :disabled="!canSubmit"
                                />
                            </UFormGroup>
                        </div>

                        <div class="min-w-0">
                            <UFormGroup name="users" class="flex-1" :label="$t('common.recipient', 2)">
                                <ClientOnly>
                                    <USelectMenu
                                        v-model="state.users"
                                        :placeholder="$t('common.recipient')"
                                        block
                                        multiple
                                        trailing
                                        :searchable="
                                            async (query: string): Promise<UserShort[]> => {
                                                usersLoading = true;
                                                const users = await completorStore.completeCitizens({
                                                    search: query,
                                                });
                                                usersLoading = false;
                                                return users;
                                            }
                                        "
                                        searchable-lazy
                                        :searchable-placeholder="$t('common.search_field')"
                                        :search-attributes="['firstname', 'lastname']"
                                        :disabled="!canSubmit"
                                    >
                                        <template #label>
                                            {{
                                                state.users.length > 0
                                                    ? $t('common.recipients', state.users.length)
                                                    : $t('common.none_selected', [$t('common.recipient', 2)])
                                            }}
                                        </template>

                                        <template #option="{ option: user }">
                                            {{ `${user?.firstname} ${user?.lastname} (${user?.dateofbirth})` }}
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
                    </div>
                </div>

                <UDivider class="my-2" />

                <UFormGroup
                    :label="$t('common.message', 1)"
                    name="content"
                    :ui="{ wrapper: 'flex flex-1 flex-col', container: 'flex flex-1 flex-col' }"
                >
                    <ClientOnly>
                        <DocEditor v-model="state.content" class="h-full w-full flex-1" :disabled="!canSubmit" />
                    </ClientOnly>
                </UFormGroup>

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
