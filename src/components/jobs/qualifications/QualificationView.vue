<script lang="ts" setup>
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue';
import { useConfirmDialog } from '@vueuse/core';
import {
    AccountIcon,
    AccountSchoolIcon,
    CalendarEditIcon,
    CalendarIcon,
    CalendarRemoveIcon,
    ChevronDownIcon,
    ListStatusIcon,
    LockIcon,
    LockOpenVariantIcon,
    MailIcon,
    PencilIcon,
    SigmaIcon,
    TestTubeIcon,
    TrashCanIcon,
} from 'mdi-vue3';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { AccessLevel, RequestStatus, ResultStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import type { DeleteQualificationResponse, GetQualificationResponse } from '~~/gen/ts/services/qualifications/qualifications';
import { checkQualificationAccess } from '~/components/jobs/qualifications/helpers';
import ConfirmDialog from '~/components/partials/ConfirmDialog.vue';
import QualificationRequestUserModal from '~/components/jobs/qualifications/QualificationRequestUserModal.vue';
import QualificationsRequestsList from '~/components/jobs/qualifications/tutor/QualificationsRequestsList.vue';
import QualificationsResultsList from '~/components/jobs/qualifications/tutor/QualificationsResultsList.vue';

const props = defineProps<{
    id: string;
}>();

const { $grpc } = useNuxtApp();

const { data, pending, refresh, error } = useLazyAsyncData(`qualification-${props.id}`, () => getQualification(props.id));

async function getQualification(qualificationId: string): Promise<GetQualificationResponse> {
    try {
        const call = $grpc.getQualificationsClient().getQualification({
            qualificationId,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

async function deleteQualification(qualificationId: string): Promise<DeleteQualificationResponse> {
    try {
        const call = $grpc.getQualificationsClient().deleteQualification({
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
const canDo = computed(() => ({
    take: checkQualificationAccess(quali.value?.access, quali.value?.creator, AccessLevel.TAKE),
    request: checkQualificationAccess(quali.value?.access, quali.value?.creator, AccessLevel.REQUEST),
    grade: checkQualificationAccess(quali.value?.access, quali.value?.creator, AccessLevel.GRADE),
    edit: checkQualificationAccess(quali.value?.access, quali.value?.creator, AccessLevel.EDIT),
}));

const { isRevealed, reveal, confirm, cancel, onConfirm } = useConfirmDialog();
onConfirm(async (id: string) => deleteQualification(id));

const openRequest = ref(false);
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
                <ConfirmDialog :open="isRevealed" :cancel="cancel" :confirm="() => confirm(quali!.id)" />

                <QualificationRequestUserModal
                    :qualification-id="quali.id"
                    :open="openRequest"
                    @close="openRequest = false"
                    @updated-request="quali.request = $event"
                />

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
                                    <template v-if="!canDo.edit">
                                        <UButton
                                            v-if="canDo.take"
                                            class="inline-flex items-center gap-x-1.5 rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold hover:bg-primary-400"
                                            @click="openRequest = true"
                                        >
                                            <TestTubeIcon class="-ml-0.5 h-auto w-5" />
                                            {{ $t('components.qualifications.take_test') }}
                                        </UButton>
                                        <UButton
                                            v-else-if="canDo.request"
                                            class="inline-flex items-center gap-x-1.5 rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold hover:bg-primary-400"
                                            :disabled="quali.request?.status === RequestStatus.PENDING"
                                            @click="openRequest = true"
                                        >
                                            <AccountSchoolIcon class="-ml-0.5 h-auto w-5" />
                                            {{ $t('common.request') }}
                                        </UButton>
                                    </template>
                                    <NuxtLink
                                        v-if="can('QualificationsService.UpdateQualification') && canDo.edit"
                                        :to="{
                                            name: 'jobs-qualifications-id-edit',
                                            params: { id: quali.id },
                                        }"
                                        class="inline-flex items-center gap-x-1.5 rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold hover:bg-primary-400"
                                    >
                                        <PencilIcon class="-ml-0.5 h-auto w-5" />
                                        {{ $t('common.edit') }}
                                    </NuxtLink>
                                    <UButton
                                        v-if="can('QualificationsService.DeleteQualification') && canDo.edit"
                                        class="inline-flex items-center gap-x-1.5 rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold hover:bg-primary-400"
                                        @click="reveal(quali.id)"
                                    >
                                        <TrashCanIcon class="-ml-0.5 h-auto w-5" />
                                        {{ $t('common.delete') }}
                                    </UButton>
                                </div>
                            </div>

                            <div class="my-4">
                                <h1 class="break-words px-0.5 py-1 text-4xl font-bold sm:pl-1">
                                    {{ quali.abbreviation }}: {{ quali.title }}
                                </h1>
                                <p v-if="quali.description" class="break-words px-0.5 py-1 text-base font-bold sm:pl-1">
                                    {{ quali.description }}
                                </p>
                            </div>

                            <div class="mb-2 flex gap-2">
                                <div
                                    v-if="quali.closed"
                                    class="flex flex-initial flex-row gap-1 rounded-full bg-error-100 px-2 py-1"
                                >
                                    <LockIcon class="size-5 text-error-400" />
                                    <span class="text-sm font-medium text-error-700">
                                        {{ $t('common.close', 2) }}
                                    </span>
                                </div>
                                <div v-else class="flex flex-initial flex-row gap-1 rounded-full bg-success-100 px-2 py-1">
                                    <LockOpenVariantIcon class="size-5 text-success-500" />
                                    <span class="text-sm font-medium text-success-700">
                                        {{ $t('common.open', 2) }}
                                    </span>
                                </div>

                                <div
                                    v-if="quali.request?.status"
                                    class="flex flex-initial flex-row gap-1 rounded-full bg-info-100 px-2 py-1"
                                >
                                    <MailIcon class="size-5 text-info-400" />
                                    <span class="text-sm font-medium text-info-700">
                                        <span c lass="font-semibold">{{ $t('common.request') }}:</span>
                                        {{
                                            $t(
                                                `enums.qualifications.RequestStatus.${RequestStatus[quali.request?.status ?? 0]}`,
                                            )
                                        }}
                                    </span>
                                </div>

                                <div
                                    v-if="quali.result?.status"
                                    class="flex flex-initial flex-row gap-1 rounded-full bg-info-100 px-2 py-1"
                                >
                                    <ListStatusIcon class="size-5 text-info-400" />
                                    <span class="text-sm font-medium text-info-700">
                                        <span class="font-semibold">{{ $t('common.result') }}:</span>
                                        {{ $t(`enums.qualifications.ResultStatus.${ResultStatus[quali.result?.status ?? 0]}`) }}
                                    </span>
                                </div>
                            </div>

                            <div class="flex snap-x flex-row flex-wrap gap-2 overflow-x-auto pb-3 sm:pb-0">
                                <div class="flex flex-initial flex-row gap-1 rounded-full bg-base-100 px-2 py-1 text-base-500">
                                    <AccountIcon class="h-auto w-5" />
                                    <span class="inline-flex items-center text-sm font-medium text-base-700">
                                        {{ $t('common.created_by') }}
                                        <CitizenInfoPopover
                                            :user="quali?.creator"
                                            class="ml-1 font-medium text-primary-600 hover:text-primary-400"
                                        />
                                    </span>
                                </div>

                                <div class="flex flex-initial flex-row gap-1 rounded-full bg-base-100 px-2 py-1 text-base-500">
                                    <CalendarIcon class="h-auto w-5" />
                                    <span class="text-sm font-medium text-base-700">
                                        {{ $t('common.created_at') }}
                                        <GenericTime :value="quali?.createdAt" type="long" />
                                    </span>
                                </div>
                                <div
                                    v-if="quali?.updatedAt"
                                    class="flex flex-initial flex-row gap-1 rounded-full bg-base-100 px-2 py-1 text-base-500"
                                >
                                    <CalendarEditIcon class="h-auto w-5" />
                                    <span class="text-sm font-medium text-base-700">
                                        {{ $t('common.updated_at') }}
                                        <GenericTime :value="quali?.updatedAt" type="long" />
                                    </span>
                                </div>
                                <div
                                    v-if="quali?.deletedAt"
                                    class="flex flex-initial flex-row gap-1 rounded-full bg-base-100 px-2 py-1 text-base-500"
                                >
                                    <CalendarRemoveIcon class="h-auto w-5" />
                                    <span class="text-sm font-medium text-base-700">
                                        {{ $t('common.deleted') }}
                                        <GenericTime :value="quali?.deletedAt" type="long" />
                                    </span>
                                </div>
                            </div>

                            <div class="mt-2 w-full">
                                <h3 class="inline-flex items-center text-base font-semibold leading-7">
                                    {{ $t('common.requirements', 2) }}:
                                </h3>

                                <div class="flex flex-row flex-wrap gap-1 pb-2">
                                    <template v-if="!quali.requirements || quali.requirements.length === 0">
                                        <p class="text-base">
                                            {{ $t('common.not_found', [$t('common.requirements', 2)]) }}
                                        </p>
                                    </template>

                                    <template v-else>
                                        <NuxtLink
                                            v-for="entry in quali.requirements"
                                            :key="entry.id"
                                            :to="{
                                                name: 'jobs-qualifications-id',
                                                params: { id: entry.targetQualificationId },
                                            }"
                                            class="flex flex-initial snap-x snap-start items-center gap-1 overflow-x-auto whitespace-nowrap rounded-full px-2 py-1"
                                            :class="
                                                entry.targetQualification?.result?.status === ResultStatus.SUCCESSFUL
                                                    ? 'bg-success-100'
                                                    : 'bg-error-100'
                                            "
                                        >
                                            <span
                                                class="size-2 rounded-full"
                                                :class="
                                                    entry.targetQualification?.result?.status === ResultStatus.SUCCESSFUL
                                                        ? 'bg-success-500'
                                                        : 'bg-error-500'
                                                "
                                            />
                                            <span class="text-sm font-medium text-info-800"
                                                >{{ entry.targetQualification?.abbreviation }}:
                                                {{ entry.targetQualification?.title }}
                                            </span>
                                        </NuxtLink>
                                    </template>
                                </div>
                            </div>

                            <div v-if="!!quali?.content" class="my-2">
                                <h2 class="sr-only">
                                    {{ $t('common.content') }}
                                </h2>
                                <div class="break-words rounded-lg bg-base-800">
                                    <!-- eslint-disable vue/no-v-html -->
                                    <div
                                        class="prose prose-invert min-w-full rounded-md bg-base-900 p-4"
                                        v-html="quali?.content"
                                    ></div>
                                </div>
                            </div>

                            <div v-if="quali.result && quali.result.id !== '0'" class="mt-2 w-full">
                                <Disclosure v-slot="{ open }" as="div" class="w-full border-neutral/20 hover:border-neutral/70">
                                    <DisclosureButton
                                        :class="[
                                            open ? 'rounded-t-lg border-b-0' : 'rounded-lg',
                                            'flex w-full items-start justify-between border-2 border-inherit p-2 text-left transition-colors',
                                        ]"
                                    >
                                        <span class="inline-flex items-center text-base font-semibold leading-7">
                                            <ListStatusIcon class="mr-2 h-auto w-5" />
                                            {{ $t('common.result') }}
                                        </span>
                                        <span class="ml-6 flex h-7 items-center">
                                            <ChevronDownIcon
                                                :class="[open ? 'upsidedown' : '', 'h-auto w-5 transition-transform']"
                                            />
                                        </span>
                                    </DisclosureButton>
                                    <DisclosurePanel class="rounded-b-lg border-2 border-t-0 border-inherit transition-colors">
                                        <div class="mx-4 flex flex-col gap-1 pb-2">
                                            <div>
                                                <span class="font-semibold">{{ $t('common.result') }}:</span>
                                                {{
                                                    $t(
                                                        `enums.qualifications.ResultStatus.${ResultStatus[quali.result?.status ?? 0]}`,
                                                    )
                                                }}
                                            </div>
                                            <div>
                                                <span class="font-semibold">{{ $t('common.summary') }}:</span>
                                                {{ quali.result?.summary }}
                                            </div>
                                            <div class="inline-flex gap-1">
                                                <span class="font-semibold">{{ $t('common.created_by') }}:</span>
                                                <CitizenInfoPopover :user="quali.result?.creator" />
                                            </div>
                                        </div>
                                    </DisclosurePanel>
                                </Disclosure>
                            </div>

                            <div class="mt-2 w-full">
                                <Disclosure v-slot="{ open }" as="div" class="w-full border-neutral/20 hover:border-neutral/70">
                                    <DisclosureButton
                                        :class="[
                                            open ? 'rounded-t-lg border-b-0' : 'rounded-lg',
                                            'flex w-full items-start justify-between border-2 border-inherit p-2 text-left transition-colors',
                                        ]"
                                    >
                                        <span class="inline-flex items-center text-base font-semibold leading-7">
                                            <LockIcon class="mr-2 h-auto w-5" />
                                            {{ $t('common.access') }}
                                        </span>
                                        <span class="ml-6 flex h-7 items-center">
                                            <ChevronDownIcon
                                                :class="[open ? 'upsidedown' : '', 'h-auto w-5 transition-transform']"
                                            />
                                        </span>
                                    </DisclosureButton>
                                    <DisclosurePanel class="rounded-b-lg border-2 border-t-0 border-inherit transition-colors">
                                        <div class="mx-4 flex flex-row flex-wrap gap-1 pb-2">
                                            <DataNoDataBlock
                                                v-if="!quali.access || quali.access?.jobs.length === 0"
                                                icon="i-mdi-file-search"
                                                :message="$t('common.not_found', [$t('common.access', 2)])"
                                            />

                                            <template v-else>
                                                <div
                                                    v-for="entry in quali.access?.jobs"
                                                    :key="entry.id"
                                                    class="flex flex-initial snap-x snap-start items-center gap-1 overflow-x-auto whitespace-nowrap rounded-full bg-info-100 px-2 py-1"
                                                >
                                                    <span class="size-2 rounded-full bg-info-500" />
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
                                                            $t(`enums.qualifications.AccessLevel.${AccessLevel[entry.access]}`)
                                                        }}
                                                    </span>
                                                </div>
                                            </template>
                                        </div>
                                    </DisclosurePanel>
                                </Disclosure>
                            </div>

                            <div v-if="canDo.grade" class="mt-2 w-full">
                                <Disclosure
                                    v-slot="{ open }"
                                    as="div"
                                    class="w-full border-neutral/20 hover:border-neutral/70"
                                    :default-open="true"
                                >
                                    <DisclosureButton
                                        :class="[
                                            open ? 'rounded-t-lg border-b-0' : 'rounded-lg',
                                            'flex w-full items-start justify-between border-2 border-inherit p-2 text-left transition-colors',
                                        ]"
                                    >
                                        <span class="inline-flex items-center text-base font-semibold leading-7">
                                            <AccountSchoolIcon class="mr-2 h-auto w-5" />
                                            {{ $t('common.request', 2) }}
                                        </span>
                                        <span class="ml-6 flex h-7 items-center">
                                            <ChevronDownIcon
                                                :class="[open ? 'upsidedown' : '', 'h-auto w-5 transition-transform']"
                                            />
                                        </span>
                                    </DisclosureButton>
                                    <DisclosurePanel class="rounded-b-lg border-2 border-t-0 border-inherit transition-colors">
                                        <div class="mx-4 pb-2">
                                            <QualificationsRequestsList :qualification-id="quali.id" />
                                        </div>
                                    </DisclosurePanel>
                                </Disclosure>
                            </div>

                            <div v-if="canDo.grade" class="mt-2 w-full">
                                <Disclosure v-slot="{ open }" as="div" class="w-full border-neutral/20 hover:border-neutral/70">
                                    <DisclosureButton
                                        :class="[
                                            open ? 'rounded-t-lg border-b-0' : 'rounded-lg',
                                            'flex w-full items-start justify-between border-2 border-inherit p-2 text-left transition-colors',
                                        ]"
                                    >
                                        <span class="inline-flex items-center text-base font-semibold leading-7">
                                            <SigmaIcon class="mr-2 h-auto w-5" />
                                            {{ $t('common.result', 2) }}
                                        </span>
                                        <span class="ml-6 flex h-7 items-center">
                                            <ChevronDownIcon
                                                :class="[open ? 'upsidedown' : '', 'h-auto w-5 transition-transform']"
                                            />
                                        </span>
                                    </DisclosureButton>
                                    <DisclosurePanel class="rounded-b-lg border-2 border-t-0 border-inherit transition-colors">
                                        <div class="mx-4 pb-2">
                                            <QualificationsResultsList :qualification-id="quali.id" />
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
