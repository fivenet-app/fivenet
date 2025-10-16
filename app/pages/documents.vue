<script lang="ts" setup>
import type { NavigationMenuItem } from '@nuxt/ui';
import type { Perms } from '~~/gen/ts/perms';

const { t } = useI18n();

const itemsLeft = computed<(NavigationMenuItem & { permission?: Perms })[]>(() => [
    {
        label: t('common.document', 2),
        to: '/documents',
        icon: 'i-mdi-file-document-multiple',
        exact: true,
        permission: 'documents.DocumentsService/ListDocuments',
    },
    {
        label: t('common.approvals'),
        to: '/documents/approvals',
        icon: 'i-mdi-approval',
        permission: 'documents.DocumentsService/ListDocuments',
    },
]);

const itemsRight = computed<(NavigationMenuItem & { permission?: Perms })[]>(() => [
    {
        label: t('common.stamp', 2),
        to: '/documents/stamps',
        icon: 'i-mdi-stamper',
        permission: 'documents.StampsService/ListUsableStamps' as Perms,
    },
]);
</script>

<template>
    <NuxtPage :items-left="itemsLeft" :items-right="itemsRight" />
</template>
