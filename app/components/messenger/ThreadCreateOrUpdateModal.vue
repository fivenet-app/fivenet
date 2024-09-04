<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { useCompletorStore } from '~/store/completor';
import { useMessengerStore } from '~/store/messenger';
import { AccessLevel } from '~~/gen/ts/resources/messenger/access';
import type { Thread } from '~~/gen/ts/resources/messenger/thread';
import type { UserShort } from '~~/gen/ts/resources/users/users';

const props = defineProps<{
    thread?: Thread;
}>();

const { isOpen } = useModal();

const completorStore = useCompletorStore();

const messengerStore = useMessengerStore();

const schema = z.object({
    title: z.string().min(3).max(255),
    archived: z.boolean(),
    users: z
        .object({
            id: z.string().optional(),
            user: z.custom<UserShort>().optional(),
            access: z.nativeEnum(AccessLevel),
        })
        .array()
        .min(1)
        .max(10),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    title: '',
    archived: false,
    users: [
        {
            access: AccessLevel.MESSAGE,
        },
    ],
});

const usersLoading = ref(false);

const accessLevels = ref<{ mode: AccessLevel; selected?: boolean }[]>([
    { mode: AccessLevel.VIEW },
    { mode: AccessLevel.MESSAGE },
    { mode: AccessLevel.MANAGE },
    { mode: AccessLevel.ADMIN },
]);

watch(props, () => setFromProps());
setFromProps();

function setFromProps(): void {
    if (!props.thread) {
        return;
    }

    state.title = props.thread.title;
    state.archived = props.thread.archived;
    state.users = props.thread.access?.users.map((u) => ({ access: u.access, user: u.user })) ?? [];
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;

    const data = event.data;

    await messengerStore
        .createOrUpdateThread({
            thread: {
                id: props.thread?.id ?? '0',
                title: data.title,
                archived: data.archived,
                creatorJob: '',
                access: {
                    jobs: [],
                    users:
                        data.users.map((u) => ({
                            id: u.id ?? '0',
                            threadId: props.thread?.id ?? '0',
                            access: u.access,
                            userId: u.user!.userId,
                        })) ?? [],
                },
            },
        })
        .then(() => (isOpen.value = false))
        .finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl font-semibold leading-6">
                            {{
                                thread
                                    ? $t('components.messenger.ThreadCreateOrUpdateModal.update.title')
                                    : $t('components.messenger.ThreadCreateOrUpdateModal.create.title')
                            }}
                        </h3>

                        <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                    </div>
                </template>

                <div>
                    <UFormGroup name="title" :label="$t('common.title')" class="flex-1" required>
                        <UInput v-model="state.title" type="text" name="title" />
                    </UFormGroup>

                    <UFormGroup name="users" class="flex-1" required>
                        <div class="flex flex-col gap-1">
                            <div v-for="(_, idx) in state.users" :key="idx" class="flex flex-1 items-center gap-1">
                                <div class="flex flex-1 items-center gap-2">
                                    <UFormGroup
                                        :name="`users.${idx}.user`"
                                        :label="idx === 0 ? $t('common.user') : undefined"
                                        class="flex-1"
                                        required
                                    >
                                        <USelectMenu
                                            v-model="state.users[idx]!.user"
                                            :placeholder="$t('common.user')"
                                            block
                                            trailing
                                            by="userId"
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
                                        >
                                            <template #label>
                                                <template v-if="state.users[idx]!.user">
                                                    {{ usersToLabel([state.users[idx]!.user!]) }}
                                                </template>
                                            </template>
                                            <template #option="{ option: user }">
                                                {{ `${user?.firstname} ${user?.lastname} (${user?.dateofbirth})` }}
                                            </template>
                                            <template #option-empty="{ query: search }">
                                                <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                            </template>
                                            <template #empty>
                                                {{ $t('common.not_found', [$t('common.creator', 2)]) }}
                                            </template>
                                        </USelectMenu>
                                    </UFormGroup>

                                    <UFormGroup
                                        :name="`users.${idx}.access`"
                                        :label="idx === 0 ? $t('common.access') : undefined"
                                        class="flex-1"
                                        required
                                    >
                                        <USelectMenu
                                            v-model="state.users[idx]!.access"
                                            :options="accessLevels"
                                            value-attribute="mode"
                                            :searchable-placeholder="$t('common.search_field')"
                                        >
                                            <template #label>
                                                <span class="truncate">{{
                                                    $t(
                                                        `enums.messenger.AccessLevel.${AccessLevel[state.users[idx]!.access ?? 0]}`,
                                                    )
                                                }}</span>
                                            </template>
                                            <template #option="{ option }">
                                                <span class="truncate">{{
                                                    $t(`enums.messenger.AccessLevel.${AccessLevel[option.mode ?? 0]}`)
                                                }}</span>
                                            </template>
                                        </USelectMenu>
                                    </UFormGroup>
                                </div>

                                <UButton
                                    :ui="{ rounded: 'rounded-full' }"
                                    icon="i-mdi-close"
                                    :disabled="!canSubmit || state.users.length <= 1"
                                    @click="state.users.splice(idx, 1)"
                                />
                            </div>
                        </div>

                        <UButton
                            :ui="{ rounded: 'rounded-full' }"
                            icon="i-mdi-plus"
                            :disabled="!canSubmit || state.users.length >= 10"
                            :class="state.users.length ? 'mt-2' : ''"
                            @click="
                                state.users.push({
                                    access: AccessLevel.MESSAGE,
                                })
                            "
                        />
                    </UFormGroup>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton color="black" block class="flex-1" @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton type="submit" block class="flex-1" :disabled="!canSubmit" :loading="!canSubmit">
                            {{ thread ? $t('common.update') : $t('common.create') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>
