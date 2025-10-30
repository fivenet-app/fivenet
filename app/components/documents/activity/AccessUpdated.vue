<script lang="ts" setup>
import AccessBadges from '~/components/partials/access/AccessBadges.vue';
import { AccessLevel } from '~~/gen/ts/resources/documents/access';
import type { DocAccessUpdated } from '~~/gen/ts/resources/documents/activity';

defineProps<{
    data: DocAccessUpdated;
}>();
</script>

<template>
    <div class="flex flex-col gap-1">
        <AccessBadges
            v-if="(data.users?.toCreate ?? []).length > 0 || (data.jobs?.toCreate ?? []).length > 0"
            :access-level="AccessLevel"
            :jobs="data.jobs?.toCreate"
            :users="data.users?.toCreate"
            i18n-key="enums.documents"
            color="success"
        />

        <AccessBadges
            v-if="(data.users?.toDelete?.length ?? 0) > 0 || (data.jobs?.toDelete?.length ?? 0) > 0"
            :access-level="AccessLevel"
            :jobs="data.jobs?.toDelete"
            :users="data.users?.toDelete"
            i18n-key="enums.documents"
            color="error"
        />

        <AccessBadges
            v-if="(data.users?.toUpdate?.length ?? 0) > 0 || (data.jobs?.toUpdate?.length ?? 0) > 0"
            :access-level="AccessLevel"
            :jobs="data.jobs?.toUpdate"
            :users="data.users?.toUpdate"
            i18n-key="enums.documents"
            color="info"
        />

        <span class="inline-flex gap-2">
            <span class="text-base font-semibold text-white">{{ $t('common.legend') }}:</span>
            <span class="bg-success-600 text-white">{{ $t('components.documents.activity_list.legend.added') }}</span>
            <span class="bg-error-600 text-white">{{ $t('components.documents.activity_list.legend.removed') }}</span>
            <span class="bg-info-600 text-white">{{ $t('components.documents.activity_list.legend.changed') }}</span>
        </span>
    </div>
</template>
