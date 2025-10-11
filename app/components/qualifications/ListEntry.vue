<script lang="ts" setup>
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import OpenClosedBadge from '~/components/partials/OpenClosedBadge.vue';
import { type Qualification, ResultStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import { resultStatusToBadgeColor } from './helpers';

defineProps<{
    qualification: Qualification;
}>();
</script>

<template>
    <li
        class="relative flex justify-between border-default p-2 hover:border-primary-500/25 hover:bg-primary-100/50 sm:px-4 dark:hover:border-primary-400/25 dark:hover:bg-primary-900/10"
    >
        <div class="flex min-w-0 gap-x-2">
            <div class="min-w-0 flex-auto">
                <p class="text-sm leading-6 font-semibold text-toned">
                    <ULink
                        class="inline-flex items-center gap-2"
                        :to="{ name: 'qualifications-id', params: { id: qualification.id } }"
                    >
                        <span class="absolute inset-x-0 -top-px bottom-0" />
                        <span class="text-highlighted"
                            ><template v-if="qualification.abbreviation">{{ qualification.abbreviation }}: </template>
                            {{ !qualification.title ? $t('common.untitled') : qualification.title }}</span
                        >

                        <UBadge
                            v-if="qualification.draft"
                            class="inline-flex gap-1"
                            color="info"
                            size="xs"
                            icon="i-mdi-pencil"
                            :label="$t('common.draft')"
                        />

                        <UBadge
                            v-if="qualification.public"
                            class="inline-flex gap-1"
                            color="neutral"
                            size="xs"
                            icon="i-mdi-earth"
                            :label="$t('common.public')"
                        />

                        <UBadge
                            v-if="qualification?.deletedAt"
                            class="inline-flex gap-1"
                            color="warning"
                            size="xs"
                            icon="i-mdi-calendar-remove"
                        >
                            {{ $t('common.deleted') }}
                            <GenericTime :value="qualification?.deletedAt" type="long" />
                        </UBadge>
                    </ULink>
                </p>
                <p class="mt-1 flex gap-1 text-xs leading-5">
                    <span class="font-semibold">{{ $t('common.description') }}:</span>
                    {{ qualification.description ?? $t('common.na') }}
                </p>
            </div>
        </div>

        <div class="flex shrink-0 items-center gap-x-2">
            <div class="hidden sm:flex sm:flex-col sm:items-end">
                <div class="flex flex-row gap-1">
                    <UBadge
                        v-if="qualification.result?.status"
                        class="inline-flex gap-1"
                        :color="resultStatusToBadgeColor(qualification.result?.status ?? 0)"
                        size="sm"
                        icon="i-mdi-list-status"
                    >
                        {{ $t('common.result') }}:
                        {{ $t(`enums.qualifications.ResultStatus.${ResultStatus[qualification.result?.status ?? 0]}`) }}
                    </UBadge>

                    <OpenClosedBadge :closed="qualification.closed" />
                </div>

                <p v-if="qualification.createdAt" class="mt-1 text-xs leading-5">
                    {{ $t('common.created_at') }} <GenericTime :value="qualification.createdAt" />
                </p>
            </div>

            <UIcon class="size-5 flex-none" name="i-mdi-chevron-right" />
        </div>
    </li>
</template>
