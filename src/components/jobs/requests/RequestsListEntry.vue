<script lang="ts" setup>
import {
    AccountIcon,
    AsteriskIcon,
    CalendarIcon,
    CancelIcon,
    CheckIcon,
    ClockEndIcon,
    ClockStartIcon,
    CloseIcon,
    LockIcon,
    LockOpenVariantIcon,
    TrashCanIcon,
    UpdateIcon,
} from 'mdi-vue3';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { Request } from '~~/gen/ts/resources/jobs/requests';
import RequestsComments from '~/components/jobs/requests/RequestsComments.vue';

defineProps<{
    request: Request;
}>();

const open = ref(false);
</script>

<template>
    <li
        :key="request.id"
        :class="[
            request.deletedAt ? 'bg-warn-800 hover:bg-warn-700' : 'bg-base-800 hover:bg-base-700',
            'my-1 flex-initial rounded-lg',
        ]"
    >
        <div class="mx-2 mb-4 mt-1" @click="open = !open">
            <div class="flex flex-row">
                <p class="py-2 pl-4 pr-3 text-lg font-medium text-neutral sm:pl-1">
                    <span
                        v-if="request.type"
                        class="inline-flex items-center rounded-md bg-primary-400/10 px-2 py-1 text-xs font-medium text-primary-400 ring-1 ring-inset ring-primary-400/30"
                    >
                        {{ request.type.name }}
                    </span>
                    {{ request.title }}
                </p>
                <p
                    v-if="request.status"
                    class="my-auto inline-flex rounded-full bg-primary-100 px-2 text-xs font-semibold leading-5 text-primary-700"
                >
                    {{ request.status }}
                </p>
                <div
                    v-if="request?.approved !== undefined"
                    class="flex flex-1 flex-row items-center justify-center text-base-200"
                >
                    <div
                        v-if="request?.approved"
                        class="flex flex-initial flex-row gap-1 rounded-full bg-success-100 px-2 py-1"
                    >
                        <CheckIcon class="h-5 w-5 text-success-500" aria-hidden="true" />
                        <span class="text-sm font-medium text-success-700">
                            {{ $t('common.approve', 2) }}
                        </span>
                    </div>
                    <div
                        v-else-if="request.approved"
                        class="flex flex-initial flex-row gap-1 rounded-full bg-error-100 px-2 py-1"
                    >
                        <CancelIcon class="h-5 w-5 text-error-700" aria-hidden="true" />
                        <span class="text-sm font-medium text-success-700">
                            {{ $t('common.decline', 2) }}
                        </span>
                    </div>
                </div>
                <div class="flex flex-1 flex-row items-center justify-end text-base-200">
                    <div v-if="request?.closed" class="flex flex-initial flex-row gap-1 rounded-full bg-error-100 px-2 py-1">
                        <LockIcon class="h-5 w-5 text-success-700" aria-hidden="true" />
                        <span class="text-sm font-medium text-error-700">
                            {{ $t('common.close', 2) }}
                        </span>
                    </div>
                    <div v-else class="flex flex-initial flex-row gap-1 rounded-full bg-success-100 px-2 py-1">
                        <LockOpenVariantIcon class="h-5 w-5 text-success-500" aria-hidden="true" />
                        <span class="text-sm font-medium text-success-700">
                            {{ $t('common.open', 2) }}
                        </span>
                    </div>
                </div>
            </div>
            <div class="flex flex-row gap-2 truncate text-base-300">
                <div v-if="request.beginsAt" class="flex flex-1 flex-row items-center justify-start">
                    <ClockStartIcon class="mr-1.5 h-5 w-5 flex-shrink-0 text-base-400" aria-hidden="true" />
                    <p>
                        {{ $t('common.begins_at') }}
                        <GenericTime :value="request.beginsAt" />
                    </p>
                </div>
                <div v-if="request.endsAt" class="flex flex-1 flex-row items-center justify-end">
                    <ClockEndIcon class="mr-1.5 h-5 w-5 flex-shrink-0 text-base-400" aria-hidden="true" />
                    <p>
                        {{ $t('common.ends_at') }}
                        <GenericTime :value="request.endsAt" />
                    </p>
                </div>
            </div>
            <div class="flex flex-row gap-2 truncate text-base-300">
                <div
                    v-if="request.deletedAt"
                    type="button"
                    class="flex flex-1 flex-row items-center justify-center font-bold text-base-100"
                >
                    <TrashCanIcon class="mr-1.5 h-5 w-5 flex-shrink-0 text-base-400" aria-hidden="true" />
                    {{ $t('common.deleted') }}
                </div>
                <div v-if="request.updatedAt" class="flex flex-1 flex-row items-center justify-end">
                    <UpdateIcon class="mr-1.5 h-5 w-5 flex-shrink-0 text-base-400" aria-hidden="true" />
                    <p>
                        {{ $t('common.updated') }}
                        <GenericTime :value="request.updatedAt" :ago="true" />
                    </p>
                </div>
            </div>
            <div class="mt-2 flex flex-row gap-2 text-base-200">
                <div class="flex flex-1 flex-row items-center justify-start gap-2">
                    <CitizenInfoPopover :user="request.creator">
                        <template #before>
                            <AccountIcon class="mr-1.5 h-5 w-5 flex-shrink-0 text-base-400" aria-hidden="true" />
                        </template>
                    </CitizenInfoPopover>

                    <CitizenInfoPopover v-if="request.approverUser" :user="request.approverUser">
                        <template #before>
                            <AsteriskIcon class="mr-1.5 h-5 w-5 flex-shrink-0 text-base-400" aria-hidden="true" />
                        </template>
                    </CitizenInfoPopover>
                </div>
                <div class="flex flex-1 flex-row items-center justify-end">
                    <CalendarIcon class="mr-1.5 h-5 w-5 flex-shrink-0 text-base-400" aria-hidden="true" />
                    <p>
                        {{ $t('common.created_at') }}
                        <GenericTime :value="request.createdAt" />
                    </p>
                </div>
            </div>
        </div>
        <div v-if="open" class="m-2 flex flex-col gap-4">
            <div class="flex flex-row items-center justify-center gap-2">
                <button
                    type="button"
                    class="inline-flex items-center justify-center gap-x-1.5 rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-primary-400"
                >
                    <template v-if="request.closed">
                        <LockOpenVariantIcon class="h-5 w-5 text-success-500" aria-hidden="true" />
                        {{ $t('common.open', 1) }}
                    </template>
                    <template v-else>
                        <LockIcon class="h-5 w-5 text-error-400" aria-hidden="true" />
                        {{ $t('common.close', 1) }}
                    </template>
                </button>
                <button
                    v-if="request.approved !== undefined"
                    type="button"
                    class="inline-flex justify-center rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                >
                    <template v-if="request.approved">
                        <CancelIcon class="h-5 w-5 text-error-400" aria-hidden="true" />
                        {{ $t('common.decline', 1) }}
                    </template>
                    <template v-else>
                        <CheckIcon class="h-5 w-5 text-success-500" aria-hidden="true" />
                        {{ $t('common.approve', 1) }}
                    </template>
                </button>
                <span v-else class="isolate inline-flex rounded-md shadow-sm">
                    <button
                        type="button"
                        class="inline-flex justify-center rounded-l-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                    >
                        <CheckIcon class="h-5 w-5 text-success-500" aria-hidden="true" />
                        {{ $t('common.approve', 1) }}
                    </button>
                    <button
                        type="button"
                        class="inline-flex justify-center rounded-r-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                    >
                        <CancelIcon class="h-5 w-5 text-error-400" aria-hidden="true" />
                        {{ $t('common.decline', 1) }}
                    </button>
                </span>
                <button
                    type="button"
                    class="inline-flex justify-center rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                >
                    <TrashCanIcon class="h-5 w-5" />
                    {{ $t('common.delete') }}
                </button>
                <button
                    type="button"
                    class="inline-flex items-center justify-center gap-x-1.5 rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-primary-400"
                    @click="open = false"
                >
                    <CloseIcon class="h-5 w-5" />
                    {{ $t('common.hide') }}
                </button>
            </div>

            <div>
                <p class="text-base text-gray-400">{{ $t('common.message') }}</p>
                <p class="text-sm text-neutral">
                    {{ request.message }}
                </p>
            </div>

            <div>
                <RequestsComments :request-id="request.id" />
            </div>
        </div>
    </li>
</template>
