<script lang="ts" setup>
import type { NavigationMenuItem } from '@nuxt/ui';
import type { Perms } from '~~/gen/ts/perms';

const { t } = useI18n();

const { can } = useAuth();

const itemsLeft = computed<NavigationMenuItem[]>(() =>
    [
        {
            label: t('common.document', 2),
            to: '/documents',
            icon: 'i-mdi-file-document-multiple',
            exact: true,
            permission: 'documents.DocumentsService/ListDocuments' as Perms,
        },
        {
            label: t('common.approvals', 2),
            to: '/documents/approvals',
            icon: 'i-mdi-approval',
            permission: 'documents.DocumentsService/ListDocuments' as Perms,
        },
    ].filter((t) => t.permission === undefined || can(t.permission).value),
);

const itemsRight = computed<NavigationMenuItem[]>(() =>
    [
        {
            label: t('common.stamp', 2),
            to: '/documents/stamps',
            icon: 'i-mdi-stamper',
            permission: 'documents.StampsService/ListUsableStamps' as Perms,
        },
    ].filter((t) => t.permission === undefined || can(t.permission).value),
);
</script>

<template>
    <NuxtPage :items-left="itemsLeft" :items-right="itemsRight" />
</template>
