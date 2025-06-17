<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import type { Account } from '~~/gen/ts/resources/accounts/accounts';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { UpdateAccountResponse } from '~~/gen/ts/services/settings/accounts';
import AccountOAuth2Connection from './AccountOAuth2Connection.vue';

const props = defineProps<{
    account: Account;
}>();

const emit = defineEmits<{
    (e: 'update:account', account: Account | undefined): void;
}>();

const account = useVModel(props, 'account', emit);

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

const notifications = useNotificationsStore();

const schema = z.object({
    enabled: z.boolean(),
    lastChar: z.number().optional(),
});

type Schema = z.output<typeof schema>;

const state = reactive({
    enabled: true,
});

async function updateAccount(values: Schema): Promise<UpdateAccountResponse | undefined> {
    try {
        const call = $grpc.settings.accounts.updateAccount({
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

        isOpen.value = false;

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
</script>

<template>
    <UModal>
        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl font-semibold leading-6">
                            {{ $t('components.settings.accounts.edit_account') }}: {{ account.username }} ({{ account.id }})
                        </h3>

                        <UButton class="-my-1" color="gray" variant="ghost" icon="i-mdi-window-close" @click="isOpen = false" />
                    </div>
                </template>

                <div>
                    <div>
                        <UFormGroup class="flex-1" name="enabled" :label="$t('common.enabled')" required>
                            <UToggle v-model="state.enabled" name="enabled" />
                        </UFormGroup>
                    </div>

                    <div>
                        <UFormGroup class="flex-1" name="oauth2Accounts" :label="$t('components.auth.OAuth2Connections.title')">
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
                        </UFormGroup>
                    </div>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton class="flex-1" block color="black" @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton class="flex-1" type="submit" block :disabled="!canSubmit" :loading="!canSubmit">
                            {{ $t('common.save') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>
