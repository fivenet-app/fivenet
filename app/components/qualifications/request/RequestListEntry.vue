<script lang="ts" setup>
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useAuthStore } from '~/stores/auth';
import { type QualificationRequest, RequestStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import { requestStatusToBadgeColor } from '../helpers';

defineProps<{
    request: QualificationRequest;
}>();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);
</script>

<template>
    <li
        class="relative flex justify-between border-default p-2 hover:border-primary-500/25 hover:bg-primary-100/50 sm:px-4 dark:hover:border-primary-400/25 dark:hover:bg-primary-900/10"
    >
        <div class="flex min-w-0 gap-x-2">
            <div class="min-w-0 flex-auto">
                <p class="text-sm leading-6 font-semibold text-toned">
                    <ULink
                        class="text-highlighted"
                        :to="{ name: 'qualifications-id', params: { id: request.qualificationId } }"
                    >
                        <span class="absolute inset-x-0 -top-px bottom-0" />
                        {{ request.qualification?.abbreviation }}:
                        {{ !request.qualification?.title ? $t('common.untitled') : request.qualification?.title }}
                    </ULink>
                </p>

                <p class="mt-1 flex text-xs leading-5">
                    <span class="inline-flex gap-1 truncate">
                        <span v-if="request.userComment">{{ $t('common.summary') }}: {{ request.userComment }}</span>
                    </span>
                </p>
            </div>
        </div>
        <div class="flex shrink-0 items-center gap-x-2">
            <div class="hidden sm:flex sm:flex-col sm:items-end">
                <UBadge
                    v-if="request.status !== undefined"
                    class="inline-flex gap-1"
                    :color="requestStatusToBadgeColor(request?.status ?? 0)"
                >
                    <UIcon class="size-5" name="i-mdi-list-status" />
                    <span>
                        {{ $t(`enums.qualifications.RequestStatus.${RequestStatus[request.status]}`) }}
                    </span>
                </UBadge>

                <p v-if="request.createdAt" class="mt-1 text-xs leading-5">
                    {{ $t('common.created_at') }} <GenericTime :value="request.createdAt" />
                </p>

                <p v-if="request.userId !== activeChar?.userId" class="mt-1 inline-flex gap-1 text-xs leading-5">
                    {{ $t('common.created_by') }} <CitizenInfoPopover :user="request.user" />
                </p>
            </div>

            <UIcon class="size-5 flex-none" name="i-mdi-chevron-right" />
        </div>
    </li>
</template>
