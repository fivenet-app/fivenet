<script lang="ts" setup>
import DocUpdatedDiff from '~/components/documents/activity/DocUpdatedDiff.vue';

import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import ResponsiveDrawerModal from '~/components/partials/ResponsiveDrawerModal.vue';
import { AccessLevel } from '~~/gen/ts/resources/documents/access/access';
import { type DocActivity, DocActivityType } from '~~/gen/ts/resources/documents/activity/activity';
import AccessUpdated from './AccessUpdated.vue';

const props = defineProps<{
    entry: DocActivity;
}>();

const detailsOpen = ref(false);
const hasDetails = computed(
    () =>
        (props.entry.activityType === DocActivityType.UPDATED && props.entry.data?.data.oneofKind === 'updated') ||
        (props.entry.activityType === DocActivityType.ACCESS_UPDATED && props.entry.data?.data.oneofKind === 'accessUpdated'),
);
</script>

<template>
    <div class="flex min-w-0 flex-auto justify-between gap-2 text-xs leading-5 text-muted">
        <div class="flex min-w-0 flex-col gap-0.5">
            <span class="text-sm font-medium text-highlighted">
                {{ $t(`enums.documents.DocActivityType.${DocActivityType[entry.activityType]}`) }}
                <template v-if="entry.data?.data.oneofKind === 'accessRequested'">
                    ({{ $t('common.access') }}:
                    {{ $t(`enums.documents.AccessLevel.${AccessLevel[entry.data.data.accessRequested.level ?? 0]}`) }})
                </template>
            </span>

            <span v-if="entry.reason" class="inline-flex gap-1 text-sm">
                <span class="font-semibold">{{ $t('common.reason', 1) }}:</span>
                <span>{{ entry.reason }}</span>
            </span>

            <UButton
                v-if="hasDetails"
                class="w-fit"
                color="neutral"
                variant="link"
                size="xs"
                icon="i-mdi-file-compare"
                :label="$t('common.view_changes')"
                @click="detailsOpen = true"
            />
        </div>

        <div class="flex flex-none flex-col items-end gap-0.5 text-xs">
            <span class="text-xs/5 text-dimmed"><GenericTime :value="entry.createdAt" type="compact" /></span>
            <span class="inline-flex items-center gap-1">
                <span class="text-dimmed">{{ $t('common.created_by') }}:</span>
                <CitizenInfoPopover :user="entry.creator" />
            </span>
        </div>
    </div>

    <ResponsiveDrawerModal
        v-model:open="detailsOpen"
        :title="$t('components.documents.activity_list.difference')"
        :ui="{ content: 'max-w-4xl' }"
    >
        <template #body>
            <DocUpdatedDiff v-if="entry.data?.data.oneofKind === 'updated'" :update="entry.data.data.updated" />
            <AccessUpdated v-else-if="entry.data?.data.oneofKind === 'accessUpdated'" :data="entry.data.data.accessUpdated" />
        </template>
    </ResponsiveDrawerModal>
</template>
