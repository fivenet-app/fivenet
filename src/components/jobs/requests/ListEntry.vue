<script lang="ts" setup>
import {
    AccountIcon,
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
import Time from '~/components/partials/elements/Time.vue';
import { Request } from '~~/gen/ts/resources/jobs/requests';
import Comments from './Comments.vue';

defineProps<{
    request: Request;
}>();

const open = ref(false);
</script>

<template>
    <li
        :key="request.id.toString()"
        :class="[
            request.deletedAt ? 'hover:bg-warn-700 bg-warn-800' : 'hover:bg-base-700 bg-base-800',
            'flex-initial my-1 rounded-lg',
        ]"
    >
        <div class="mx-2 mt-1 mb-4" @click="open = !open">
            <div class="flex flex-row">
                <p class="py-2 pl-4 pr-3 text-lg font-medium text-neutral sm:pl-0">
                    <span
                        v-if="request.type"
                        class="inline-flex items-center rounded-md bg-primary-400/10 px-2 py-1 text-xs font-medium text-primary-400 ring-1 ring-inset ring-primary-400/30"
                    >
                        {{ request.type.name }}
                    </span>
                    {{ request.title }}
                </p>
                <p
                    class="inline-flex px-2 text-xs font-semibold leading-5 rounded-full bg-primary-100 text-primary-700 my-auto"
                    v-if="request.status"
                >
                    {{ request.status }}
                </p>
                <div
                    class="flex flex-row items-center justify-center flex-1 text-base-200"
                    v-if="request?.approved !== undefined"
                >
                    <div
                        v-if="request?.approved"
                        class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-success-100"
                    >
                        <CheckIcon class="w-5 h-5 text-success-500" aria-hidden="true" />
                        <span class="text-sm font-medium text-success-700">
                            {{ $t('common.approve', 2) }}
                        </span>
                    </div>
                    <div v-else="request.approved" class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-error-100">
                        <CancelIcon class="w-5 h-5 text-error-700" aria-hidden="true" />
                        <span class="text-sm font-medium text-success-700">
                            {{ $t('common.decline', 2) }}
                        </span>
                    </div>
                </div>
                <div class="flex flex-row items-center justify-end flex-1 text-base-200">
                    <div v-if="request?.closed" class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-error-100">
                        <LockIcon class="w-5 h-5 text-success-700" aria-hidden="true" />
                        <span class="text-sm font-medium text-error-700">
                            {{ $t('common.close', 2) }}
                        </span>
                    </div>
                    <div v-else class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-success-100">
                        <LockOpenVariantIcon class="w-5 h-5 text-success-500" aria-hidden="true" />
                        <span class="text-sm font-medium text-success-700">
                            {{ $t('common.open', 2) }}
                        </span>
                    </div>
                </div>
            </div>
            <div class="flex flex-row gap-2 text-base-300 truncate">
                <div v-if="request.beginsAt" class="flex flex-row items-center justify-start flex-1">
                    <ClockStartIcon class="mr-1.5 h-5 w-5 flex-shrink-0 text-base-400" aria-hidden="true" />
                    <p>
                        {{ $t('common.begins_at') }}
                        <Time :value="request.beginsAt" />
                    </p>
                </div>
                <div v-if="request.endsAt" class="flex flex-row items-center justify-end flex-1">
                    <ClockEndIcon class="mr-1.5 h-5 w-5 flex-shrink-0 text-base-400" aria-hidden="true" />
                    <p>
                        {{ $t('common.ends_at') }}
                        <Time :value="request.endsAt" />
                    </p>
                </div>
            </div>
            <div class="flex flex-row gap-2 text-base-300 truncate">
                <div
                    v-if="request.deletedAt"
                    type="button"
                    class="flex flex-row items-center justify-center flex-1 text-base-100 font-bold"
                >
                    <TrashCanIcon class="mr-1.5 h-5 w-5 flex-shrink-0 text-base-400" aria-hidden="true" />
                    {{ $t('common.deleted') }}
                </div>
                <div v-if="request.updatedAt" class="flex flex-row items-center justify-end flex-1">
                    <UpdateIcon class="mr-1.5 h-5 w-5 flex-shrink-0 text-base-400" aria-hidden="true" />
                    <p>
                        {{ $t('common.updated') }}
                        <Time :value="request.updatedAt" :ago="true" />
                    </p>
                </div>
            </div>
            <div class="mt-2 flex flex-row gap-2 text-base-200">
                <div class="flex flex-row items-center justify-start flex-1">
                    <CitizenInfoPopover :user="request.creator">
                        <template v-slot:before>
                            <AccountIcon class="mr-1.5 h-5 w-5 flex-shrink-0 text-base-400" aria-hidden="true" />
                        </template>
                    </CitizenInfoPopover>
                </div>
                <div class="flex flex-row items-center justify-end flex-1">
                    <CalendarIcon class="mr-1.5 h-5 w-5 flex-shrink-0 text-base-400" aria-hidden="true" />
                    <p>
                        {{ $t('common.created_at') }}
                        <Time :value="request.createdAt" />
                    </p>
                </div>
            </div>
        </div>
        <div v-if="open" class="flex flex-col gap-4 m-2">
            <div class="flex flex-row gap-2 items-center justify-center">
                <button
                    type="button"
                    class="inline-flex items-center justify-center gap-x-1.5 rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-primary-400"
                >
                    <template v-if="request.closed">
                        <LockOpenVariantIcon class="w-5 h-5 text-success-500" aria-hidden="true" />
                        {{ $t('common.open', 2) }}
                    </template>
                    <template v-else>
                        <LockIcon class="w-5 h-5 text-error-400" aria-hidden="true" />
                        {{ $t('common.close', 1) }}
                    </template>
                </button>
                <button
                    type="button"
                    class="inline-flex justify-center px-3 py-2 text-sm font-semibold rounded-md bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                    v-if="request.approved !== undefined"
                >
                    <template v-if="request.approved">
                        <CheckIcon class="w-5 h-5 text-success-500" aria-hidden="true" />
                        {{ $t('common.approve', 1) }}
                    </template>
                    <template v-else>
                        <CancelIcon class="w-5 h-5 text-error-400" aria-hidden="true" />
                        {{ $t('common.decline', 1) }}
                    </template>
                </button>
                <span v-else class="isolate inline-flex rounded-md shadow-sm">
                    <button
                        type="button"
                        class="inline-flex justify-center px-3 py-2 text-sm font-semibold rounded-l-md bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                    >
                        <CheckIcon class="w-5 h-5 text-success-500" aria-hidden="true" />
                        {{ $t('common.approve', 1) }}
                    </button>
                    <button
                        type="button"
                        class="inline-flex justify-center px-3 py-2 text-sm font-semibold rounded-r-md bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                    >
                        <CancelIcon class="w-5 h-5 text-error-400" aria-hidden="true" />
                        {{ $t('common.decline', 1) }}
                    </button>
                </span>
                <button
                    type="button"
                    class="inline-flex justify-center px-3 py-2 text-sm font-semibold rounded-md bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                >
                    <TrashCanIcon class="w-5 h-5" />
                    {{ $t('common.delete') }}
                </button>
                <button
                    type="button"
                    @click="open = false"
                    class="inline-flex items-center justify-center gap-x-1.5 rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-primary-400"
                >
                    <CloseIcon class="h-5 w-5" />
                    {{ $t('common.hide') }}
                </button>
            </div>

            <div>
                <p class="text-gray-400 text-base">{{ $t('common.message') }}</p>
                <p class="text-white text-sm">
                    {{ request.message }}
                </p>
            </div>

            <div>
                <Comments :request-id="request.id" />
            </div>
        </div>
    </li>
</template>
