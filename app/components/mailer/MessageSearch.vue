<script lang="ts" setup>
import type { CommandPaletteGroup, CommandPaletteItem } from '@nuxt/ui';
import { useMailerStore } from '~/stores/mailer';
import { getMailerMailerClient } from '~~/gen/ts/clients';

const { t, d } = useI18n();

const appConfig = useAppConfig();

const isOpen = ref(false);

const mailerStore = useMailerStore();
const { selectedEmail, selectedThread } = storeToRefs(mailerStore);

const mailerMailerClient = await getMailerMailerClient();

const searchTerm = ref('');
const searchTermDebounced = debouncedRef(searchTerm, 200);

const { data: threads, status } = useLazyAsyncData(
    () => `mailer-threads-search-${searchTermDebounced.value}`,
    () => searchThreads(searchTerm.value),
    {
        watch: [searchTermDebounced],
    },
);

async function searchThreads(q: string): Promise<CommandPaletteItem[]> {
    if (q.length < 3) return [];

    try {
        const call = mailerMailerClient.searchThreads({
            pagination: {
                offset: 0,
                pageSize: 6,
            },
            search: q.trim().substring(0, 64),
        });
        const { response } = await call;

        return response.messages.flatMap((message) => ({
            id: message.id,
            label: message.title,
            suffix: `${t('common.sender')}: ${message.sender?.email} - ${t('common.sent_at')}: ${d(toDate(message.createdAt), 'compact')}`,
            to: `/mail/${message.threadId}?email=${selectedEmail.value?.id}&msg=${message.id}`,
            active: selectedThread.value?.id === message.threadId,
            onSelect: () => {
                // Close the search modal when selecting an item
                isOpen.value = false;
            },
        }));
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const groups = computed(
    () =>
        [
            {
                id: 'messages',
                label: searchTerm.value ? `${t('common.search')}: ${searchTerm.value}...` : t('common.search'),
                items: threads.value || [],
                ignoreFilter: true,
            },
        ] as CommandPaletteGroup<CommandPaletteItem>[],
);
</script>

<template>
    <UButton
        class="w-full"
        :icon="appConfig.ui.icons.search"
        color="neutral"
        variant="outline"
        truncate
        aria-label="Search"
        v-bind="$attrs"
        @click="isOpen = !isOpen"
    />

    <UModal v-model:open="isOpen">
        <template #content>
            <ClientOnly>
                <UCommandPalette
                    v-model:search-term="searchTerm"
                    :loading="isRequestPending(status)"
                    :color-mode="false"
                    :groups="groups"
                    :empty-state="{
                        icon: 'i-mdi-email',
                        label: $t('commandpalette.empty.title'),
                        queryLabel: $t('commandpalette.empty.title'),
                    }"
                    :placeholder="`${$t('common.search_field')}`"
                    :fuse="{ resultLimit: 6, fuseOptions: { threshold: 0.1 } }"
                />
            </ClientOnly>
        </template>
    </UModal>
</template>
