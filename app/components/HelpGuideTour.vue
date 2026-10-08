<script setup lang="ts">
import { HELP_GUIDE_ACTION_EVENT, type HelpGuideActionPhase, type HelpGuideStep } from '~/composables/useHelpGuides';

const { currentGuide, activeGuide, finishGuide } = useHelpGuides();
const route = useRoute();

const steps = computed(() => currentGuide.value?.steps ?? []);
const tour = useTour(steps, { scrollIntoView: false });

const tourOpen = computed(() => tour.open.value);
const currentStep = computed<HelpGuideStep | undefined>(() => tour.current.value as HelpGuideStep | undefined);
const showOverlay = computed(() => tourOpen.value && currentStep.value?.showOverlay !== false);
const spotlight = ref<{ top: number; right: number; bottom: number; left: number }>();
const centerBottomReference = {
    getBoundingClientRect() {
        const x = typeof window === 'undefined' ? 0 : window.innerWidth / 2;
        const y = typeof window === 'undefined' ? 0 : window.innerHeight - 24;
        return { x, y, top: y, left: x, right: x, bottom: y, width: 0, height: 0, toJSON: () => ({}) };
    },
};
const tourReference = computed(() => {
    // useTour resolves string targets as element IDs. Help guides use CSS selectors,
    // so resolve the visible target ourselves for correct popover positioning.
    spotlight.value;
    if (currentStep.value?.popoverPosition === 'center-bottom') return centerBottomReference;

    return getTargetElement(currentStep.value?.popoverTarget ?? currentStep.value?.target) ?? tour.reference.value;
});
let activeStep: HelpGuideStep | undefined;
let targetObserver: MutationObserver | undefined;

function dispatchAction(actionId: string, phase: HelpGuideActionPhase, step: HelpGuideStep): void {
    const guide = currentGuide.value;
    if (!guide) return;

    window.dispatchEvent(
        new CustomEvent(HELP_GUIDE_ACTION_EVENT, {
            detail: {
                guideId: guide.id,
                stepId: step.id,
                actionId,
                phase,
            },
        }),
    );
}

function leaveStep(step = activeStep): void {
    if (!step) return;
    if (step.onLeaveAction) dispatchAction(step.onLeaveAction, 'leave', step);
    activeStep = undefined;
}

function getTargetElements(target?: unknown): HTMLElement[] {
    if (typeof target !== 'string') return [];

    return Array.from(document.querySelectorAll<HTMLElement>(target)).filter((element) => {
        const style = window.getComputedStyle(element);
        return element.getClientRects().length > 0 && style.display !== 'none' && style.visibility !== 'hidden';
    });
}

function getTargetElement(target?: unknown): HTMLElement | undefined {
    return getTargetElements(target)[0];
}

function setTargetHighlight(step?: HelpGuideStep): void {
    document.querySelectorAll<HTMLElement>('[data-tour-active]').forEach((element) => {
        element.removeAttribute('data-tour-active');
    });
    document.querySelectorAll<HTMLElement>('[data-tour-disabled]').forEach((element) => {
        element.removeAttribute('data-tour-disabled');
    });

    getTargetElements(step?.target).forEach((target) => {
        target.setAttribute('data-tour-active', 'true');
        if (step?.disableInteraction) target.setAttribute('data-tour-disabled', 'true');
    });
}

function updateSpotlight(target?: unknown): void {
    const elements = getTargetElements(target);
    if (elements.length === 0) {
        spotlight.value = undefined;
        return;
    }

    const padding = 8;
    const rects = elements.map((element) => element.getBoundingClientRect());
    const rect = {
        top: Math.min(...rects.map((item) => item.top)),
        right: Math.max(...rects.map((item) => item.right)),
        bottom: Math.max(...rects.map((item) => item.bottom)),
        left: Math.min(...rects.map((item) => item.left)),
    };
    spotlight.value = {
        top: Math.max(0, rect.top - padding),
        right: Math.max(0, window.innerWidth - rect.right - padding),
        bottom: Math.max(0, window.innerHeight - rect.bottom - padding),
        left: Math.max(0, rect.left - padding),
    };
}

function observeLateTarget(step?: HelpGuideStep): void {
    targetObserver?.disconnect();
    targetObserver = undefined;

    if (!step?.target || getTargetElement(step.target) || typeof document === 'undefined' || !document.body) return;

    targetObserver = new MutationObserver(() => {
        if (!getTargetElement(step.target)) return;

        targetObserver?.disconnect();
        targetObserver = undefined;
        setTargetHighlight(step);
        updateSpotlight(step.target);
    });
    targetObserver.observe(document.body, {
        childList: true,
        subtree: true,
        attributes: true,
        attributeFilter: ['class', 'style'],
    });
}

function refreshSpotlight(): void {
    updateSpotlight(currentStep.value?.target);
}

function moveToTarget(target?: unknown): void {
    const element = getTargetElement(target);
    element?.scrollIntoView({ behavior: 'smooth', block: 'center' });
    requestAnimationFrame(() => updateSpotlight(target));
}

watch(activeGuide, async (guideId) => {
    if (guideId === undefined) {
        leaveStep();
        setTargetHighlight();
        tour.finish();
        return;
    }

    await nextTick();
    tour.start(0);
});

watch(steps, (nextSteps) => {
    const resumeToStepId = activeStep?.resumeToStepId;
    if (!resumeToStepId || nextSteps.some((step) => step.id === activeStep?.id)) return;

    const resumeIndex = nextSteps.findIndex((step) => step.id === resumeToStepId);
    if (resumeIndex > -1 && tour.index.value !== resumeIndex) tour.goTo(resumeIndex);
});

watch(
    currentStep,
    async (step) => {
        if (activeStep?.id !== step?.id) leaveStep();

        setTargetHighlight(step);
        updateSpotlight(step?.target);
        observeLateTarget(step);
        if (step && activeStep?.id !== step.id) {
            activeStep = step;
            if (step.onEnterAction) {
                dispatchAction(step.onEnterAction, 'enter', step);
                if (step.advanceAfterEnterAction) {
                    await nextTick();
                    await new Promise<void>((resolve) => requestAnimationFrame(() => resolve()));
                    if (activeStep?.id === step.id) tour.next();
                    return;
                }
            }
            await nextTick();
            updateSpotlight(step.target);
            observeLateTarget(step);
        }
        if (step?.moveToTarget) {
            await nextTick();
            moveToTarget(step.target);
        }
    },
    { flush: 'post' },
);

watch(
    () => route.fullPath,
    async () => {
        await nextTick();
        setTargetHighlight(currentStep.value);
        updateSpotlight(currentStep.value?.target);
        observeLateTarget(currentStep.value);
    },
);

onMounted(() => {
    window.addEventListener('resize', refreshSpotlight);
    window.addEventListener('scroll', refreshSpotlight, true);
});

onBeforeUnmount(() => {
    window.removeEventListener('resize', refreshSpotlight);
    window.removeEventListener('scroll', refreshSpotlight, true);
    targetObserver?.disconnect();
    setTargetHighlight();
});

async function runAction(): Promise<void> {
    const action = currentStep.value?.action;
    const step = currentStep.value;
    if (!action || !step) return;

    dispatchAction(action.id, 'trigger', step);

    if (action.advanceAfterAction) {
        await nextTick();
        await new Promise<void>((resolve) => requestAnimationFrame(() => resolve()));
        tour.next();
    }
}

function finish(): void {
    leaveStep();
    setTargetHighlight();
    tour.finish();
    finishGuide();
}
</script>

<template>
    <template v-if="showOverlay">
        <template v-if="spotlight">
            <div
                class="pointer-events-none fixed inset-x-0 top-0 z-40 bg-black/50 backdrop-blur-[1px]"
                :style="{ height: `${spotlight.top}px` }"
            />
            <div
                class="pointer-events-none fixed inset-x-0 bottom-0 z-40 bg-black/50 backdrop-blur-[1px]"
                :style="{ height: `${spotlight.bottom}px` }"
            />
            <div
                class="pointer-events-none fixed left-0 z-40 bg-black/50 backdrop-blur-[1px]"
                :style="{
                    top: `${spotlight.top}px`,
                    bottom: `${spotlight.bottom}px`,
                    width: `${spotlight.left}px`,
                }"
            />
            <div
                class="pointer-events-none fixed right-0 z-40 bg-black/50 backdrop-blur-[1px]"
                :style="{
                    top: `${spotlight.top}px`,
                    bottom: `${spotlight.bottom}px`,
                    width: `${spotlight.right}px`,
                }"
            />
        </template>
        <div v-else class="pointer-events-none fixed inset-0 z-40 bg-black/50 backdrop-blur-[1px]" />
    </template>

    <UPopover
        v-if="currentStep"
        :open="tourOpen"
        :reference="tourReference"
        :dismissible="false"
        :content="{
            side: currentStep.popoverPosition === 'center-bottom' ? 'top' : (currentStep.side ?? 'bottom'),
            sideOffset: 12,
        }"
        :ui="{ content: 'pointer-events-auto relative z-[100] min-w-80 max-w-120 p-4 ring-primary ring-2' }"
    >
        <template #content>
            <div class="space-y-4">
                <div class="space-y-1">
                    <p class="text-xs font-medium text-muted">
                        {{ $t('help_guides.step_count', { current: tour.index.value + 1, total: tour.total.value }) }}
                    </p>
                    <h2 class="font-semibold text-highlighted">{{ currentStep.title }}</h2>
                    <p class="text-sm text-muted">{{ currentStep.body }}</p>
                </div>

                <div class="flex items-center justify-between gap-2">
                    <UButton color="neutral" variant="ghost" size="sm" :label="$t('help_guides.skip')" @click="finish" />

                    <div class="flex flex-wrap justify-end gap-2">
                        <UButton
                            v-if="currentStep.action"
                            color="primary"
                            variant="soft"
                            size="sm"
                            :label="currentStep.action.label"
                            @click.stop="runAction"
                        />
                        <UButton
                            v-if="tour.hasPrev.value"
                            color="neutral"
                            variant="soft"
                            size="sm"
                            :label="$t('common.go_back')"
                            @click.stop="tour.prev()"
                        />
                        <UButton
                            size="sm"
                            :label="tour.hasNext.value ? $t('common.next') : $t('help_guides.finish')"
                            @click.stop="tour.hasNext.value ? tour.next() : finish()"
                        />
                    </div>
                </div>
            </div>
        </template>
    </UPopover>
</template>

<style>
[data-tour-active='true'] {
    position: relative;
    z-index: 50;
    outline: 3px solid var(--ui-primary);
    outline-offset: 3px;
    border-radius: var(--ui-radius);
    box-shadow: 0 0 0 6px color-mix(in srgb, var(--ui-primary) 20%, transparent);
}

[data-tour-disabled='true'] {
    pointer-events: none;
}
</style>
