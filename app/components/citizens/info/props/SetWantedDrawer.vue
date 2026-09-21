<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import { getCitizensCitizensClient } from '~~/gen/ts/clients';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { UserProps } from '~~/gen/ts/resources/users/props/props';
import type { User } from '~~/gen/ts/resources/users/user';

const props = defineProps<{
    user: User;
}>();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
    (e: 'update:wantedStatus', value: boolean): void;
}>();

const notifications = useNotificationsStore();

const citizensCitizensClient = await getCitizensCitizensClient();

const schema = z.object({
    reason: z.coerce.string().min(3).max(255),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    reason: '',
});

const { hasUnsavedChanges, confirmLeave } = useSnapshotChanges(state);

async function setWantedState(values: Schema): Promise<void> {
    const userProps: UserProps = {
        userId: props.user.userId,
        wanted: props.user.props ? !props.user.props.wanted : true,
    };

    try {
        const call = citizensCitizensClient.setUserProps({
            props: userProps,
            reason: values.reason,
        });
        const { response } = await call;

        emit('update:wantedStatus', response.props?.wanted ?? false);

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
    <UDrawer
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
    </UDrawer>
</template>
