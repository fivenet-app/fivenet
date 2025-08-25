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

const overlay = useOverlay();

const { attr, can, activeChar } = useAuth();

const { game } = useAppConfig();

const selfServicePropsAbsenceDateModal = overlay.create(SelfServicePropsAbsenceDateModal);
</script>

<template>
    <div class="mb-4 flex items-center gap-2 px-4">
        <ProfilePictureImg
            :src="colleague.avatar"
            :name="`${colleague.firstname} ${colleague.lastname}`"
            :enable-popup="true"
            size="3xl"
        />

        <div class="w-full flex-1">
            <div class="flex snap-x flex-row flex-wrap justify-between gap-2 overflow-x-auto">
                <h1 class="flex-1 px-0.5 py-1 text-4xl font-bold break-words sm:pl-1">
                    <ColleagueName :colleague="colleague" />
                </h1>
            </div>

            <div class="inline-flex flex-col gap-2 lg:flex-row">
                <UBadge class="truncate">
                    {{ colleague.jobLabel }}
                    <template v-if="colleague.job !== game.unemployedJobName">
                        ({{ $t('common.rank') }}: {{ colleague.jobGradeLabel }})
                    </template>
                </UBadge>

                <UBadge
                    v-if="colleague.props?.absenceEnd && isFuture(toDate(colleague.props?.absenceEnd))"
                    class="inline-flex items-center gap-1"
                >
                    <UIcon class="size-5" name="i-mdi-island" />
                    <GenericTime :value="colleague.props?.absenceBegin" type="date" />
                    <span>{{ $t('common.to') }}</span>
                    <GenericTime :value="colleague.props?.absenceEnd" type="date" />
                </UBadge>
            </div>
        </div>

        <div class="inline-flex flex-initial flex-col gap-1 sm:flex-row">
            <PartialsBackButton fallback-to="/jobs/colleagues" />

            <UButton
                v-if="
                    can('jobs.JobsService/SetColleagueProps').value &&
                    (colleague.userId === activeChar!.userId ||
                        attr('jobs.JobsService/SetColleagueProps', 'Types', 'AbsenceDate').value) &&
                    checkIfCanAccessColleague(colleague, 'jobs.JobsService/SetColleagueProps')
                "
                icon="i-mdi-island"
                size="md"
                @click="
                    selfServicePropsAbsenceDateModal.open({
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
