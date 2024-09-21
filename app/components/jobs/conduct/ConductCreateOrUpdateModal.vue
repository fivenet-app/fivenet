<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import DatePickerPopoverClient from '~/components/partials/DatePickerPopover.client.vue';
import { useAuthStore } from '~/store/auth';
import { useCompletorStore } from '~/store/completor';
import { useNotificatorStore } from '~/store/notificator';
import type { ConductEntry } from '~~/gen/ts/resources/jobs/conduct';
import { ConductType } from '~~/gen/ts/resources/jobs/conduct';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { UserShort } from '~~/gen/ts/resources/users/users';
import { conductTypesToBGColor } from './helpers';

const props = defineProps<{
    entry?: ConductEntry;
    userId?: number;
}>();

const emit = defineEmits<{
    (e: 'created', entry: ConductEntry): void;
    (e: 'updated', entry: ConductEntry): void;
}>();

const { isOpen } = useModal();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const completorStore = useCompletorStore();

const notifications = useNotificatorStore();

const usersLoading = ref(false);

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

async function conductCreateOrUpdateEntry(values: Schema, id?: string): Promise<void> {
    try {
        const req = {
            entry: {
                id: id ?? '0',
                job: '',
                creatorId: activeChar.value?.userId ?? 1,
                type: values.type,
                message: values.message,
                targetUserId: values.targetUser.userId,
                expiresAt: values.expiresAt ? toTimestamp(values.expiresAt) : undefined,
            },
        };

        if (id === undefined) {
            const call = getGRPCJobsConductClient().createConductEntry(req);
            const { response } = await call;

            emit('created', response.entry!);
        } else {
            const call = getGRPCJobsConductClient().updateConductEntry(req);
            const { response } = await call;

            emit('updated', response.entry!);
        }

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        isOpen.value = false;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function setFormFromProps(): Promise<void> {
    state.targetUser = props.entry?.targetUser;
    state.type = props.entry?.type ?? ConductType.NOTE;
    state.message = props.entry?.message ?? '';
    state.expiresAt = props.entry?.expiresAt ? toDate(props.entry?.expiresAt) : undefined;
}

watch(props, () => setFormFromProps());

onMounted(() => setFormFromProps());

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await conductCreateOrUpdateEntry(event.data, props.entry?.id).finally(() =>
        useTimeoutFn(() => (canSubmit.value = true), 400),
    );
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
                                entry === undefined
                                    ? $t('components.jobs.conduct.CreateOrUpdateModal.create.title')
                                    : $t('components.jobs.conduct.CreateOrUpdateModal.update.title')
                            }}
                        </h3>

                        <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                    </div>
                </template>

                <div>
                    <dl class="divide-neutral/10 divide-y">
                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm font-medium leading-6">
                                <label for="type" class="block text-sm font-medium leading-6">
                                    {{ $t('common.type') }}
                                </label>
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <UFormGroup name="type">
                                    <USelectMenu
                                        v-model="state.type"
                                        :options="cTypes"
                                        value-attribute="status"
                                        :searchable-placeholder="$t('common.search_field')"
                                    >
                                        <template #label>
                                            <span class="truncate">{{
                                                $t(`enums.jobs.ConductType.${ConductType[state.type ?? 0]}`)
                                            }}</span>
                                        </template>
                                        <template #option="{ option }">
                                            <span class="truncate" :class="conductTypesToBGColor(option.status)">{{
                                                $t(`enums.jobs.ConductType.${ConductType[option.status ?? 0]}`)
                                            }}</span>
                                        </template>
                                    </USelectMenu>
                                </UFormGroup>
                            </dd>
                        </div>
                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm font-medium leading-6">
                                <label for="targetUser" class="block text-sm font-medium leading-6">
                                    {{ $t('common.target') }}
                                </label>
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <UFormGroup name="targetUserId">
                                    <USelectMenu
                                        v-model="state.targetUser"
                                        :searchable="
                                            async (query: string) => {
                                                usersLoading = true;
                                                const colleagues = await completorStore.listColleagues({
                                                    search: query,
                                                });
                                                usersLoading = false;
                                                return colleagues;
                                            }
                                        "
                                        searchable-lazy
                                        :searchable-placeholder="$t('common.search_field')"
                                        :search-attributes="['firstname', 'lastname']"
                                        block
                                        :placeholder="$t('common.colleague')"
                                        trailing
                                        by="userId"
                                    >
                                        <template #label>
                                            <template v-if="state.targetUser">
                                                {{ userToLabel(state.targetUser) }}
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
                            </dd>
                        </div>
                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm font-medium leading-6">
                                <label for="message" class="block text-sm font-medium leading-6">
                                    {{ $t('common.message') }}
                                </label>
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <UFormGroup name="message">
                                    <UTextarea
                                        v-model="state.message"
                                        name="message"
                                        :rows="6"
                                        :placeholder="$t('common.message')"
                                    />
                                </UFormGroup>
                            </dd>
                        </div>
                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm font-medium leading-6">
                                <label for="expiresAt" class="block text-sm font-medium leading-6">
                                    {{ $t('common.expires_at') }}?
                                </label>
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <UFormGroup name="expiresAt">
                                    <DatePickerPopoverClient
                                        v-model="state.expiresAt"
                                        :popover="{ popper: { placement: 'bottom-start' } }"
                                        :date-picker="{ clearable: true }"
                                    />
                                </UFormGroup>
                            </dd>
                        </div>
                    </dl>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton color="black" block class="flex-1" @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton type="submit" block class="flex-1" :disabled="!canSubmit" :loading="!canSubmit">
                            {{ entry?.id === undefined ? $t('common.create') : $t('common.update') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>
