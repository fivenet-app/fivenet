import type { RoutePathSchema } from '@typed-router';
import type { AppFeature } from '~/composables/useAppFeatures';
import { isRoute } from '~/utils/route';
import { getDocumentsDocumentsClient } from '~~/gen/ts/clients';
import type { Perms } from '~~/gen/ts/perms';

export const HELP_GUIDE_ACTION_EVENT = 'fivenet:help-guide-action';

export type HelpGuideActionPhase = 'enter' | 'leave' | 'trigger';

export type HelpGuideActionDetail = {
    guideId: string;
    stepId: string;
    actionId: string;
    phase: HelpGuideActionPhase;
};

export type HelpGuideStep = {
    id: string;
    title: string;
    body: string;
    target?: string;
    side?: 'top' | 'right' | 'bottom' | 'left';
    /** CSS selector for the element the instructional popover should anchor to. */
    popoverTarget?: string;
    /** Position the instructional popover independently from the highlighted target. */
    popoverPosition?: 'center-bottom';
    /** Scroll this step's target into view when it becomes active. */
    moveToTarget?: boolean;
    /** Dim the rest of the page while this step is active. Defaults to true. */
    showOverlay?: boolean;
    /** Prevent interaction with the highlighted target while this step is active. */
    disableInteraction?: boolean;
    /** Resume at this step when a reactive guide branch replaces the current step. */
    resumeToStepId?: string;
    permission?: Perms | Perms[];
    /** Action dispatched automatically as this step becomes active. */
    onEnterAction?: string;
    /** Advance after an automatic enter action has been dispatched. */
    advanceAfterEnterAction?: boolean;
    /** Action dispatched before leaving this step or ending the guide. */
    onLeaveAction?: string;
    action?: {
        id: string;
        label: string;
        advanceAfterAction?: boolean;
    };
};

export type HelpGuide = {
    id: string;
    featureId: AppFeature['id'];
    title: string;
    description: string;
    icon: string;
    route: RoutePathSchema;
    permission?: Perms | Perms[];
    steps: HelpGuideStep[];
};

const activeGuideId = () => useState<string | undefined>('help-active-guide', () => undefined);

export const useHelpGuides = () => {
    const { t } = useI18n();
    const { can } = useAuth();
    const { isHelpSlideoverOpen } = useDashboard();
    const route = useRoute();
    const router = useRouter();
    const mailerStore = useMailerStore();
    const { loaded: mailerLoaded, hasPrivateEmail, threads } = storeToRefs(mailerStore);
    const canListDocuments = can('documents.DocumentsService/ListDocuments' as Perms);
    const approvalGuideEligible = useState<boolean | undefined>('help-documents-approval-eligible', () => undefined);
    const approvalGuideDocumentId = useState<number | undefined>('help-documents-approval-document-id', () => undefined);
    const approvalGuideCheckStarted = useState('help-documents-approval-check-started', () => false);
    const activeGuide = activeGuideId();

    async function checkApprovalGuideEligibility(): Promise<void> {
        if (!canListDocuments.value || approvalGuideCheckStarted.value) return;

        approvalGuideCheckStarted.value = true;
        try {
            const documentsClient = await getDocumentsDocumentsClient();
            const { response } = await documentsClient.listDocuments({
                pagination: { offset: 0, pageSize: 1 },
                categoryIds: [],
                creatorIds: [],
                documentIds: [],
                onlyDrafts: false,
            });
            const document = response.documents[0];
            approvalGuideDocumentId.value = document?.id;
            approvalGuideEligible.value = document !== undefined;
        } catch {
            approvalGuideDocumentId.value = undefined;
            approvalGuideEligible.value = false;
        }
    }

    if (import.meta.client) {
        watch(canListDocuments, () => void checkApprovalGuideEligibility(), { immediate: true });
    }

    const guides = computed<HelpGuide[]>(() => {
        const available: HelpGuide[] = [];
        const addFeatureGuide = (guide: HelpGuide): void => {
            if (guide.permission === undefined || can(guide.permission, 'all').value) available.push(guide);
        };

        available.push({
            id: 'overview-quick-access',
            featureId: 'personal-productivity',
            title: t('help_guides.overview_quick_access.title'),
            description: t('help_guides.overview_quick_access.description'),
            icon: 'i-mdi-pin-outline',
            route: '/overview',
            steps: [
                {
                    id: 'intro',
                    title: t('help_guides.overview_quick_access.steps.intro.title'),
                    body: t('help_guides.overview_quick_access.steps.intro.body'),
                    target: '[data-tour="overview-feature-cards"]',
                    side: 'top',
                    showOverlay: false,
                    disableInteraction: true,
                },
                {
                    id: 'feature-card',
                    title: t('help_guides.overview_quick_access.steps.feature_card.title'),
                    body: t('help_guides.overview_quick_access.steps.feature_card.body'),
                    target: '[data-tour="overview-feature-card"]',
                    side: 'right',
                    showOverlay: false,
                },
                {
                    id: 'quick-access',
                    title: t('help_guides.overview_quick_access.steps.quick_access.title'),
                    body: t('help_guides.overview_quick_access.steps.quick_access.body'),
                },
                {
                    id: 'finish',
                    title: t('help_guides.overview_quick_access.steps.finish.title'),
                    body: t('help_guides.overview_quick_access.steps.finish.body'),
                    showOverlay: false,
                },
            ],
        });

        addFeatureGuide({
            id: 'clipboard-basics',
            featureId: 'personal-productivity',
            title: t('help_guides.clipboard_basics.title'),
            description: t('help_guides.clipboard_basics.description'),
            icon: 'i-mdi-clipboard-list-outline',
            route: '/overview',
            permission: 'documents.DocumentsService/UpdateDocument' as Perms,
            steps: [
                {
                    id: 'intro',
                    title: t('help_guides.clipboard_basics.steps.intro.title'),
                    body: t('help_guides.clipboard_basics.steps.intro.body'),
                    showOverlay: false,
                },
                {
                    id: 'open',
                    title: t('help_guides.clipboard_basics.steps.open.title'),
                    body: t('help_guides.clipboard_basics.steps.open.body'),
                    target: '[data-tour="clipboard-open"]',
                    side: 'bottom',
                    popoverPosition: 'center-bottom',
                    moveToTarget: true,
                },
                {
                    id: 'modal',
                    title: t('help_guides.clipboard_basics.steps.modal.title'),
                    body: t('help_guides.clipboard_basics.steps.modal.body'),
                    target: '[data-tour="clipboard-modal"]',
                    side: 'bottom',
                    popoverPosition: 'center-bottom',
                    showOverlay: false,
                    onEnterAction: 'clipboard-open-modal',
                    onLeaveAction: 'clipboard-close-modal',
                },
                {
                    id: 'categories',
                    title: t('help_guides.clipboard_basics.steps.categories.title'),
                    body: t('help_guides.clipboard_basics.steps.categories.body'),
                    target: '[data-tour="clipboard-tabs"]',
                    side: 'bottom',
                    showOverlay: false,
                },
                {
                    id: 'remove',
                    title: t('help_guides.clipboard_basics.steps.remove.title'),
                    body: t('help_guides.clipboard_basics.steps.remove.body'),
                    target: '[data-tour="clipboard-clear"]',
                    side: 'top',
                    showOverlay: false,
                },
                {
                    id: 'templates',
                    title: t('help_guides.clipboard_basics.steps.templates.title'),
                    body: t('help_guides.clipboard_basics.steps.templates.body'),
                    showOverlay: false,
                },
                {
                    id: 'finish',
                    title: t('help_guides.clipboard_basics.steps.finish.title'),
                    body: t('help_guides.clipboard_basics.steps.finish.body'),
                    showOverlay: false,
                },
            ],
        });

        addFeatureGuide({
            id: 'documents-workflow',
            featureId: 'documents',
            title: t('help_guides.documents_workflow.title'),
            description: t('help_guides.documents_workflow.description'),
            icon: 'i-mdi-file-document-multiple-outline',
            route: '/documents/',
            permission: 'documents.DocumentsService/ListDocuments' as Perms,
            steps: [
                {
                    id: 'intro',
                    title: t('help_guides.documents_workflow.steps.intro.title'),
                    body: t('help_guides.documents_workflow.steps.intro.body'),
                    showOverlay: false,
                },
                {
                    id: 'search',
                    title: t('help_guides.documents_workflow.steps.search.title'),
                    body: t('help_guides.documents_workflow.steps.search.body'),
                    target: '[data-tour="documents-search"]',
                    side: 'bottom',
                    moveToTarget: true,
                },
                {
                    id: 'filters',
                    title: t('help_guides.documents_workflow.steps.filters.title'),
                    body: t('help_guides.documents_workflow.steps.filters.body'),
                    target: '[data-tour="documents-filters"]',
                    side: 'bottom',
                },
                {
                    id: 'pinned',
                    title: t('help_guides.documents_workflow.steps.pinned.title'),
                    body: t('help_guides.documents_workflow.steps.pinned.body'),
                    target: '[data-tour="documents-pinned"]',
                    side: 'left',
                    popoverPosition: 'center-bottom',
                    showOverlay: false,
                    permission: 'documents.DocumentsService/ToggleDocumentPin' as Perms,
                    onEnterAction: 'documents-open-pinned',
                    onLeaveAction: 'documents-close-pinned',
                },
                {
                    id: 'actions',
                    title: t('help_guides.documents_workflow.steps.actions.title'),
                    body: t('help_guides.documents_workflow.steps.actions.body'),
                    target: '[data-tour="documents-list"]',
                    side: 'top',
                },
                {
                    id: 'finish',
                    title: t('help_guides.documents_workflow.steps.finish.title'),
                    body: t('help_guides.documents_workflow.steps.finish.body'),
                    showOverlay: false,
                },
            ],
        });

        if (approvalGuideEligible.value === true) {
            addFeatureGuide({
                id: 'documents-approval',
                featureId: 'documents',
                title: t('help_guides.documents_approval.title'),
                description: t('help_guides.documents_approval.description'),
                icon: 'i-mdi-file-check-outline',
                route: '/documents/',
                permission: 'documents.DocumentsService/ListDocuments' as Perms,
                steps: [
                    {
                        id: 'intro',
                        title: t('help_guides.documents_approval.steps.intro.title'),
                        body: t('help_guides.documents_approval.steps.intro.body'),
                        showOverlay: false,
                    },
                    {
                        id: 'actions',
                        title: t('help_guides.documents_approval.steps.actions.title'),
                        body: t('help_guides.documents_approval.steps.actions.body'),
                        target: '[data-tour="documents-action-toolbar"]',
                        side: 'bottom',
                        action: {
                            id: 'documents-open-approval-drawer',
                            label: t('help_guides.documents_approval.steps.actions.action'),
                            advanceAfterAction: true,
                        },
                        onLeaveAction: 'documents-close-approval-drawer',
                    },
                    {
                        id: 'approval-panel',
                        title: t('help_guides.documents_approval.steps.approval_panel.title'),
                        body: t('help_guides.documents_approval.steps.approval_panel.body'),
                        target: '[data-tour="documents-approval-drawer"]',
                        side: 'left',
                        popoverPosition: 'center-bottom',
                        showOverlay: false,
                        onEnterAction: 'documents-open-approval-drawer',
                    },
                    {
                        id: 'decision-buttons',
                        title: t('help_guides.documents_approval.steps.decision_buttons.title'),
                        body: t('help_guides.documents_approval.steps.decision_buttons.body'),
                        target: '[data-tour="documents-approval-decisions"]',
                        side: 'top',
                        showOverlay: false,
                        onLeaveAction: 'documents-close-approval-drawer',
                    },
                    {
                        id: 'finish',
                        title: t('help_guides.documents_approval.steps.finish.title'),
                        body: t('help_guides.documents_approval.steps.finish.body'),
                        showOverlay: false,
                    },
                ],
            });
        }

        if (
            can(
                ['documents.DocumentsService/ListDocuments' as Perms, 'documents.DocumentsService/UpdateDocument' as Perms],
                'all',
            ).value
        ) {
            available.push({
                id: 'documents-create',
                featureId: 'documents',
                title: t('help_guides.documents_create.title'),
                description: t('help_guides.documents_create.description'),
                icon: 'i-mdi-file-document-edit-outline',
                route: '/documents/',
                permission: [
                    'documents.DocumentsService/ListDocuments' as Perms,
                    'documents.DocumentsService/UpdateDocument' as Perms,
                ],
                steps: [
                    {
                        id: 'intro',
                        title: t('help_guides.documents_create.steps.intro.title'),
                        body: t('help_guides.documents_create.steps.intro.body'),
                        showOverlay: false,
                    },
                    {
                        id: 'open-chooser',
                        title: t('help_guides.documents_create.steps.open_chooser.title'),
                        body: t('help_guides.documents_create.steps.open_chooser.body'),
                        target: '[data-tour="documents-create"]',
                        side: 'bottom',
                        moveToTarget: true,
                    },
                    {
                        id: 'blank',
                        title: t('help_guides.documents_create.steps.blank.title'),
                        body: t('help_guides.documents_create.steps.blank.body'),
                        target: '[data-tour="documents-create-blank"]',
                        side: 'top',
                        disableInteraction: true,
                        onEnterAction: 'documents-open-template-chooser',
                        onLeaveAction: 'documents-close-template-chooser',
                    },
                    {
                        id: 'template',
                        title: t('help_guides.documents_create.steps.template.title'),
                        body: t('help_guides.documents_create.steps.template.body'),
                        target: '[data-tour="documents-create-template"]',
                        side: 'top',
                        disableInteraction: true,
                        permission: 'documents.TemplatesService/ListTemplates' as Perms,
                        onEnterAction: 'documents-open-template-chooser',
                        onLeaveAction: 'documents-close-template-chooser',
                    },
                    {
                        id: 'template-requirements',
                        title: t('help_guides.documents_create.steps.template_requirements.title'),
                        body: t('help_guides.documents_create.steps.template_requirements.body'),
                        target: '[data-tour="documents-create-template"]',
                        side: 'top',
                        disableInteraction: true,
                        permission: 'documents.TemplatesService/ListTemplates' as Perms,
                        onEnterAction: 'documents-open-template-chooser-first-template',
                        onLeaveAction: 'documents-close-template-chooser',
                    },
                    {
                        id: 'finish',
                        title: t('help_guides.documents_create.steps.finish.title'),
                        body: t('help_guides.documents_create.steps.finish.body'),
                    },
                ],
            });
        }

        if (can('livemap.LivemapService/Stream' as Perms).value) {
            available.push({
                id: 'livemap-basics',
                featureId: 'livemap',
                title: t('help_guides.livemap_basics.title'),
                description: t('help_guides.livemap_basics.description'),
                icon: 'i-mdi-map-outline',
                route: '/livemap',
                permission: 'livemap.LivemapService/Stream' as Perms,
                steps: [
                    {
                        id: 'intro',
                        title: t('help_guides.livemap_basics.steps.intro.title'),
                        body: t('help_guides.livemap_basics.steps.intro.body'),
                        showOverlay: false,
                    },
                    {
                        id: 'zoom',
                        title: t('help_guides.livemap_basics.steps.zoom.title'),
                        body: t('help_guides.livemap_basics.steps.zoom.body'),
                        target: '[data-tour="livemap-zoom"]',
                        side: 'right',
                    },
                    {
                        id: 'layers',
                        title: t('help_guides.livemap_basics.steps.layers.title'),
                        body: t('help_guides.livemap_basics.steps.layers.body'),
                        target: '[data-tour="livemap-layers"]',
                        side: 'left',
                        onEnterAction: 'livemap-open-layers',
                        onLeaveAction: 'livemap-close-layers',
                    },
                    {
                        id: 'postal-search',
                        title: t('help_guides.livemap_basics.steps.postal_search.title'),
                        body: t('help_guides.livemap_basics.steps.postal_search.body'),
                        target: '[data-tour="livemap-postal-search"]',
                        side: 'top',
                    },
                    {
                        id: 'employee-search',
                        title: t('help_guides.livemap_basics.steps.employee_search.title'),
                        body: t('help_guides.livemap_basics.steps.employee_search.body'),
                        target: '[data-tour="livemap-employee-search"]',
                        side: 'top',
                    },
                    {
                        id: 'dispatch-search',
                        title: t('help_guides.livemap_basics.steps.dispatch_search.title'),
                        body: t('help_guides.livemap_basics.steps.dispatch_search.body'),
                        target: '[data-tour="livemap-dispatch-search"]',
                        side: 'top',
                        permission: 'centrum.CentrumService/Stream' as Perms,
                    },
                    {
                        id: 'settings',
                        title: t('help_guides.livemap_basics.steps.settings.title'),
                        body: t('help_guides.livemap_basics.steps.settings.body'),
                        target: '[data-tour="livemap-settings"]',
                        side: 'left',
                        onEnterAction: 'livemap-open-settings',
                        onLeaveAction: 'livemap-close-settings',
                    },
                    {
                        id: 'finish',
                        title: t('help_guides.livemap_basics.steps.finish.title'),
                        body: t('help_guides.livemap_basics.steps.finish.body'),
                        showOverlay: false,
                    },
                ],
            });
        }

        if (can(['livemap.LivemapService/Stream' as Perms, 'centrum.CentrumService/Stream' as Perms], 'all').value) {
            available.push({
                id: 'centrum-basics',
                featureId: 'livemap',
                title: t('help_guides.centrum_basics.title'),
                description: t('help_guides.centrum_basics.description'),
                icon: 'i-mdi-radio-tower',
                route: '/livemap',
                permission: ['livemap.LivemapService/Stream' as Perms, 'centrum.CentrumService/Stream' as Perms],
                steps: [
                    {
                        id: 'intro',
                        title: t('help_guides.centrum_basics.steps.intro.title'),
                        body: t('help_guides.centrum_basics.steps.intro.body'),
                        showOverlay: false,
                    },
                    {
                        id: 'membership-button',
                        title: t('help_guides.centrum_basics.steps.membership_button.title'),
                        body: t('help_guides.centrum_basics.steps.membership_button.body'),
                        target: '[data-tour="centrum-unit-membership"]',
                        side: 'right',
                    },
                    {
                        id: 'membership-dialog',
                        title: t('help_guides.centrum_basics.steps.membership_dialog.title'),
                        body: t('help_guides.centrum_basics.steps.membership_dialog.body'),
                        target: '[data-tour="centrum-demo-unit"]',
                        side: 'left',
                        showOverlay: false,
                        onEnterAction: 'centrum-open-unit-selector',
                        onLeaveAction: 'centrum-close-unit-selector',
                    },
                    {
                        id: 'unit-status',
                        title: t('help_guides.centrum_basics.steps.unit_status.title'),
                        body: t('help_guides.centrum_basics.steps.unit_status.body'),
                        target: '[data-tour="centrum-unit-status"]',
                        side: 'right',
                    },
                    {
                        id: 'take-dispatch',
                        title: t('help_guides.centrum_basics.steps.take_dispatch.title'),
                        body: t('help_guides.centrum_basics.steps.take_dispatch.body'),
                        target: '[data-tour="centrum-take-dispatch-footer"]',
                        side: 'top',
                        showOverlay: false,
                        onEnterAction: 'centrum-open-take-dispatches',
                        onLeaveAction: 'centrum-close-take-dispatches',
                    },
                    {
                        id: 'own-dispatches',
                        title: t('help_guides.centrum_basics.steps.own_dispatches.title'),
                        body: t('help_guides.centrum_basics.steps.own_dispatches.body'),
                        target: '[data-tour="centrum-own-dispatches"]',
                        side: 'left',
                    },
                    {
                        id: 'dispatch-status',
                        title: t('help_guides.centrum_basics.steps.dispatch_status.title'),
                        body: t('help_guides.centrum_basics.steps.dispatch_status.body'),
                        target: '[data-tour="centrum-dispatch-status"]',
                        side: 'right',
                        onEnterAction: 'centrum-select-demo-dispatch',
                        onLeaveAction: 'centrum-select-demo-dispatch',
                    },
                    {
                        id: 'finish',
                        title: t('help_guides.centrum_basics.steps.finish.title'),
                        body: t('help_guides.centrum_basics.steps.finish.body'),
                        showOverlay: false,
                    },
                ],
            });
        }

        available.push({
            id: 'calendar-basics',
            featureId: 'calendar',
            title: t('help_guides.calendar_basics.title'),
            description: t('help_guides.calendar_basics.description'),
            icon: 'i-mdi-calendar-outline',
            route: '/calendar',
            steps: [
                {
                    id: 'intro',
                    title: t('help_guides.calendar_basics.steps.intro.title'),
                    body: t('help_guides.calendar_basics.steps.intro.body'),
                    showOverlay: false,
                },
                {
                    id: 'navigation',
                    title: t('help_guides.calendar_basics.steps.navigation.title'),
                    body: t('help_guides.calendar_basics.steps.navigation.body'),
                    target: '[data-tour="calendar-navigation"]',
                    side: 'bottom',
                },
                {
                    id: 'date-picker',
                    title: t('help_guides.calendar_basics.steps.date_picker.title'),
                    body: t('help_guides.calendar_basics.steps.date_picker.body'),
                    target: '[data-tour="calendar-date-picker"]',
                    side: 'right',
                },
                {
                    id: 'view',
                    title: t('help_guides.calendar_basics.steps.view.title'),
                    body: t('help_guides.calendar_basics.steps.view.body'),
                    target: '[data-tour="calendar-view"]',
                    side: 'right',
                },
                {
                    id: 'calendar-list',
                    title: t('help_guides.calendar_basics.steps.calendar_list.title'),
                    body: t('help_guides.calendar_basics.steps.calendar_list.body'),
                    target: '[data-tour="calendar-list"]',
                    side: 'right',
                },
                {
                    id: 'public-calendars',
                    title: t('help_guides.calendar_basics.steps.public_calendars.title'),
                    body: t('help_guides.calendar_basics.steps.public_calendars.body'),
                    target: '[data-tour="calendar-selector"]',
                    side: 'right',
                    onEnterAction: 'calendar-open-public-calendars',
                    onLeaveAction: 'calendar-close-public-calendars',
                },
                {
                    id: 'reminders',
                    title: t('help_guides.calendar_basics.steps.reminders.title'),
                    body: t('help_guides.calendar_basics.steps.reminders.body'),
                },
                {
                    id: 'finish',
                    title: t('help_guides.calendar_basics.steps.finish.title'),
                    body: t('help_guides.calendar_basics.steps.finish.body'),
                    showOverlay: false,
                },
            ],
        });

        if (can('mailer.MailerService/ListEmails' as Perms).value) {
            const mailerNeedsEmailSetup = mailerLoaded.value && !hasPrivateEmail.value;
            const mailerHasThreads = (threads.value?.threads.length ?? 0) > 0;

            available.push({
                id: 'mailer-basics',
                featureId: 'mail',
                title: t('help_guides.mailer_basics.title'),
                description: t('help_guides.mailer_basics.description'),
                icon: 'i-mdi-email-outline',
                // The generated optional `:thread?` route expects the trailing slash.
                route: '/mail/' as RoutePathSchema,
                permission: 'mailer.MailerService/ListEmails' as Perms,
                steps: [
                    {
                        id: 'intro',
                        title: t('help_guides.mailer_basics.steps.intro.title'),
                        body: t('help_guides.mailer_basics.steps.intro.body'),
                        showOverlay: false,
                    },
                    ...(mailerNeedsEmailSetup
                        ? [
                              {
                                  id: 'email-name',
                                  title: t('help_guides.mailer_basics.steps.email_name.title'),
                                  body: t('help_guides.mailer_basics.steps.email_name.body'),
                                  target: '[data-tour="mailer-email-name"]',
                                  side: 'right' as const,
                                  moveToTarget: true,
                                  showOverlay: false,
                              },
                              {
                                  id: 'create-email',
                                  title: t('help_guides.mailer_basics.steps.create_email.title'),
                                  body: t('help_guides.mailer_basics.steps.create_email.body'),
                                  target: '[data-tour="mailer-create-email"]',
                                  side: 'top' as const,
                                  showOverlay: false,
                                  resumeToStepId: 'mailbox',
                              },
                          ]
                        : [
                              {
                                  id: 'mailbox',
                                  title: t('help_guides.mailer_basics.steps.mailbox.title'),
                                  body: t('help_guides.mailer_basics.steps.mailbox.body'),
                                  target: '[data-tour="mailer-mailbox"]',
                                  side: 'bottom' as const,
                              },
                              {
                                  id: 'filters',
                                  title: t('help_guides.mailer_basics.steps.filters.title'),
                                  body: t('help_guides.mailer_basics.steps.filters.body'),
                                  target: '[data-tour="mailer-filters"]',
                                  side: 'bottom' as const,
                              },
                              {
                                  id: 'search',
                                  title: t('help_guides.mailer_basics.steps.search.title'),
                                  body: t('help_guides.mailer_basics.steps.search.body'),
                                  target: '[data-tour="mailer-search"]',
                                  side: 'bottom' as const,
                              },
                              {
                                  id: 'compose',
                                  title: t('help_guides.mailer_basics.steps.compose.title'),
                                  body: t('help_guides.mailer_basics.steps.compose.body'),
                                  target: '[data-tour="mailer-compose"]',
                                  side: 'bottom' as const,
                              },
                              {
                                  id: 'templates',
                                  title: t('help_guides.mailer_basics.steps.templates.title'),
                                  body: t('help_guides.mailer_basics.steps.templates.body'),
                                  target: '[data-tour="mailer-templates"]',
                                  side: 'top' as const,
                              },
                              {
                                  id: 'threads',
                                  title: t('help_guides.mailer_basics.steps.threads.title'),
                                  body: t('help_guides.mailer_basics.steps.threads.body'),
                                  target: '[data-tour="mailer-thread-list"]',
                                  side: 'right' as const,
                                  action: mailerHasThreads
                                      ? {
                                            id: 'mailer-select-first-thread',
                                            label: t('help_guides.mailer_basics.steps.threads.action'),
                                            advanceAfterAction: true,
                                        }
                                      : undefined,
                              },
                              ...(mailerHasThreads
                                  ? [
                                        {
                                            id: 'thread-actions',
                                            title: t('help_guides.mailer_basics.steps.thread_actions.title'),
                                            body: t('help_guides.mailer_basics.steps.thread_actions.body'),
                                            target: '[data-tour="mailer-thread-actions"]',
                                            side: 'bottom' as const,
                                        },
                                    ]
                                  : []),
                              {
                                  id: 'finish',
                                  title: t('help_guides.mailer_basics.steps.finish.title'),
                                  body: t('help_guides.mailer_basics.steps.finish.body'),
                                  showOverlay: false,
                              },
                          ]),
                ],
            });
        }

        if (
            can(['mailer.MailerService/ListEmails' as Perms, 'mailer.MailerService/CreateOrUpdateEmail' as Perms], 'all').value
        ) {
            available.push({
                id: 'mailer-advanced',
                featureId: 'mail',
                title: t('help_guides.mailer_advanced.title'),
                description: t('help_guides.mailer_advanced.description'),
                icon: 'i-mdi-email-edit-outline',
                route: '/mail/manage',
                permission: ['mailer.MailerService/ListEmails' as Perms, 'mailer.MailerService/CreateOrUpdateEmail' as Perms],
                steps: [
                    {
                        id: 'intro',
                        title: t('help_guides.mailer_advanced.steps.intro.title'),
                        body: t('help_guides.mailer_advanced.steps.intro.body'),
                        showOverlay: false,
                    },
                    {
                        id: 'select-email',
                        title: t('help_guides.mailer_advanced.steps.select_email.title'),
                        body: t('help_guides.mailer_advanced.steps.select_email.body'),
                        target: '[data-tour="mailer-addresses"]',
                        side: 'right',
                        action: {
                            id: 'mailer-select-first-email',
                            label: t('help_guides.mailer_advanced.steps.select_email.action'),
                            advanceAfterAction: true,
                        },
                    },
                    {
                        id: 'addresses',
                        title: t('help_guides.mailer_advanced.steps.addresses.title'),
                        body: t('help_guides.mailer_advanced.steps.addresses.body'),
                        target: '[data-tour="mailer-addresses"]',
                        side: 'right',
                    },
                    {
                        id: 'access',
                        title: t('help_guides.mailer_advanced.steps.access.title'),
                        body: t('help_guides.mailer_advanced.steps.access.body'),
                        target: '[data-tour="mailer-access"]',
                        side: 'top',
                    },
                    {
                        id: 'blocked',
                        title: t('help_guides.mailer_advanced.steps.blocked.title'),
                        body: t('help_guides.mailer_advanced.steps.blocked.body'),
                        target: '[data-tour="mailer-blocked"]',
                        side: 'top',
                    },
                    {
                        id: 'address-book',
                        title: t('help_guides.mailer_advanced.steps.address_book.title'),
                        body: t('help_guides.mailer_advanced.steps.address_book.body'),
                        target: '[data-tour="mailer-address-book"]',
                        side: 'top',
                    },
                    {
                        id: 'signature',
                        title: t('help_guides.mailer_advanced.steps.signature.title'),
                        body: t('help_guides.mailer_advanced.steps.signature.body'),
                        target: '[data-tour="mailer-signature"]',
                        side: 'top',
                    },
                    {
                        id: 'finish',
                        title: t('help_guides.mailer_advanced.steps.finish.title'),
                        body: t('help_guides.mailer_advanced.steps.finish.body'),
                        showOverlay: false,
                    },
                ],
            });
        }

        addFeatureGuide({
            id: 'dispatch-workflow',
            featureId: 'dispatch',
            title: t('help_guides.dispatch_workflow.title'),
            description: t('help_guides.dispatch_workflow.description'),
            icon: 'i-mdi-radio-tower',
            route: '/dispatch',
            permission: 'centrum.CentrumService/TakeControl' as Perms,
            steps: [
                {
                    id: 'intro',
                    title: t('help_guides.dispatch_workflow.steps.intro.title'),
                    body: t('help_guides.dispatch_workflow.steps.intro.body'),
                    showOverlay: false,
                },
                {
                    id: 'layout',
                    title: t('help_guides.dispatch_workflow.steps.layout.title'),
                    body: t('help_guides.dispatch_workflow.steps.layout.body'),
                    target: '[data-tour="dispatch-center-layout"], [data-tour="dispatch-center-layout-popover"]',
                    popoverTarget: '[data-tour="dispatch-center-layout-popover"]',
                    side: 'bottom',
                    showOverlay: false,
                    onEnterAction: 'dispatch-center-open-layout',
                    onLeaveAction: 'dispatch-center-close-layout',
                },
                {
                    id: 'units',
                    title: t('help_guides.dispatch_workflow.steps.units.title'),
                    body: t('help_guides.dispatch_workflow.steps.units.body'),
                    target: '[data-tour="dispatch-units"]',
                    side: 'right',
                    moveToTarget: true,
                },
                {
                    id: 'dispatches',
                    title: t('help_guides.dispatch_workflow.steps.dispatches.title'),
                    body: t('help_guides.dispatch_workflow.steps.dispatches.body'),
                    target: '[data-tour="dispatch-list"]',
                    side: 'left',
                },
                {
                    id: 'dispatcher',
                    title: t('help_guides.dispatch_workflow.steps.dispatcher.title'),
                    body: t('help_guides.dispatch_workflow.steps.dispatcher.body'),
                    target: '[data-tour="dispatch-center-join"]',
                    side: 'bottom',
                    permission: 'centrum.CentrumService/TakeControl' as Perms,
                },
                {
                    id: 'finish',
                    title: t('help_guides.dispatch_workflow.steps.finish.title'),
                    body: t('help_guides.dispatch_workflow.steps.finish.body'),
                    showOverlay: false,
                },
            ],
        });

        addFeatureGuide({
            id: 'livemap-dispatch',
            featureId: 'livemap',
            title: t('help_guides.livemap_dispatch.title'),
            description: t('help_guides.livemap_dispatch.description'),
            icon: 'i-mdi-map-marker-radius-outline',
            route: '/livemap',
            permission: 'livemap.LivemapService/Stream' as Perms,
            steps: [
                {
                    id: 'intro',
                    title: t('help_guides.livemap_dispatch.steps.intro.title'),
                    body: t('help_guides.livemap_dispatch.steps.intro.body'),
                    showOverlay: false,
                },
                {
                    id: 'markers',
                    title: t('help_guides.livemap_dispatch.steps.markers.title'),
                    body: t('help_guides.livemap_dispatch.steps.markers.body'),
                    target: '[data-tour="livemap-dispatch-markers"]',
                    side: 'left',
                    moveToTarget: true,
                    permission: 'centrum.CentrumService/Stream' as Perms,
                },
                {
                    id: 'panel',
                    title: t('help_guides.livemap_dispatch.steps.panel.title'),
                    body: t('help_guides.livemap_dispatch.steps.panel.body'),
                    target: '[data-tour="livemap-side-panel"]',
                    side: 'left',
                },
                {
                    id: 'finish',
                    title: t('help_guides.livemap_dispatch.steps.finish.title'),
                    body: t('help_guides.livemap_dispatch.steps.finish.body'),
                    showOverlay: false,
                },
            ],
        });

        addFeatureGuide({
            id: 'citizen-lookup',
            featureId: 'citizens',
            title: t('help_guides.citizen_lookup.title'),
            description: t('help_guides.citizen_lookup.description'),
            icon: 'i-mdi-account-search-outline',
            route: '/citizens',
            permission: 'citizens.CitizensService/ListCitizens' as Perms,
            steps: [
                {
                    id: 'intro',
                    title: t('help_guides.citizen_lookup.steps.intro.title'),
                    body: t('help_guides.citizen_lookup.steps.intro.body'),
                    showOverlay: false,
                },
                {
                    id: 'search',
                    title: t('help_guides.citizen_lookup.steps.search.title'),
                    body: t('help_guides.citizen_lookup.steps.search.body'),
                    target: '[data-tour="citizens-search"]',
                    side: 'bottom',
                    moveToTarget: true,
                },
                {
                    id: 'advanced',
                    title: t('help_guides.citizen_lookup.steps.advanced.title'),
                    body: t('help_guides.citizen_lookup.steps.advanced.body'),
                    target: '[data-tour="citizens-advanced-search"]',
                    side: 'bottom',
                },
                {
                    id: 'open-profile',
                    title: t('help_guides.citizen_lookup.steps.open_profile.title'),
                    body: t('help_guides.citizen_lookup.steps.open_profile.body'),
                    target: '[data-tour="citizens-list"]',
                    side: 'top',
                    onEnterAction: 'citizens-open-first-profile',
                    advanceAfterEnterAction: true,
                },
                {
                    id: 'profile',
                    title: t('help_guides.citizen_lookup.steps.profile.title'),
                    body: t('help_guides.citizen_lookup.steps.profile.body'),
                    target: '[data-tour="citizen-profile"]',
                    side: 'top',
                    popoverPosition: 'center-bottom',
                },
                {
                    id: 'profile-tabs',
                    title: t('help_guides.citizen_lookup.steps.profile_tabs.title'),
                    body: t('help_guides.citizen_lookup.steps.profile_tabs.body'),
                    target: '[data-tour="citizen-profile-tabs"]',
                    side: 'bottom',
                    popoverPosition: 'center-bottom',
                },
                {
                    id: 'profile-actions',
                    title: t('help_guides.citizen_lookup.steps.profile_actions.title'),
                    body: t('help_guides.citizen_lookup.steps.profile_actions.body'),
                    target: '[data-tour="citizen-profile-actions"]',
                    side: 'bottom',
                    popoverPosition: 'center-bottom',
                },
                {
                    id: 'finish',
                    title: t('help_guides.citizen_lookup.steps.finish.title'),
                    body: t('help_guides.citizen_lookup.steps.finish.body'),
                    showOverlay: false,
                },
            ],
        });

        addFeatureGuide({
            id: 'job-management',
            featureId: 'jobs',
            title: t('help_guides.job_management.title'),
            description: t('help_guides.job_management.description'),
            icon: 'i-mdi-briefcase-account-outline',
            route: '/jobs/overview',
            permission: 'jobs.ColleaguesService/ListColleagues' as Perms,
            steps: [
                {
                    id: 'intro',
                    title: t('help_guides.job_management.steps.intro.title'),
                    body: t('help_guides.job_management.steps.intro.body'),
                    showOverlay: false,
                },
                {
                    id: 'navigation',
                    title: t('help_guides.job_management.steps.navigation.title'),
                    body: t('help_guides.job_management.steps.navigation.body'),
                    target: '[data-tour="jobs-navigation"]',
                    side: 'bottom',
                    moveToTarget: true,
                },
                {
                    id: 'absence',
                    title: t('help_guides.job_management.steps.absence.title'),
                    body: t('help_guides.job_management.steps.absence.body'),
                    target: '[data-tour="jobs-self-service-absence"]',
                    side: 'top',
                    permission: 'jobs.ColleaguesService/SetColleagueProps' as Perms,
                },
                {
                    id: 'profile-picture',
                    title: t('help_guides.job_management.steps.profile_picture.title'),
                    body: t('help_guides.job_management.steps.profile_picture.body'),
                    target: '[data-tour="jobs-self-service-avatar"]',
                    side: 'top',
                },
                {
                    id: 'finish',
                    title: t('help_guides.job_management.steps.finish.title'),
                    body: t('help_guides.job_management.steps.finish.body'),
                    showOverlay: false,
                },
            ],
        });

        addFeatureGuide({
            id: 'calendar-workflow',
            featureId: 'calendar',
            title: t('help_guides.calendar_workflow.title'),
            description: t('help_guides.calendar_workflow.description'),
            icon: 'i-mdi-calendar-edit-outline',
            route: '/calendar',
            steps: [
                {
                    id: 'intro',
                    title: t('help_guides.calendar_workflow.steps.intro.title'),
                    body: t('help_guides.calendar_workflow.steps.intro.body'),
                    showOverlay: false,
                },
                {
                    id: 'create',
                    title: t('help_guides.calendar_workflow.steps.create.title'),
                    body: t('help_guides.calendar_workflow.steps.create.body'),
                    target: '[data-tour="calendar-create"]',
                    side: 'bottom',
                    moveToTarget: true,
                    action: {
                        id: 'calendar-open-create-entry',
                        label: t('help_guides.calendar_workflow.steps.create.action'),
                        advanceAfterAction: true,
                    },
                },
                {
                    id: 'invitees',
                    title: t('help_guides.calendar_workflow.steps.invitees.title'),
                    body: t('help_guides.calendar_workflow.steps.invitees.body'),
                    target: '[data-tour="calendar-entry-form"]',
                    side: 'left',
                },
                {
                    id: 'finish',
                    title: t('help_guides.calendar_workflow.steps.finish.title'),
                    body: t('help_guides.calendar_workflow.steps.finish.body'),
                    showOverlay: false,
                },
            ],
        });

        addFeatureGuide({
            id: 'qualifications-workflow',
            featureId: 'qualifications',
            title: t('help_guides.qualifications_workflow.title'),
            description: t('help_guides.qualifications_workflow.description'),
            icon: 'i-mdi-school-outline',
            route: '/qualifications',
            permission: 'qualifications.QualificationsService/ListQualifications' as Perms,
            steps: [
                {
                    id: 'intro',
                    title: t('help_guides.qualifications_workflow.steps.intro.title'),
                    body: t('help_guides.qualifications_workflow.steps.intro.body'),
                    showOverlay: false,
                },
                {
                    id: 'your-qualifications',
                    title: t('help_guides.qualifications_workflow.steps.your_qualifications.title'),
                    body: t('help_guides.qualifications_workflow.steps.your_qualifications.body'),
                    target: '[data-tour="qualifications-your-tab"]',
                    side: 'bottom',
                },
                {
                    id: 'all-qualifications',
                    title: t('help_guides.qualifications_workflow.steps.all_qualifications.title'),
                    body: t('help_guides.qualifications_workflow.steps.all_qualifications.body'),
                    target: '[data-tour="qualifications-all-tab"]',
                    side: 'bottom',
                    showOverlay: false,
                },
                {
                    id: 'find',
                    title: t('help_guides.qualifications_workflow.steps.find.title'),
                    body: t('help_guides.qualifications_workflow.steps.find.body'),
                    target: '[data-tour="qualifications-all-search"]',
                    side: 'bottom',
                    moveToTarget: true,
                },
                {
                    id: 'review',
                    title: t('help_guides.qualifications_workflow.steps.review.title'),
                    body: t('help_guides.qualifications_workflow.steps.review.body'),
                    target: '[data-tour="qualifications-list"]',
                    side: 'top',
                },
                {
                    id: 'finish',
                    title: t('help_guides.qualifications_workflow.steps.finish.title'),
                    body: t('help_guides.qualifications_workflow.steps.finish.body'),
                    showOverlay: false,
                },
            ],
        });

        addFeatureGuide({
            id: 'wiki-workflow',
            featureId: 'wiki',
            title: t('help_guides.wiki_workflow.title'),
            description: t('help_guides.wiki_workflow.description'),
            icon: 'i-mdi-book-open-page-variant-outline',
            route: '/wiki',
            permission: 'wiki.WikiService/ListPages' as Perms,
            steps: [
                {
                    id: 'intro',
                    title: t('help_guides.wiki_workflow.steps.intro.title'),
                    body: t('help_guides.wiki_workflow.steps.intro.body'),
                    showOverlay: false,
                },
                {
                    id: 'search',
                    title: t('help_guides.wiki_workflow.steps.search.title'),
                    body: t('help_guides.wiki_workflow.steps.search.body'),
                    target: '[data-tour="wiki-search"]',
                    side: 'bottom',
                    moveToTarget: true,
                },
                {
                    id: 'select-wiki',
                    title: t('help_guides.wiki_workflow.steps.select_wiki.title'),
                    body: t('help_guides.wiki_workflow.steps.select_wiki.body'),
                    target: '[data-tour="wiki-pages"]',
                    side: 'top',
                    popoverPosition: 'center-bottom',
                    onEnterAction: 'wiki-select-first-wiki',
                    advanceAfterEnterAction: true,
                },
                {
                    id: 'sidebar',
                    title: t('help_guides.wiki_workflow.steps.sidebar.title'),
                    body: t('help_guides.wiki_workflow.steps.sidebar.body'),
                    target: '[data-tour="wiki-sidebar"]',
                    side: 'right',
                    popoverPosition: 'center-bottom',
                },
                {
                    id: 'content',
                    title: t('help_guides.wiki_workflow.steps.content.title'),
                    body: t('help_guides.wiki_workflow.steps.content.body'),
                    target: '[data-tour="wiki-content"]',
                    side: 'left',
                    popoverPosition: 'center-bottom',
                },
                {
                    id: 'finish',
                    title: t('help_guides.wiki_workflow.steps.finish.title'),
                    body: t('help_guides.wiki_workflow.steps.finish.body'),
                    showOverlay: false,
                },
            ],
        });

        addFeatureGuide({
            id: 'vehicle-lookup',
            featureId: 'vehicles',
            title: t('help_guides.vehicle_lookup.title'),
            description: t('help_guides.vehicle_lookup.description'),
            icon: 'i-mdi-car-search-outline',
            route: '/vehicles',
            permission: 'vehicles.VehiclesService/ListVehicles' as Perms,
            steps: [
                {
                    id: 'intro',
                    title: t('help_guides.vehicle_lookup.steps.intro.title'),
                    body: t('help_guides.vehicle_lookup.steps.intro.body'),
                    showOverlay: false,
                },
                {
                    id: 'search',
                    title: t('help_guides.vehicle_lookup.steps.search.title'),
                    body: t('help_guides.vehicle_lookup.steps.search.body'),
                    target: '[data-tour="vehicles-search"]',
                    side: 'bottom',
                    moveToTarget: true,
                },
                {
                    id: 'profile',
                    title: t('help_guides.vehicle_lookup.steps.profile.title'),
                    body: t('help_guides.vehicle_lookup.steps.profile.body'),
                    target: '[data-tour="vehicles-list"]',
                    side: 'top',
                },
                {
                    id: 'finish',
                    title: t('help_guides.vehicle_lookup.steps.finish.title'),
                    body: t('help_guides.vehicle_lookup.steps.finish.body'),
                    showOverlay: false,
                },
            ],
        });

        return available
            .filter((guide) => guide.permission === undefined || can(guide.permission, 'all').value)
            .map((guide) => ({
                ...guide,
                steps: guide.steps.filter((step) => step.permission === undefined || can(step.permission).value),
            }));
    });

    const currentGuide = computed(() => guides.value.find((guide) => guide.id === activeGuide.value));

    watch(
        () => route.path,
        (path) => {
            if (activeGuide.value === undefined) return;

            const guide = currentGuide.value;
            if (guide === undefined || !isRoute(path, guide.route)) {
                finishGuide();
            }
        },
    );

    watch(currentGuide, (guide) => {
        if (activeGuide.value !== undefined && guide === undefined) {
            finishGuide();
        }
    });

    async function startGuide(id: string): Promise<void> {
        const guide = guides.value.find((item) => item.id === id);
        if (!guide) return;

        isHelpSlideoverOpen.value = false;
        activeGuide.value = undefined;
        if (id === 'documents-approval' && approvalGuideDocumentId.value !== undefined) {
            await navigateTo({
                name: 'documents-id',
                params: { id: approvalGuideDocumentId.value },
            });
        } else {
            await navigateTo(guide.route);
        }

        // Let route middleware and the destination page finish their final
        // updates before mounting the tour. Otherwise the popover can open
        // briefly and then be closed by the route watcher on the next update.
        await nextTick();
        await new Promise<void>((resolve) => requestAnimationFrame(() => resolve()));
        await nextTick();

        if (!isRoute(router.currentRoute.value.path, guide.route)) return;
        activeGuide.value = guide.id;
    }

    function finishGuide(): void {
        activeGuide.value = undefined;
    }

    return {
        guides,
        currentGuide,
        activeGuide,
        startGuide,
        finishGuide,
    };
};
