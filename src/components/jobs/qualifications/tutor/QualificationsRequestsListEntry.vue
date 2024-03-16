<script lang="ts" setup>
import { ChevronRightIcon, ListStatusIcon } from 'mdi-vue3';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { QualificationRequest, RequestStatus } from '~~/gen/ts/resources/qualifications/qualifications';

defineProps<{
    request: QualificationRequest;
}>();

defineEmits<{
    (e: 'selected'): void;
}>();
</script>

<template>
    <li class="relative flex justify-between px-4 py-5">
        <div class="flex min-w-0 gap-x-4">
            <div class="min-w-0 flex-auto">
                <p class="text-sm font-semibold leading-6 text-gray-100">
                    <NuxtLink :to="{ name: 'jobs-qualifications-id', params: { id: request.qualificationId } }">
                        {{ request.qualification?.abbreviation }}: {{ request.qualification?.title }}
                    </NuxtLink>
                </p>
                <p class="mt-1 flex text-xs leading-5 text-gray-300">
                    <span class="inline-flex gap-1">
                        <span v-if="request.userComment"> ({{ $t('common.summary') }}: {{ request.userComment }})</span>
                    </span>
                </p>
            </div>
        </div>
        <div class="flex shrink-0 items-center gap-x-4">
            <div class="hidden sm:flex sm:flex-col sm:items-end">
                <div class="flex flex-row gap-1">
                    <div class="flex flex-initial flex-row gap-1 rounded-full bg-info-100 px-2 py-1">
                        <ListStatusIcon class="h-5 w-5 text-info-400" aria-hidden="true" />
                        <template v-if="request.status !== undefined">
                            <span class="text-sm font-medium text-info-700">
                                <span class="font-semibold">{{
                                    $t(`enums.qualifications.RequestStatus.${RequestStatus[request.status]}`)
                                }}</span>
                            </span>
                        </template>
                    </div>
                </div>
                <p v-if="request.createdAt" class="mt-1 text-xs leading-5 text-gray-300">
                    {{ $t('common.created_at') }} <GenericTime :value="request.createdAt" />
                </p>
            </div>
            <button type="button" @click="$emit('selected')">
                <ChevronRightIcon class="h-5 w-5 flex-none text-gray-300" aria-hidden="true" />
            </button>
        </div>
    </li>
</template>
