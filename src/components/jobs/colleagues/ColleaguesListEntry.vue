<script lang="ts" setup>
import { EyeIcon, IslandIcon } from 'mdi-vue3';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import type { Colleague } from '~~/gen/ts/resources/jobs/colleagues';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import SelfServicePropsAbsenceDateModal from '~/components/jobs/colleagues/SelfServicePropsAbsenceDateModal.vue';
import { useAuthStore } from '~/store/auth';
import { checkIfCanAccessColleague } from '~/components/jobs/colleagues/helpers';

defineProps<{
    user: Colleague;
}>();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const absenceDateModal = ref(false);
</script>

<template>
    <tr :key="user.userId" class="transition-colors even:bg-base-800 hover:bg-neutral/5">
        <SelfServicePropsAbsenceDateModal
            :open="absenceDateModal"
            :user-id="user.userId"
            :user-props="user.props"
            @close="absenceDateModal = false"
        />
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-base font-medium text-neutral sm:pl-1">
            <ProfilePictureImg
                :url="user.avatar?.url"
                :name="`${user.firstname} ${user.lastname}`"
                size="sm"
                :rounded="false"
            />
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-base font-medium text-neutral sm:pl-1">
            {{ user.firstname }} {{ user.lastname }}
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-accent-200">{{ user.jobGradeLabel }} ({{ user.jobGrade }})</td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-base font-medium text-neutral sm:pl-1">
            <GenericTime v-if="user.props?.absenceDate" :value="user.props?.absenceDate" type="date" />
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-accent-200">
            <PhoneNumberBlock :number="user.phoneNumber" />
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-accent-200">
            {{ user.dateofbirth }}
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-accent-200">
            <div class="flex flex-row justify-end">
                <button
                    v-if="
                        can('JobsService.SetJobsUserProps') &&
                        checkIfCanAccessColleague(activeChar!, user, 'JobsService.SetJobsUserProps')
                    "
                    type="button"
                    class="flex-initial text-primary-500 hover:text-primary-400"
                    @click="absenceDateModal = true"
                >
                    <IslandIcon class="mr-2.5 w-5 h-auto" />
                </button>
                <NuxtLink
                    v-if="
                        can('JobsService.GetColleague') &&
                        checkIfCanAccessColleague(activeChar!, user, 'JobsService.GetColleague')
                    "
                    :to="{
                        name: 'jobs-colleagues-id',
                        params: { id: user.userId ?? 0 },
                    }"
                    class="flex-initial text-primary-500 hover:text-primary-400"
                >
                    <EyeIcon class="mr-2.5 w-5 h-auto" />
                </NuxtLink>
            </div>
        </td>
    </tr>
</template>
