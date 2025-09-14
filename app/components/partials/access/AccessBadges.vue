<script
    lang="ts"
    setup
    generic="JobsT extends JobAccessEntry, UsersT extends UserAccessEntry, QualiT extends QualificationAccessEntry"
>
import type { z } from 'zod';
import JobInfoPopover from '../JobInfoPopover.vue';
import type { JobAccessEntry, QualificationAccessEntry, UserAccessEntry } from './helpers';

withDefaults(
    defineProps<{
        accessLevel: z.util.EnumLike;
        jobs?: JobsT[];
        users?: UsersT[];
        qualifications?: QualiT[];
        i18nKey: string;
        i18nAccessLevelKey?: string;
    }>(),
    {
        jobs: () => [],
        users: () => [],
        qualifications: () => [],
        i18nAccessLevelKey: 'AccessLevel',
    },
);

defineOptions({
    inheritAttrs: false,
});
</script>

<template>
    <div class="flex flex-col gap-2">
        <div v-if="jobs.length > 0" class="flex flex-row flex-wrap gap-1">
            <UBadge v-for="entry in jobs" :key="entry.id" class="inline-flex gap-1" color="neutral" size="md" v-bind="$attrs">
                <span class="size-2 rounded-full bg-info-500" />
                <span class="inline-flex gap-1">
                    <span v-if="entry.jobLabel">
                        {{ entry.jobLabel }}
                        <span
                            v-if="entry.minimumGrade !== undefined && entry.minimumGrade >= 0"
                            :title="`${entry.jobLabel} - ${$t('common.rank')} ${entry.minimumGrade}`"
                        >
                            ({{ entry.jobGradeLabel ?? entry.minimumGrade }})</span
                        >
                    </span>
                    <JobInfoPopover v-else :job="entry.job" :job-label="entry.jobLabel" :grade="entry.minimumGrade" />
                    -
                    <span>{{ $t(`${i18nKey}.${i18nAccessLevelKey}.${accessLevel[entry.access]}`) }}</span>
                </span>
            </UBadge>
        </div>

        <div v-if="users.length > 0" class="flex flex-row flex-wrap gap-1">
            <UBadge v-for="entry in users" :key="entry.id" class="inline-flex gap-1" color="neutral" size="md" v-bind="$attrs">
                <span class="size-2 rounded-full bg-amber-500" />

                <span class="inline-flex gap-1" :title="`${$t('common.id')} ${entry.userId}`">
                    <span v-if="entry.user">
                        {{ entry.user?.firstname }}
                        {{ entry.user?.lastname }}
                    </span>
                    <span v-else>
                        {{ entry.userId }}
                    </span>
                    -
                    <span>{{ $t(`${i18nKey}.${i18nAccessLevelKey}.${accessLevel[entry.access]}`) }}</span>
                </span>
            </UBadge>
        </div>

        <div v-if="qualifications.length > 0" class="flex flex-row flex-wrap gap-1">
            <UBadge
                v-for="entry in qualifications"
                :key="entry.id"
                class="inline-flex gap-1"
                color="neutral"
                size="md"
                v-bind="$attrs"
            >
                <span class="size-2 rounded-full bg-amber-500" />

                <span class="inline-flex gap-1">
                    <span v-if="entry.qualification">
                        {{ entry.qualification.abbreviation }}: {{ entry.qualification.title }}
                    </span>
                    <span v-else>
                        {{ entry.qualificationId }}
                    </span>
                    -
                    <span>{{ $t(`${i18nKey}.${i18nAccessLevelKey}.${accessLevel[entry.access]}`) }}</span>
                </span>
            </UBadge>
        </div>
    </div>
</template>
