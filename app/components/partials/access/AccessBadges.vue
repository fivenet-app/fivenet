<script lang="ts" setup generic="JobsT extends JobAccessEntry, UsersT extends UserAccessEntry">
import JobInfoPopover from '../JobInfoPopover.vue';
import type { JobAccessEntry, UserAccessEntry } from './helpers';

withDefaults(
    defineProps<{
        accessLevel: Zod.EnumLike;
        jobs?: JobsT[];
        users?: UsersT[];
        i18nKey: string;
        i18nAccessLevelKey?: string;
    }>(),
    {
        jobs: () => [],
        users: () => [],
        i18nAccessLevelKey: 'AccessLevel',
    },
);

defineOptions({
    inheritAttrs: false,
});
</script>

<template>
    <div class="flex flex-col gap-2">
        <div class="flex flex-row flex-wrap gap-1">
            <UBadge v-for="entry in jobs" :key="entry.id" class="inline-flex gap-1" color="black" size="md" v-bind="$attrs">
                <span class="size-2 rounded-full border bg-info-500" />
                <span>
                    <template v-if="entry.jobLabel">
                        {{ entry.jobLabel }}
                        <span
                            v-if="entry.minimumGrade !== undefined && entry.minimumGrade >= 0"
                            :title="`${entry.jobLabel} - ${$t('common.rank')} ${entry.minimumGrade}`"
                        >
                            ({{ entry.jobGradeLabel ?? entry.minimumGrade }})</span
                        >
                    </template>
                    <JobInfoPopover v-else :job="entry.job" :job-label="entry.jobLabel" :grade="entry.minimumGrade" />
                    -
                    {{ $t(`${i18nKey}.${i18nAccessLevelKey}.${accessLevel[entry.access]}`) }}
                </span>
            </UBadge>
        </div>

        <div class="flex flex-row flex-wrap gap-1">
            <UBadge v-for="entry in users" :key="entry.id" class="inline-flex gap-1" color="black" size="md" v-bind="$attrs">
                <span class="size-2 rounded-full bg-amber-500" />

                <span :title="`${$t('common.id')} ${entry.userId}`">
                    <template v-if="entry.user">
                        {{ entry.user?.firstname }}
                        {{ entry.user?.lastname }}
                    </template>
                    <template v-else>
                        {{ entry.userId }}
                    </template>
                    - {{ $t(`${i18nKey}.${i18nAccessLevelKey}.${accessLevel[entry.access]}`) }}
                </span>
            </UBadge>
        </div>
    </div>
</template>
