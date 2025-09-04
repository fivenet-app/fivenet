<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import InputDatePicker from '~/components/partials/InputDatePicker.vue';
import SelectMenu from '~/components/partials/SelectMenu.vue';
import { useAuthStore } from '~/stores/auth';
import { useCompletorStore } from '~/stores/completor';
import { getJobsConductClient } from '~~/gen/ts/clients';
import { type ConductEntry, ConductType } from '~~/gen/ts/resources/jobs/conduct';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { UserShort } from '~~/gen/ts/resources/users/users';
import ColleagueName from '../colleagues/ColleagueName.vue';
import { conductTypesToBadgeColor } from './helpers';

const props = defineProps<{
    entry?: ConductEntry;
    userId?: number;
}>();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
    (e: 'created', entry: ConductEntry): void;
    (e: 'updated', entry: ConductEntry): void;
}>();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const completorStore = useCompletorStore();

const notifications = useNotificationsStore();

const jobsConductClient = await getJobsConductClient();

const cTypes = ref<{ status: ConductType }[]>([
    { status: ConductType.NOTE },
    { status: ConductType.NEUTRAL },
    { status: ConductType.POSITIVE },
    { status: ConductType.NEGATIVE },
    { status: ConductType.WARNING },
    { status: ConductType.SUSPENSION },
]);

const schema = z.object({
    targetUser: z.custom<UserShort>(),
    type: z.nativeEnum(ConductType),
    message: z.string().min(3).max(2000),
    expiresAt: z.date().optional(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Partial<Schema>>({
    targetUser: undefined,
    type: ConductType.NOTE,
    message: '',
    expiresAt: undefined,
});

async function conductCreateOrUpdateEntry(values: Schema, id?: number): Promise<void> {
    try {
        const req = {
            entry: {
                id: id ?? 0,
                job: '',
                creatorId: activeChar.value?.userId ?? 0,
                type: values.type,
                message: values.message,
                targetUserId: values.targetUser.userId,
                expiresAt: values.expiresAt ? toTimestamp(values.expiresAt) : undefined,
            },
        };

        if (id === undefined) {
            const call = jobsConductClient.createConductEntry(req);
            const { response } = await call;

            emit('created', response.entry!);
        } else {
            const call = jobsConductClient.updateConductEntry(req);
            const { response } = await call;

            emit('updated', response.entry!);
        }

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        emit('close', false);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function setFromProps(): Promise<void> {
    state.targetUser = props.entry?.targetUser;
    state.type = props.entry?.type ?? ConductType.NOTE;
    state.message = props.entry?.message ?? '';
    state.expiresAt = props.entry?.expiresAt ? toDate(props.entry?.expiresAt) : undefined;
}

watch(props, () => setFromProps());

onBeforeMount(() => setFromProps());

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await conductCreateOrUpdateEntry(event.data, props.entry?.id).finally(() =>
        useTimeoutFn(() => (canSubmit.value = true), 400),
    );
}, 1000);

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UModal
        :title="
            entry === undefined
                ? $t('components.jobs.conduct.CreateOrUpdateModal.create.title')
                : $t('components.jobs.conduct.CreateOrUpdateModal.update.title')
        "
    >
        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <dl class="divide-neutral/10 divide-y">
                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm leading-6 font-medium">
                            <label class="block text-sm leading-6 font-medium" for="type">
                                {{ $t('common.type') }}
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <UFormField name="type">
                                <ClientOnly>
                                    <USelectMenu
                                        v-model="state.type"
                                        :items="cTypes"
                                        value-key="status"
                                        :search-input="{ placeholder: $t('common.search_field') }"
                                    >
                                        <template #item-label="{ item }">
                                            <UBadge :color="conductTypesToBadgeColor(item.status)" truncate>
                                                {{ $t(`enums.jobs.ConductType.${ConductType[item.status ?? 0]}`) }}
                                            </UBadge>
                                        </template>

                                        <template #item="{ item }">
                                            <UBadge :color="conductTypesToBadgeColor(item.status)" truncate>
                                                {{ $t(`enums.jobs.ConductType.${ConductType[item.status ?? 0]}`) }}
                                            </UBadge>
                                        </template>

                                        <template #empty>
                                            {{ $t('common.not_found', [$t('common.type', 2)]) }}
                                        </template>
                                    </USelectMenu>
                                </ClientOnly>
                            </UFormField>
                        </dd>
                    </div>

                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm leading-6 font-medium">
                            <label class="block text-sm leading-6 font-medium" for="targetUser">
                                {{ $t('common.target') }}
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <UFormField name="targetUserId">
                                <SelectMenu
                                    v-model="state.targetUser"
                                    :searchable="
                                        async (q: string) =>
                                            await completorStore.listColleagues({
                                                search: q,
                                                labelIds: [],
                                                userIds: [],
                                            })
                                    "
                                    searchable-key="completor-colleagues"
                                    :search-input="{ placeholder: $t('common.search_field') }"
                                    :filter-fields="['firstname', 'lastname']"
                                    block
                                    :placeholder="$t('common.colleague')"
                                    trailing
                                >
                                    <template #item-label>
                                        <template v-if="state.targetUser">
                                            {{ userToLabel(state.targetUser) }}
                                        </template>
                                    </template>

                                    <template #item="{ item }">
                                        <ColleagueName class="truncate" :colleague="item" birthday />
                                    </template>

                                    <template #empty>
                                        {{ $t('common.not_found', [$t('common.creator', 2)]) }}
                                    </template>
                                </SelectMenu>
                            </UFormField>
                        </dd>
                    </div>

                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm leading-6 font-medium">
                            <label class="block text-sm leading-6 font-medium" for="message">
                                {{ $t('common.message') }}
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <UFormField name="message">
                                <UTextarea
                                    v-model="state.message"
                                    name="message"
                                    :rows="6"
                                    :placeholder="$t('common.message')"
                                />
                            </UFormField>
                        </dd>
                    </div>

                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm leading-6 font-medium">
                            <label class="block text-sm leading-6 font-medium" for="expiresAt">
                                {{ $t('common.expires_at') }}?
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <UFormField name="expiresAt">
                                <InputDatePicker v-model="state.expiresAt" clearable time />
                            </UFormField>
                        </dd>
                    </div>
                </dl>
            </UForm>
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="$emit('close', false)" />

                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    :label="entry?.id === undefined ? $t('common.create') : $t('common.update')"
                    @click="() => formRef?.submit()"
                />
            </UButtonGroup>
        </template>
    </UModal>
</template>
