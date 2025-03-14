<script lang="ts" setup>
import type { Group } from '#ui/types';
import { useMailerStore } from '~/stores/mailer';

const { $grpc } = useNuxtApp();

const { t, d } = useI18n();

const appConfig = useAppConfig();

const isOpen = ref(false);

const mailerStore = useMailerStore();
const { selectedEmail } = storeToRefs(mailerStore);

const groups = [
    {
        key: 'pages',
        label: (q: string | undefined) => q && `${t('common.search')}: ${q}`,
        search: async (q: string) => {
            try {
                const call = $grpc.mailer.mailer.searchThreads({
                    pagination: {
                        offset: 0,
                    },
                    search: q.trim().substring(0, 64),
                });
                const { response } = await call;

                return response.messages.flatMap((message) => ({
                    id: message.id,
                    label: message.title,
                    suffix: `${t('common.sender')}: ${message.sender?.email} - ${t('common.sent_at')}: ${d(toDate(message.createdAt), 'compact')}`,
                    to: `/mail?email=${selectedEmail.value?.id}&thread=${message.threadId}&msg=${message.id}`,
                }));
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        },
    },
] as Group[];
</script>

<template>
    <UButton
        v-bind="$attrs"
        :icon="appConfig.ui.icons.search"
        color="gray"
        class="w-full"
        :label="$t('common.search_field')"
        truncate
        aria-label="Search"
        @click="isOpen = !isOpen"
    />

    <ClientOnly>
        <UDashboardSearch
            v-model="isOpen"
            hide-color-mode
            :groups="groups"
            :empty-state="{
                icon: 'i-mdi-email',
                label: $t('commandpalette.empty.title'),
                queryLabel: $t('commandpalette.empty.title'),
            }"
            :placeholder="`${$t('common.search_field')}`"
            :autoselect="false"
            :fuse="{ resultLimit: 6, fuseOptions: { threshold: 0.1 } }"
        />
    </ClientOnly>
</template>
