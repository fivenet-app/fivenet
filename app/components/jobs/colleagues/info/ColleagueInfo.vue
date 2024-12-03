<script lang="ts" setup>
import { isFuture } from 'date-fns';
import SelfServicePropsAbsenceDateModal from '~/components/jobs/colleagues/SelfServicePropsAbsenceDateModal.vue';
import { checkIfCanAccessColleague } from '~/components/jobs/colleagues/helpers';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import type { Colleague } from '~~/gen/ts/resources/jobs/colleagues';
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';
import ColleagueName from '../ColleagueName.vue';

defineProps<{
    colleague: Colleague;
}>();

defineEmits<{
    (e: 'update:absenceDates', value: { userId: number; absenceBegin?: Timestamp; absenceEnd?: Timestamp }): void;
}>();

const modal = useModal();

const { attr, can, activeChar } = useAuth();
</script>

<template>
    <div class="mb-4 flex items-center gap-2 px-4">
        <ProfilePictureImg
            :src="colleague.avatar?.url"
            :name="`${colleague.firstname} ${colleague.lastname}`"
            :enable-popup="true"
            size="3xl"
        />

        <div class="w-full flex-1">
            <div class="flex snap-x flex-row flex-wrap justify-between gap-2 overflow-x-auto">
                <h1 class="flex-1 break-words px-0.5 py-1 text-4xl font-bold sm:pl-1">
                    <ColleagueName :colleague="colleague" />
                </h1>
            </div>

            <div class="inline-flex flex-col gap-2 lg:flex-row">
                <UBadge>
                    {{ colleague.jobLabel }}
                    <span v-if="colleague.jobGrade > 0" class="ml-1 truncate">
                        ({{ $t('common.rank') }}: {{ colleague.jobGradeLabel }})</span
                    >
                </UBadge>

                <UBadge
                    v-if="colleague.props?.absenceEnd && isFuture(toDate(colleague.props?.absenceEnd))"
                    class="inline-flex items-center gap-1"
                >
                    <UIcon name="i-mdi-island" class="size-5" />
                    <GenericTime :value="colleague.props?.absenceBegin" type="date" />
                    <span>{{ $t('common.to') }}</span>
                    <GenericTime :value="colleague.props?.absenceEnd" type="date" />
                </UBadge>
            </div>
        </div>

        <div class="inline-flex flex-initial flex-col gap-1 sm:flex-row">
            <UButton color="black" icon="i-mdi-arrow-back" to="/jobs/colleagues">
                {{ $t('common.back') }}
            </UButton>

            <UButton
                v-if="
                    can('JobsService.SetJobsUserProps').value &&
                    (colleague.userId === activeChar!.userId ||
                        attr('JobsService.SetJobsUserProps', 'Types', 'AbsenceDate').value) &&
                    checkIfCanAccessColleague(colleague, 'JobsService.SetJobsUserProps')
                "
                icon="i-mdi-island"
                size="md"
                @click="
                    modal.open(SelfServicePropsAbsenceDateModal, {
                        userId: colleague.userId,
                        userProps: colleague.props,
                        'onUpdate:absenceDates': ($event) => $emit('update:absenceDates', $event),
                    })
                "
            >
                {{ $t('components.jobs.self_service.set_absence_date') }}
            </UButton>
        </div>
    </div>
</template>
