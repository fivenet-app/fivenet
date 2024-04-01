<script lang="ts" setup>
import { EyeIcon, IslandIcon } from 'mdi-vue3';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import type { Colleague } from '~~/gen/ts/resources/jobs/colleagues';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import SelfServicePropsAbsenceDateModal from '~/components/jobs/colleagues/SelfServicePropsAbsenceDateModal.vue';
import { useAuthStore } from '~/store/auth';
import { checkIfCanAccessColleague } from '~/components/jobs/colleagues/helpers';
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';

const props = defineProps<{
    colleague: Colleague;
}>();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

defineEmits<{
    (e: 'update:absenceDates', value: { userId: number; absenceBegin?: Timestamp; absenceEnd?: Timestamp }): void;
}>();

const absenceDateModal = ref(false);

const today = new Date();
today.setHours(0);
today.setMinutes(0);
today.setSeconds(0);
today.setMilliseconds(0);

const showAbsence = props.colleague.props?.absenceEnd && toDate(props.colleague.props?.absenceEnd).getTime() >= today.getTime();
</script>

<template>
    <tr :key="colleague.userId" class="transition-colors even:bg-base-800 hover:bg-neutral/5">
        <SelfServicePropsAbsenceDateModal
            :open="absenceDateModal"
            :user-id="colleague.userId"
            :user-props="colleague.props"
            @close="absenceDateModal = false"
            @update:absence-dates="$emit('update:absenceDates', $event)"
        />

        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-base font-medium text-neutral sm:pl-1">
            <ProfilePictureImg
                :url="colleague.avatar?.url"
                :name="`${colleague.firstname} ${colleague.lastname}`"
                size="sm"
                :rounded="false"
                :enable-popup="true"
            />
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-base font-medium text-neutral sm:pl-1">
            {{ colleague.firstname }} {{ colleague.lastname }}
            <dl class="font-normal lg:hidden">
                <dt class="sr-only">{{ $t('common.job_grade') }}</dt>
                <dd class="mt-1 truncate text-accent-200">
                    {{ colleague.jobGradeLabel }}<span v-if="colleague.jobGrade > 0"> ({{ colleague.jobGrade }})</span>
                </dd>
            </dl>
        </td>
        <td class="hidden whitespace-nowrap p-1 text-left text-accent-200 lg:table-cell">
            {{ colleague.jobGradeLabel }}<span v-if="colleague.jobGrade > 0"> ({{ colleague.jobGrade }})</span>
        </td>
        <td class="hidden whitespace-nowrap py-2 pl-4 pr-3 text-base font-medium text-neutral sm:block sm:pl-1">
            <dl v-if="showAbsence" class="font-normal">
                <dd class="truncate text-accent-200">
                    {{ $t('common.from') }}:
                    <GenericTime :value="colleague.props?.absenceBegin" type="date" />
                </dd>
                <dd class="truncate text-accent-200">
                    {{ $t('common.to') }}: <GenericTime :value="colleague.props?.absenceEnd" type="date" />
                </dd>
            </dl>
        </td>
        <td class="whitespace-nowrap p-1 text-left text-accent-200">
            <PhoneNumberBlock :number="colleague.phoneNumber" />
        </td>
        <td class="whitespace-nowrap p-1 text-left text-accent-200">
            {{ colleague.dateofbirth }}
            <dl v-if="showAbsence" class="block font-normal sm:hidden">
                <dd class="truncate text-accent-200">
                    {{ $t('common.from') }}:
                    <GenericTime :value="colleague.props?.absenceBegin" type="date" />
                </dd>
                <dd class="truncate text-accent-200">
                    {{ $t('common.to') }}: <GenericTime :value="colleague.props?.absenceEnd" type="date" />
                </dd>
            </dl>
        </td>
        <td class="whitespace-nowrap p-1 text-left text-accent-200">
            <div class="flex flex-col justify-end md:flex-row">
                <span
                    v-if="
                        can('JobsService.SetJobsUserProps') &&
                        checkIfCanAccessColleague(activeChar!, colleague, 'JobsService.SetJobsUserProps')
                    "
                    class="flex-initial text-primary-500 hover:text-primary-400"
                    @click="absenceDateModal = true"
                >
                    <IslandIcon class="mr-2.5 h-auto w-5" aria-hidden="true" />
                </span>

                <NuxtLink
                    v-if="
                        can('JobsService.GetColleague') &&
                        checkIfCanAccessColleague(activeChar!, colleague, 'JobsService.GetColleague')
                    "
                    :to="{
                        name: 'jobs-colleagues-id',
                        params: { id: colleague.userId ?? 0 },
                    }"
                    class="flex-initial text-primary-500 hover:text-primary-400"
                >
                    <EyeIcon class="mr-2.5 h-auto w-5" aria-hidden="true" />
                </NuxtLink>
            </div>
        </td>
    </tr>
</template>
