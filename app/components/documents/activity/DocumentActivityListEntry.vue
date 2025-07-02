<script lang="ts" setup>
import ActivityDocUpdatedDiff from '~/components/documents/activity/ActivityDocUpdatedDiff.vue';

import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { AccessLevel } from '~~/gen/ts/resources/documents/access';
import { type DocActivity, DocActivityType } from '~~/gen/ts/resources/documents/activity';
import { getDocAtivityIcon } from '../helpers';
import ActivityAccessUpdated from './ActivityAccessUpdated.vue';

defineProps<{
    entry: DocActivity;
}>();

function spoilerNeeded(activityType: DocActivityType): boolean {
    switch (activityType) {
        case DocActivityType.ACCESS_UPDATED:
        case DocActivityType.UPDATED:
            return true;

        default:
            return false;
    }
}
</script>

<template>
    <li
        class="hover:border-primary-500/25 dark:hover:border-primary-400/25 hover:bg-primary-100/50 dark:hover:bg-primary-900/10 border-white px-2 py-2 dark:border-gray-900"
    >
        <div v-if="!spoilerNeeded(entry.activityType)" class="flex space-x-3">
            <div class="my-auto flex size-10 items-center justify-center rounded-full">
                <UIcon class="size-7" :name="getDocAtivityIcon(entry.activityType)" />
            </div>

            <div class="flex-1 space-y-1">
                <div class="flex items-center justify-between">
                    <h3 class="inline-flex items-center gap-2 text-sm font-medium">
                        <span class="font-bold text-gray-900 dark:text-white">
                            {{ $t(`enums.documents.DocActivityType.${DocActivityType[entry.activityType]}`) }}
                        </span>
                        <span v-if="entry.data">
                            <span v-if="entry.data?.data.oneofKind === 'accessRequested'">
                                ({{ $t('common.access') }}:
                                {{
                                    $t(
                                        `enums.documents.AccessLevel.${AccessLevel[entry.data?.data.accessRequested.level ?? 0]}`,
                                    )
                                }})
                            </span>
                        </span>
                    </h3>

                    <p class="text-sm text-gray-400">
                        <GenericTime :value="entry.createdAt" type="long" />
                    </p>
                </div>

                <div class="flex items-center justify-between">
                    <p class="inline-flex gap-1 text-sm">
                        <template v-if="!!entry.reason">
                            <span class="font-semibold">{{ $t('common.reason', 1) }}:</span>
                            <span>
                                {{ entry.reason }}
                            </span>
                        </template>
                    </p>

                    <p class="inline-flex text-sm">
                        {{ $t('common.created_by') }}
                        <CitizenInfoPopover class="ml-1" :user="entry.creator" />
                    </p>
                </div>
            </div>
        </div>

        <UAccordion v-else :items="[{}]">
            <template #default="{ open }">
                <div class="flex space-x-3">
                    <div class="my-auto flex size-10 items-center justify-center rounded-full">
                        <UIcon class="size-7" :name="getDocAtivityIcon(entry.activityType)" />
                    </div>

                    <div class="flex-1 space-y-1">
                        <div class="flex items-center justify-between">
                            <h3 class="inline-flex items-center text-sm font-medium">
                                <span class="font-bold text-gray-900 dark:text-white">
                                    {{ $t(`enums.documents.DocActivityType.${DocActivityType[entry.activityType]}`) }}
                                </span>
                                <span class="ml-6 flex h-7 items-center">
                                    <UIcon
                                        :class="[open ? '!rotate-180' : '', 'size-5 transition-transform']"
                                        name="i-mdi-chevron-down"
                                    />
                                </span>
                            </h3>
                            <p class="text-sm text-gray-400">
                                <GenericTime :value="entry.createdAt" type="long" />
                            </p>
                        </div>

                        <div class="flex items-center justify-between">
                            <p class="inline-flex gap-1 text-sm">
                                <template v-if="!!entry.reason">
                                    <span class="font-semibold">{{ $t('common.reason', 1) }}:</span>
                                    <span>
                                        {{ entry.reason }}
                                    </span>
                                </template>
                            </p>

                            <p class="inline-flex text-sm">
                                {{ $t('common.created_by') }}
                                <CitizenInfoPopover class="ml-1" :user="entry.creator" />
                            </p>
                        </div>
                    </div>
                </div>
            </template>

            <template #item>
                <template v-if="entry.activityType === DocActivityType.UPDATED">
                    <div class="bg-background rounded-md p-2">
                        <ActivityDocUpdatedDiff
                            v-if="entry.data?.data.oneofKind === 'updated'"
                            :update="entry.data?.data.updated"
                        />
                    </div>
                </template>
                <template
                    v-else-if="
                        entry.activityType === DocActivityType.ACCESS_UPDATED && entry.data?.data.oneofKind === 'accessUpdated'
                    "
                >
                    <div class="bg-background rounded-md p-2">
                        <ActivityAccessUpdated :data="entry.data?.data.accessUpdated" />
                    </div>
                </template>
            </template>
        </UAccordion>
    </li>
</template>
