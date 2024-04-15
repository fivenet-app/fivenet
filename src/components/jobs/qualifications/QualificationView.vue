<script lang="ts" setup>
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { AccessLevel, RequestStatus, ResultStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import type { DeleteQualificationResponse, GetQualificationResponse } from '~~/gen/ts/services/qualifications/qualifications';
import { checkQualificationAccess } from '~/components/jobs/qualifications/helpers';
import QualificationRequestUserModal from '~/components/jobs/qualifications/QualificationRequestUserModal.vue';
import QualificationsRequestsList from '~/components/jobs/qualifications/tutor/QualificationsRequestsList.vue';
import QualificationsResultsList from '~/components/jobs/qualifications/tutor/QualificationsResultsList.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import type { AccordionItem } from '#ui/types';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';

const props = defineProps<{
    qualificationId: string;
}>();

const { t } = useI18n();

const { $grpc } = useNuxtApp();

const modal = useModal();

const { data, pending, refresh, error } = useLazyAsyncData(`qualification-${props.qualificationId}`, () =>
    getQualification(props.qualificationId),
);

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

const qualification = computed(() => data.value?.qualification);
const canDo = computed(() => ({
    take: checkQualificationAccess(qualification.value?.access, qualification.value?.creator, AccessLevel.TAKE),
    request: checkQualificationAccess(qualification.value?.access, qualification.value?.creator, AccessLevel.REQUEST),
    grade: checkQualificationAccess(qualification.value?.access, qualification.value?.creator, AccessLevel.GRADE),
    edit: checkQualificationAccess(qualification.value?.access, qualification.value?.creator, AccessLevel.EDIT),
}));

const accordionItems = computed(() => {
    return [
        qualification.value?.result && parseInt(qualification.value?.result.id) > 0
            ? { slot: 'result', label: t('common.result', 1), icon: 'i-mdi-list-status', defaultOpen: true }
            : undefined,
        { slot: 'access', label: t('common.access'), icon: 'i-mdi-lock', defaultOpen: true },
        { slot: 'requests', label: t('common.request', 2), icon: 'i-mdi-account-school' },
        { slot: 'results', label: t('common.result', 2), icon: 'i-mdi-sigma' },
    ].filter((item) => item !== undefined) as AccordionItem[];
});
</script>

<template>
    <UDashboardNavbar :title="$t('pages.qualifications.single.title')">
        <template #right>
            <UButtonGroup class="inline-flex">
                <UButton color="black" icon="i-mdi-arrow-back" to="/jobs/qualifications">
                    {{ $t('common.back') }}
                </UButton>

                <IDCopyBadge
                    :id="qualification?.id ?? 0"
                    prefix="QUAL"
                    :title="{ key: 'notifications.quali?.ment_view.copy_quali?.ment_id.title', parameters: {} }"
                    :content="{
                        key: 'notifications.quali?.ment_view.copy_quali?.ment_id.content',
                        parameters: {},
                    }"
                />
            </UButtonGroup>
        </template>
    </UDashboardNavbar>

    <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.qualifications', 1)])" />
    <DataErrorBlock v-else-if="error" :title="$t('common.unable_to_load', [$t('common.qualifications', 1)])" :retry="refresh" />
    <DataNoDataBlock
        v-else-if="!qualification"
        icon="i-mdi-account-school"
        :message="$t('common.not_found', [$t('common.qualification', 1)])"
    />

    <template v-else>
        <UDashboardToolbar>
            <template #default>
                <div class="flex flex-1 snap-x flex-row flex-wrap justify-between gap-2 overflow-x-auto">
                    <template v-if="!canDo.edit">
                        <UButton
                            v-if="canDo.take"
                            icon="i-mdi-test-tube"
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
                            :disabled="qualification.request?.status === RequestStatus.PENDING"
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
                            name: 'jobs-qualifications-id-edit',
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
                    <UBadge v-if="qualification.closed" color="red" class="inline-flex gap-1" size="md">
                        <UIcon name="i-mdi-lock" color="red" class="h-auto w-5" />
                        <span>
                            {{ $t('common.close', 2) }}
                        </span>
                    </UBadge>
                    <UBadge v-else color="green" class="inline-flex gap-1" size="md">
                        <UIcon name="i-mdi-lock-open-variant" color="green" class="h-auto w-5" />
                        <span>
                            {{ $t('common.open', 2) }}
                        </span>
                    </UBadge>

                    <UBadge
                        v-if="qualification.request?.status"
                        class="flex flex-initial flex-row gap-1 rounded-full bg-info-100 px-2 py-1"
                    >
                        <UIcon name="i-mdi-mail" class="size-5" />
                        <span>
                            <span class="font-medium">{{ $t('common.request') }}:</span>
                            {{ $t(`enums.qualifications.RequestStatus.${RequestStatus[qualification.request?.status ?? 0]}`) }}
                        </span>
                    </UBadge>

                    <UBadge
                        v-if="qualification.result?.status"
                        class="flex flex-initial flex-row gap-1 rounded-full bg-info-100 px-2 py-1"
                    >
                        <UIcon name="i-mdi-list-status" class="size-5" />
                        <span>
                            <span class="font-medium">{{ $t('common.result') }}:</span>
                            {{ $t(`enums.qualifications.ResultStatus.${ResultStatus[qualification.result?.status ?? 0]}`) }}
                        </span>
                    </UBadge>
                </div>

                <div class="flex snap-x flex-row flex-wrap gap-2 overflow-x-auto pb-3 sm:pb-0">
                    <UBadge color="black" class="inline-flex gap-1" size="md">
                        <UIcon name="i-mdi-account" class="h-auto w-5" />
                        <span class="inline-flex items-center gap-1">
                            <span class="text-sm font-medium">{{ $t('common.created_by') }}</span>
                            <CitizenInfoPopover :user="qualification?.creator" />
                        </span>
                    </UBadge>

                    <UBadge color="black" class="inline-flex gap-1" size="md">
                        <UIcon name="i-mdi-calendar" class="h-auto w-5" />

                        <span>
                            {{ $t('common.created_at') }}
                            <GenericTime :value="qualification?.createdAt" type="long" />
                        </span>
                    </UBadge>

                    <UBadge v-if="qualification?.updatedAt" color="black" class="inline-flex gap-1" size="md">
                        <UIcon name="i-mdi-calendar-edit" class="h-auto w-5" />
                        <span>
                            {{ $t('common.updated_at') }}
                            <GenericTime :value="qualification?.updatedAt" type="long" />
                        </span>
                    </UBadge>

                    <UBadge v-if="qualification?.deletedAt" color="amber" class="inline-flex gap-1" size="md">
                        <UIcon name="i-mdi-calendar-remove" class="h-auto w-5" />
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
                            <NuxtLink
                                v-for="entry in qualification.requirements"
                                :key="entry.id"
                                :to="{
                                    name: 'jobs-qualifications-id',
                                    params: { id: entry.targetQualificationId },
                                }"
                            >
                                <UBadge
                                    :color="
                                        entry.targetQualification?.result?.status === ResultStatus.SUCCESSFUL ? 'green' : 'red'
                                    "
                                >
                                    <span>
                                        {{ entry.targetQualification?.abbreviation }}:
                                        {{ entry.targetQualification?.title }}
                                    </span>
                                </UBadge>
                            </NuxtLink>
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
                    <div ref="contentRef" class="prose prose-invert min-w-full px-4 py-2" v-html="qualification.content"></div>
                </div>
            </div>

            <template #footer>
                <UAccordion :items="accordionItems" multiple :unmount="true">
                    <template v-if="qualification.result && qualification.result.id !== '0'" #result>
                        <UContainer>
                            <div class="flex flex-col gap-1">
                                <div>
                                    <span class="font-semibold">{{ $t('common.result') }}:</span>
                                    {{
                                        $t(
                                            `enums.qualifications.ResultStatus.${ResultStatus[qualification.result?.status ?? 0]}`,
                                        )
                                    }}
                                </div>
                                <div>
                                    <span class="font-semibold">{{ $t('common.summary') }}:</span>
                                    {{ qualification.result?.summary }}
                                </div>
                                <div class="inline-flex gap-1">
                                    <span class="font-semibold">{{ $t('common.created_by') }}:</span>
                                    <CitizenInfoPopover :user="qualification.result?.creator" />
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

                    <template v-if="canDo.grade" #requests>
                        <UContainer>
                            <QualificationsRequestsList :qualification-id="qualification.id" />
                        </UContainer>
                    </template>

                    <template v-if="canDo.grade" #results>
                        <UContainer>
                            <QualificationsResultsList :qualification-id="qualification.id" />
                        </UContainer>
                    </template>
                </UAccordion>
            </template>
        </UCard>
    </template>
</template>
