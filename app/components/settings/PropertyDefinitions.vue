<script lang="ts" setup>
import type { Form } from '@nuxt/ui';
import { z } from 'zod';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import ReorderButtons from '~/components/partials/ReorderButtons.vue';
import PropertyDefinitionEditorDrawer from '~/components/settings/PropertyDefinitionEditorDrawer.vue';
import PropertyGroupEditorDrawer from '~/components/settings/PropertyGroupEditorDrawer.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import { getSettingsPropertyClient } from '~~/gen/ts/clients';
import {
    PropertyBehavior,
    PropertyDefinitionSource,
    PropertyDisplayMode,
    PropertyRenderer,
    PropertyScope,
    PropertyValueType,
    type PropertyDefinition,
    type PropertyGroup,
} from '~~/gen/ts/resources/properties/properties';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const { can } = useAuth();
const { t, te } = useI18n();

const overlay = useOverlay();

const notifications = useNotificationsStore();
const settingsPropertyClient = await getSettingsPropertyClient();

const canEdit = can('internal.Superuser/ConfigAdmin');
const sortMode = ref(false);

const schema = z.object({
    scope: z.coerce.number().pipe(z.enum(PropertyScope)).default(PropertyScope.VEHICLE),
    showDeleted: z.boolean().default(false),
});
type Schema = z.output<typeof schema>;

const query = useSearchForm('settings-property-definitions', schema);
const formRef = useTemplateRef<Form<typeof schema>>('formRef');
const { validatedQuery, commitValidatedQuery } = useFormSearchValidation<typeof schema>(query, formRef);

type PropertyDefinitionsData = {
    scopes: PropertyScope[];
    groups: PropertyGroup[];
    definitions: PropertyDefinition[];
};

async function getPropertyDefinitions(values: Schema, signal: AbortSignal): Promise<PropertyDefinitionsData> {
    try {
        const [{ response: scopesResponse }, { response: definitionsResponse }] = await Promise.all([
            settingsPropertyClient.listPropertyScopes({}, { abort: signal }),
            settingsPropertyClient.listPropertyDefinitions(
                { scope: values.scope, includeDeleted: values.showDeleted },
                { abort: signal },
            ),
        ]);

        return {
            scopes: scopesResponse.scopes,
            definitions: definitionsResponse.definitions,
            groups: definitionsResponse.groups,
        };
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const { data, status, error, refresh } = useAuthedLazyAsyncData(
    'capabilities',
    computed(() => `settings-property-definitions:${JSON.stringify(validatedQuery.value)}`),
    ({ signal }) => getPropertyDefinitions(validatedQuery.value, signal),
    {
        default: (): PropertyDefinitionsData => ({ scopes: [], definitions: [], groups: [] }),
    },
);

const scopes = computed(() => data.value?.scopes ?? []);
const definitions = computed(() => data.value?.definitions ?? []);
const groups = computed(() => data.value?.groups.filter((group) => group.scope === validatedQuery.value.scope) ?? []);
const pinnedCount = computed(
    () => definitions.value.filter((definition) => definition.pinned && !isDeletedDefinition(definition)).length,
);

function groupLabel(group: PropertyGroup): string {
    if (group.label) return group.label;
    const key = group.displayKey || `components.settings.property_definitions.groups.${group.key}`;
    return te(key) ? t(key) : group.key;
}

function isSystemGroup(group?: PropertyGroup): boolean {
    return Boolean(group?.displayKey);
}

function isDeletedGroup(group?: PropertyGroup): boolean {
    return group?.deletedAt !== undefined;
}

const definitionGroups = computed(() => {
    type DefinitionGroup = {
        key: string;
        label: string;
        definitions: PropertyDefinition[];
        group?: PropertyGroup;
    };
    const definitionsByGroup = new Map<string, PropertyDefinition[]>();
    for (const definition of definitions.value) {
        const group = definitionsByGroup.get(definition.groupKey) ?? [];
        group.push(definition);
        definitionsByGroup.set(definition.groupKey, group);
    }

    const result: DefinitionGroup[] = groups.value.map((group) => ({
        key: group.key,
        label: groupLabel(group),
        definitions: definitionsByGroup.get(group.key) ?? [],
        group,
    }));
    const ungrouped = definitionsByGroup.get('') ?? [];
    if (ungrouped.length > 0) {
        result.push({
            key: 'ungrouped',
            label: t('components.settings.property_definitions.ungrouped'),
            definitions: ungrouped,
            group: undefined,
        });
    }
    return result;
});

function moveProperty(groupKey: string, index: number, direction: -1 | 1): void {
    if (!data.value || !sortMode.value) return;
    const group = definitionGroups.value.find((entry) => entry.key === groupKey);
    if (!group) return;
    const target = index + direction;
    if (target < 0 || target >= group.definitions.length) return;
    const current = group.definitions[index]!;
    const next = group.definitions[target]!;
    const currentIndex = data.value.definitions.indexOf(current);
    const nextIndex = data.value.definitions.indexOf(next);
    if (currentIndex < 0 || nextIndex < 0) return;
    [data.value.definitions[currentIndex], data.value.definitions[nextIndex]] = [next, current];
    group.definitions.forEach((definition, order) => (definition.sortOrder = (order + 1) * 100));
    void savePropertyOrder(groupKey, group.definitions);
}

async function savePropertyOrder(groupKey: string, ordered: PropertyDefinition[]): Promise<void> {
    try {
        await settingsPropertyClient.reorderPropertyDefinitions({
            scope: query.scope,
            groupKey: groupKey === 'ungrouped' ? '' : groupKey,
            propertyKeys: ordered.map((definition) => definition.key),
        });
    } catch (e) {
        handleGRPCError(e as RpcError);
        await refresh();
    }
}

function moveGroup(groupKey: string, direction: -1 | 1): void {
    if (!data.value || groupKey === 'ungrouped' || !sortMode.value) return;
    const ordered = [...groups.value];
    const index = ordered.findIndex((group) => group.key === groupKey);
    const target = index + direction;
    if (index < 0 || target < 0 || target >= ordered.length) return;
    [ordered[index], ordered[target]] = [ordered[target]!, ordered[index]!];
    ordered.forEach((group, order) => (group.sortOrder = (order + 1) * 100));
    data.value.groups = ordered;
    void saveGroupOrder(ordered);
}

async function saveGroupOrder(ordered: PropertyGroup[]): Promise<void> {
    try {
        await settingsPropertyClient.reorderPropertyGroups({
            scope: query.scope,
            groupKeys: ordered.map((group) => group.key),
        });
    } catch (e) {
        handleGRPCError(e as RpcError);
        await refresh();
    }
}

const propertyDefinitionEditorDrawer = overlay.create(PropertyDefinitionEditorDrawer);
const propertyGroupEditorDrawer = overlay.create(PropertyGroupEditorDrawer);
const confirmModal = overlay.create(ConfirmModal);

function openGroupEditor(group?: PropertyGroup): void {
    if (group && isSystemGroup(group)) return;
    propertyGroupEditorDrawer.open({
        group,
        scope: query.scope,
        onSave: (value: PropertyGroup) => saveCustomGroup(value, group !== undefined),
    });
}

async function saveCustomGroup(group: PropertyGroup, isEditing: boolean): Promise<void> {
    if (isSystemGroup(group)) return;
    try {
        const { response } = isEditing
            ? await settingsPropertyClient.updatePropertyGroup({ group: { ...group, displayKey: '' } })
            : await settingsPropertyClient.createPropertyGroup({ group: { ...group, displayKey: '' } });
        if (response.group) {
            const index = data.value?.groups.findIndex((item) => item.key === response.group?.key) ?? -1;
            if (data.value && index === -1) data.value.groups.push(response.group);
            else if (data.value && index !== -1) data.value.groups[index] = response.group;
        }
        propertyGroupEditorDrawer.close();
        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

function deleteCustomGroup(group: PropertyGroup): void {
    if (isSystemGroup(group)) return;
    confirmModal.open({
        title: t('components.settings.property_definitions.delete_group_title'),
        description: t('components.settings.property_definitions.delete_group_description'),
        confirm: async () => {
            try {
                await settingsPropertyClient.deletePropertyGroup({ scope: query.scope, key: group.key });
                if (!query.showDeleted && data.value) {
                    data.value.groups = data.value.groups.filter((item) => item.key !== group.key);
                }
                await refresh();
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        },
    });
}

async function restoreCustomGroup(group: PropertyGroup): Promise<void> {
    if (isSystemGroup(group)) return;
    try {
        const { response } = await settingsPropertyClient.restorePropertyGroup({ scope: query.scope, key: group.key });
        if (response.group) Object.assign(group, response.group);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

function scopeLabel(scope: PropertyScope): string {
    if (scope === PropertyScope.VEHICLE) return t('components.settings.property_definitions.scopes.vehicle');
    if (scope === PropertyScope.USER) return t('components.settings.property_definitions.scopes.user');
    return PropertyScope[scope];
}

function scopeIcon(scope: PropertyScope): string {
    if (scope === PropertyScope.VEHICLE) return 'i-mdi-car-outline';
    if (scope === PropertyScope.USER) return 'i-mdi-account-group-outline';
    return 'i-mdi-help-circle-outline';
}

const scopeItems = computed(() =>
    scopes.value.map((scope) => ({ label: scopeLabel(scope), value: scope, icon: scopeIcon(scope) })),
);

function isSystemDefinition(definition: PropertyDefinition): boolean {
    return definition.source === PropertyDefinitionSource.SYSTEM;
}

function isCustomDefinition(definition: PropertyDefinition): boolean {
    return definition.source === PropertyDefinitionSource.CUSTOM;
}

function isDeletedDefinition(definition: PropertyDefinition): boolean {
    return definition.deletedAt !== undefined;
}

function isManagedDefinition(definition: PropertyDefinition): boolean {
    return definition.behavior === PropertyBehavior.MANAGED;
}

function canDisableDefinition(definition: PropertyDefinition): boolean {
    return !isManagedDefinition(definition);
}

function displayModeLabel(definition: PropertyDefinition): string {
    switch (definition.displayMode) {
        case PropertyDisplayMode.HIDDEN:
            return t('components.settings.property_definitions.visibility_hidden');
        case PropertyDisplayMode.PROFILE_BADGE:
            return t('components.settings.property_definitions.visibility_profile_badge');
        default:
            return t('components.settings.property_definitions.visibility_default');
    }
}

function rendererLabel(definition: PropertyDefinition): string {
    if (definition.renderer === PropertyRenderer.UNSPECIFIED) {
        return t('components.settings.property_definitions.renderer_default');
    }
    const key = `components.settings.property_definitions.renderer_${PropertyRenderer[definition.renderer].toLowerCase()}`;
    return te(key) ? t(key) : PropertyRenderer[definition.renderer];
}

function openEditor(definition?: PropertyDefinition): void {
    propertyDefinitionEditorDrawer.open({
        definition,
        groups: groups.value,
        scope: query.scope,
        onSave: (value: PropertyDefinition) => saveCustomDefinition(value, definition !== undefined),
    });
}

async function saveCustomDefinition(definition: PropertyDefinition, isEditing: boolean): Promise<void> {
    try {
        const payload = { definition };
        const { response } = isEditing
            ? await settingsPropertyClient.updatePropertyDefinition(payload)
            : await settingsPropertyClient.createPropertyDefinition(payload);
        const savedDefinition = response.definition;
        if (savedDefinition) {
            const index = data.value?.definitions.findIndex((item) => item.key === savedDefinition.key) ?? -1;
            if (data.value && index === -1) data.value.definitions.push(savedDefinition);
            else if (data.value && index !== -1) data.value.definitions[index] = savedDefinition;
        }
        propertyDefinitionEditorDrawer.close();
        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function toggleSystemDefinition(definition: PropertyDefinition, enabled: boolean): Promise<void> {
    try {
        const { response } = await settingsPropertyClient.setSystemPropertyEnabled({
            scope: query.scope,
            key: definition.key,
            enabled,
        });
        if (response.definition) Object.assign(definition, response.definition);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function setCustomDefinitionEnabled(definition: PropertyDefinition, enabled: boolean): Promise<void> {
    try {
        const { response } = await settingsPropertyClient.updatePropertyDefinition({
            definition: { ...definition, enabled },
        });
        if (response.definition) Object.assign(definition, response.definition);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function deleteCustomDefinition(definition: PropertyDefinition): Promise<void> {
    try {
        await settingsPropertyClient.deletePropertyDefinition({ scope: query.scope, key: definition.key });
        if (data.value) data.value.definitions = data.value.definitions.filter((item) => item.key !== definition.key);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function restoreCustomDefinition(definition: PropertyDefinition): Promise<void> {
    try {
        const { response } = await settingsPropertyClient.restorePropertyDefinition({
            scope: query.scope,
            key: definition.key,
        });
        if (response.definition) Object.assign(definition, response.definition);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}
</script>

<template>
    <UDashboardPanel :ui="{ root: 'pb-(--page-content-bottom-offset)' }">
        <template #header>
            <UDashboardNavbar :title="t('components.settings.property_definitions.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <PartialsBackButton fallback-to="/settings" />

                    <UTooltip v-if="canEdit" :text="$t('common.change_order')">
                        <UButton
                            color="neutral"
                            variant="outline"
                            :icon="sortMode ? 'i-mdi-check' : 'i-mdi-sort'"
                            @click="sortMode = !sortMode"
                        />
                    </UTooltip>

                    <UButton
                        v-if="canEdit"
                        color="primary"
                        icon="i-mdi-plus"
                        :label="$t('components.settings.property_definitions.add')"
                        @click="openEditor()"
                    />
                    <UButton
                        v-if="canEdit"
                        color="neutral"
                        variant="outline"
                        icon="i-mdi-folder-plus-outline"
                        :label="$t('components.settings.property_definitions.add_group')"
                        @click="openGroupEditor()"
                    />
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar>
                <UForm
                    ref="formRef"
                    data-tour="documents-search"
                    class="my-2 flex w-full flex-1 flex-row gap-2"
                    :schema="schema"
                    :state="query"
                    @submit="commitValidatedQuery"
                >
                    <UFormField name="scope" class="flex-1" :label="t('components.settings.property_definitions.scope')">
                        <USelectMenu
                            v-model="query.scope"
                            class="w-full"
                            :items="scopeItems"
                            label-key="label"
                            value-key="value"
                            :disabled="isRequestPending(status)"
                            :loading="isRequestPending(status)"
                        >
                            <template #leading>
                                <UIcon class="size-4" :name="scopeIcon(query.scope)" />
                            </template>
                        </USelectMenu>
                    </UFormField>

                    <UFormField
                        name="showDeleted"
                        class="shrink-0"
                        :label="t('components.settings.property_definitions.show_deleted')"
                    >
                        <USwitch v-model="query.showDeleted" />
                    </UFormField>
                </UForm>
            </UDashboardToolbar>
        </template>

        <template #body>
            <DataPendingBlock
                v-if="isRequestPending(status)"
                :message="$t('common.loading', [$t('components.settings.property_definitions.title')])"
            />
            <DataErrorBlock
                v-else-if="error"
                :title="t('components.settings.property_definitions.load_error')"
                :error="error"
                :retry="refresh"
            />
            <template v-else>
                <UAlert
                    v-if="pinnedCount >= 6"
                    color="warning"
                    icon="i-mdi-pin"
                    :title="t('components.settings.property_definitions.pinned_warning', { count: pinnedCount })"
                />

                <DataNoDataBlock v-if="definitions.length === 0" :type="t('components.settings.property_definitions.title')" />

                <div v-for="group in definitionGroups" :key="group.key" class="space-y-3">
                    <div class="flex items-center justify-between gap-2">
                        <ReorderButtons
                            v-if="sortMode && group.key !== 'ungrouped'"
                            :idx="groups.findIndex((item) => item.key === group.key)"
                            :move-up="() => moveGroup(group.key, -1)"
                            :move-down="() => moveGroup(group.key, 1)"
                            :disable-up="groups.findIndex((item) => item.key === group.key) <= 0"
                            :disable-down="
                                groups.findIndex((item) => item.key === group.key) < 0 ||
                                groups.findIndex((item) => item.key === group.key) >= groups.length - 1
                            "
                        />
                        <div class="flex min-w-0 items-center gap-2">
                            <h2 class="text-lg font-semibold text-highlighted">{{ group.label }}</h2>
                            <UBadge
                                v-if="isSystemGroup(group.group)"
                                color="primary"
                                :label="$t('components.settings.property_definitions.system')"
                            />
                            <UBadge
                                v-else-if="isDeletedGroup(group.group)"
                                color="error"
                                :label="$t('components.settings.property_definitions.deleted')"
                            />
                        </div>
                        <div v-if="group.group && canEdit && !isSystemGroup(group.group)" class="flex shrink-0 gap-1">
                            <UButton
                                v-if="!isDeletedGroup(group.group)"
                                color="neutral"
                                variant="ghost"
                                icon="i-mdi-pencil"
                                :aria-label="$t('common.update')"
                                @click="openGroupEditor(group.group)"
                            />
                            <UButton
                                :color="isDeletedGroup(group.group) ? 'success' : 'error'"
                                variant="ghost"
                                :icon="isDeletedGroup(group.group) ? 'i-mdi-restore' : 'i-mdi-delete'"
                                :aria-label="isDeletedGroup(group.group) ? $t('common.restore') : $t('common.delete')"
                                @click="
                                    isDeletedGroup(group.group)
                                        ? restoreCustomGroup(group.group)
                                        : deleteCustomGroup(group.group)
                                "
                            />
                        </div>
                    </div>
                    <div class="grid gap-3 xl:grid-cols-2">
                        <UCard v-for="(definition, definitionIndex) in group.definitions" :key="definition.key">
                            <div class="flex flex-col gap-3">
                                <div class="flex items-start justify-between gap-3">
                                    <div class="min-w-0">
                                        <div class="flex items-center gap-2">
                                            <ReorderButtons
                                                v-if="sortMode"
                                                :idx="definitionIndex"
                                                :move-up="() => moveProperty(group.key, definitionIndex, -1)"
                                                :move-down="() => moveProperty(group.key, definitionIndex, 1)"
                                                :disable-up="definitionIndex === 0"
                                                :disable-down="definitionIndex === group.definitions.length - 1"
                                            />
                                            <div class="font-semibold text-highlighted">
                                                {{
                                                    definition.label || (definition.displayKey ? t(definition.displayKey) : '')
                                                }}
                                            </div>
                                        </div>
                                        <div class="text-sm text-muted">
                                            <span class="font-mono">{{ definition.key }}</span> ·
                                            {{
                                                t(
                                                    `enums.properties.PropertyValueType.${PropertyValueType[definition.valueType]}`,
                                                )
                                            }}
                                        </div>
                                        <p v-if="definition.description" class="mt-1 text-sm text-muted">
                                            {{ definition.description }}
                                        </p>
                                    </div>

                                    <USwitch
                                        v-if="!isDeletedDefinition(definition)"
                                        class="shrink-0"
                                        :model-value="definition.enabled"
                                        :disabled="!canEdit || !canDisableDefinition(definition)"
                                        @update:model-value="
                                            (enabled) =>
                                                isSystemDefinition(definition)
                                                    ? toggleSystemDefinition(definition, enabled)
                                                    : setCustomDefinitionEnabled(definition, enabled)
                                        "
                                    />
                                </div>

                                <div class="flex flex-wrap items-center gap-2">
                                    <UBadge
                                        :color="
                                            isDeletedDefinition(definition)
                                                ? 'error'
                                                : isSystemDefinition(definition)
                                                  ? 'primary'
                                                  : 'neutral'
                                        "
                                        :label="
                                            isDeletedDefinition(definition)
                                                ? $t('components.settings.property_definitions.deleted')
                                                : isSystemDefinition(definition)
                                                  ? $t('components.settings.property_definitions.system')
                                                  : $t('components.settings.property_definitions.custom')
                                        "
                                    />
                                    <UBadge
                                        :color="definition.pinned ? 'warning' : 'neutral'"
                                        :icon="definition.pinned ? 'i-mdi-pin' : 'i-mdi-pin-off-outline'"
                                        :label="
                                            definition.pinned
                                                ? $t('components.settings.property_definitions.pinned')
                                                : $t('components.settings.property_definitions.not_pinned')
                                        "
                                    />
                                    <UBadge color="neutral" :label="displayModeLabel(definition)" />
                                    <UBadge
                                        :label="
                                            $t('components.settings.property_definitions.renderer_label', {
                                                renderer: rendererLabel(definition),
                                            })
                                        "
                                    />
                                    <UBadge
                                        v-if="definition.searchable"
                                        color="primary"
                                        :label="$t('components.settings.property_definitions.searchable')"
                                    />
                                    <UBadge
                                        v-if="isManagedDefinition(definition)"
                                        color="warning"
                                        :label="$t('components.settings.property_definitions.managed')"
                                    />
                                </div>

                                <div
                                    v-if="canEdit && (isCustomDefinition(definition) || isDeletedDefinition(definition))"
                                    class="flex justify-end gap-2"
                                >
                                    <template v-if="isCustomDefinition(definition) && !isDeletedDefinition(definition)">
                                        <UButton
                                            icon="i-mdi-pencil"
                                            color="neutral"
                                            variant="ghost"
                                            @click="openEditor(definition)"
                                        />
                                        <UButton
                                            icon="i-mdi-delete"
                                            color="error"
                                            variant="ghost"
                                            @click="deleteCustomDefinition(definition)"
                                        />
                                    </template>
                                    <UButton
                                        v-if="isCustomDefinition(definition) && isDeletedDefinition(definition)"
                                        icon="i-mdi-restore"
                                        color="primary"
                                        variant="ghost"
                                        @click="restoreCustomDefinition(definition)"
                                    />
                                </div>
                            </div>
                        </UCard>
                    </div>
                </div>
            </template>
        </template>

        <template #footer>
            <Pagination :status="status" :refresh="refresh" hide-buttons hide-text />
        </template>
    </UDashboardPanel>
</template>
