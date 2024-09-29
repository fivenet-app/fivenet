<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import { useAuthStore } from '~/store/auth';
import { useCompletorStore } from '~/store/completor';
import type { UserShort } from '~~/gen/ts/resources/users/users';
import type { ShareCalendarEntryResponse } from '~~/gen/ts/services/calendar/calendar';

const props = defineProps<{
    entryId: string;
}>();

const emits = defineEmits<{
    (e: 'close'): void;
    (e: 'refresh'): void;
}>();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const completorStore = useCompletorStore();

const usersLoading = ref(false);

const schema = z.object({
    users: z.custom<UserShort>().array().max(20),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    users: [],
});

async function shareCalendarEntry(values: Schema): Promise<undefined | ShareCalendarEntryResponse> {
    if (values.users.length === 0) {
        emits('close');
        return;
    }

    const call = getGRPCCalendarClient().shareCalendarEntry({
        entryId: props.entryId,
        userIds: values.users.map((u) => u.userId),
    });
    const { response } = await call;

    emits('refresh');
    emits('close');

    values.users.length = 0;

    return response;
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await shareCalendarEntry(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
        <UCard
            :ui="{
                ring: '',
                divide: 'divide-y divide-gray-100 dark:divide-gray-800',
                body: { padding: 'px-2 py-3 sm:p-3' },
                header: { padding: 'px-2 py-3 sm:p-3' },
                footer: { padding: 'px-2 py-2 sm:p-3' },
            }"
        >
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-xl font-semibold leading-6">
                        {{ $t('components.calendar.EntryShareModal.title') }}
                    </h3>
                </div>
            </template>

            <div>
                <UFormGroup name="participants" :label="$t('common.guest', 2)" class="flex-1">
                    <ClientOnly>
                        <USelectMenu
                            v-model="state.users"
                            multiple
                            :searchable="
                                async (query: string) => {
                                    usersLoading = true;
                                    const users = await completorStore.completeCitizens({
                                        search: query,
                                    });
                                    usersLoading = false;
                                    return users.filter((u) => u.userId !== activeChar?.userId);
                                }
                            "
                            searchable-lazy
                            :searchable-placeholder="$t('common.search_field')"
                            :search-attributes="['firstname', 'lastname']"
                            block
                            :placeholder="$t('common.citizen', 2)"
                            trailing
                            by="userId"
                        >
                            <template #option="{ option: user }">
                                {{ `${user?.firstname} ${user?.lastname} (${user?.dateofbirth})` }}
                            </template>
                            <template #option-empty="{ query: search }">
                                <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                            </template>
                            <template #empty> {{ $t('common.not_found', [$t('common.citizen', 2)]) }} </template>
                        </USelectMenu>
                    </ClientOnly>
                </UFormGroup>

                <div class="mt-2 overflow-hidden rounded-md bg-base-900">
                    <ul role="list" class="grid grid-cols-2 text-sm font-medium text-gray-100 lg:grid-cols-3">
                        <li
                            v-for="user in state.users"
                            :key="user.userId"
                            class="flex items-center border-b border-gray-100 px-4 py-2 dark:border-gray-800"
                        >
                            <CitizenInfoPopover :user="user" show-avatar show-avatar-in-name />
                        </li>
                    </ul>
                </div>
            </div>

            <template #footer>
                <UButtonGroup class="inline-flex w-full">
                    <UButton color="black" block class="flex-1" @click="$emit('close')">
                        {{ $t('common.cancel', 1) }}
                    </UButton>

                    <UButton type="submit" block class="flex-1" :disabled="!canSubmit" :loading="!canSubmit">
                        {{ $t('common.save') }}
                    </UButton>
                </UButtonGroup>
            </template>
        </UCard>
    </UForm>
</template>
