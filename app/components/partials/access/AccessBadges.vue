<script lang="ts" setup generic="JobsT extends JobAccessEntry, UsersT extends UserAccessEntry">
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
</script>

<template>
    <div class="flex flex-col gap-2">
        <div class="flex flex-row flex-wrap gap-1">
            <UBadge v-for="entry in jobs" :key="entry.id" color="black" class="inline-flex gap-1" size="md">
                <span class="size-2 rounded-full bg-info-500" />
                <span>
                    {{ entry.jobLabel
                    }}<span
                        v-if="entry.minimumGrade !== undefined && entry.minimumGrade > 0"
                        :title="`${entry.jobLabel} - ${$t('common.rank')} ${entry.minimumGrade}`"
                    >
                        ({{ entry.jobGradeLabel }})</span
                    >
                    -
                    {{ $t(`${i18nKey}.${i18nAccessLevelKey}.${accessLevel[entry.access]}`) }}
                </span>
            </UBadge>
        </div>

        <div class="flex flex-row flex-wrap gap-1">
            <UBadge v-for="entry in users" :key="entry.id" color="black" class="inline-flex gap-1" size="md">
                <span class="size-2 rounded-full bg-amber-500" />
                <span :title="`${$t('common.id')} ${entry.userId}`">
                    {{ entry.user?.firstname }}
                    {{ entry.user?.lastname }} -
                    {{ $t(`${i18nKey}.${i18nAccessLevelKey}.${accessLevel[entry.access]}`) }}
                </span>
            </UBadge>
        </div>
    </div>
</template>
