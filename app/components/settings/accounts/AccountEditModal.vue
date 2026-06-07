<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import { getSettingsAccountsClient } from '~~/gen/ts/clients';
import type { Account } from '~~/gen/ts/resources/accounts/accounts';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { UpdateAccountResponse } from '~~/gen/ts/services/settings/accounts';
import AccountSocialLogin from './AccountSocialLogin.vue';

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const account = defineModel<Account>('account', { required: true });

const notifications = useNotificationsStore();

const settingsAccountsClient = await getSettingsAccountsClient();

const schema = z.object({
    enabled: z.coerce.boolean().default(true),
    lastChar: z.coerce.number().optional(),
    groups: z.string().array().max(5).optional(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    enabled: true,
});

const { hasUnsavedChanges, confirmLeave, syncSnapshot } = useSnapshotChanges(state);

async function updateAccount(values: Schema): Promise<UpdateAccountResponse | undefined> {
    try {
        const call = settingsAccountsClient.updateAccount({
            id: account.value.id,

            enabled: values.enabled,
        });
        const { response } = await call;

        if (response.account) {
            account.value = response.account;

            notifications.add({
                title: { key: 'notifications.action_successful.title', parameters: {} },
                description: { key: 'notifications.action_successful.content', parameters: {} },
                type: NotificationType.SUCCESS,
            });
        }

        emit('close', false);

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

function setFromProps(): void {
    if (!account.value) return;

    state.enabled = account.value.enabled;
    state.groups = account.value.groups?.groups ?? [];
    syncSnapshot();
}

setFromProps();
watch(account, () => setFromProps());

const canSubmit = ref<boolean>(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await updateAccount(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const formRef = useTemplateRef('formRef');

async function closeModal(): Promise<void> {
    if (!canSubmit.value) return;

    if (hasUnsavedChanges.value && !(await confirmLeave())) return;

    emit('close', false);
}
</script>

<template>
    <UModal
        :title="`${$t('components.settings.accounts.edit_account')}: ${account.username} (${$t('common.id')}: ${account.id})`"
        :close="false"
        :dismissible="!hasUnsavedChanges && canSubmit"
    >
        <template #header>
            <div class="flex w-full items-center justify-between gap-2">
                <h3 class="font-semibold text-highlighted">
                    {{ $t('components.settings.accounts.edit_account') }}: {{ account.username }} ({{ $t('common.id') }}:
                    {{ account.id }})
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
            <UForm ref="formRef" class="space-y-4" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <UFormField class="flex-1" name="enabled" :label="$t('common.enabled')" required>
                    <USwitch v-model="state.enabled" name="enabled" />
                </UFormField>

                <UFormField class="flex-1" name="groups" :label="$t('common.group', 2)">
                    <UInputTags v-model="state.groups" class="w-full" disabled />
                </UFormField>

                <UFormField class="flex-1" name="oauth2Accounts" :label="$t('components.auth.SocialLogins.title')">
                    <div class="flex flex-col gap-2">
                        <DataNoDataBlock
                            v-if="account.oauth2Accounts.length === 0"
                            :type="$t('components.auth.SocialLogins.title')"
                        />

                        <template v-else>
                            <AccountSocialLogin
                                v-for="connection in account.oauth2Accounts"
                                :key="connection.providerName"
                                :account-id="account.id"
                                :connection="connection"
                                @deleted="
                                    account.oauth2Accounts = account.oauth2Accounts.filter(
                                        (c) => c.providerName !== connection.providerName,
                                    )
                                "
                            />
                        </template>
                    </div>
                </UFormField>
            </UForm>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton
                    class="flex-1"
                    block
                    color="neutral"
                    :disabled="!canSubmit"
                    :label="$t('common.close', 1)"
                    @click="closeModal"
                />

                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    :label="$t('common.save')"
                    @click="() => formRef?.submit()"
                />
            </UFieldGroup>
        </template>
    </UModal>
</template>
