<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import { getCitizensCitizensClient } from '~~/gen/ts/clients';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { PropertyMutation } from '~~/gen/ts/resources/properties/properties';
import type { UserProps } from '~~/gen/ts/resources/users/props/props';
import type { User } from '~~/gen/ts/resources/users/user';
import WantedTillPicker from '~/components/partials/WantedTillPicker.vue';
import { toTimestamp } from '~/utils/time';
import ResponsiveDrawerModal from '~/components/partials/ResponsiveDrawerModal.vue';

const props = defineProps<{
    user: User;
}>();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
    (e: 'update:wantedStatus', value: UserProps): void;
}>();

const notifications = useNotificationsStore();

const citizensCitizensClient = await getCitizensCitizensClient();

const schema = z.object({
    reason: z.coerce.string().min(3).max(255),
    wantedTill: z.date().optional(),
    permanent: z.boolean(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    reason: '',
    wantedTill: undefined,
    permanent: true,
});

const { hasUnsavedChanges, confirmLeave } = useSnapshotChanges(state);

async function setWantedState(values: Schema): Promise<void> {
    const wanted = !props.user.props?.wanted;
    let propertyMutations: PropertyMutation[];
    if (wanted) {
        const mutation: PropertyMutation = {
            key: 'user.wanted',
            valueChange: {
                oneofKind: 'setValue',
                setValue: { value: { oneofKind: 'boolValue', boolValue: true } },
            },
            reason: values.reason,
        };
        if (!values.permanent && values.wantedTill) {
            const expiresAt = toTimestamp(values.wantedTill);
            if (expiresAt) mutation.state = { expiryChange: { oneofKind: 'expiresAt', expiresAt } };
        }
        propertyMutations = [mutation];
    } else {
        propertyMutations = [
            { key: 'user.wanted', reason: values.reason, valueChange: { oneofKind: 'clearValue', clearValue: {} } },
        ];
    }

    try {
        const call = citizensCitizensClient.setUserProps({
            props: { userId: props.user.userId },
            reason: values.reason,
            propertyMutations,
        });
        const { response } = await call;

        if (response.props) {
            emit('update:wantedStatus', response.props);
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

const { submit, isSubmitting, canSubmit } = useSubmitGuard(async (event: FormSubmitEvent<Schema>) => {
    await setWantedState(event.data);
});

const formRef = useTemplateRef('formRef');

async function closeModal(): Promise<void> {
    if (!canSubmit.value) return;

    if (hasUnsavedChanges.value && !(await confirmLeave())) return;

    emit('close', false);
}
</script>

<template>
    <ResponsiveDrawerModal
        :title="
            user.props?.wanted
                ? $t('components.citizens.CitizenInfoProfile.revoke_wanted')
                : $t('components.citizens.CitizenInfoProfile.set_wanted')
        "
        :close="false"
        :dismissible="!hasUnsavedChanges && canSubmit"
        :ui="{ body: 'mx-auto w-full max-w-xl' }"
    >
        <template #header>
            <div class="flex w-full items-center justify-between gap-2">
                <h3 class="font-semibold text-highlighted">
                    {{
                        user.props?.wanted
                            ? $t('components.citizens.CitizenInfoProfile.revoke_wanted')
                            : $t('components.citizens.CitizenInfoProfile.set_wanted')
                    }}
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
            <UForm ref="formRef" :schema="schema" :state="state" @submit="submit">
                <UAlert
                    class="mb-4"
                    :color="user.props?.wanted ? 'success' : 'error'"
                    :icon="user.props?.wanted ? 'i-mdi-account-check' : 'i-mdi-account-alert'"
                    :title="
                        user.props?.wanted
                            ? $t('components.citizens.CitizenInfoProfile.revoke_wanted')
                            : $t('components.citizens.CitizenInfoProfile.set_wanted')
                    "
                    :description="
                        user.props?.wanted
                            ? $t('components.citizens.CitizenInfoProfile.revoke_wanted_description')
                            : $t('components.citizens.CitizenInfoProfile.set_wanted_description')
                    "
                />

                <UFormField class="flex-1" name="reason" :label="$t('common.reason')" required>
                    <UInput v-model="state.reason" class="w-full" type="text" :placeholder="$t('common.reason')" />
                </UFormField>

                <template v-if="!user.props?.wanted">
                    <UFormField name="wantedTill" :label="$t('common.duration')">
                        <WantedTillPicker v-model="state.wantedTill" v-model:permanent="state.permanent" />
                    </UFormField>
                </template>
            </UForm>
        </template>

        <template #footer>
            <div class="flex w-full flex-col-reverse gap-2 sm:flex-row">
                <UButton
                    class="flex-1"
                    color="neutral"
                    block
                    :disabled="!canSubmit"
                    :label="$t('common.close', 1)"
                    @click="closeModal"
                />

                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="isSubmitting"
                    :label="$t('common.save')"
                    @click="formRef?.submit()"
                />
            </div>
        </template>
    </ResponsiveDrawerModal>
</template>
