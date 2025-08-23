<script lang="ts" setup>
import type { TypedRouteFromName } from '@typed-router';
import ColleagueSetLabels from '~/components/jobs/colleagues/ColleagueSetLabels.vue';
import ColleagueSetName from '~/components/jobs/colleagues/ColleagueSetName.vue';
import ColleagueSetNote from '~/components/jobs/colleagues/ColleagueSetNote.vue';
import EmailInfoPopover from '~/components/mailer/EmailInfoPopover.vue';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import type { Colleague } from '~~/gen/ts/resources/jobs/colleagues';

useHead({
    title: 'pages.jobs.colleagues.single.title',
});

definePageMeta({
    title: 'pages.jobs.colleagues.single.title',
    requiresAuth: true,
    permission: 'jobs.JobsService/GetColleague',
    validate: async (route) => {
        route = route as TypedRouteFromName<'jobs-colleagues-id-info'>;
        // Check if the id is made up of digits
        if (typeof route.params.id !== 'string') {
            return false;
        }
        return !!(route.params.id && !isNaN(Number(route.params.id))) && Number(route.params.id) > -1;
    },
});

defineProps<{
    colleague: Colleague;
}>();

defineEmits<{
    (e: 'refresh'): void;
}>();

const { attr, can } = useAuth();
</script>

<template>
    <UContainer class="w-full">
        <div class="w-full grow lg:flex lg:flex-col">
            <div class="flex-1 px-4 py-5 sm:p-0">
                <dl class="space-y-4 sm:space-y-0 xl:grid xl:grid-cols-2">
                    <div class="border-b border-gray-100 sm:flex sm:px-5 sm:py-4 dark:border-gray-800">
                        <dt class="text-sm font-medium sm:w-40 sm:shrink-0 lg:w-48">
                            {{ $t('common.date_of_birth') }}
                        </dt>
                        <dd class="text-base-800 dark:text-base-300 mt-1 text-sm sm:col-span-2 sm:mt-0 sm:ml-6">
                            {{ colleague.dateofbirth }}
                        </dd>
                    </div>

                    <div class="border-b border-gray-100 sm:flex sm:px-5 sm:py-4 dark:border-gray-800">
                        <dt class="text-sm font-medium sm:w-40 sm:shrink-0 lg:w-48">
                            {{ $t('common.phone_number') }}
                        </dt>
                        <dd class="text-base-800 dark:text-base-300 mt-1 text-sm sm:col-span-2 sm:mt-0 sm:ml-6">
                            <PhoneNumberBlock :number="colleague.phoneNumber" />
                        </dd>
                    </div>

                    <div class="border-b border-gray-100 sm:flex sm:px-5 sm:py-4 dark:border-gray-800">
                        <dt class="text-sm font-medium sm:w-40 sm:shrink-0 lg:w-48">
                            {{ $t('common.mail') }}
                        </dt>
                        <dd class="text-base-800 dark:text-base-300 mt-1 text-sm sm:col-span-2 sm:mt-0 sm:ml-6">
                            <EmailInfoPopover :email="colleague.email" />
                        </dd>
                    </div>

                    <div class="border-b border-gray-100 sm:flex sm:px-5 sm:py-4 dark:border-gray-800">
                        <dt class="text-sm font-medium sm:w-40 sm:shrink-0 lg:w-48">
                            {{ $t('common.name') }}
                        </dt>
                        <dd class="text-base-800 dark:text-base-300 mt-1 text-sm sm:col-span-2 sm:mt-0 sm:ml-6">
                            <ColleagueSetName
                                v-if="
                                    can('jobs.JobsService/SetColleagueProps').value &&
                                    attr('jobs.JobsService/SetColleagueProps', 'Types', 'Name').value
                                "
                                v-model:name-prefix="colleague.props!.namePrefix"
                                v-model:name-suffix="colleague.props!.nameSuffix"
                                :user-id="colleague.userId"
                                @refresh="$emit('refresh')"
                            />
                            <div v-else class="text-sm leading-6">
                                <span>{{ $t('common.prefix') }}: {{ colleague?.props?.namePrefix ?? $t('common.na') }}</span>
                                <span>{{ $t('common.suffix') }}: {{ colleague?.props?.nameSuffix ?? $t('common.na') }}</span>
                            </div>
                        </dd>
                    </div>

                    <!-- Labels -->
                    <div
                        v-if="attr('jobs.JobsService/GetColleague', 'Types', 'Labels').value"
                        class="border-b border-gray-100 py-1 hover:bg-primary-100/50 sm:flex sm:px-5 sm:py-4 dark:border-gray-800 dark:hover:bg-primary-900/10"
                    >
                        <dt class="text-sm font-medium sm:w-40 sm:shrink-0 lg:w-48">
                            {{ $t('common.label', 2) }}
                        </dt>
                        <dd class="text-base-800 dark:text-base-300 mt-1 text-sm sm:col-span-2 sm:mt-0 sm:ml-6">
                            <ColleagueSetLabels
                                v-if="
                                    can('jobs.JobsService/SetColleagueProps').value &&
                                    attr('jobs.JobsService/SetColleagueProps', 'Types', 'Labels').value
                                "
                                v-model="colleague.props!.labels"
                                :user-id="colleague.userId"
                                @refresh="$emit('refresh')"
                            />
                            <template v-else>
                                <p v-if="!colleague?.props?.labels?.list.length" class="text-sm leading-6">
                                    {{ $t('common.none', [$t('common.label', 2)]) }}
                                </p>

                                <template v-else>
                                    <div class="flex max-w-80 flex-row flex-wrap gap-1">
                                        <UBadge
                                            v-for="label in colleague?.props?.labels?.list"
                                            :key="label.name"
                                            class="justify-between gap-2"
                                            :class="
                                                isColorBright(hexToRgb(label.color, RGBBlack)!) ? 'text-black!' : 'text-white!'
                                            "
                                            :style="{ backgroundColor: label.color }"
                                            size="md"
                                        >
                                            {{ label.name }}
                                        </UBadge>
                                    </div>
                                </template>
                            </template>
                        </dd>
                    </div>

                    <!-- Note -->
                    <div
                        v-if="attr('jobs.JobsService/GetColleague', 'Types', 'Note').value"
                        class="border-b border-gray-100 py-1 hover:bg-primary-100/50 sm:flex sm:px-5 sm:py-4 dark:border-gray-800 dark:hover:bg-primary-900/10"
                    >
                        <dt class="text-sm font-medium sm:w-40 sm:shrink-0 lg:w-48">
                            {{ $t('common.note') }}
                        </dt>
                        <dd
                            class="text-base-800 dark:text-base-300 mt-1 flex w-full flex-1 text-sm sm:col-span-2 sm:mt-0 sm:ml-6"
                        >
                            <ColleagueSetNote
                                v-model="colleague.props!.note"
                                :user-id="colleague.userId"
                                @refresh="$emit('refresh')"
                            />
                        </dd>
                    </div>
                </dl>
            </div>
        </div>
    </UContainer>
</template>
