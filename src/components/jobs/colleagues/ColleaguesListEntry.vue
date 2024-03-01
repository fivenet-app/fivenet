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

defineProps<{
    colleague: Colleague;
}>();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

defineEmits<{
    (e: 'update:absenceDate', value: { userId: number; absenceDate?: Timestamp }): void;
}>();

const absenceDateModal = ref(false);
</script>

<template>
    <tr :key="colleague.userId" class="transition-colors even:bg-base-800 hover:bg-neutral/5">
        <SelfServicePropsAbsenceDateModal
            :open="absenceDateModal"
            :user-id="colleague.userId"
            :user-props="colleague.props"
            @close="absenceDateModal = false"
            @update:absence-date="$emit('update:absenceDate', $event)"
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
        <td class="whitespace-nowrap px-1 py-1 text-left text-accent-200 hidden lg:table-cell">
            {{ colleague.jobGradeLabel }}<span v-if="colleague.jobGrade > 0"> ({{ colleague.jobGrade }})</span>
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-base font-medium text-neutral sm:pl-1">
            <GenericTime
                v-if="colleague.props?.absenceDate && toDate(colleague.props?.absenceDate).getTime() > new Date().getTime()"
                :value="colleague.props?.absenceDate"
                type="date"
            />
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-accent-200">
            <PhoneNumberBlock :number="colleague.phoneNumber" />
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-accent-200">
            {{ colleague.dateofbirth }}
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-accent-200">
            <div class="flex flex-col md:flex-row justify-end">
                <button
                    v-if="
                        can('JobsService.SetJobsUserProps') &&
                        checkIfCanAccessColleague(activeChar!, colleague, 'JobsService.SetJobsUserProps')
                    "
                    type="button"
                    class="flex-initial text-primary-500 hover:text-primary-400"
                    @click="absenceDateModal = true"
                >
                    <IslandIcon class="mr-2.5 w-5 h-auto" aria-hidden="true" />
                </button>

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
                    <EyeIcon class="mr-2.5 w-5 h-auto" aria-hidden="true" />
                </NuxtLink>
            </div>
        </td>
    </tr>
</template>
