<script lang="ts" setup>
import { isFuture } from 'date-fns';
import EmailBlock from '~/components/partials/citizens/EmailBlock.vue';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import type { ClassProp } from '~/utils/types';
import { getJobsJobsClient } from '~~/gen/ts/clients';
import type { Colleague } from '~~/gen/ts/resources/jobs/colleagues/colleagues';
import ColleagueName from './ColleagueName.vue';

const props = withDefaults(
    defineProps<{
        userId?: number;
        user?: Colleague;
        textClass?: ClassProp;
        showAvatar?: boolean;
        trailing?: boolean;
        hideProps?: boolean;
    }>(),
    {
        userId: undefined,
        user: undefined,
        textClass: '',
        showAvatar: undefined,
        trailing: true,
        hideProps: false,
    },
);

const { can, activeChar } = useAuth();

const { game, popover } = useAppConfig();

const jobsJobsClient = await getJobsJobsClient();

const userId = computed(() => props.userId ?? props.user?.userId ?? 0);

const { data, refresh, status, error } = useLazyAsyncData(`colleague-info-${userId.value}`, () => getCitizen(userId.value), {
    immediate: !props.user,
});

async function getCitizen(id: number): Promise<Colleague> {
    const call = jobsJobsClient.getColleague({
        userId: id,
        infoOnly: true,
    });
    const { response } = await call;

    return response.colleague!;
}

const user = computed(() =>
    data.value || props.user
        ? ({
              ...data.value,
              ...props.user,
          } as Colleague)
        : undefined,
);

const opened = ref(false);
watchOnce(opened, async () => {
    if (props.user) {
        useTimeoutFn(async () => refresh(), popover.waitTime);
    }
});
</script>

<template>
    <template v-if="!user && !userId">
        <span class="inline-flex items-center">
            <slot name="before" />
            <span>{{ $t('common.na') }}</span>
            <slot name="after" />
        </span>
    </template>

    <UPopover v-else>
        <UButton
            class="inline-flex items-center gap-1 p-px"
            variant="link"
            :trailing-icon="trailing ? 'i-mdi-chevron-down' : undefined"
            v-bind="$attrs"
            @click.prevent="opened = true"
        >
            <template v-if="showAvatar" #leading>
                <slot name="before" />
                <USkeleton v-if="!user && isRequestPending(status)" class="h-6 w-6" />
                <ProfilePictureImg
                    v-else
                    :src="user?.profilePicture"
                    :name="`${user?.firstname} ${user?.lastname}`"
                    size="3xs"
                />
            </template>

            <USkeleton v-if="!user && isRequestPending(status)" class="h-8 w-[125px]" />
            <span v-else class="truncate" :class="textClass"> <ColleagueName v-if="user" :colleague="user" /> </span>
            <slot name="after" />
        </UButton>

        <template #content>
            <div class="flex flex-col gap-2 p-4">
                <UFieldGroup class="inline-flex w-full">
                    <UButton
                        v-if="can('citizens.CitizensService/ListCitizens').value"
                        variant="link"
                        icon="i-mdi-account"
                        :label="$t('common.profile')"
                        :to="`/citizens/${userId}`"
                    />

                    <UButton
                        v-if="can('jobs.JobsService/GetColleague').value && user?.job === activeChar?.job"
                        variant="link"
                        icon="i-mdi-briefcase"
                        :label="$t('common.colleague')"
                        :to="`/jobs/colleagues/${userId}/info`"
                    />

                    <PhoneNumberBlock v-if="user?.phoneNumber" :number="user?.phoneNumber" hide-number show-label />

                    <EmailBlock v-if="user?.email" :email="user.email" hide-email />
                </UFieldGroup>

                <div v-if="error">
                    <DataErrorBlock
                        :title="$t('common.unable_to_load', [$t('common.colleague', 2)])"
                        :error="error"
                        :retry="refresh"
                    />
                </div>

                <div v-else-if="isRequestPending(status) && !user" class="flex flex-col gap-2 text-highlighted">
                    <USkeleton class="h-8 w-[250px]" />

                    <div class="flex flex-row items-center gap-2">
                        <USkeleton class="h-7 w-[60px]" />
                        <USkeleton class="h-6 w-[215px]" />
                    </div>
                </div>

                <div v-else-if="user" class="flex flex-row gap-2 text-highlighted">
                    <div v-if="showAvatar === undefined || showAvatar">
                        <ProfilePictureImg :src="user.profilePicture" :name="`${user.firstname} ${user.lastname}`" />
                    </div>
                    <div>
                        <UButton
                            v-if="activeChar?.job === user.job && can('jobs.JobsService/GetColleague').value"
                            variant="link"
                            :to="{
                                name: 'jobs-colleagues-id',
                                params: { id: user.userId ?? 0 },
                            }"
                        >
                            <ColleagueName :colleague="user" />
                        </UButton>
                        <UButton
                            v-else-if="can('citizens.CitizensService/ListCitizens').value"
                            variant="link"
                            :to="{
                                name: 'citizens-id',
                                params: { id: user.userId ?? 0 },
                            }"
                        >
                            <ColleagueName :colleague="user" />
                        </UButton>
                        <UButton v-else variant="link">
                            <ColleagueName :colleague="user" />
                        </UButton>

                        <p v-if="user.jobLabel" class="text-sm font-normal">
                            <span class="font-semibold">{{ $t('common.job') }}:</span>
                            {{ user.jobLabel }}
                            <template v-if="user.jobGradeLabel && user.job !== game.unemployedJobName">
                                ({{ $t('common.rank') }}: {{ user.jobGradeLabel }})
                            </template>
                        </p>

                        <p v-if="user.dateofbirth" class="text-sm font-normal">
                            <span class="font-semibold">{{ $t('common.date_of_birth') }}:</span>
                            {{ user.dateofbirth }}
                        </p>

                        <template v-if="!hideProps">
                            <div
                                v-if="user.props?.absenceEnd && isFuture(toDate(user.props?.absenceEnd))"
                                class="text-sm font-normal"
                            >
                                <span class="font-semibold">{{ $t('common.absent') }}:</span>
                                <dl class="text-sm font-normal">
                                    <dd class="truncate">
                                        {{ $t('common.from') }}:
                                        <GenericTime :value="user.props?.absenceBegin" type="date" />
                                    </dd>
                                    <dd class="truncate">
                                        {{ $t('common.to') }}: <GenericTime :value="user.props?.absenceEnd" type="date" />
                                    </dd>
                                </dl>
                            </div>
                        </template>
                    </div>
                </div>
            </div>
        </template>
    </UPopover>
</template>
