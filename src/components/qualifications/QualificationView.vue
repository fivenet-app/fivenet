<script lang="ts" setup>
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { AccessLevel, RequestStatus, ResultStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import type { DeleteQualificationResponse, GetQualificationResponse } from '~~/gen/ts/services/qualifications/qualifications';
import {
    checkQualificationAccess,
    requestStatusToBadgeColor,
    requestStatusToTextColor,
    requirementsFullfilled,
    resultStatusToBadgeColor,
    resultStatusToTextColor,
} from '~/components/qualifications/helpers';
import QualificationRequestUserModal from '~/components/qualifications/QualificationRequestUserModal.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import QualificationTutorView from './tutor/QualificationTutorView.vue';
import OpenClosedBadge from '../partials/OpenClosedBadge.vue';

const props = defineProps<{
    qualificationId: string;
}>();

const { t } = useI18n();

const modal = useModal();

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`qualification-${props.qualificationId}`, () => getQualification(props.qualificationId));

async function getQualification(qualificationId: string): Promise<GetQualificationResponse> {
    try {
        const call = getGRPCQualificationsClient().getQualification({
            qualificationId,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function deleteQualification(qualificationId: string): Promise<DeleteQualificationResponse> {
    try {
        const call = getGRPCQualificationsClient().deleteQualification({
            qualificationId,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const qualification = computed(() => data.value?.qualification);
const canDo = computed(() => ({
    request: checkQualificationAccess(qualification.value?.access, qualification.value?.creator, AccessLevel.REQUEST),
    take: checkQualificationAccess(qualification.value?.access, qualification.value?.creator, AccessLevel.TAKE),
    grade: checkQualificationAccess(qualification.value?.access, qualification.value?.creator, AccessLevel.GRADE),
    edit: checkQualificationAccess(qualification.value?.access, qualification.value?.creator, AccessLevel.EDIT),
}));

const accordionItems = computed(() =>
    [
        qualification.value?.result
            ? { slot: 'result', label: t('common.result', 1), icon: 'i-mdi-list-status', defaultOpen: true }
            : qualification.value?.request
              ? { slot: 'request', label: t('common.request'), icon: 'i-mdi-mail', defaultOpen: true }
              : undefined,
        { slot: 'access', label: t('common.access'), icon: 'i-mdi-lock', defaultOpen: true },
        canDo.value.grade ? { slot: 'tutor', label: t('common.tutor'), icon: 'i-mdi-sigma' } : undefined,
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);
</script>

<template>
    <UDashboardNavbar :title="$t('pages.qualifications.single.title')">
        <template #right>
            <UButtonGroup class="inline-flex">
                <UButton color="black" icon="i-mdi-arrow-back" to="/qualifications">
                    {{ $t('common.back') }}
                </UButton>

                <IDCopyBadge
                    :id="qualification?.id ?? 0"
                    prefix="QUAL"
                    :title="{ key: 'notifications.qualifications.copy_qualification.title', parameters: {} }"
                    :content="{
                        key: 'notifications.qualifications.copy_qualification.content',
                        parameters: {},
                    }"
                />
            </UButtonGroup>
        </template>
    </UDashboardNavbar>

    <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.qualifications', 1)])" />
    <DataErrorBlock v-else-if="error" :title="$t('common.unable_to_load', [$t('common.qualifications', 1)])" :retry="refresh" />
    <DataNoDataBlock
        v-else-if="!qualification"
        icon="i-mdi-account-school"
        :message="$t('common.not_found', [$t('common.qualification', 1)])"
    />

    <template v-else>
        <UDashboardToolbar v-if="canDo.edit || qualification.result?.status !== ResultStatus.SUCCESSFUL">
            <template #default>
                <div class="flex flex-1 snap-x flex-row flex-wrap justify-between gap-2 overflow-x-auto">
                    <template v-if="!canDo.edit">
                        <!-- TODO Enable when the exam logic is ready -->
                        <UButton
                            v-if="false || canDo.take"
                            icon="i-mdi-test-tube"
                            :disabled="qualification.closed || !requirementsFullfilled(qualification.requirements)"
                            @click="
                                modal.open(QualificationRequestUserModal, {
                                    qualificationId: qualification!.id,
                                    onUpdatedRequest: ($event) => (qualification!.request = $event),
                                })
                            "
                        >
                            {{ $t('components.qualifications.take_test') }}
                        </UButton>

                        <UButton
                            v-else-if="canDo.request"
                            :disabled="
                                qualification.closed ||
                                !requirementsFullfilled(qualification.requirements) ||
                                qualification.request?.status === RequestStatus.PENDING ||
                                qualification.request?.status === RequestStatus.ACCEPTED
                            "
                            icon="i-mdi-account-school"
                            @click="
                                modal.open(QualificationRequestUserModal, {
                                    qualificationId: qualification!.id,
                                    onUpdatedRequest: ($event) => (qualification!.request = $event),
                                })
                            "
                        >
                            {{ $t('common.request') }}
                        </UButton>
                    </template>

                    <UButton
                        v-if="can('QualificationsService.UpdateQualification') && canDo.edit"
                        :to="{
                            name: 'qualifications-id-edit',
                            params: { id: qualification.id },
                        }"
                        icon="i-mdi-pencil"
                    >
                        {{ $t('common.edit') }}
                    </UButton>

                    <UButton
                        v-if="can('QualificationsService.DeleteQualification') && canDo.edit"
                        icon="i-mdi-trash-can"
                        @click="
                            modal.open(ConfirmModal, {
                                confirm: async () => deleteQualification(qualification!.id),
                            })
                        "
                    >
                        {{ $t('common.delete') }}
                    </UButton>
                </div>
            </template>
        </UDashboardToolbar>

        <UCard>
            <template #header>
                <div class="mb-4">
                    <h1 class="break-words px-0.5 py-1 text-4xl font-bold sm:pl-1">
                        {{ qualification.abbreviation }}: {{ qualification.title }}
                    </h1>

                    <p v-if="qualification.description" class="break-words px-0.5 py-1 text-base font-bold sm:pl-1">
                        {{ qualification.description }}
                    </p>
                </div>

                <div class="mb-2 flex gap-2">
                    <OpenClosedBadge :closed="qualification.closed" />

                    <UBadge
                        v-if="qualification.result?.status"
                        :color="resultStatusToBadgeColor(qualification.result?.status ?? 0)"
                        class="inline-flex gap-1"
                        size="md"
                    >
                        <UIcon name="i-mdi-list-status" class="size-5" />
                        <span>
                            {{ $t('common.result') }}:
                            {{ $t(`enums.qualifications.ResultStatus.${ResultStatus[qualification.result?.status ?? 0]}`) }}
                        </span>
                    </UBadge>
                    <UBadge
                        v-else-if="qualification.request?.status"
                        :color="requestStatusToBadgeColor(qualification.request?.status ?? 0)"
                        class="inline-flex gap-1"
                        size="md"
                    >
                        <UIcon name="i-mdi-mail" class="size-5" />
                        <span>
                            {{ $t('common.request') }}:
                            {{ $t(`enums.qualifications.RequestStatus.${RequestStatus[qualification.request?.status ?? 0]}`) }}
                        </span>
                    </UBadge>
                </div>

                <div class="flex snap-x flex-row flex-wrap gap-2 overflow-x-auto pb-3 sm:pb-0">
                    <UBadge color="black" class="inline-flex gap-1" size="md">
                        <UIcon name="i-mdi-account" class="size-5" />
                        <span class="inline-flex items-center gap-1">
                            <span class="text-sm font-medium">{{ $t('common.created_by') }}</span>
                            <CitizenInfoPopover :user="qualification?.creator" />
                        </span>
                    </UBadge>

                    <UBadge color="black" class="inline-flex gap-1" size="md">
                        <UIcon name="i-mdi-calendar" class="size-5" />

                        <span>
                            {{ $t('common.created_at') }}
                            <GenericTime :value="qualification?.createdAt" type="long" />
                        </span>
                    </UBadge>

                    <UBadge v-if="qualification?.updatedAt" color="black" class="inline-flex gap-1" size="md">
                        <UIcon name="i-mdi-calendar-edit" class="size-5" />
                        <span>
                            {{ $t('common.updated_at') }}
                            <GenericTime :value="qualification?.updatedAt" type="long" />
                        </span>
                    </UBadge>

                    <UBadge v-if="qualification?.deletedAt" color="amber" class="inline-flex gap-1" size="md">
                        <UIcon name="i-mdi-calendar-remove" class="size-5" />
                        <span>
                            {{ $t('common.deleted') }}
                            <GenericTime :value="qualification?.deletedAt" type="long" />
                        </span>
                    </UBadge>
                </div>

                <div class="mt-2 w-full">
                    <h3>{{ $t('common.requirements', 2) }}:</h3>

                    <div class="flex flex-row flex-wrap gap-1 pb-2">
                        <template v-if="!qualification.requirements || qualification.requirements.length === 0">
                            <p class="text-base">
                                {{ $t('common.not_found', [$t('common.requirements', 2)]) }}
                            </p>
                        </template>

                        <template v-else>
                            <ULink
                                v-for="requirement in qualification.requirements"
                                :key="requirement.id"
                                :to="{
                                    name: 'qualifications-id',
                                    params: { id: requirement.targetQualificationId },
                                }"
                            >
                                <UBadge
                                    :color="
                                        requirement.targetQualification?.result?.status === ResultStatus.SUCCESSFUL
                                            ? 'green'
                                            : 'red'
                                    "
                                >
                                    <span>
                                        {{ requirement.targetQualification?.abbreviation }}:
                                        {{ requirement.targetQualification?.title }}
                                    </span>
                                </UBadge>
                            </ULink>
                        </template>
                    </div>
                </div>
            </template>

            <div>
                <h2 class="sr-only">
                    {{ $t('common.content') }}
                </h2>

                <div class="mx-auto max-w-screen-xl break-words rounded-lg bg-base-900">
                    <!-- eslint-disable vue/no-v-html -->
                    <div
                        v-if="qualification.content"
                        v-html="qualification.content"
                        class="prose dark:prose-invert min-w-full px-4 py-2"
                    ></div>
                    <p v-else>
                        {{ $t('components.qualifications.content_unavailable') }}
                    </p>
                </div>
            </div>

            <template #footer>
                <UAccordion :items="accordionItems" multiple :unmount="true">
                    <template v-if="qualification.result" #result>
                        <UContainer>
                            <div class="flex flex-col gap-1">
                                <div class="inline-flex gap-1">
                                    <span class="font-semibold">{{ $t('common.result') }}:</span>
                                    <span :class="resultStatusToTextColor(qualification.result?.status ?? 0)">
                                        {{
                                            $t(
                                                `enums.qualifications.ResultStatus.${ResultStatus[qualification.result?.status ?? 0]}`,
                                            )
                                        }}
                                    </span>
                                </div>
                                <div>
                                    <span class="font-semibold">{{ $t('common.summary') }}:</span>
                                    {{ qualification.result?.summary }}
                                </div>
                                <div>
                                    <span class="font-semibold">{{ $t('common.score') }}:</span>
                                    {{ $t('common.point', qualification.result?.score ?? 0) }}
                                </div>
                                <div class="inline-flex gap-1">
                                    <span class="font-semibold">{{ $t('common.created_by') }}:</span>
                                    <CitizenInfoPopover :user="qualification.result?.creator" />
                                </div>
                            </div>
                        </UContainer>
                    </template>

                    <template v-if="qualification.request" #request>
                        <UContainer>
                            <div class="flex flex-col gap-1">
                                <div class="inline-flex gap-1">
                                    <span class="font-semibold">{{ $t('common.status') }}:</span>
                                    <span :class="requestStatusToTextColor(qualification.request?.status ?? 0)">
                                        {{
                                            $t(
                                                `enums.qualifications.RequestStatus.${RequestStatus[qualification.request?.status ?? 0]}`,
                                            )
                                        }}
                                    </span>
                                </div>
                                <div>
                                    <span class="font-semibold">{{ $t('common.request') }} {{ $t('common.message') }}:</span>
                                    {{ qualification.request?.userComment }}
                                </div>
                                <div>
                                    <span class="font-semibold">{{ $t('common.comment') }}:</span>
                                    {{ qualification.request?.approverComment ?? $t('common.na') }}
                                </div>
                                <div v-if="qualification.request.approvedAt" class="inline-flex gap-1">
                                    <span class="font-semibold">{{ $t('common.approved_at') }}:</span>
                                    <span class="inline-flex gap-1">
                                        <GenericTime :value="qualification.request?.approvedAt" />
                                    </span>
                                </div>
                                <div v-if="qualification.request.approver" class="inline-flex gap-1">
                                    <span class="font-semibold">{{ $t('common.approved_by') }}:</span>
                                    <CitizenInfoPopover :user="qualification.request?.approver" />
                                </div>
                            </div>
                        </UContainer>
                    </template>

                    <template #access>
                        <UContainer>
                            <DataNoDataBlock
                                v-if="!qualification.access || qualification.access?.jobs.length === 0"
                                icon="i-mdi-file-search"
                                :message="$t('common.not_found', [$t('common.access', 2)])"
                            />
                            <div v-else class="mx-4 flex flex-col gap-2">
                                <div class="flex flex-row flex-wrap gap-1">
                                    <UBadge
                                        v-for="entry in qualification.access?.jobs"
                                        :key="entry.id"
                                        color="black"
                                        class="inline-flex gap-1"
                                        size="md"
                                    >
                                        <span class="size-2 rounded-full bg-info-500" />
                                        <span>
                                            {{ entry.jobLabel
                                            }}<span
                                                v-if="entry.minimumGrade > 0"
                                                :title="`${entry.jobLabel} - ${$t('common.rank')} ${entry.minimumGrade}`"
                                            >
                                                ({{ entry.jobGradeLabel }})</span
                                            >
                                            -
                                            {{ $t(`enums.qualifications.AccessLevel.${AccessLevel[entry.access]}`) }}
                                        </span>
                                    </UBadge>
                                </div>
                            </div>
                        </UContainer>
                    </template>

                    <template v-if="canDo.grade" #tutor>
                        <QualificationTutorView :qualification="qualification" />
                    </template>
                </UAccordion>
            </template>
        </UCard>
    </template>
</template>
