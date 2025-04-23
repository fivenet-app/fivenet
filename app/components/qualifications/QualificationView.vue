<script lang="ts" setup>
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import OpenClosedBadge from '~/components/partials/OpenClosedBadge.vue';
import AccessBadges from '~/components/partials/access/AccessBadges.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import HTMLContent from '~/components/partials/content/HTMLContent.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import QualificationRequestUserModal from '~/components/qualifications/QualificationRequestUserModal.vue';
import {
    checkQualificationAccess,
    requestStatusToBadgeColor,
    requestStatusToTextColor,
    requirementsFullfilled,
    resultStatusToBadgeColor,
    resultStatusToTextColor,
} from '~/components/qualifications/helpers';
import { useNotificatorStore } from '~/stores/notificator';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import { AccessLevel } from '~~/gen/ts/resources/qualifications/access';
import { QualificationExamMode, RequestStatus, ResultStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import type { DeleteQualificationResponse, GetQualificationResponse } from '~~/gen/ts/services/qualifications/qualifications';
import QualificationTutorView from './tutor/QualificationTutorView.vue';

const props = defineProps<{
    qualificationId: number;
}>();

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const { can } = useAuth();

const modal = useModal();

const notifications = useNotificatorStore();

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`qualification-${props.qualificationId}`, () => getQualification(props.qualificationId));

async function getQualification(qualificationId: number): Promise<GetQualificationResponse> {
    try {
        const call = $grpc.qualifications.qualifications.getQualification({
            qualificationId: qualificationId,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function deleteQualification(qualificationId: number): Promise<DeleteQualificationResponse> {
    try {
        const call = $grpc.qualifications.qualifications.deleteQualification({
            qualificationId: qualificationId,
        });
        const { response } = await call;

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        navigateTo({
            name: 'qualifications',
            query: {
                tab: 'tab=all',
            },
            hash: '#',
        });

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
    delete: checkQualificationAccess(qualification.value?.access, qualification.value?.creator, AccessLevel.EDIT),
}));

watchOnce(data, async () => {
    if (!data.value?.qualification?.request) {
        return;
    }

    if (data.value?.qualification?.request.status === RequestStatus.EXAM_STARTED) {
        await navigateTo({
            name: 'qualifications-id-exam',
            params: { id: props.qualificationId },
        });
    }
});

const items = computed(() =>
    [
        { slot: 'info', label: t('common.content'), icon: 'i-mdi-info' },
        canDo.value.grade ? { slot: 'tutor', label: t('common.tutor'), icon: 'i-mdi-sigma' } : undefined,
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);

const route = useRoute();
const router = useRouter();

const selectedTab = computed({
    get() {
        const index = items.value.findIndex((item) => item.slot === route.query.tab);
        if (index === -1) {
            return 0;
        }

        return index;
    },
    set(value) {
        // Hash is specified here to prevent the page from scrolling to the top
        router.replace({ query: { tab: items.value[value]?.slot }, hash: '#' });
    },
});

const accordionItems = computed(() =>
    [
        qualification.value?.result
            ? { slot: 'result', label: t('common.result', 1), icon: 'i-mdi-list-status', defaultOpen: true }
            : qualification.value?.request
              ? { slot: 'request', label: t('common.request'), icon: 'i-mdi-mail', defaultOpen: true }
              : undefined,
        { slot: 'access', label: t('common.access'), icon: 'i-mdi-lock', defaultOpen: true },
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);
</script>

<template>
    <UDashboardNavbar :title="$t('pages.qualifications.single.title')">
        <template #right>
            <PartialsBackButton to="/qualifications" />

            <UButtonGroup class="inline-flex">
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
    <DataErrorBlock
        v-else-if="error"
        :title="$t('common.unable_to_load', [$t('common.qualifications', 1)])"
        :error="error"
        :retry="refresh"
    />
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
                        <UButton
                            v-if="qualification.examMode !== QualificationExamMode.ENABLED"
                            :disabled="
                                qualification.closed ||
                                !requirementsFullfilled(qualification.requirements) ||
                                qualification.request?.status === RequestStatus.PENDING ||
                                qualification.request?.status === RequestStatus.ACCEPTED ||
                                qualification.request?.status === RequestStatus.EXAM_STARTED ||
                                qualification.request?.status === RequestStatus.EXAM_GRADING
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

                        <UButton
                            v-if="
                                canDo.take &&
                                qualification.examMode !== QualificationExamMode.DISABLED &&
                                qualification.result?.status !== ResultStatus.SUCCESSFUL
                            "
                            icon="i-mdi-test-tube"
                            :disabled="
                                qualification.closed ||
                                !requirementsFullfilled(qualification.requirements) ||
                                qualification.request?.status === RequestStatus.EXAM_GRADING ||
                                (qualification.examMode === QualificationExamMode.REQUEST_NEEDED &&
                                    qualification.request?.status !== RequestStatus.ACCEPTED)
                            "
                            :to="{ name: 'qualifications-id-exam', params: { id: qualification.id } }"
                        >
                            {{ $t('components.qualifications.take_test') }}
                        </UButton>
                    </template>

                    <UButton
                        v-if="can('QualificationsService.UpdateQualification').value && canDo.edit"
                        :to="{
                            name: 'qualifications-id-edit',
                            params: { id: qualification.id },
                        }"
                        icon="i-mdi-pencil"
                    >
                        {{ $t('common.edit') }}
                    </UButton>

                    <UButton
                        v-if="can('QualificationsService.DeleteQualification').value && canDo.delete"
                        :color="!qualification.deletedAt ? 'red' : 'green'"
                        :icon="!qualification.deletedAt ? 'i-mdi-trash-can' : 'i-mdi-restore'"
                        :label="!qualification.deletedAt ? $t('common.delete') : $t('common.restore')"
                        @click="
                            modal.open(ConfirmModal, {
                                confirm: async () => deleteQualification(qualification!.id),
                            })
                        "
                    />
                </div>
            </template>
        </UDashboardToolbar>

        <UDashboardPanelContent class="p-0 sm:pb-0">
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

                        <UBadge v-if="qualification.public" color="black" class="inline-flex gap-1" size="md">
                            <UIcon name="i-mdi-earth" class="size-5" />
                            <span>
                                {{ $t('common.public') }}
                            </span>
                        </UBadge>

                        <UBadge class="inline-flex gap-1" size="md">
                            <UIcon name="i-mdi-test-tube" class="size-5" />
                            <span>
                                {{ $t('common.exam', 1) }}:
                                {{
                                    $t(
                                        `enums.qualifications.QualificationExamMode.${QualificationExamMode[qualification.examMode]}`,
                                    )
                                }}
                            </span>
                        </UBadge>

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
                                {{
                                    $t(
                                        `enums.qualifications.RequestStatus.${RequestStatus[qualification.request?.status ?? 0]}`,
                                    )
                                }}
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

                    <div v-if="qualification.requirements && qualification.requirements.length > 0" class="mt-2 w-full">
                        <h3>{{ $t('common.requirements', 2) }}:</h3>

                        <div class="flex flex-row flex-wrap gap-1 pb-2">
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
                        </div>
                    </div>
                </template>

                <div>
                    <UAlert
                        v-if="!canDo.grade"
                        icon="i-mdi-info"
                        :description="$t('components.qualifications.content_unavailable')"
                    />

                    <UTabs v-else v-model="selectedTab" :items="items" class="w-full">
                        <template #info>
                            <h2 class="sr-only">
                                {{ $t('common.content') }}
                            </h2>

                            <div class="mx-auto w-full max-w-screen-xl !break-words rounded-lg bg-neutral-100 dark:bg-base-900">
                                <HTMLContent
                                    v-if="qualification.content?.content"
                                    class="px-4 py-2"
                                    :value="qualification.content.content"
                                />
                                <UAlert
                                    v-else
                                    icon="i-mdi-info"
                                    :description="$t('components.qualifications.content_unavailable')"
                                />
                            </div>
                        </template>

                        <template v-if="canDo.grade" #tutor>
                            <QualificationTutorView :qualification="qualification" />
                        </template>
                    </UTabs>
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
                                        <span class="font-semibold"
                                            >{{ $t('common.request') }} {{ $t('common.message') }}:</span
                                        >
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
                                <AccessBadges
                                    v-else
                                    :access-level="AccessLevel"
                                    :jobs="qualification?.access.jobs"
                                    i18n-key="enums.qualifications"
                                />
                            </UContainer>
                        </template>
                    </UAccordion>
                </template>
            </UCard>
        </UDashboardPanelContent>
    </template>
</template>
