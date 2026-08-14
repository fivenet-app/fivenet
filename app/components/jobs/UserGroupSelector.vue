<script lang="ts" setup>
import SelectMenu from '~/components/partials/SelectMenu.vue';
import ColleagueName from '~/components/jobs/colleagues/ColleagueName.vue';
import { useCompletorStore } from '~/stores/completor';
import type { Colleague } from '~~/gen/ts/resources/jobs/colleagues/colleagues';
import type { UserSelector } from '~~/gen/ts/resources/jobs/user_selector';

const props = withDefaults(
    defineProps<{
        modelValue?: UserSelector;
        groupsOnly?: boolean;
    }>(),
    {
        modelValue: () => ({ userIds: [] }),
        groupsOnly: false,
    },
);

const emit = defineEmits<{
    'update:modelValue': [value: UserSelector];
}>();

const completorStore = useCompletorStore();

const USER_PREFIX = 'u:';
const GROUP_PREFIX = 'g:';

function toUserValue(id: number): string {
    return `${USER_PREFIX}${id}`;
}

function toGroupValue(id: number): string {
    return `${GROUP_PREFIX}${id}`;
}

const selectedValues = ref<Set<string>>(new Set());
let syncingFromProps = false;

function setFromProps() {
    const vals = new Set<string>();
    if (!props.groupsOnly) {
        for (const uid of props.modelValue?.userIds ?? []) {
            vals.add(toUserValue(uid));
        }
    }
    for (const gid of props.modelValue?.groups?.groupIds ?? []) {
        vals.add(toGroupValue(gid));
    }
    selectedValues.value = vals;
}

const includeLeaders = ref<boolean>(props.modelValue?.groups?.includeLeaders ?? false);

function emitValue() {
    const userIds: number[] = [];
    const groupIds: number[] = [];

    for (const val of selectedValues.value) {
        if (val.startsWith(USER_PREFIX)) {
            const id = parseInt(val.slice(USER_PREFIX.length), 10);
            if (!isNaN(id) && !props.groupsOnly) userIds.push(id);
        } else if (val.startsWith(GROUP_PREFIX)) {
            const id = parseInt(val.slice(GROUP_PREFIX.length), 10);
            if (!isNaN(id)) groupIds.push(id);
        }
    }

    const selector: UserSelector = {
        userIds,
    };
    if (groupIds.length > 0 || includeLeaders.value) {
        selector.groups = {
            groupIds,
            includeLeaders: includeLeaders.value,
            includeExcluded: false,
        };
    }
    emit('update:modelValue', selector);
}

watch(
    () => props.modelValue,
    () => {
        syncingFromProps = true;
        setFromProps();
        includeLeaders.value = props.modelValue?.groups?.includeLeaders ?? false;
        syncingFromProps = false;

        if (props.groupsOnly && (props.modelValue?.userIds?.length ?? 0) > 0) {
            emitValue();
        }
    },
    {
        deep: true,
        immediate: true,
    },
);

watch(
    selectedValues,
    () => {
        if (!syncingFromProps) emitValue();
    },
    { flush: 'sync' },
);
watch(
    includeLeaders,
    () => {
        if (!syncingFromProps) emitValue();
    },
    { flush: 'sync' },
);

const groupIds = computed(() => {
    const ids: number[] = [];
    for (const val of selectedValues.value) {
        if (val.startsWith(GROUP_PREFIX)) {
            const id = parseInt(val.slice(GROUP_PREFIX.length), 10);
            if (!isNaN(id)) ids.push(id);
        }
    }
    return ids;
});

function handleUpdate(vals: unknown) {
    if (Array.isArray(vals)) {
        selectedValues.value = new Set(vals as string[]);
    } else if (typeof vals === 'string') {
        selectedValues.value = new Set([vals]);
    } else {
        selectedValues.value = new Set();
    }
}

type SearchUserItem = Colleague & {
    value: string;
    kind: 'user';
    label: string;
};

interface SearchGroupItem {
    value: string;
    kind: 'group';
    label: string;
    id: number;
}

type SearchItem = SearchUserItem | SearchGroupItem;

async function searchItems(q: string): Promise<SearchItem[]> {
    const restoredItems = await loadRestoredItems();

    if (!q) return [...restoredItems];

    const items: SearchItem[] = [];

    if (props.groupsOnly || q === '@' || q.startsWith('@')) {
        const query = q.startsWith('@') ? q.slice(1).trim() : q.trim();
        const groups = await completorStore.completeGroups(query || '');

        items.push(
            ...(groups ?? []).map((g) => ({
                id: g.id,
                value: toGroupValue(g.id),
                kind: 'group' as const,
                label: g.shortName ? `${g.shortName}: ${g.name}` : g.name || String(g.id),
            })),
        );
    } else {
        const users = await completorStore.completeColleagues(q);
        items.push(
            ...(users ?? []).map((u) => ({
                ...u,
                value: toUserValue(u.userId),
                kind: 'user' as const,
                label: userToLabel(u),
            })),
        );
    }

    for (const restoredItem of restoredItems) {
        if (!items.some((item) => item.value === restoredItem.value)) items.push(restoredItem);
    }

    return items;
}

const restoredItems = ref<SearchItem[]>([]);
let restoredItemsKey = '';
let restoredItemsLoadingKey = '';
let restoredItemsPromise: Promise<SearchItem[]> | undefined;

function selectedValuesKey(): string {
    return [...selectedValues.value].sort().join('\u0000');
}

async function loadRestoredItems(): Promise<SearchItem[]> {
    const key = selectedValuesKey();
    if (key === restoredItemsKey) return restoredItems.value;
    if (restoredItemsPromise && key === restoredItemsLoadingKey) return restoredItemsPromise;

    restoredItemsLoadingKey = key;
    const promise = (async () => {
        const vals = new Set(selectedValues.value);
        if (vals.size === 0) return [];

        const userIds: number[] = [];
        const targetGroupIds = new Set<number>();

        for (const val of vals) {
            if (val.startsWith(USER_PREFIX)) {
                const id = parseInt(val.slice(USER_PREFIX.length), 10);
                if (!isNaN(id) && !props.groupsOnly) userIds.push(id);
            } else if (val.startsWith(GROUP_PREFIX)) {
                const id = parseInt(val.slice(GROUP_PREFIX.length), 10);
                if (!isNaN(id)) targetGroupIds.add(id);
            }
        }

        const items: SearchItem[] = [];
        if (userIds.length > 0) {
            const users = await completorStore.completeColleagues('', userIds, true);
            for (const u of users ?? []) {
                items.push({
                    ...u,
                    value: toUserValue(u.userId),
                    kind: 'user' as const,
                    label: userToLabel(u),
                });
            }
        }

        if (targetGroupIds.size > 0) {
            const groups = await completorStore.completeGroups('', [...targetGroupIds.values()]);
            for (const g of groups ?? []) {
                if (targetGroupIds.has(g.id)) {
                    items.push({
                        id: g.id,
                        value: toGroupValue(g.id),
                        kind: 'group' as const,
                        label: g.shortName ? `${g.shortName}: ${g.name}` : g.name || String(g.id),
                    });
                }
            }
        }

        if (selectedValuesKey() === key) {
            restoredItems.value = items;
            restoredItemsKey = key;
        }

        return items;
    })();
    restoredItemsPromise = promise;

    try {
        return await promise;
    } finally {
        if (restoredItemsPromise === promise) {
            restoredItemsPromise = undefined;
            restoredItemsLoadingKey = '';
        }
    }
}
</script>

<template>
    <div class="flex flex-col gap-2">
        <div class="flex items-end gap-2">
            <UFormField class="flex-1" name="users" v-bind="$attrs">
                <UFieldGroup v-bind="$attrs">
                    <SelectMenu
                        :model-value="[...selectedValues]"
                        :searchable="searchItems"
                        multiple
                        :search-input="{
                            placeholder: props.groupsOnly
                                ? $t('common.group', 2)
                                : $t('components.jobs.user_group_selector.search_placeholder'),
                        }"
                        :placeholder="props.groupsOnly ? $t('common.group', 2) : $t('common.colleague', 2)"
                        label-key="label"
                        :searchable-key="props.groupsOnly ? 'user-group-select-groups' : 'user-group-select'"
                        value-key="value"
                        v-bind="$attrs"
                        @update:model-value="handleUpdate"
                    >
                        <template #item-label="{ item }">
                            <ColleagueName v-if="item.kind === 'user'" class="truncate" :colleague="item" birthday />
                            <span v-else class="truncate">@{{ item.label }}</span>
                        </template>

                        <template #item-leading="{ item }">
                            <UIcon v-if="item.kind === 'group'" class="size-5" name="i-mdi-account-group" />
                        </template>

                        <template #empty>
                            {{ $t('common.not_found', [props.groupsOnly ? $t('common.group', 2) : $t('common.colleague', 2)]) }}
                        </template>
                    </SelectMenu>

                    <UTooltip
                        v-if="groupIds.length > 0"
                        :text="
                            includeLeaders
                                ? $t('components.jobs.user_group_selector.leaders.exclude')
                                : $t('components.jobs.user_group_selector.leaders.include')
                        "
                    >
                        <UButton
                            :color="includeLeaders ? 'warning' : 'neutral'"
                            trailing-icon="i-mdi-account-star"
                            :label="$t('components.jobs.groups.leaders')"
                            @click="includeLeaders = !includeLeaders"
                        />
                    </UTooltip>
                </UFieldGroup>
            </UFormField>
        </div>
    </div>
</template>
