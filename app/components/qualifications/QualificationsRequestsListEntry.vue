<script lang="ts" setup>
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useAuthStore } from '~/stores/auth';
import type { QualificationRequest } from '~~/gen/ts/resources/qualifications/qualifications';
import { RequestStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import { requestStatusToBadgeColor } from './helpers';

defineProps<{
    request: QualificationRequest;
}>();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);
</script>

<template>
    <li
        class="hover:border-primary-500/25 dark:hover:border-primary-400/25 hover:bg-primary-100/50 dark:hover:bg-primary-900/10 relative flex justify-between border-white px-2 py-2 sm:px-4 dark:border-gray-900"
    >
        <div class="flex min-w-0 gap-x-2">
            <div class="min-w-0 flex-auto">
                <p class="text-sm font-semibold leading-6 text-gray-100">
                    <ULink :to="{ name: 'qualifications-id', params: { id: request.qualificationId } }">
                        <span class="absolute inset-x-0 -top-px bottom-0" />
                        {{ request.qualification?.abbreviation }}: {{ request.qualification?.title }}
                    </ULink>
                </p>

                <p class="mt-1 flex text-xs leading-5">
                    <span class="inline-flex gap-1">
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
