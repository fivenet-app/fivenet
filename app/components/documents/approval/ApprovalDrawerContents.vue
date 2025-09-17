<script setup lang="ts">
import { computed, reactive, ref } from 'vue';
import { ApprovalTaskStatus, type ApprovalTask } from '~~/gen/ts/resources/documents/approval';
import type { ApprovalPanelSnapshot } from '~~/gen/ts/services/documents/approval';

const props = defineProps<{
    panel: ApprovalPanelSnapshot;
    canManage?: boolean;
    me?: number;
}>();

const emit = defineEmits<{
    (e: 'refresh'): void;
    (e: 'close'): void;
    (e: 'approve', taskId: number, comment?: string): void;
    (e: 'decline', taskId: number, reason: string): void;
    (e: 'cancelTask', taskId: number): void;
    (e: 'reassign', taskId: number): void;
    (e: 'addReviewers', payload: { groupKey: string; orderNo: number; selector: any }): void;
}>();

// Derived data
const versionLabel = computed(() => 'v0');

const sortedGroups = computed(() => {
    return [...(props.panel?.groups ?? [])].sort((a, b) => a.orderNo - b.orderNo || a.groupKey.localeCompare(b.groupKey));
});

const tasksByGroup = computed<Record<string, ApprovalTask[]>>(() => {
    const map: Record<string, ApprovalTask[]> = {};
    // Ensure myPending are always visible in their groups (even if not full list)
    // In a real app you'd fetch group tasks separately; for now we render only myPending inside the groups if present
    for (const t of props.panel?.myPending ?? []) {
        const k = `${t.groupKey}:${t.orderNo}`;
        if (!map[k]) map[k] = [];
        map[k].push(t);
    }
    return map;
});

function gKey(g: GroupSummary) {
    return `${g.groupKey}:${g.orderNo}`;
}
function prettyGroupName(key: string) {
    return key.replaceAll('_', ' ').replace(/\b\w/g, (c) => c.toUpperCase());
}
function isMyPending(t: ApprovalTask) {
    return t.status === ApprovalTaskStatus.PENDING;
}

// Leader manage menu
const canManage = computed(() => !!props.canManage);
function manageMenu(t: ApprovalTask) {
    return [
        [{ label: 'Cancel task', icon: 'i-heroicons-no-symbol', click: () => emit('cancelTask', t.id) }],
        [{ label: 'Reassign…', icon: 'i-heroicons-user-plus', click: () => emit('reassign', t.id) }],
    ];
}

// Add reviewers (inline)
const addingToGroup = ref<string | null>(null);
const addingBusy = ref(false);
const selector = reactive<{ kind: any; value: string; minRank?: number }>({
    kind: { label: 'User', value: 'user' },
    value: '',
});
const selectorKinds = [
    { label: 'User', value: 'user' },
    { label: 'Role', value: 'role' },
    { label: 'Faction ≥ Rank', value: 'faction' },
];
function toggleAddReviewer(g: GroupSummary) {
    const k = gKey(g);
    addingToGroup.value = addingToGroup.value === k ? null : k;
}
async function emitAddReviewer(g: GroupSummary) {
    addingBusy.value = true;
    try {
        emit('addReviewers', {
            groupKey: g.groupKey,
            orderNo: g.orderNo,
            selector: { kind: selector.kind.value, value: selector.value, minRank: selector.minRank },
        });
        addingToGroup.value = null;
    } finally {
        addingBusy.value = false;
    }
}

// Decide modal
const decideOpen = ref(false);
const decideMode = ref<'approve' | 'decline'>('approve');
const decideTask = ref<ApprovalTask | null>(null);
const decideComment = ref('');

const decideTaskTitle = computed(() =>
    decideTask.value ? `${prettyGroupName(decideTask.value.groupKey)} · order ${decideTask.value.orderNo}` : '',
);

function openDecide(t: ApprovalTask, mode: 'approve' | 'decline') {
    decideTask.value = t;
    decideMode.value = mode;
    decideComment.value = '';
    decideOpen.value = true;
}

function submitDecision() {
    if (!decideTask.value) return;
    const id = decideTask.value.id;
    if (decideMode.value === 'approve') {
        emit('approve', id, decideComment.value || undefined);
    } else {
        const reason = decideComment.value?.trim();
        if (!reason) return; // could show toast
        emit('decline', id, reason);
    }
    decideOpen.value = false;
}
</script>

<template>
    <div class="flex h-full flex-col gap-3 p-4">
        <!-- Header -->
        <header class="flex items-start justify-between gap-3">
            <div class="min-w-0">
                <p class="mt-0.5 text-sm text-gray-500">
                    Version {{ versionLabel }} ·
                    <span v-if="panel.allGroupsComplete" class="text-green-600">All stages complete</span>
                    <span v-else class="text-amber-600">Awaiting order {{ panel.currentOrder }}</span>
                </p>
            </div>
            <div class="flex items-center gap-2">
                <UButton icon="i-heroicons-arrow-path" size="sm" @click="$emit('refresh')">Refresh</UButton>
            </div>
        </header>

        <USeparator />

        <!-- Progress summary by group (stage) -->
        <section class="grid grid-cols-1 gap-3 md:grid-cols-2">
            <UCard v-for="g in sortedGroups" :key="g.groupKey + ':' + g.orderNo" :ui="{ body: { padding: 'p-3' } }">
                <div class="flex items-start justify-between gap-3">
                    <div class="min-w-0">
                        <div class="flex items-center gap-2">
                            <h3 class="truncate font-medium">{{ prettyGroupName(g.groupKey) }}</h3>
                            <UBadge v-if="g.complete" color="green" variant="subtle">Complete</UBadge>
                            <UBadge v-else-if="g.orderNo === panel.currentOrder" color="amber" variant="subtle">Current</UBadge>
                            <UBadge v-else color="gray" variant="subtle">Queued</UBadge>
                        </div>
                        <p class="mt-1 text-xs text-gray-500">
                            Order {{ g.orderNo }} ·
                            <span v-if="g.quorumAny">Quorum: any {{ g.quorumAny }}</span>
                            <span v-else>Require all</span>
                        </p>
                    </div>
                    <div class="text-right text-sm">
                        <div class="flex items-center justify-end gap-2">
                            <UBadge color="primary" variant="soft">{{ g.approved }}/{{ g.assigned }}</UBadge>
                            <UBadge v-if="g.declined" color="red" variant="soft">{{ g.declined }} declined</UBadge>
                        </div>
                    </div>
                </div>

                <USeparator class="my-3" />

                <!-- Group tasks table (compact) -->
                <div class="space-y-2">
                    <div
                        v-for="t in tasksByGroup[gKey(g)]"
                        :key="t.id"
                        class="flex items-start justify-between gap-2 rounded-lg border p-2"
                    >
                        <div class="flex min-w-0 items-center gap-2">
                            <UAvatar :src="t.user?.avatarUrl" :alt="t.user?.displayName" size="xs" />
                            <div class="min-w-0">
                                <p class="truncate text-sm font-medium">{{ t.user?.displayName || 'User #' + t.userId }}</p>
                                <p class="truncate text-xs text-gray-500">
                                    <span v-if="t.role?.name">{{ t.role.name }}</span>
                                    <span v-if="t.faction?.name"> · {{ t.faction.name }}</span>
                                </p>
                            </div>
                        </div>

                        <div class="flex items-center gap-2">
                            <UBadge v-if="t.status === ApprovalTaskStatus.APPROVED" color="green" variant="solid"
                                >Approved</UBadge
                            >
                            <UBadge v-else-if="t.status === ApprovalTaskStatus.DECLINED" color="red" variant="solid"
                                >Declined</UBadge
                            >
                            <UBadge v-else-if="t.status === ApprovalTaskStatus.EXPIRED" color="orange" variant="soft"
                                >Expired</UBadge
                            >
                            <UBadge v-else-if="t.status === ApprovalTaskStatus.CANCELLED" color="gray" variant="soft"
                                >Cancelled</UBadge
                            >
                            <UBadge v-else color="blue" variant="soft">Pending</UBadge>

                            <!-- Quick actions for my pending tasks only -->
                            <template v-if="isMyPending(t)">
                                <UButton
                                    size="xs"
                                    color="green"
                                    variant="soft"
                                    icon="i-heroicons-check-20-solid"
                                    @click="openDecide(t, 'approve')"
                                    >Approve</UButton
                                >
                                <UButton
                                    size="xs"
                                    color="red"
                                    variant="soft"
                                    icon="i-heroicons-x-mark-20-solid"
                                    @click="openDecide(t, 'decline')"
                                    >Decline</UButton
                                >
                            </template>

                            <!-- Leader tools: cancel / reassign -->
                            <UDropdownMenu v-if="canManage" :items="manageMenu(t)" :popper="{ placement: 'bottom-end' }">
                                <UButton icon="i-heroicons-ellipsis-vertical" color="gray" variant="ghost" size="xs" />
                            </UDropdownMenu>
                        </div>
                    </div>
                </div>

                <!-- Add reviewers into this group (leader-only) -->
                <div v-if="canManage" class="mt-3 rounded-md bg-gray-50 p-2">
                    <div class="flex items-center justify-between">
                        <p class="text-xs font-medium text-gray-700">Add reviewers to this stage</p>
                        <UButton size="xs" color="gray" variant="ghost" icon="i-heroicons-plus" @click="toggleAddReviewer(g)"
                            >Add</UButton
                        >
                    </div>
                    <div v-if="addingToGroup === gKey(g)" class="mt-2 grid grid-cols-1 gap-2 md:grid-cols-3">
                        <USelectMenu v-model="selector.kind" :options="selectorKinds" option-attribute="label" />
                        <UInput
                            v-if="selector.kind.value === 'user'"
                            v-model="selector.value"
                            placeholder="Search user id or name"
                        />
                        <UInput v-else-if="selector.kind.value === 'role'" v-model="selector.value" placeholder="Role ID" />
                        <div v-else class="flex gap-2">
                            <UInput v-model="selector.value" class="flex-1" placeholder="Faction ID" />
                            <UInput v-model.number="selector.minRank" class="w-28" placeholder="Min rank" />
                        </div>
                        <UButton class="md:col-span-3" size="xs" :loading="addingBusy" @click="emitAddReviewer(g)"
                            >Add to stage</UButton
                        >
                    </div>
                </div>
            </UCard>
        </section>

        <USeparator />

        <!-- My pending tasks (quick act) -->
        <section>
            <h3 class="mb-2 text-sm font-semibold">Your pending approvals</h3>
            <div v-if="panel.myPending?.length" class="space-y-2">
                <div
                    v-for="t in panel.myPending"
                    :key="'my-' + t.id"
                    class="flex items-center justify-between rounded-lg border p-2"
                >
                    <div class="truncate text-sm">
                        <span class="font-medium">{{ prettyGroupName(t.groupKey) }}</span>
                        <span class="text-gray-500"> · order {{ t.orderNo }}</span>
                    </div>
                    <div class="flex items-center gap-2">
                        <UButton size="xs" color="green" variant="soft" @click="openDecide(t, 'approve')">Approve</UButton>
                        <UButton size="xs" color="red" variant="soft" @click="openDecide(t, 'decline')">Decline</UButton>
                    </div>
                </div>
            </div>
            <UAlert v-else color="gray" variant="soft" title="No pending approvals" class="mt-2" />
        </section>

        <!-- Decide modal -->
        <UModal v-model="decideOpen" :ui="{ content: 'max-w-md' }">
            <template #title>
                <div class="flex items-start justify-between">
                    <p class="text-sm text-gray-500">{{ decideMode === 'approve' ? 'Approve' : 'Decline' }} task</p>
                    <h3 class="font-semibold">{{ decideTaskTitle }}</h3>
                </div>
            </template>

            <template #body>
                <div class="space-y-3">
                    <div v-if="decideMode === 'approve'">
                        <UTextarea v-model="decideComment" placeholder="Optional comment" :rows="3" />
                    </div>
                    <div v-else>
                        <UTextarea v-model="decideComment" placeholder="Reason (required)" :rows="3" />
                        <p class="text-xs text-gray-500">A short reason helps the author understand the change.</p>
                    </div>
                </div>
            </template>

            <template #footer>
                <div class="flex w-full items-center justify-end gap-2">
                    <UButton color="gray" variant="ghost" @click="decideOpen = false">Cancel</UButton>
                    <UButton :color="decideMode === 'approve' ? 'green' : 'red'" :loading="decideBusy" @click="submitDecision">
                        {{ decideMode === 'approve' ? 'Approve' : 'Decline' }}
                    </UButton>
                </div>
            </template>
        </UModal>
    </div>
</template>
