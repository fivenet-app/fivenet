<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import SelectMenu from '~/components/partials/SelectMenu.vue';
import { useAuthStore } from '~/stores/auth';
import { useCompletorStore } from '~/stores/completor';
import { getCalendarCalendarClient } from '~~/gen/ts/clients';
import type { UserShort } from '~~/gen/ts/resources/users/users';
import type { ShareCalendarEntryResponse } from '~~/gen/ts/services/calendar/calendar';

const props = defineProps<{
    entryId: number;
}>();

const emit = defineEmits<{
    (e: 'close'): void;
    (e: 'refresh'): void;
}>();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const completorStore = useCompletorStore();

const calendarCalendarClient = await getCalendarCalendarClient();

const schema = z.object({
    users: z.custom<UserShort>().array().max(20).default([]),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    users: [],
});

async function shareCalendarEntry(values: Schema): Promise<undefined | ShareCalendarEntryResponse> {
    if (values.users.length === 0) {
        emit('close');
        return;
    }

    const call = calendarCalendarClient.shareCalendarEntry({
        entryId: props.entryId,
        userIds: values.users.map((u) => u.userId),
    });
    const { response } = await call;

    emit('refresh');
    emit('close');

    values.users.length = 0;

    return response;
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await shareCalendarEntry(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UCard>
        <template #header>
            <div class="flex items-center justify-between">
                <h3 class="text-xl leading-6 font-semibold">
                    {{ $t('components.calendar.EntryShareModal.title') }}
                </h3>
            </div>
        </template>

        <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UFormField class="flex-1" name="participants" :label="$t('common.guest', 2)">
                <SelectMenu
                    v-model="state.users"
                    multiple
                    :searchable="
                        async (q: string) => {
                            const users = await completorStore.completeCitizens({
                                search: q,
                                userIds: state.users.map((u) => u.userId),
                            });
                            return users.filter((u) => u.userId !== activeChar?.userId);
                        }
                    "
                    searchable-key="completor-citizens"
                    :search-input="{ placeholder: $t('common.search_field') }"
                    :filter-fields="['firstname', 'lastname']"
                    block
                    :placeholder="$t('common.citizen', 2)"
                    trailing
                >
                    <template #item-label>
                        {{ $t('common.selected', state.users.length) }}
                    </template>

                    <template #item="{ item }">
                        {{ `${item?.firstname} ${item?.lastname} (${item?.dateofbirth})` }}
                    </template>

                    <template #empty> {{ $t('common.not_found', [$t('common.citizen', 2)]) }} </template>
                </SelectMenu>
            </UFormField>

            <div class="dark:bg-base-900 mt-2 overflow-hidden rounded-md bg-neutral-100">
                <ul class="grid grid-cols-2 text-sm font-medium text-gray-100 lg:grid-cols-3" role="list">
                    <li
                        v-for="user in state.users"
                        :key="user.userId"
                        class="flex items-center border-b border-neutral-100 px-4 py-2 dark:border-neutral-800"
                    >
                        <CitizenInfoPopover :user="user" show-avatar show-avatar-in-name />
                    </li>
                </ul>
            </div>
        </UForm>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" block :label="$t('common.cancel', 1)" @click="$emit('close')" />

                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    :label="$t('common.save')"
                    @click="formRef?.submit()"
                />
            </UButtonGroup>
        </template>
    </UCard>
</template>
