<script lang="ts" setup>
import type { DocumentMeta } from '~~/gen/ts/resources/documents/documents';

defineProps<{
    meta: DocumentMeta | undefined;
}>();
</script>

<template>
    <UBadge
        v-if="meta?.apPoliciesActive"
        class="inline-flex gap-1"
        size="md"
        :color="meta?.approved ? 'green' : meta?.apDeclinedCount && meta?.apDeclinedCount > 0 ? 'error' : 'warning'"
        icon="i-mdi-approval"
        :label="
            meta?.approved
                ? $t('common.approved')
                : meta?.apDeclinedCount && meta?.apDeclinedCount > 0
                  ? $t('common.declined')
                  : $t('common.unapproved')
        "
    />
</template>
