<script lang="ts" setup>
import type { Colleague } from '~~/gen/ts/resources/jobs/colleagues';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import { useAuthStore } from '~/store/auth';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { isFuture } from 'date-fns';

withDefaults(
    defineProps<{
        user: Colleague | undefined;
        noPopover?: boolean;
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
</script>

<template>
    <template v-if="!user">
        <span class="inline-flex items-center">
            <slot name="before" />
            <span>N/A</span>
            <slot name="after" />
        </span>
    </template>
    <template v-else-if="noPopover">
        <span class="inline-flex items-center">
            <slot name="before" />
            <UButton
                variant="link"
                :padded="false"
                :to="{
                    name: activeChar?.job === user.job ? 'jobs-colleagues-id' : 'citizens-id',
                    params: { id: user.userId ?? 0 },
                }"
            >
                {{ user.firstname }} {{ user.lastname }}
            </UButton>
            <span v-if="user.phoneNumber">
                <PhoneNumberBlock v-if="user.phoneNumber" :number="user.phoneNumber" :hide-number="true" :show-label="false" />
            </span>
            <slot name="after" />
        </span>
    </template>
    <UPopover v-else>
        <UButton
            variant="link"
            :padded="false"
            class="inline-flex items-center"
            :trailing-icon="trailing ? 'i-mdi-chevron-down' : undefined"
        >
            <slot name="before" />
            <span class="truncate" :class="textClass"> {{ user.firstname }} {{ user.lastname }} </span>
            <slot name="after" />
        </UButton>

        <template #panel>
            <div class="flex flex-col gap-2 p-4">
                <UButtonGroup class="inline-flex w-full">
                    <UButton
                        v-if="can('JobsService.GetColleague') && activeChar?.job === user.job"
                        variant="link"
                        icon="i-mdi-account"
                        :to="{
                            name: 'jobs-colleagues-id',
                            params: { id: user.userId ?? 0 },
                        }"
                    >
                        {{ $t('common.profile') }}
                    </UButton>
                    <UButton
                        v-else-if="can('CitizenStoreService.ListCitizens')"
                        variant="link"
                        icon="i-mdi-account"
                        :to="{
                            name: 'citizens-id',
                            params: { id: user.userId ?? 0 },
                        }"
                    >
                        {{ $t('common.profile') }}
                    </UButton>

                    <PhoneNumberBlock
                        v-if="user.phoneNumber"
                        :number="user.phoneNumber"
                        :hide-number="true"
                        :show-label="true"
                    />
                </UButtonGroup>

                <div class="inline-flex flex-row gap-2 text-gray-900 dark:text-white">
                    <div v-if="showAvatar === undefined || showAvatar">
                        <ProfilePictureImg :src="user.avatar?.url" :name="`${user.firstname} ${user.lastname}`" />
                    </div>
                    <div>
                        <UButton
                            v-if="can('JobsService.GetColleague') && activeChar?.job === user.job"
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
                            v-else-if="can('CitizenStoreService.ListCitizens')"
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
