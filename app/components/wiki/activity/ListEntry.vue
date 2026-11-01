<script lang="ts" setup>
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import ResponsiveDrawerModal from '~/components/partials/ResponsiveDrawerModal.vue';
import PageUpdatedDiff from '~/components/wiki/activity/PageUpdatedDiff.vue';
import { type PageActivity, PageActivityType } from '~~/gen/ts/resources/wiki/activity/activity';

const props = defineProps<{
    entry: PageActivity;
}>();

const detailsOpen = ref(false);
const hasUpdateDiff = computed(
    () => props.entry.activityType === PageActivityType.UPDATED && props.entry.data?.data.oneofKind === 'updated',
);
</script>

<template>
    <div class="flex min-w-0 flex-auto justify-between gap-2 text-xs leading-5 text-muted">
        <div class="flex min-w-0 flex-col gap-0.5">
            <span class="text-sm font-medium text-highlighted">
                {{ $t(`enums.wiki.PageActivityType.${PageActivityType[entry.activityType]}`) }}
            </span>

            <UButton
                v-if="hasUpdateDiff"
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
            <span class="text-xs/5 text-dimmed"><GenericTime :value="entry.createdAt" type="short" /></span>

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
            <PageUpdatedDiff v-if="entry.data?.data.oneofKind === 'updated'" :update="entry.data.data.updated" />
        </template>
    </ResponsiveDrawerModal>
</template>
