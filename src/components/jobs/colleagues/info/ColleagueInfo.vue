<script lang="ts" setup>
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import SelfServicePropsAbsenceDateModal from '~/components/jobs/colleagues/SelfServicePropsAbsenceDateModal.vue';
import { checkIfCanAccessColleague } from '~/components/jobs/colleagues/helpers';
import { useAuthStore } from '~/store/auth';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import type { Colleague } from '~~/gen/ts/resources/jobs/colleagues';
import { isFuture } from 'date-fns';

defineProps<{
    colleague: Colleague;
}>();

const modal = useModal();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);
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
                    {{ colleague.firstname }} {{ colleague.lastname }}
                </h1>
            </div>

            <div class="inline-flex gap-2">
                <UBadge>
                    {{ colleague.jobLabel }}
                    <span v-if="colleague.jobGrade > 0" class="ml-1">
                        ({{ $t('common.rank') }}: {{ colleague.jobGradeLabel }})</span
                    >
                </UBadge>

                <UBadge
                    v-if="colleague.props?.absenceEnd && isFuture(toDate(colleague.props?.absenceEnd))"
                    class="inline-flex items-center gap-1 rounded-full bg-base-100 px-2.5 py-0.5 text-sm font-medium text-base-800"
                >
                    <UIcon name="i-mdi-island" class="size-5" />
                    <GenericTime :value="colleague.props?.absenceBegin" type="date" />
                    <span>{{ $t('common.to') }}</span>
                    <GenericTime :value="colleague.props?.absenceEnd" type="date" />
                </UBadge>
            </div>
        </div>

        <UButtonGroup class="inline-flex flex-initial">
            <UButton color="black" icon="i-mdi-arrow-back" to="/jobs/colleagues">
                {{ $t('common.back') }}
            </UButton>

            <UButton
                v-if="
                    can('JobsService.SetJobsUserProps').value &&
                    (colleague.userId === activeChar!.userId ||
                        attr('JobsService.SetJobsUserProps', 'Types', 'AbsenceDate').value) &&
                    checkIfCanAccessColleague(activeChar!, colleague, 'JobsService.SetJobsUserProps')
                "
                icon="i-mdi-island"
                size="md"
                @click="
                    modal.open(SelfServicePropsAbsenceDateModal, {
                        userId: colleague.userId,
                        userProps: colleague.props,
                    })
                "
            >
                {{ $t('components.jobs.self_service.set_absence_date') }}
            </UButton>
        </UButtonGroup>
    </div>
</template>
