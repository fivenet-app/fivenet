<script lang="ts" setup>
import { NuxtImg } from '#components';
import type { OAuth2Account } from '~~/gen/ts/resources/accounts/oauth2';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = defineProps<{
    accountId: number;
    connection: OAuth2Account;
}>();

const emit = defineEmits<{
    (e: 'deleted', providerName: string): void;
}>();

const { $grpc } = useNuxtApp();

const notifications = useNotificationsStore();

async function disconnectOAuth2Connection(accountId: number, providerName: string): Promise<void> {
    try {
        await $grpc.settings.accounts.disconnectOAuth2Connection({
            id: accountId,
            providerName: providerName,
        });

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        emit('deleted', providerName);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const { login } = useAppConfig();
const provider = computed(() => login.providers.find((p) => p.name === props.connection.providerName));
</script>

<template>
    <UPageCard
        :ui="{
            body: {
                padding: 'px-4 py-4 sm:p-4',
            },
            icon: { wrapper: 'mb-1' },
        }"
    >
        <template #title>
            <div class="flex flex-1 gap-2">
                <UButton class="inline-flex flex-1 gap-2" variant="ghost" external :to="provider?.homepage" target="_blank">
                    <NuxtImg
                        v-if="!provider?.icon?.startsWith('i-')"
                        class="size-10"
                        :src="provider?.icon"
                        :alt="provider?.name"
                        placeholder-class="size-10"
                        loading="lazy"
                    />
                    <UIcon
                        v-else
                        class="size-10"
                        :name="provider.icon"
                        :style="provider.name === 'discord' && { color: '#7289da' }"
                    />

                    <div class="flex items-center gap-1.5 text-base font-semibold text-gray-900 dark:text-white">
                        {{ provider?.label }}
                    </div>
                </UButton>

                <div class="flex items-center justify-between">
                    <UButton
                        icon="i-mdi-close-circle"
                        color="error"
                        @click="disconnectOAuth2Connection(accountId, connection.providerName)"
                    >
                        {{ $t('common.disconnect') }}
                    </UButton>
                </div>
            </div>
        </template>

        <template v-if="connection" #footer>
            <div class="inline-flex items-center gap-4">
                <template v-if="connection">
                    <UAvatar :as="NuxtImg" size="md" :src="connection.avatar" :alt="$t('common.image')" loading="lazy" />

                    <UTooltip :text="`ID: ${connection.externalId}`">
                        <span class="text-left">
                            {{ connection.username }}
                        </span>
                    </UTooltip>
                </template>
            </div>
        </template>
    </UPageCard>
</template>
