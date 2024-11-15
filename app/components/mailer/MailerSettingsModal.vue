<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { useCompletorStore } from '~/store/completor';
import { useMailerStore } from '~/store/mailer';
import type { UserShort } from '~~/gen/ts/resources/users/users';

const { isOpen } = useModal();

const { activeChar } = useAuth();

const completorStore = useCompletorStore();

const mailerStore = useMailerStore();

const schema = z.object({
    users: z.custom<UserShort>().array().max(25),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    users: [],
});

const usersLoading = ref(false);

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;

    const values = event.data;
    await mailerStore
        .setEmailSettings({
            settings: {
                userId: activeChar.value?.userId ?? 0,
                blockedUsers: values.users.map((u) => ({
                    userId: u.userId,
                    user: u,
                })),
            },
        })
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
                            {{ $t('common.settings') }}
                        </h3>

                        <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                    </div>
                </template>

                <div>
                    <UFormGroup name="users" class="flex-1" :label="$t('common.blocklist')">
                        <ClientOnly>
                            <USelectMenu
                                v-model="state.users"
                                :placeholder="$t('common.citizen')"
                                block
                                multiple
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
                                    {{
                                        state.users.length > 0
                                            ? $t('common.citizens', state.users.length)
                                            : $t('common.none_selected', [$t('common.citizen', 2)])
                                    }}
                                </template>

                                <template #option="{ option: user }">
                                    {{ `${user?.firstname} ${user?.lastname} (${user?.dateofbirth})` }}
                                </template>

                                <template #option-empty="{ query: search }">
                                    <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                </template>

                                <template #empty>
                                    {{ $t('common.not_found', [$t('common.citizen', 2)]) }}
                                </template>
                            </USelectMenu>
                        </ClientOnly>
                    </UFormGroup>

                    <div class="grid grid-cols-2 gap-2 lg:grid-cols-4">
                        <div v-for="user in state.users" :key="user.userId">
                            {{ `${user?.firstname} ${user?.lastname} (${user?.dateofbirth})` }}
                        </div>
                    </div>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton color="black" block class="flex-1" @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton type="submit" block class="flex-1" :disabled="!canSubmit" :loading="!canSubmit">
                            {{ $t('common.save') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>
