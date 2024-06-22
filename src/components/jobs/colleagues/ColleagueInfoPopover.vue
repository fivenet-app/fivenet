<script lang="ts" setup>
import type { Colleague } from '~~/gen/ts/resources/jobs/colleagues';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import { useAuthStore } from '~/store/auth';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { isFuture } from 'date-fns';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';

const props = withDefaults(
    defineProps<{
        userId?: number | string;
        user?: Colleague;
        textClass?: unknown;
        showAvatar?: boolean;
        trailing?: boolean;
        hideProps?: boolean;
    }>(),
    {
        textClass: '' as any,
        showAvatar: undefined,
        trailing: true,
        hideProps: false,
    },
);

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const { popover } = useAppConfig();

const userId = computed(() => {
    if (typeof props.userId === 'string') {
        return parseInt(props.userId);
    }

    return props.userId ?? props.user?.userId ?? 0;
});

const {
    data,
    refresh,
    pending: loading,
    error,
} = useLazyAsyncData(`colleague-info-${userId.value}`, () => getCitizen(userId.value), { immediate: !props.user });

async function getCitizen(id: number): Promise<Colleague> {
    try {
        const call = getGRPCJobsClient().getColleague({
            userId: id,
            infoOnly: true,
        });
        const { response } = await call;

        return response.colleague!;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const user = computed(() => data.value || props.user);

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
            v-bind="$attrs"
            variant="link"
            :padded="false"
            class="inline-flex items-center gap-1 p-px"
            :trailing-icon="trailing ? 'i-mdi-chevron-down' : undefined"
            @click="opened = true"
        >
            <slot name="before" />
            <template #leading v-if="showAvatar">
                <USkeleton v-if="!user && loading" class="h-6 w-6" :ui="{ rounded: 'rounded-full' }" />
                <ProfilePictureImg v-else :src="user?.avatar?.url" :name="`${user?.firstname} ${user?.lastname}`" size="3xs" />
            </template>

            <USkeleton v-if="!user && loading" class="h-8 w-[125px]" />
            <span v-else class="truncate" :class="textClass"> {{ user?.firstname }} {{ user?.lastname }} </span>
            <slot name="after" />
        </UButton>

        <template #panel>
            <div class="flex flex-col gap-2 p-4">
                <UButtonGroup class="inline-flex w-full">
                    <UButton
                        v-if="can('JobsService.GetColleague').value && activeChar?.job === user?.job"
                        variant="link"
                        icon="i-mdi-account"
                        :to="{
                            name: 'jobs-colleagues-id',
                            params: { id: user?.userId ?? 0 },
                        }"
                    >
                        {{ $t('common.profile') }}
                    </UButton>
                    <UButton
                        v-else-if="can('CitizenStoreService.ListCitizens').value"
                        variant="link"
                        icon="i-mdi-account"
                        :to="{
                            name: 'citizens-id',
                            params: { id: user?.userId ?? 0 },
                        }"
                    >
                        {{ $t('common.profile') }}
                    </UButton>

                    <PhoneNumberBlock
                        v-if="user?.phoneNumber"
                        :number="user?.phoneNumber"
                        :hide-number="true"
                        :show-label="true"
                    />
                </UButtonGroup>

                <div v-if="error">
                    <DataErrorBlock :title="$t('common.unable_to_load', [$t('common.colleague', 2)])" :retry="refresh" />
                </div>

                <div v-else-if="loading && !user" class="flex flex-col gap-2 text-gray-900 dark:text-white">
                    <USkeleton class="h-8 w-[250px]" />

                    <div class="flex flex-row items-center gap-2">
                        <USkeleton class="h-7 w-[60px]" />
                        <USkeleton class="h-6 w-[215px]" />
                    </div>
                </div>

                <div v-else-if="user" class="inline-flex flex-row gap-2 text-gray-900 dark:text-white">
                    <div v-if="showAvatar === undefined || showAvatar">
                        <ProfilePictureImg :src="user.avatar?.url" :name="`${user.firstname} ${user.lastname}`" />
                    </div>
                    <div>
                        <UButton
                            v-if="can('JobsService.GetColleague').value && activeChar?.job === user.job"
                            variant="link"
                            :padded="false"
                            :to="{
                                name: 'jobs-colleagues-id',
                                params: { id: user.userId ?? 0 },
                            }"
                        >
                            {{ user.firstname }} {{ user.lastname }}
                        </UButton>
                        <UButton
                            v-else-if="can('CitizenStoreService.ListCitizens').value"
                            variant="link"
                            :padded="false"
                            :to="{
                                name: 'citizens-id',
                                params: { id: user.userId ?? 0 },
                            }"
                        >
                            {{ user.firstname }} {{ user.lastname }}
                        </UButton>
                        <UButton v-else variant="link" :padded="false"> {{ user.firstname }} {{ user.lastname }} </UButton>

                        <p v-if="user.jobLabel" class="text-sm font-normal">
                            <span class="font-semibold">{{ $t('common.job') }}:</span>
                            {{ user.jobLabel }}
                            <span v-if="user.jobGrade > 0 && user.jobGradeLabel"> ({{ user.jobGradeLabel }})</span>
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
