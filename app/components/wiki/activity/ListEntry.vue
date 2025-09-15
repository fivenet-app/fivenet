<script lang="ts" setup>
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import PageUpdatedDiff from '~/components/wiki/activity/PageUpdatedDiff.vue';
import { type PageActivity, PageActivityType } from '~~/gen/ts/resources/wiki/activity';
import { getPageAtivityIcon } from '../helpers';

defineProps<{
    entry: PageActivity;
}>();

function spoilerNeeded(activityType: PageActivityType): boolean {
    switch (activityType) {
        case PageActivityType.UPDATED:
            return true;

        default:
            return false;
    }
}
</script>

<template>
    <li
        class="border-default p-2 hover:border-primary-500/25 hover:bg-primary-100/50 dark:hover:border-primary-400/25 dark:hover:bg-primary-900/10"
    >
        <div v-if="!spoilerNeeded(entry.activityType)" class="flex space-x-3">
            <div class="my-auto flex size-10 items-center justify-center rounded-full">
                <UIcon class="size-7" :name="getPageAtivityIcon(entry.activityType)" />
            </div>
            <div class="flex-1 space-y-1">
                <div class="flex items-center justify-between">
                    <h3 class="inline-flex items-center gap-2 text-sm font-medium">
                        <span class="font-bold">
                            {{ $t(`enums.wiki.PageActivityType.${PageActivityType[entry.activityType]}`) }}
                        </span>
                    </h3>

                    <p class="text-sm text-dimmed">
                        <GenericTime :value="entry.createdAt" type="long" />
                    </p>
                </div>

                <p class="inline-flex text-sm">
                    {{ $t('common.created_by') }}
                    <CitizenInfoPopover class="ml-1" :user="entry.creator" />
                </p>
            </div>
        </div>

        <UCollapsible v-else>
            <template #default>
                <div class="flex space-x-3">
                    <div class="my-auto flex size-10 items-center justify-center rounded-full">
                        <UIcon class="size-7" :name="getPageAtivityIcon(entry.activityType)" />
                    </div>
                    <div class="flex-1 space-y-1">
                        <div class="flex items-center justify-between">
                            <h3 class="inline-flex items-center text-sm font-medium">
                                <span class="font-bold">
                                    {{ $t(`enums.wiki.PageActivityType.${PageActivityType[entry.activityType]}`) }}
                                </span>

                                <span class="flex items-center">
                                    <UIcon
                                        name="i-mdi-chevron-down"
                                        class="size-5 transition-transform duration-200 group-data-[state=open]:rotate-180"
                                    />
                                </span>
                            </h3>

                            <p class="text-sm text-dimmed">
                                <GenericTime :value="entry.createdAt" type="long" />
                            </p>
                        </div>

                        <p class="inline-flex text-sm">
                            {{ $t('common.created_by') }}
                            <CitizenInfoPopover class="ml-1" :user="entry.creator" />
                        </p>
                    </div>
                </div>
            </template>

            <template v-if="entry.activityType === PageActivityType.UPDATED" #content>
                <div class="rounded-md bg-default p-2">
                    <PageUpdatedDiff v-if="entry.data?.data.oneofKind === 'updated'" :update="entry.data?.data.updated" />
                </div>
            </template>
        </UCollapsible>
    </li>
</template>
