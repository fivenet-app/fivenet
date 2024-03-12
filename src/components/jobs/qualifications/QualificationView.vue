<script lang="ts" setup>
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc';
import {
    AccountIcon,
    CalendarEditIcon,
    CalendarIcon,
    CalendarRemoveIcon,
    ChevronDownIcon,
    FileSearchIcon,
    LockIcon,
    LockOpenVariantIcon,
    NoteCheckIcon,
    PencilIcon,
    TrashCanIcon,
} from 'mdi-vue3';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { AccessLevel } from '~~/gen/ts/resources/jobs/qualifications';
import type { GetQualificationResponse } from '~~/gen/ts/services/jobs/qualifications';

const props = defineProps<{
    id: string;
}>();

const { $grpc } = useNuxtApp();

const { data, pending, refresh, error } = useLazyAsyncData(`jobs-qualification-${props.id}`, () => getQualification(props.id));

async function getQualification(qualificationId: string): Promise<GetQualificationResponse> {
    try {
        const call = $grpc.getJobsQualificationsClient().getQualification({
            qualificationId,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const quali = computed(() => data.value?.qualification);
</script>

<template>
    <div class="py-2 pb-14">
        <div class="px-1 sm:px-2 lg:px-4">
            <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.qualifications', 1)])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.qualifications', 1)])"
                :retry="refresh"
            />
            <DataNoDataBlock v-else-if="!quali" />

            <div v-else class="rounded-lg bg-base-700">
                <div class="h-full px-4 py-6 sm:px-6 lg:px-8">
                    <div>
                        <div>
                            <div class="flex snap-x flex-row flex-wrap justify-between gap-2 overflow-x-auto">
                                <IDCopyBadge
                                    :id="quali?.id"
                                    prefix="QUAL"
                                    :title="{ key: 'notifications.quali?.ment_view.copy_quali?.ment_id.title', parameters: {} }"
                                    :content="{
                                        key: 'notifications.quali?.ment_view.copy_quali?.ment_id.content',
                                        parameters: {},
                                    }"
                                />

                                <div class="flex space-x-2 self-end">
                                    <button
                                        v-if="can('DocStoreService.ToggleDocument')"
                                        type="button"
                                        class="inline-flex items-center gap-x-1.5 rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-primary-400"
                                    >
                                        <template v-if="!quali.closed">
                                            <LockOpenVariantIcon class="h-5 w-5 text-success-500" aria-hidden="true" />
                                            {{ $t('common.open', 1) }}
                                        </template>
                                        <template v-else>
                                            <LockIcon class="h-5 w-5 text-error-400" aria-hidden="true" />
                                            {{ $t('common.close', 1) }}
                                        </template>
                                    </button>
                                    <NuxtLink
                                        v-if="can('DocStoreService.UpdateDocument')"
                                        :to="{
                                            name: 'jobs-qualifications-id-edit',
                                            params: { id: quali.id },
                                        }"
                                        type="button"
                                        class="inline-flex items-center gap-x-1.5 rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-primary-400"
                                    >
                                        <PencilIcon class="-ml-0.5 w-5 h-auto" aria-hidden="true" />
                                        {{ $t('common.edit') }}
                                    </NuxtLink>
                                    <button
                                        v-if="can('DocStoreService.DeleteDocument')"
                                        type="button"
                                        class="inline-flex items-center gap-x-1.5 rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-primary-400"
                                    >
                                        <TrashCanIcon class="-ml-0.5 w-5 h-auto" aria-hidden="true" />
                                        {{ $t('common.delete') }}
                                    </button>
                                </div>
                            </div>

                            <div class="my-4">
                                <h1 class="break-words py-1 pl-0.5 pr-0.5 text-4xl font-bold text-neutral sm:pl-1">
                                    {{ quali.abbreviation }}: {{ quali.title }}
                                </h1>
                            </div>

                            <div class="mb-2 flex gap-2">
                                <div
                                    v-if="quali.closed"
                                    class="flex flex-initial flex-row gap-1 rounded-full bg-error-100 px-2 py-1"
                                >
                                    <LockIcon class="h-5 w-5 text-error-400" aria-hidden="true" />
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

                                <div
                                    v-if="quali.abbreviation"
                                    class="flex flex-initial flex-row gap-1 rounded-full bg-info-100 px-2 py-1 text-info-500"
                                >
                                    <NoteCheckIcon class="w-5 h-auto" aria-hidden="true" />
                                    <span class="text-sm font-medium text-info-800">
                                        {{ quali.abbreviation }}
                                    </span>
                                </div>
                            </div>

                            <div class="flex snap-x flex-row flex-wrap gap-2 overflow-x-auto pb-3 sm:pb-0">
                                <div class="flex flex-initial flex-row gap-1 rounded-full bg-base-100 px-2 py-1 text-base-500">
                                    <AccountIcon class="w-5 h-auto" aria-hidden="true" />
                                    <span class="inline-flex items-center text-sm font-medium text-base-700">
                                        {{ $t('common.created_by') }}
                                        <CitizenInfoPopover
                                            :user="quali?.creator"
                                            class="ml-1 font-medium text-primary-600 hover:text-primary-400"
                                        />
                                    </span>
                                </div>

                                <div class="flex flex-initial flex-row gap-1 rounded-full bg-base-100 px-2 py-1 text-base-500">
                                    <CalendarIcon class="w-5 h-auto" aria-hidden="true" />
                                    <span class="text-sm font-medium text-base-700">
                                        {{ $t('common.created_at') }}
                                        <GenericTime :value="quali?.createdAt" type="long" />
                                    </span>
                                </div>
                                <div
                                    v-if="quali?.updatedAt"
                                    class="flex flex-initial flex-row gap-1 rounded-full bg-base-100 px-2 py-1 text-base-500"
                                >
                                    <CalendarEditIcon class="w-5 h-auto" aria-hidden="true" />
                                    <span class="text-sm font-medium text-base-700">
                                        {{ $t('common.updated_at') }}
                                        <GenericTime :value="quali?.updatedAt" type="long" />
                                    </span>
                                </div>
                                <div
                                    v-if="quali?.deletedAt"
                                    class="flex flex-initial flex-row gap-1 rounded-full bg-base-100 px-2 py-1 text-base-500"
                                >
                                    <CalendarRemoveIcon class="w-5 h-auto" aria-hidden="true" />
                                    <span class="text-sm font-medium text-base-700">
                                        {{ $t('common.deleted') }}
                                        <GenericTime :value="quali?.deletedAt" type="long" />
                                    </span>
                                </div>
                            </div>

                            <div class="my-2">
                                <h2 class="sr-only">
                                    {{ $t('common.content') }}
                                </h2>
                                <div class="break-words rounded-lg bg-base-800 text-neutral">
                                    <!-- eslint-disable vue/no-v-html -->
                                    <div
                                        class="prose prose-invert min-w-full rounded-md bg-base-900 px-4 py-4"
                                        v-html="quali?.description"
                                    ></div>
                                </div>
                            </div>

                            <div class="w-full">
                                <Disclosure
                                    v-slot="{ open }"
                                    as="div"
                                    class="w-full border-neutral/20 text-neutral hover:border-neutral/70"
                                    :default-open="true"
                                >
                                    <DisclosureButton
                                        :class="[
                                            open ? 'rounded-t-lg border-b-0' : 'rounded-lg',
                                            'flex w-full items-start justify-between border-2 border-inherit p-2 text-left transition-colors',
                                        ]"
                                    >
                                        <span class="inline-flex items-center text-base font-semibold leading-7">
                                            <LockIcon class="mr-2 w-5 h-auto" aria-hidden="true" />
                                            {{ $t('common.access') }}
                                        </span>
                                        <span class="ml-6 flex h-7 items-center">
                                            <ChevronDownIcon
                                                :class="[open ? 'upsidedown' : '', 'h-autotransition-transform w-5']"
                                                aria-hidden="true"
                                            />
                                        </span>
                                    </DisclosureButton>
                                    <DisclosurePanel class="rounded-b-lg border-2 border-t-0 border-inherit transition-colors">
                                        <div class="mx-4 flex flex-row flex-wrap gap-1 pb-2">
                                            <template v-if="!quali.access || quali.access?.jobs.length === 0">
                                                <DataNoDataBlock
                                                    :icon="FileSearchIcon"
                                                    :message="$t('common.not_found', [$t('common.access', 2)])"
                                                />
                                            </template>
                                            <div
                                                v-for="entry in quali.access?.jobs"
                                                :key="entry.id"
                                                class="flex flex-initial snap-x snap-start items-center gap-1 overflow-x-auto whitespace-nowrap rounded-full bg-info-100 px-2 py-1"
                                            >
                                                <span class="h-2 w-2 rounded-full bg-info-500" aria-hidden="true" />
                                                <span class="text-sm font-medium text-info-800"
                                                    >{{ entry.jobLabel
                                                    }}<span
                                                        v-if="entry.minimumGrade > 0"
                                                        :title="`${entry.jobLabel} - ${$t('common.rank')} ${entry.minimumGrade}`"
                                                    >
                                                        ({{ entry.jobGradeLabel }})</span
                                                    >
                                                    -
                                                    {{
                                                        $t(`enums.jobs.qualifications.AccessLevel.${AccessLevel[entry.access]}`)
                                                    }}
                                                </span>
                                            </div>
                                        </div>
                                    </DisclosurePanel>
                                </Disclosure>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
