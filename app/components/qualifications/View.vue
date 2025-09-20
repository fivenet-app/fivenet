<script lang="ts" setup>
import type { TabsItem } from '@nuxt/ui';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import {
    checkQualificationAccess,
    requestStatusToBadgeColor,
    requestStatusToTextColor,
    requirementsFullfilled,
    resultStatusToBadgeColor,
    resultStatusToTextColor,
} from '~/components/qualifications/helpers';
import RequestUserModal from '~/components/qualifications/request/RequestUserModal.vue';
import { getQualificationsQualificationsClient } from '~~/gen/ts/clients';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import { AccessLevel } from '~~/gen/ts/resources/qualifications/access';
import { QualificationExamMode, RequestStatus, ResultStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import type { DeleteQualificationResponse, GetQualificationResponse } from '~~/gen/ts/services/qualifications/qualifications';
import AccessBadges from '../partials/access/AccessBadges.vue';
import CitizenInfoPopover from '../partials/citizens/CitizenInfoPopover.vue';
import HTMLContent from '../partials/content/HTMLContent.vue';
import DataErrorBlock from '../partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '../partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '../partials/data/DataPendingBlock.vue';
import GenericTime from '../partials/elements/GenericTime.vue';
import IDCopyBadge from '../partials/IDCopyBadge.vue';
import OpenClosedBadge from '../partials/OpenClosedBadge.vue';
import TutorView from './tutor/TutorView.vue';

const props = defineProps<{
    qualificationId: number;
}>();

const { t } = useI18n();

const { can } = useAuth();

const overlay = useOverlay();

const notifications = useNotificationsStore();

const qualificationsQualificationsClient = await getQualificationsQualificationsClient();

const { data, status, refresh, error } = useLazyAsyncData(`qualification-${props.qualificationId}`, () =>
    getQualification(props.qualificationId),
);

async function getQualification(qualificationId: number): Promise<GetQualificationResponse> {
    try {
        const call = qualificationsQualificationsClient.getQualification({
            qualificationId: qualificationId,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

useHead({
    title: () =>
        data.value?.qualification?.title
            ? `${data.value.qualification.abbreviation}: ${data.value.qualification.title} - ${t('pages.qualifications.id.title')}`
            : t('pages.qualifications.id.title'),
});

async function deleteQualification(qualificationId: number): Promise<DeleteQualificationResponse> {
    try {
        const call = qualificationsQualificationsClient.deleteQualification({
            qualificationId: qualificationId,
        });
        const { response } = await call;

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
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
    request: checkQualificationAccess(
        qualification.value?.access,
        qualification.value?.creator,
        AccessLevel.REQUEST,
        undefined,
        qualification.value?.creatorJob,
    ),
    take: checkQualificationAccess(
        qualification.value?.access,
        qualification.value?.creator,
        AccessLevel.TAKE,
        undefined,
        qualification.value?.creatorJob,
    ),
    grade: checkQualificationAccess(
        qualification.value?.access,
        qualification.value?.creator,
        AccessLevel.GRADE,
        undefined,
        qualification.value?.creatorJob,
    ),
    edit: checkQualificationAccess(
        qualification.value?.access,
        qualification.value?.creator,
        AccessLevel.EDIT,
        undefined,
        qualification.value?.creatorJob,
    ),
    delete: checkQualificationAccess(
        qualification.value?.access,
        qualification.value?.creator,
        AccessLevel.EDIT,
        undefined,
        qualification.value?.creatorJob,
    ),
}));

watchOnce(data, async () => {
    if (!data.value?.qualification?.request) return;

    if (data.value?.qualification?.request.status === RequestStatus.EXAM_STARTED) {
        await navigateTo({
            name: 'qualifications-id-exam',
            params: { id: props.qualificationId },
        });
    }
});

const items = computed<TabsItem[]>(() =>
    [
        { slot: 'info' as const, label: t('common.content'), icon: 'i-mdi-info', value: 'info' },
        canDo.value.grade
            ? { slot: 'tutor' as const, label: t('common.tutor'), icon: 'i-mdi-sigma', value: 'tutor' }
            : undefined,
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);

const route = useRoute();
const router = useRouter();

const selectedTab = computed({
    get() {
        return (route.query.tab as string) || 'info';
    },
    set(tab) {
        // Hash is specified here to prevent the page from scrolling to the top
        router.push({ query: { tab: tab }, hash: '#control-active-item' });
    },
});

const accordionItems = computed(() =>
    [
        qualification.value?.result
            ? { slot: 'result' as const, label: t('common.result', 1), icon: 'i-mdi-list-status', defaultOpen: true }
            : qualification.value?.request
              ? { slot: 'request' as const, label: t('common.request'), icon: 'i-mdi-mail', defaultOpen: true }
              : undefined,
        { slot: 'access' as const, label: t('common.access'), icon: 'i-mdi-lock', defaultOpen: true },
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);

const confirmModal = overlay.create(ConfirmModal);
const requestUserModal = overlay.create(RequestUserModal);
</script>

<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0 overflow-y-hidden' }">
        <template #header>
            <UDashboardNavbar :title="$t('pages.qualifications.id.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <PartialsBackButton to="/qualifications" />

                    <UButton
                        icon="i-mdi-refresh"
                        :label="$t('common.refresh')"
                        :loading="isRequestPending(status)"
                        :ui="{ label: 'hidden sm:inline-flex' }"
                        @click="() => refresh()"
                    />

                    <UButtonGroup class="inline-flex">
                        <IDCopyBadge
                            :id="qualification?.id ?? qualificationId ?? 0"
                            prefix="QUAL"
                            size="md"
                            :title="{ key: 'notifications.qualifications.copy_qualification.title', parameters: {} }"
                            :content="{
                                key: 'notifications.qualifications.copy_qualification.content',
                                parameters: {},
                            }"
                        />
                    </UButtonGroup>
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar
                v-if="
                    qualification &&
                    (canDo.edit ||
                        (qualification.result !== undefined && qualification.result?.status !== ResultStatus.SUCCESSFUL) ||
                        canDo.request)
                "
            >
                <template #default>
                    <div
                        class="mx-auto flex w-full max-w-(--breakpoint-xl) flex-1 snap-x flex-row flex-wrap justify-between gap-2 overflow-x-auto"
                    >
                        <template v-if="!canDo.edit">
                            <UTooltip
                                v-if="canDo.request && qualification.examMode !== QualificationExamMode.ENABLED"
                                class="flex-1"
                                :text="$t('common.request')"
                            >
                                <UButton
                                    :disabled="
                                        qualification.closed ||
                                        !requirementsFullfilled(qualification.requirements) ||
                                        qualification.request?.status === RequestStatus.PENDING ||
                                        qualification.request?.status === RequestStatus.ACCEPTED ||
                                        qualification.request?.status === RequestStatus.EXAM_STARTED ||
                                        qualification.request?.status === RequestStatus.EXAM_GRADING
                                    "
                                    block
                                    icon="i-mdi-account-school"
                                    :label="$t('common.request')"
                                    variant="ghost"
                                    @click="
                                        requestUserModal.open({
                                            qualificationId: qualification!.id,
                                            onUpdatedRequest: ($event) => (qualification!.request = $event),
                                        })
                                    "
                                />
                            </UTooltip>

                            <UTooltip
                                v-if="
                                    canDo.take &&
                                    qualification.examMode !== QualificationExamMode.DISABLED &&
                                    qualification.result?.status !== ResultStatus.SUCCESSFUL
                                "
                                class="flex-1"
                                :text="$t('components.qualifications.take_test')"
                            >
                                <UButton
                                    block
                                    icon="i-mdi-test-tube"
                                    :disabled="
                                        qualification.closed ||
                                        !requirementsFullfilled(qualification.requirements) ||
                                        qualification.request?.status === RequestStatus.EXAM_GRADING ||
                                        (qualification.examMode === QualificationExamMode.REQUEST_NEEDED &&
                                            qualification.request?.status !== RequestStatus.ACCEPTED)
                                    "
                                    :to="
                                        qualification.closed ||
                                        !requirementsFullfilled(qualification.requirements) ||
                                        qualification.request?.status === RequestStatus.EXAM_GRADING ||
                                        (qualification.examMode === QualificationExamMode.REQUEST_NEEDED &&
                                            qualification.request?.status !== RequestStatus.ACCEPTED)
                                            ? undefined
                                            : { name: 'qualifications-id-exam', params: { id: qualification.id } }
                                    "
                                    :label="$t('components.qualifications.take_test')"
                                    variant="ghost"
                                />
                            </UTooltip>
                        </template>

                        <UTooltip
                            v-if="can('qualifications.QualificationsService/UpdateQualification').value && canDo.edit"
                            class="flex-1"
                            :text="$t('common.edit')"
                        >
                            <UButton
                                block
                                :to="{
                                    name: 'qualifications-id-edit',
                                    params: { id: qualification.id },
                                }"
                                color="neutral"
                                icon="i-mdi-pencil"
                                :label="$t('common.edit')"
                                variant="ghost"
                            />
                        </UTooltip>

                        <UTooltip
                            v-if="can('qualifications.QualificationsService/DeleteQualification').value && canDo.delete"
                            class="flex-1"
                            :text="!qualification.deletedAt ? $t('common.delete') : $t('common.restore')"
                        >
                            <UButton
                                block
                                :color="!qualification.deletedAt ? 'error' : 'success'"
                                :icon="!qualification.deletedAt ? 'i-mdi-delete' : 'i-mdi-restore'"
                                :label="!qualification.deletedAt ? $t('common.delete') : $t('common.restore')"
                                variant="ghost"
                                @click="
                                    confirmModal.open({
                                        confirm: async () => deleteQualification(qualification!.id),
                                    })
                                "
                            />
                        </UTooltip>
                    </div>
                </template>
            </UDashboardToolbar>

            <UDashboardToolbar v-if="qualification" class="print:hidden">
                <div class="mx-auto my-2 w-full max-w-(--breakpoint-xl)">
                    <div class="mb-4">
                        <h1
                            class="px-0.5 py-1 text-4xl font-bold break-words sm:pl-1"
                            :class="!qualification.title ? 'italic' : ''"
                        >
                            <template v-if="qualification.abbreviation">{{ qualification.abbreviation }}: </template>
                            {{ !qualification.title ? $t('common.untitled') : qualification.title }}
                        </h1>

                        <p v-if="qualification.description" class="px-0.5 py-1 text-base font-bold break-words sm:pl-1">
                            {{ qualification.description }}
                        </p>
                    </div>

                    <div class="mb-2 flex gap-2">
                        <OpenClosedBadge :closed="qualification.closed" />

                        <UBadge
                            v-if="qualification.draft"
                            class="inline-flex gap-1"
                            icon="i-mdi-pencil"
                            color="info"
                            size="md"
                            :label="$t('common.draft')"
                        />

                        <UBadge
                            v-if="qualification.public"
                            class="inline-flex gap-1"
                            icon="i-mdi-earth"
                            color="neutral"
                            size="md"
                            :label="$t('common.public')"
                        />

                        <UBadge class="inline-flex gap-1" icon="i-mdi-test-tube" size="md">
                            {{ $t('common.exam', 1) }}:
                            {{
                                $t(
                                    `enums.qualifications.QualificationExamMode.${QualificationExamMode[qualification.examMode]}`,
                                )
                            }}
                        </UBadge>

                        <UBadge
                            v-if="qualification.result?.status"
                            class="inline-flex gap-1"
                            icon="i-mdi-list-status"
                            :color="resultStatusToBadgeColor(qualification.result?.status ?? 0)"
                            size="md"
                        >
                            {{ $t('common.result') }}:
                            {{ $t(`enums.qualifications.ResultStatus.${ResultStatus[qualification.result?.status ?? 0]}`) }}
                        </UBadge>
                        <UBadge
                            v-else-if="qualification.request?.status"
                            class="inline-flex gap-1"
                            icon="i-mdi-mail"
                            :color="requestStatusToBadgeColor(qualification.request?.status ?? 0)"
                            size="md"
                        >
                            {{ $t('common.request') }}:
                            {{ $t(`enums.qualifications.RequestStatus.${RequestStatus[qualification.request?.status ?? 0]}`) }}
                        </UBadge>
                    </div>

                    <div class="flex snap-x flex-row flex-wrap gap-2 overflow-x-auto pb-3 sm:pb-0">
                        <UBadge class="inline-flex gap-1" icon="i-mdi-account" color="neutral" size="md">
                            <span class="inline-flex items-center gap-1">
                                <span class="font-medium">{{ $t('common.created_by') }}</span>
                                <CitizenInfoPopover :user="qualification?.creator" text-class="text-xs" />
                            </span>
                        </UBadge>

                        <UBadge class="inline-flex gap-1" icon="i-mdi-calendar" color="neutral" size="md">
                            {{ $t('common.created_at') }}
                            <GenericTime :value="qualification?.createdAt" type="long" />
                        </UBadge>

                        <UBadge
                            v-if="qualification?.updatedAt"
                            class="inline-flex gap-1"
                            icon="i-mdi-calendar-edit"
                            color="neutral"
                            size="md"
                        >
                            {{ $t('common.updated_at') }}
                            <GenericTime :value="qualification?.updatedAt" type="long" />
                        </UBadge>

                        <UBadge
                            v-if="qualification?.deletedAt"
                            class="inline-flex gap-1"
                            icon="i-mdi-calendar-remove"
                            color="warning"
                            size="md"
                        >
                            {{ $t('common.deleted') }}
                            <GenericTime :value="qualification?.deletedAt" type="long" />
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
                                            ? 'success'
                                            : 'error'
                                    "
                                >
                                    <span>
                                        {{ requirement.targetQualification?.abbreviation }}:
                                        {{
                                            !requirement.targetQualification?.title
                                                ? $t('common.untitled')
                                                : requirement.targetQualification?.title
                                        }}
                                    </span>
                                </UBadge>
                            </ULink>
                        </div>
                    </div>
                </div>
            </UDashboardToolbar>
        </template>

        <template #body>
            <DataPendingBlock
                v-if="isRequestPending(status)"
                :message="$t('common.loading', [$t('common.qualifications', 1)])"
            />
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

            <div v-else class="flex min-h-full w-full max-w-full flex-1 flex-col overflow-y-hidden">
                <div v-if="!canDo.grade && qualification.content?.content === undefined" class="p-4 sm:p-4">
                    <UAlert
                        icon="i-mdi-information"
                        color="info"
                        variant="subtle"
                        :description="$t('components.qualifications.content_unavailable')"
                    />
                </div>

                <UTabs
                    v-else
                    v-model="selectedTab"
                    class="flex-1 flex-col overflow-y-hidden"
                    :items="items"
                    variant="link"
                    :ui="{
                        content: 'flex flex-col flex-1 min-h-0 max-h-full overflow-y-hidden',
                        list: 'mx-auto max-w-(--breakpoint-xl)',
                    }"
                >
                    <template #info>
                        <UDashboardPanel :ui="{ root: 'h-full min-h-0 overflow-y-auto' }">
                            <template #body>
                                <div class="mx-auto w-full max-w-(--breakpoint-xl)">
                                    <div
                                        v-if="qualification.content?.content"
                                        class="w-full rounded-lg bg-neutral-100 p-4 break-words dark:bg-neutral-800"
                                    >
                                        <HTMLContent :value="qualification.content.content" />
                                    </div>
                                    <UAlert
                                        v-else
                                        icon="i-mdi-info"
                                        color="error"
                                        variant="subtle"
                                        :description="$t('components.qualifications.content_unavailable')"
                                    />

                                    <UDashboardToolbar v-if="qualification" :ui="{ root: 'border-b-0' }">
                                        <UAccordion :items="accordionItems" type="multiple" class="p-2 sm:p-2">
                                            <template v-if="qualification.result" #result>
                                                <UContainer class="mb-2">
                                                    <div class="flex flex-col gap-1">
                                                        <div class="inline-flex gap-1">
                                                            <span class="font-semibold">{{ $t('common.result') }}:</span>
                                                            <span
                                                                :class="
                                                                    resultStatusToTextColor(qualification.result?.status ?? 0)
                                                                "
                                                            >
                                                                {{
                                                                    $t(
                                                                        `enums.qualifications.ResultStatus.${ResultStatus[qualification.result?.status ?? 0]}`,
                                                                    )
                                                                }}
                                                            </span>
                                                        </div>

                                                        <div class="truncate">
                                                            <span class="font-semibold">{{ $t('common.summary') }}:</span>
                                                            {{ qualification.result?.summary }}
                                                        </div>

                                                        <div>
                                                            <span class="font-semibold">{{ $t('common.score') }}:</span>
                                                            {{ $t('common.point', qualification.result?.score ?? 0) }}
                                                        </div>

                                                        <div v-if="qualification.result?.creator" class="inline-flex gap-1">
                                                            <span class="font-semibold">{{ $t('common.created_by') }}:</span>
                                                            <CitizenInfoPopover :user="qualification.result?.creator" />
                                                        </div>
                                                    </div>
                                                </UContainer>
                                            </template>

                                            <template v-if="qualification.request" #request>
                                                <UContainer class="mb-2">
                                                    <div class="flex flex-col gap-1">
                                                        <div class="inline-flex gap-1">
                                                            <span class="font-semibold">{{ $t('common.status') }}:</span>
                                                            <span
                                                                :class="
                                                                    requestStatusToTextColor(qualification.request?.status ?? 0)
                                                                "
                                                            >
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
                                                <UContainer class="mb-2">
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
                                    </UDashboardToolbar>
                                </div>
                            </template>
                        </UDashboardPanel>
                    </template>

                    <template v-if="canDo.grade" #tutor>
                        <div class="w-full overflow-y-auto">
                            <div class="mx-auto w-full max-w-(--breakpoint-xl)">
                                <TutorView :qualification="qualification" />
                            </div>
                        </div>
                    </template>
                </UTabs>
            </div>
        </template>
    </UDashboardPanel>
</template>
