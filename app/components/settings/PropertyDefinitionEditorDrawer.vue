<script lang="ts" setup>
import type { Form } from '@nuxt/ui';
import { z } from 'zod';
import {
    PropertyDefinitionSource,
    PropertyDisplayMode,
    PropertyRenderer,
    PropertyBehavior,
    PropertyScope,
    PropertyValueFormat,
    PropertyValueType,
    type PropertyRendererConfig,
    type PropertyDefinition,
    type PropertyGroup,
} from '~~/gen/ts/resources/properties/properties';
import ResponsiveDrawerModal from '~/components/partials/ResponsiveDrawerModal.vue';

const props = defineProps<{
    definition?: PropertyDefinition;
    groups: PropertyGroup[];
    scope: PropertyScope;
    onSave: (definition: PropertyDefinition) => Promise<void>;
}>();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const { t, te } = useI18n();
const isEditing = computed(() => props.definition !== undefined);

const schema = z.object({
    key: z
        .string()
        .min(1)
        .max(64)
        .refine((key) => key.includes('.') && key.slice(key.indexOf('.') + 1).length > 0),
    label: z.string().min(1).max(128),
    description: z.string().max(512).optional(),
});

const formRef = useTemplateRef<Form<typeof schema>>('formRef');

const editor = reactive<PropertyDefinition>({
    scope: props.scope,
    key: 'vehicle.',
    valueType: PropertyValueType.STRING,
    source: PropertyDefinitionSource.CUSTOM,
    renderer: PropertyRenderer.DEFAULT,
    displayMode: PropertyDisplayMode.DEFAULT,
    enabled: true,
    searchable: false,
    pinned: false,
    implicit: false,
    groupKey: '',
    groupSortOrder: 0,
    sortOrder: 0,
    valueFormat: PropertyValueFormat.PLAIN,
    behavior: PropertyBehavior.GENERIC,
    unit: undefined,
    decimalPlaces: undefined,
    currency: undefined,
    label: '',
    description: undefined,
    displayKey: '',
    rendererConfig: undefined,
});

const rendererConfigDraft = reactive<{
    prefix: string;
    suffix: string;
    truncate: boolean;
    decimalPlaces?: number;
    unit?: string;
}>({ prefix: '', suffix: '', truncate: false });

function buildRendererConfig(): PropertyRendererConfig | undefined {
    if (editor.renderer === PropertyRenderer.STRING) {
        if (!rendererConfigDraft.prefix && !rendererConfigDraft.suffix && !rendererConfigDraft.truncate) return undefined;
        return {
            config: {
                oneofKind: 'string',
                string: {
                    prefix: rendererConfigDraft.prefix,
                    suffix: rendererConfigDraft.suffix,
                    truncate: rendererConfigDraft.truncate,
                },
            },
        };
    }
    if (editor.renderer === PropertyRenderer.INTEGER || editor.renderer === PropertyRenderer.NUMBER) {
        if (
            rendererConfigDraft.decimalPlaces === undefined &&
            !rendererConfigDraft.prefix &&
            !rendererConfigDraft.suffix &&
            !rendererConfigDraft.unit
        ) {
            return undefined;
        }
        const config = {
            decimalPlaces: rendererConfigDraft.decimalPlaces,
            prefix: rendererConfigDraft.prefix,
            suffix: rendererConfigDraft.suffix,
            unit: rendererConfigDraft.unit,
        };
        return {
            config:
                editor.renderer === PropertyRenderer.INTEGER
                    ? { oneofKind: 'integer', integer: config }
                    : { oneofKind: 'number', number: config },
        };
    }
    return undefined;
}

const { submit, isSubmitting, canSubmit } = useSubmitGuard(async () =>
    props.onSave({ ...editor, rendererConfig: buildRendererConfig() }),
);

const valueTypeOptions = [
    PropertyValueType.STRING,
    PropertyValueType.INTEGER,
    PropertyValueType.NUMBER,
    PropertyValueType.BOOL,
    PropertyValueType.TIMESTAMP,
    PropertyValueType.DATE,
].map((value) => ({ label: t(`enums.properties.PropertyValueType.${PropertyValueType[value]}`), value }));

const valueFormatOptions = computed(() => {
    const formats = [PropertyValueFormat.PLAIN];
    switch (editor.valueType) {
        case PropertyValueType.INTEGER:
            formats.push(
                PropertyValueFormat.MONEY,
                PropertyValueFormat.PERCENTAGE,
                PropertyValueFormat.DURATION,
                PropertyValueFormat.DISTANCE,
            );
            break;
        case PropertyValueType.NUMBER:
            formats.push(PropertyValueFormat.PERCENTAGE, PropertyValueFormat.DISTANCE);
            break;
        case PropertyValueType.TIMESTAMP:
            formats.push(PropertyValueFormat.DATETIME);
            break;
        case PropertyValueType.DATE:
            formats.push(PropertyValueFormat.DATE);
            break;
    }
    return formats.map((value) => ({ label: t(`enums.properties.PropertyValueFormat.${PropertyValueFormat[value]}`), value }));
});

const rendererOptions = computed(() => {
    const options = [
        { label: t('components.settings.property_definitions.renderer_default'), value: PropertyRenderer.DEFAULT },
    ];
    const rendererByType: Partial<Record<PropertyValueType, PropertyRenderer>> = {
        [PropertyValueType.STRING]: PropertyRenderer.STRING,
        [PropertyValueType.INTEGER]: PropertyRenderer.INTEGER,
        [PropertyValueType.NUMBER]: PropertyRenderer.NUMBER,
        [PropertyValueType.BOOL]: PropertyRenderer.BOOLEAN,
        [PropertyValueType.DATE]: PropertyRenderer.DATE,
        [PropertyValueType.TIMESTAMP]: PropertyRenderer.DATETIME,
    };
    const renderer = rendererByType[editor.valueType];
    if (renderer !== undefined) {
        options.push({
            label: t(`components.settings.property_definitions.renderer_${PropertyRenderer[renderer].toLowerCase()}`),
            value: renderer,
        });
    }
    return options;
});

// Group keys must start with a lowercase letter, so this cannot collide with a persisted group key.
const ungroupedGroupValue = '__ungrouped__';

const selectedGroupKey = computed({
    get: () => editor.groupKey || ungroupedGroupValue,
    set: (value: string) => {
        editor.groupKey = value === ungroupedGroupValue ? '' : value;
    },
});

const groupOptions = computed(() => [
    {
        label: t('components.settings.property_definitions.ungrouped'),
        value: ungroupedGroupValue,
    },
    ...props.groups
        .filter((group) => group.deletedAt === undefined)
        .map((group) => ({
            label:
                group.label ||
                (group.displayKey && te(group.displayKey)
                    ? t(group.displayKey)
                    : te(`components.settings.property_definitions.groups.${group.key}`)
                      ? t(`components.settings.property_definitions.groups.${group.key}`)
                      : group.key),
            value: group.key,
        })),
]);

const hasRendererConfig = computed(() =>
    [PropertyRenderer.STRING, PropertyRenderer.INTEGER, PropertyRenderer.NUMBER].includes(editor.renderer),
);

const hasNumberRendererConfig = computed(() => [PropertyRenderer.INTEGER, PropertyRenderer.NUMBER].includes(editor.renderer));

function scopeKeyPrefix(scope: PropertyScope): string {
    return scope === PropertyScope.USER ? 'user.' : 'vehicle.';
}

const keySuffix = computed({
    get: () => {
        const separatorIndex = editor.key.indexOf('.');
        return separatorIndex === -1 ? editor.key : editor.key.slice(separatorIndex + 1);
    },
    set: (value: string) => {
        editor.key = `${scopeKeyPrefix(props.scope)}${value}`;
    },
});

function resetEditor(): void {
    const definition = props.definition;
    Object.assign(editor, {
        scope: props.scope,
        key: definition?.key ?? scopeKeyPrefix(props.scope),
        valueType: definition?.valueType ?? PropertyValueType.STRING,
        source: PropertyDefinitionSource.CUSTOM,
        renderer: definition?.renderer ?? PropertyRenderer.DEFAULT,
        displayMode: definition?.displayMode ?? PropertyDisplayMode.DEFAULT,
        enabled: definition?.enabled ?? true,
        searchable: definition?.searchable ?? false,
        pinned: definition?.pinned ?? false,
        implicit: definition?.implicit ?? false,
        groupKey: definition?.groupKey ?? '',
        groupSortOrder: definition?.groupSortOrder ?? 0,
        sortOrder: definition?.sortOrder ?? 0,
        valueFormat: definition?.valueFormat || PropertyValueFormat.PLAIN,
        behavior: PropertyBehavior.GENERIC,
        unit: definition?.unit,
        decimalPlaces: definition?.decimalPlaces,
        currency: definition?.currency,
        label: definition?.label ?? '',
        description: definition?.description,
        displayKey: '',
        rendererConfig: definition?.rendererConfig,
    });
    Object.assign(rendererConfigDraft, {
        prefix: '',
        suffix: '',
        truncate: false,
        decimalPlaces: undefined,
        unit: undefined,
    });
    const config = definition?.rendererConfig?.config;
    if (config?.oneofKind === 'string') Object.assign(rendererConfigDraft, config.string);
    if (config?.oneofKind === 'integer') Object.assign(rendererConfigDraft, config.integer);
    if (config?.oneofKind === 'number') Object.assign(rendererConfigDraft, config.number);
}

const { hasUnsavedChanges, confirmLeave, syncSnapshot } = useSnapshotChanges(() => ({
    editor: { ...editor },
    rendererConfigDraft: { ...rendererConfigDraft },
}));

watch([() => props.definition, () => props.scope], resetEditor, { immediate: true });
watch([() => props.definition, () => props.scope], () => syncSnapshot(), { immediate: true });
watch(
    () => editor.valueType,
    () => {
        if (!valueFormatOptions.value.some((option) => option.value === editor.valueFormat)) {
            editor.valueFormat = PropertyValueFormat.PLAIN;
        }
    },
);
watch(
    () => editor.renderer,
    () => {
        if (!rendererOptions.value.some((option) => option.value === editor.renderer))
            editor.renderer = PropertyRenderer.DEFAULT;
    },
);

async function close(): Promise<void> {
    if (!canSubmit.value) return;

    if (hasUnsavedChanges.value && !(await confirmLeave())) return;

    emit('close', false);
}
</script>

<template>
    <ResponsiveDrawerModal
        :title="t('components.settings.property_definitions.editor_title')"
        :close="false"
        :dismissible="!hasUnsavedChanges"
    >
        <template #header>
            <div class="flex w-full items-center justify-between gap-1.5">
                <h3 class="font-semibold text-highlighted">
                    {{ t('components.settings.property_definitions.editor_title') }}
                </h3>
                <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-close"
                    :disabled="!canSubmit"
                    :aria-label="$t('common.close', 1)"
                    @click="close"
                />
            </div>
        </template>

        <template #body>
            <UForm ref="formRef" class="flex flex-col gap-4" :schema="schema" :state="editor" @submit="submit">
                <UFormField
                    name="key"
                    :label="t('components.settings.property_definitions.key')"
                    :description="t('components.settings.property_definitions.key_description')"
                    required
                >
                    <div class="flex items-center gap-1">
                        <span class="shrink-0 font-mono text-muted">{{ scopeKeyPrefix(props.scope) }}</span>
                        <UInput
                            v-model="keySuffix"
                            class="min-w-0 flex-1 font-mono"
                            :maxlength="64 - scopeKeyPrefix(props.scope).length"
                            :disabled="isEditing"
                        />
                    </div>
                </UFormField>

                <UFormField
                    name="label"
                    :label="t('components.settings.property_definitions.label')"
                    :description="t('components.settings.property_definitions.label_description')"
                    required
                >
                    <UInput v-model="editor.label" class="w-full" />
                </UFormField>

                <UFormField
                    name="description"
                    :label="t('components.settings.property_definitions.description')"
                    :description="t('components.settings.property_definitions.description_description')"
                >
                    <UTextarea v-model="editor.description" class="w-full" />
                </UFormField>

                <UFormField
                    name="groupKey"
                    :label="t('components.settings.property_definitions.group')"
                    :description="t('components.settings.property_definitions.group_description')"
                >
                    <USelectMenu
                        v-model="selectedGroupKey"
                        class="w-full"
                        :items="groupOptions"
                        label-key="label"
                        value-key="value"
                    />
                </UFormField>

                <UFormField
                    name="valueType"
                    :label="t('components.settings.property_definitions.value_type')"
                    :description="t('components.settings.property_definitions.value_type_description')"
                >
                    <USelectMenu
                        v-model="editor.valueType"
                        class="w-full"
                        :items="valueTypeOptions"
                        label-key="label"
                        value-key="value"
                        :disabled="isEditing"
                    />
                </UFormField>

                <UFormField
                    name="valueFormat"
                    :label="t('components.settings.property_definitions.value_format')"
                    :description="t('components.settings.property_definitions.value_format_description')"
                >
                    <USelectMenu
                        v-model="editor.valueFormat"
                        class="w-full"
                        :items="valueFormatOptions"
                        label-key="label"
                        value-key="value"
                    />
                </UFormField>

                <UFormField
                    name="renderer"
                    :label="t('components.settings.property_definitions.renderer')"
                    :description="t('components.settings.property_definitions.renderer_description')"
                >
                    <USelectMenu
                        v-model="editor.renderer"
                        class="w-full"
                        :items="rendererOptions"
                        label-key="label"
                        value-key="value"
                    />
                </UFormField>

                <template v-if="hasRendererConfig">
                    <UFormField
                        :label="t('components.settings.property_definitions.renderer_prefix')"
                        :description="t('components.settings.property_definitions.renderer_prefix_description')"
                    >
                        <UInput v-model="rendererConfigDraft.prefix" class="w-full" />
                    </UFormField>
                    <UFormField
                        :label="t('components.settings.property_definitions.renderer_suffix')"
                        :description="t('components.settings.property_definitions.renderer_suffix_description')"
                    >
                        <UInput v-model="rendererConfigDraft.suffix" class="w-full" />
                    </UFormField>
                    <UFormField
                        v-if="hasNumberRendererConfig"
                        :label="t('components.settings.property_definitions.renderer_unit')"
                        :description="t('components.settings.property_definitions.renderer_unit_description')"
                    >
                        <UInput v-model="rendererConfigDraft.unit" class="w-full" />
                    </UFormField>
                    <UFormField
                        v-if="hasNumberRendererConfig"
                        :label="t('components.settings.property_definitions.renderer_decimal_places')"
                        :description="t('components.settings.property_definitions.renderer_decimal_places_description')"
                    >
                        <UInput
                            v-model.number="rendererConfigDraft.decimalPlaces"
                            class="w-full"
                            type="number"
                            min="0"
                            max="9"
                        />
                    </UFormField>
                    <UFormField
                        v-if="editor.renderer === PropertyRenderer.STRING"
                        :label="t('components.settings.property_definitions.renderer_truncate')"
                        :description="t('components.settings.property_definitions.renderer_truncate_description')"
                    >
                        <USwitch v-model="rendererConfigDraft.truncate" />
                    </UFormField>
                </template>

                <UFormField
                    v-if="editor.valueFormat === PropertyValueFormat.MONEY"
                    :label="t('components.settings.property_definitions.currency')"
                    :description="t('components.settings.property_definitions.currency_description')"
                >
                    <UInput v-model="editor.currency" class="w-full uppercase" maxlength="3" />
                </UFormField>

                <UFormField
                    v-if="editor.valueFormat === PropertyValueFormat.DISTANCE"
                    :label="t('components.settings.property_definitions.unit')"
                    :description="t('components.settings.property_definitions.unit_description')"
                >
                    <UInput v-model="editor.unit" class="w-full" maxlength="32" />
                </UFormField>

                <UFormField
                    v-if="editor.valueType === PropertyValueType.NUMBER || editor.valueFormat === PropertyValueFormat.MONEY"
                    :label="t('components.settings.property_definitions.decimal_places')"
                    :description="t('components.settings.property_definitions.decimal_places_description')"
                >
                    <UInput v-model.number="editor.decimalPlaces" class="w-full" type="number" min="0" max="9" />
                </UFormField>

                <div class="flex flex-col gap-2">
                    <UFormField
                        name="enabled"
                        :label="t('components.settings.property_definitions.enabled')"
                        :description="t('components.settings.property_definitions.enabled_description')"
                    >
                        <USwitch v-model="editor.enabled" />
                    </UFormField>
                    <UFormField
                        name="searchable"
                        :label="t('components.settings.property_definitions.searchable')"
                        :description="t('components.settings.property_definitions.searchable_description')"
                    >
                        <USwitch v-model="editor.searchable" />
                    </UFormField>
                    <UFormField
                        name="pinned"
                        :label="t('components.settings.property_definitions.pinned')"
                        :description="t('components.settings.property_definitions.pinned_description')"
                    >
                        <USwitch v-model="editor.pinned" />
                    </UFormField>
                    <UFormField
                        name="implicit"
                        :label="t('components.settings.property_definitions.implicit')"
                        :description="t('components.settings.property_definitions.implicit_description')"
                    >
                        <USwitch v-model="editor.implicit" />
                    </UFormField>
                </div>
            </UForm>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton
                    class="flex-1"
                    color="neutral"
                    block
                    :disabled="!canSubmit"
                    :label="$t('common.close', 1)"
                    @click="close"
                />
                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="isSubmitting"
                    :label="$t('common.save')"
                    @click="formRef?.submit()"
                />
            </UFieldGroup>
        </template>
    </ResponsiveDrawerModal>
</template>
