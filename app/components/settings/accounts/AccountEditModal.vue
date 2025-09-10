<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import { getSettingsAccountsClient } from '~~/gen/ts/clients';
import type { Account } from '~~/gen/ts/resources/accounts/accounts';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { UpdateAccountResponse } from '~~/gen/ts/services/settings/accounts';
import AccountOAuth2Connection from './AccountOAuth2Connection.vue';

const props = defineProps<{
    account: Account;
}>();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
    (e: 'update:account', account: Account | undefined): void;
}>();

const account = useVModel(props, 'account', emit);

const notifications = useNotificationsStore();

const settingsAccountsClient = await getSettingsAccountsClient();

const schema = z.object({
    enabled: z.coerce.boolean().default(true),
    lastChar: z.coerce.number().optional(),
});

type Schema = z.output<typeof schema>;

const state = reactive({
    enabled: true,
});

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
    state.enabled = account.value.enabled;
}

setFromProps();
watch(props, () => setFromProps());

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await updateAccount(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UModal :title="`${$t('components.settings.accounts.edit_account')}: ${account.username} (${account.id})`">
        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <div>
                    <UFormField class="flex-1" name="enabled" :label="$t('common.enabled')" required>
                        <USwitch v-model="state.enabled" name="enabled" />
                    </UFormField>
                </div>

                <div>
                    <UFormField class="flex-1" name="oauth2Accounts" :label="$t('components.auth.OAuth2Connections.title')">
                        <div class="flex flex-col gap-2">
                            <DataNoDataBlock
                                v-if="account.oauth2Accounts.length === 0"
                                :type="$t('components.auth.OAuth2Connections.title')"
                            />

                            <template v-else>
                                <AccountOAuth2Connection
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
                </div>
            </UForm>
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton class="flex-1" block color="neutral" :label="$t('common.close', 1)" @click="$emit('close', false)" />

                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    :label="$t('common.save')"
                    @click="() => formRef?.submit()"
                />
            </UButtonGroup>
        </template>
    </UModal>
</template>
