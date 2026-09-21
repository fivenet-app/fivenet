<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import { getCitizensCitizensClient } from '~~/gen/ts/clients';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { UserProps } from '~~/gen/ts/resources/users/props/props';
import type { User } from '~~/gen/ts/resources/users/user';

const props = defineProps<{
    user: User;
}>();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
    (e: 'update:trafficInfractionPoints', value: number): void;
}>();

const notifications = useNotificationsStore();

const { t } = useI18n();

const citizensCitizensClient = await getCitizensCitizensClient();

const overlay = useOverlay();
const resetConfirmModal = overlay.create(ConfirmModal);

const schema = z.object({
    reason: z.coerce.string().min(3).max(255),
    trafficInfractionPoints: z.coerce.number().int().nonnegative().lt(99999),
    reset: z.coerce.boolean(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    reason: '',
    trafficInfractionPoints: props.user.props?.trafficInfractionPoints ?? 0,
    reset: false,
});

const { hasUnsavedChanges, confirmLeave } = useSnapshotChanges(state);

async function setTrafficPoints(values: Schema): Promise<void> {
    const userProps: UserProps = {
        userId: props.user.userId,
        trafficInfractionPoints: values.trafficInfractionPoints,
    };

    if (values.reset) {
        userProps.trafficInfractionPoints = 0;
    }

    try {
        const call = citizensCitizensClient.setUserProps({
            props: userProps,
            reason: values.reason,
        });
        const { response } = await call;

        emit('update:trafficInfractionPoints', response.props?.trafficInfractionPoints ?? 0);

        state.reset = false;

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

watch(
    () => props.user.props?.trafficInfractionPoints,
    () => {
        state.trafficInfractionPoints = props.user.props?.trafficInfractionPoints ?? 0;
    },
    { immediate: true },
);

const { submit, isSubmitting, canSubmit } = useSubmitGuard(async (event: FormSubmitEvent<Schema>) => {
    await setTrafficPoints(event.data);
});

const formRef = useTemplateRef('formRef');

function submitChanges(): void {
    state.reset = false;
    formRef.value?.submit();
}

function confirmReset(): void {
    resetConfirmModal.open({
        title: t('components.citizens.CitizenInfoProfile.reset_traffic_points'),
        description: t('components.citizens.CitizenInfoProfile.reset_traffic_points_description'),
        confirm: () => {
            state.reset = true;
            formRef.value?.submit();
        },
    });
}

async function closeModal(): Promise<void> {
    if (!canSubmit.value) return;

    if (hasUnsavedChanges.value && !(await confirmLeave())) return;

    emit('close', false);
}
</script>

<template>
    <UDrawer
        :title="$t('components.citizens.CitizenInfoProfile.set_traffic_points')"
        :close="false"
        :dismissible="!hasUnsavedChanges && canSubmit"
        :ui="{ body: 'mx-auto w-full max-w-xl' }"
    >
        <template #header>
            <div class="flex w-full items-center justify-between gap-2">
                <h3 class="font-semibold text-highlighted">
                    {{ $t('components.citizens.CitizenInfoProfile.set_traffic_points') }}
                </h3>

                <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-close"
                    :disabled="!canSubmit"
                    :aria-label="$t('common.close', 1)"
                    @click="closeModal"
                />
            </div>
        </template>

        <template #body>
            <UForm ref="formRef" class="grid gap-4" :schema="schema" :state="state" @submit="submit">
                <UFormField name="reason" :label="$t('common.reason')" required>
                    <UInput v-model="state.reason" class="w-full" type="text" :placeholder="$t('common.reason')" />
                </UFormField>

                <UFormField
                    name="trafficInfractionPoints"
                    :label="$t('common.traffic_infraction_points')"
                    :description="
                        $t('components.citizens.CitizenInfoProfile.current_traffic_points', {
                            points: props.user.props?.trafficInfractionPoints ?? 0,
                        })
                    "
                >
                    <UInputNumber
                        v-model="state.trafficInfractionPoints"
                        class="w-full"
                        :min="0"
                        :max="99998"
                        :step="1"
                        :placeholder="$t('common.traffic_infraction_points')"
                    />
                </UFormField>
            </UForm>
        </template>

        <template #footer>
            <div class="flex w-full flex-col-reverse gap-2 sm:flex-row">
                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="isSubmitting"
                    :label="$t('common.save')"
                    @click="submitChanges"
                />

                <UButton
                    class="flex-1"
                    color="error"
                    block
                    :disabled="!canSubmit"
                    :loading="isSubmitting"
                    :label="$t('components.citizens.CitizenInfoProfile.reset_traffic_points')"
                    @click="confirmReset"
                />

                <UButton
                    class="flex-1"
                    color="neutral"
                    block
                    :disabled="!canSubmit"
                    :label="$t('common.close', 1)"
                    @click="closeModal"
                />
            </div>
        </template>
    </UDrawer>
</template>
