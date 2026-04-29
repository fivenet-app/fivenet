<script lang="ts" setup>
import type { CommandPaletteGroup, CommandPaletteItem, NavigationMenuItem } from '@nuxt/ui';
import { getCitizensCitizensClient, getDocumentsDocumentsClient } from '~~/gen/ts/clients';
import FiveNetLogo from './logos/FiveNetLogo.vue';

const props = defineProps<{
    children: NavigationMenuItem[];
}>();

const { t } = useI18n();

const { isCommandSearchOpen } = useDashboard();

const citizensCitizensClient = await getCitizensCitizensClient();
const documentsDocumentsClient = await getDocumentsDocumentsClient();

const searchTerm = ref('');

const { data: citizens, status: citizensStatus } = useLazyAsyncData(
    `citizens-search-${searchTerm.value}`,
    () => searchCitiznes(searchTerm.value),
    {
        watch: [searchTerm],
        deep: false,
    },
);

async function searchCitiznes(q: string): Promise<CommandPaletteItem[]> {
    if (!q.startsWith('@') || q.length < 3) return [];

    try {
        const call = citizensCitizensClient.listCitizens({
            pagination: {
                offset: 0,
                pageSize: 10,
            },
            search: q.trim().substring(1, 64).trim(),
        });
        const { response } = await call;

        return response.users.map((u) => ({
            id: u.userId,
            label: `${u.firstname} ${u.lastname}`,
            suffix: u.dateofbirth,
            to: `/citizens/${u.userId}`,
        }));
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const { data: documents, status: documentsStatus } = useLazyAsyncData(
    `documents-search-${searchTerm.value}`,
    () => searchDocuments(searchTerm.value),
    {
        watch: [searchTerm],
        deep: false,
    },
);

async function searchDocuments(q: string): Promise<CommandPaletteItem[]> {
    if (!q.startsWith('#') || q.length < 3) return [];

    try {
        const call = documentsDocumentsClient.listDocuments({
            pagination: {
                offset: 0,
                pageSize: 10,
            },
            search: q.trim().substring(1, 64).trim(),
            categoryIds: [],
            creatorIds: [],
            documentIds: [],
        });
        const { response } = await call;

        return response.documents.map((d) => ({
            id: d.id,
            label: d.title,
            suffix: d.meta?.state,
            to: `/documents/${d.id}`,
        }));
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const idsLink = computed<CommandPaletteItem[]>(() => {
    const q = searchTerm.value;
    const links: CommandPaletteItem[] = [
        {
            id: 'id-doc',
            label: `DOC-...`,
        },
        {
            id: 'id-citizen',
            label: `CIT-...`,
        },
    ].filter((l) => l.label.toLowerCase().includes(q.toLowerCase()));

    const id = q.substring(q.indexOf('-') + 1).trim();
    if (id.length > 0 && id !== '' && !isNaN(Number(id.toString()))) {
        if (q.startsWith('CIT')) {
            links.push({
                id: 'id-citizen',
                label: `CIT-${id}`,
                to: `/citizens/${id}`,
            });
        } else if (q.startsWith('DOC')) {
            links.push({
                id: 'id-doc',
                label: `DOC-${id}`,
                to: `/documents/${id}`,
            });
        }
    }

    return links;
});

const searchLinks = computed<CommandPaletteItem[]>(
    () =>
        [
            !searchTerm.value.startsWith('@') && !citizens.value?.length
                ? {
                      id: 'cit',
                      label: t('common.citizen', 2),
                      prefix: '@',
                      icon: 'i-mdi-account-multiple-outline',
                  }
                : undefined,
            ...(citizens.value ?? []),
            !searchTerm.value.startsWith('#') && !documents.value?.length
                ? {
                      id: 'doc',
                      label: t('common.document', 2),
                      prefix: '#',
                      icon: 'i-mdi-file-document-box-multiple-outline',
                  }
                : undefined,
            ...(documents.value ?? []),
        ].filter((i) => i !== undefined) as CommandPaletteItem[],
);

function mapNavigationItemToCommandPaletteItem(link: NavigationMenuItem): CommandPaletteItem {
    const childItems = link.children?.map(mapNavigationItemToCommandPaletteItem) ?? [];
    const hasChildren = childItems.length > 0;
    const mainItem =
        hasChildren && link.to
            ? [
                  {
                      ...link,
                      chip: typeof link.chip === 'object' ? link.chip : undefined,
                      kbds: typeof link.tooltip === 'object' ? link.tooltip.kbds : undefined,
                      children: undefined,
                  } satisfies CommandPaletteItem,
              ]
            : [];

    return {
        ...link,
        chip: typeof link.chip === 'object' ? link.chip : undefined,
        to: hasChildren ? undefined : link.to,
        kbds: typeof link.tooltip === 'object' ? link.tooltip.kbds : undefined,
        children: hasChildren ? [...mainItem, ...childItems] : undefined,
    };
}

const groups = computed<CommandPaletteGroup<CommandPaletteItem>[]>(() => [
    {
        id: 'ids',
        label: t('common.id', 2),
        items: idsLink.value,
    },
    {
        id: 'search',
        label:
            (citizens.value?.length
                ? t('common.citizen', 2) + ' '
                : documents.value?.length
                  ? t('common.document', 2) + ' '
                  : '') + t('common.search'),
        items: searchLinks.value,
        ignoreFilter: true,
    },
    {
        id: 'links',
        label: t('common.goto'),
        items: props.children.map(mapNavigationItemToCommandPaletteItem),
    },
]);
</script>

<template>
    <UDashboardSearch
        v-model:open="isCommandSearchOpen"
        v-model:search-term="searchTerm"
        class="h-80 flex-1"
        :groups="groups"
        :loading="isRequestPending(citizensStatus) || isRequestPending(documentsStatus)"
        :placeholder="`${$t('common.search_field')} (${$t('commandpalette.footer', { key1: '@', key2: '#' })})`"
        :ui="{ root: 'flex-1', content: 'flex-1' }"
    >
        <template #empty>
            {{ $t('commandpalette.empty.title') }}
        </template>

        <template #footer>
            <div class="flex items-center justify-between gap-2">
                <FiveNetLogo class="ml-1 size-5 text-dimmed" />
                <div class="flex items-center gap-1">
                    <UButton class="text-dimmed" color="neutral" variant="ghost" :label="$t('common.open_command')" size="xs">
                        <template #trailing>
                            <UKbd value="enter" />
                        </template>
                    </UButton>

                    <USeparator class="h-4" orientation="vertical" />

                    <UButton class="text-dimmed" color="neutral" variant="ghost" :label="$t('common.commandpalette')" size="xs">
                        <template #trailing>
                            <UKbd value="meta" />
                            <UKbd value="k" />
                        </template>
                    </UButton>
                </div>
            </div>
        </template>
    </UDashboardSearch>
</template>
