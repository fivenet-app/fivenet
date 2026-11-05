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
    PropertyStringAllowedValuesMode,
    PropertyExpiryAction,
    PropertyReasonMode,
    type PropertyRendererConfig,
    type PropertyDefinition,
    type PropertyGroup,
    type PropertyValueConstraints,
} from '~~/gen/ts/resources/properties/properties';
import ResponsiveDrawerModal from '~/components/partials/ResponsiveDrawerModal.vue';
import InputDatePicker from '~/components/partials/InputDatePicker.vue';
import { toDate, toTimestamp } from '~/utils/time';

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
    statePolicy: undefined,
    reasonMode: PropertyReasonMode.REQUIRED,
    valueConstraints: undefined,
});

const rendererConfigDraft = reactive<{
    prefix: string;
    suffix: string;
    truncate: boolean;
    decimalPlaces?: number;
    unit?: string;
}>({ prefix: '', suffix: '', truncate: false });

// Constraints deliberately live in a presentation draft: date pickers use
// Date while the definition API uses LocalDate/Timestamp scalar messages.
const constraintsDraft = reactive<{
    stringMinLength?: number;
    stringMaxLength?: number;
    stringAllowedValues: string;
    stringAllowedValuesMode: PropertyStringAllowedValuesMode;
    integerMin?: number;
    integerMax?: number;
    numberMin?: number;
    numberMax?: number;
    timestampMin?: Date;
    timestampMax?: Date;
    dateMin?: Date;
    dateMax?: Date;
}>({
    stringAllowedValues: '',
    stringAllowedValuesMode: PropertyStringAllowedValuesMode.STRICT,
});

function dateToLocalDate(date: Date | undefined) {
    if (!date) return undefined;
    return { year: date.getFullYear(), month: date.getMonth() + 1, day: date.getDate() };
}

function localDateToDate(date: PropertyValueConstraints['dateMin']): Date | undefined {
    if (!date) return undefined;
    return new Date(date.year, date.month - 1, date.day);
}

function allowedValues(): string[] {
    return constraintsDraft.stringAllowedValues
        .replaceAll('\r\n', '\n')
        .split('\n')
        .filter((value) => value.length > 0);
}

const allowedValuesLimit = 50;

const reasonModeOptions = [PropertyReasonMode.REQUIRED, PropertyReasonMode.OPTIONAL, PropertyReasonMode.NONE].map((value) => ({
    label: t(`components.settings.property_definitions.reason_mode_${PropertyReasonMode[value].toLowerCase()}`),
    value,
}));

const allowedValuesModeOptions = [PropertyStringAllowedValuesMode.STRICT, PropertyStringAllowedValuesMode.SUGGEST].map(
    (value) => ({
        label: t(
            `components.settings.property_definitions.constraint_allowed_values_mode_${PropertyStringAllowedValuesMode[value].toLowerCase()}`,
        ),
        value,
    }),
);

const constraintErrors = computed(() => {
    const errors: Partial<Record<'stringLength' | 'allowedValues' | 'integer' | 'number' | 'timestamp' | 'date', string>> = {};
    if (
        constraintsDraft.stringMinLength !== undefined &&
        constraintsDraft.stringMaxLength !== undefined &&
        constraintsDraft.stringMinLength > constraintsDraft.stringMaxLength
    ) {
        errors.stringLength = t('components.settings.property_definitions.constraint_min_max_error');
    }
    const values = allowedValues();
    if (new Set(values).size !== values.length) {
        errors.allowedValues = t('components.settings.property_definitions.constraint_duplicate_values_error');
    } else if (values.length > allowedValuesLimit) {
        errors.allowedValues = t('components.settings.property_definitions.constraint_allowed_values_limit_error', {
            limit: allowedValuesLimit,
        });
    }
    if (
        constraintsDraft.integerMin !== undefined &&
        constraintsDraft.integerMax !== undefined &&
        constraintsDraft.integerMin > constraintsDraft.integerMax
    ) {
        errors.integer = t('components.settings.property_definitions.constraint_min_max_error');
    }
    if (
        constraintsDraft.numberMin !== undefined &&
        constraintsDraft.numberMax !== undefined &&
        constraintsDraft.numberMin > constraintsDraft.numberMax
    ) {
        errors.number = t('components.settings.property_definitions.constraint_min_max_error');
    }
    if (
        constraintsDraft.timestampMin !== undefined &&
        constraintsDraft.timestampMax !== undefined &&
        constraintsDraft.timestampMin > constraintsDraft.timestampMax
    ) {
        errors.timestamp = t('components.settings.property_definitions.constraint_min_max_error');
    }
    if (
        constraintsDraft.dateMin !== undefined &&
        constraintsDraft.dateMax !== undefined &&
        constraintsDraft.dateMin > constraintsDraft.dateMax
    ) {
        errors.date = t('components.settings.property_definitions.constraint_min_max_error');
    }
    return errors;
});

const hasConstraintErrors = computed(() => Object.keys(constraintErrors.value).length > 0);

function buildValueConstraints(): PropertyValueConstraints | undefined {
    switch (editor.valueType) {
        case PropertyValueType.STRING: {
            const stringAllowedValues = allowedValues();
            if (
                constraintsDraft.stringMinLength === undefined &&
                constraintsDraft.stringMaxLength === undefined &&
                stringAllowedValues.length === 0
            ) {
                return undefined;
            }
            return {
                stringMinLength: constraintsDraft.stringMinLength,
                stringMaxLength: constraintsDraft.stringMaxLength,
                stringAllowedValues,
                stringAllowedValuesMode:
                    stringAllowedValues.length > 0
                        ? constraintsDraft.stringAllowedValuesMode
                        : PropertyStringAllowedValuesMode.UNSPECIFIED,
            };
        }
        case PropertyValueType.INTEGER:
            return constraintsDraft.integerMin === undefined && constraintsDraft.integerMax === undefined
                ? undefined
                : {
                      stringAllowedValues: [],
                      stringAllowedValuesMode: PropertyStringAllowedValuesMode.UNSPECIFIED,
                      integerMin: constraintsDraft.integerMin,
                      integerMax: constraintsDraft.integerMax,
                  };
        case PropertyValueType.NUMBER:
            return constraintsDraft.numberMin === undefined && constraintsDraft.numberMax === undefined
                ? undefined
                : {
                      stringAllowedValues: [],
                      stringAllowedValuesMode: PropertyStringAllowedValuesMode.UNSPECIFIED,
                      numberMin: constraintsDraft.numberMin,
                      numberMax: constraintsDraft.numberMax,
                  };
        case PropertyValueType.TIMESTAMP:
            return constraintsDraft.timestampMin === undefined && constraintsDraft.timestampMax === undefined
                ? undefined
                : {
                      stringAllowedValues: [],
                      stringAllowedValuesMode: PropertyStringAllowedValuesMode.UNSPECIFIED,
                      timestampMin: toTimestamp(constraintsDraft.timestampMin),
                      timestampMax: toTimestamp(constraintsDraft.timestampMax),
                  };
        case PropertyValueType.DATE:
            return constraintsDraft.dateMin === undefined && constraintsDraft.dateMax === undefined
                ? undefined
                : {
                      stringAllowedValues: [],
                      stringAllowedValuesMode: PropertyStringAllowedValuesMode.UNSPECIFIED,
                      dateMin: dateToLocalDate(constraintsDraft.dateMin),
                      dateMax: dateToLocalDate(constraintsDraft.dateMax),
                  };
        default:
            return undefined;
    }
}

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

const { submit, isSubmitting, canSubmit } = useSubmitGuard(async () => {
    if (hasConstraintErrors.value) return;
    await props.onSave({ ...editor, rendererConfig: buildRendererConfig(), valueConstraints: buildValueConstraints() });
});

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
        statePolicy: definition?.statePolicy,
        reasonMode: definition?.reasonMode ?? PropertyReasonMode.REQUIRED,
        valueConstraints: definition?.valueConstraints,
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
    const constraints = definition?.valueConstraints;
    Object.assign(constraintsDraft, {
        stringMinLength: constraints?.stringMinLength,
        stringMaxLength: constraints?.stringMaxLength,
        stringAllowedValues: constraints?.stringAllowedValues.join('\n') ?? '',
        stringAllowedValuesMode:
            constraints?.stringAllowedValuesMode === PropertyStringAllowedValuesMode.SUGGEST
                ? PropertyStringAllowedValuesMode.SUGGEST
                : PropertyStringAllowedValuesMode.STRICT,
        integerMin: constraints?.integerMin,
        integerMax: constraints?.integerMax,
        numberMin: constraints?.numberMin,
        numberMax: constraints?.numberMax,
        timestampMin: constraints?.timestampMin ? toDate(constraints.timestampMin) : undefined,
        timestampMax: constraints?.timestampMax ? toDate(constraints.timestampMax) : undefined,
        dateMin: localDateToDate(constraints?.dateMin),
        dateMax: localDateToDate(constraints?.dateMax),
    });
}

const { hasUnsavedChanges, confirmLeave, syncSnapshot } = useSnapshotChanges(() => ({
    editor: { ...editor },
    rendererConfigDraft: { ...rendererConfigDraft },
    constraintsDraft: { ...constraintsDraft },
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
        :title="$t('components.settings.property_definitions.editor_title')"
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
                <USeparator :label="$t('components.settings.property_definitions.section_identity')" />

                <UFormField
                    name="key"
                    :label="$t('components.settings.property_definitions.key')"
                    :description="$t('components.settings.property_definitions.key_description')"
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
                    :label="$t('components.settings.property_definitions.label')"
                    :description="$t('components.settings.property_definitions.label_description')"
                    required
                >
                    <UInput v-model="editor.label" class="w-full" />
                </UFormField>

                <UFormField
                    name="description"
                    :label="$t('components.settings.property_definitions.description')"
                    :description="$t('components.settings.property_definitions.description_description')"
                >
                    <UTextarea v-model="editor.description" class="w-full" />
                </UFormField>

                <USeparator :label="$t('components.settings.property_definitions.section_organization')" />
                <UFormField
                    name="groupKey"
                    :label="$t('components.settings.property_definitions.group')"
                    :description="$t('components.settings.property_definitions.group_description')"
                >
                    <USelectMenu
                        v-model="selectedGroupKey"
                        class="w-full"
                        :items="groupOptions"
                        label-key="label"
                        value-key="value"
                    />
                </UFormField>

                <div class="space-y-2">
                    <USeparator :label="$t('components.settings.property_definitions.section_value')" />
                </div>
                <UFormField
                    name="valueType"
                    :label="$t('components.settings.property_definitions.value_type')"
                    :description="$t('components.settings.property_definitions.value_type_description')"
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

                <template v-if="editor.valueType === PropertyValueType.STRING">
                    <div class="grid gap-3 sm:grid-cols-2">
                        <UFormField
                            :label="$t('components.settings.property_definitions.constraint_min_characters')"
                            :error="constraintErrors.stringLength"
                        >
                            <UInputNumber v-model="constraintsDraft.stringMinLength" class="w-full" :min="0" :max="2048" />
                        </UFormField>
                        <UFormField
                            :label="$t('components.settings.property_definitions.constraint_max_characters')"
                            :error="constraintErrors.stringLength"
                        >
                            <UInputNumber v-model="constraintsDraft.stringMaxLength" class="w-full" :min="0" :max="2048" />
                        </UFormField>
                    </div>
                    <UFormField
                        :label="$t('components.settings.property_definitions.constraint_allowed_values')"
                        :description="$t('components.settings.property_definitions.constraint_allowed_values_description')"
                        :error="constraintErrors.allowedValues"
                    >
                        <UTextarea v-model="constraintsDraft.stringAllowedValues" class="w-full" :rows="5" />
                    </UFormField>
                    <UFormField
                        v-if="allowedValues().length > 0"
                        :label="$t('components.settings.property_definitions.constraint_allowed_values_mode')"
                    >
                        <USelectMenu
                            v-model="constraintsDraft.stringAllowedValuesMode"
                            class="w-full"
                            :items="allowedValuesModeOptions"
                            label-key="label"
                            value-key="value"
                        />
                    </UFormField>
                </template>

                <template v-else-if="editor.valueType === PropertyValueType.INTEGER">
                    <div class="grid gap-3 sm:grid-cols-2">
                        <UFormField
                            :label="$t('components.settings.property_definitions.constraint_minimum')"
                            :error="constraintErrors.integer"
                        >
                            <UInputNumber v-model="constraintsDraft.integerMin" class="w-full" :step="1" />
                        </UFormField>
                        <UFormField
                            :label="$t('components.settings.property_definitions.constraint_maximum')"
                            :error="constraintErrors.integer"
                        >
                            <UInputNumber v-model="constraintsDraft.integerMax" class="w-full" :step="1" />
                        </UFormField>
                    </div>
                </template>

                <template v-else-if="editor.valueType === PropertyValueType.NUMBER">
                    <div class="grid gap-3 sm:grid-cols-2">
                        <UFormField
                            :label="$t('components.settings.property_definitions.constraint_minimum')"
                            :error="constraintErrors.number"
                        >
                            <UInputNumber v-model="constraintsDraft.numberMin" class="w-full" :step="0.01" />
                        </UFormField>
                        <UFormField
                            :label="$t('components.settings.property_definitions.constraint_maximum')"
                            :error="constraintErrors.number"
                        >
                            <UInputNumber v-model="constraintsDraft.numberMax" class="w-full" :step="0.01" />
                        </UFormField>
                    </div>
                </template>

                <template v-else-if="editor.valueType === PropertyValueType.DATE">
                    <div class="grid gap-3 sm:grid-cols-2">
                        <UFormField
                            :label="$t('components.settings.property_definitions.constraint_earliest_date')"
                            :error="constraintErrors.date"
                        >
                            <InputDatePicker v-model="constraintsDraft.dateMin" class="w-full" clearable />
                        </UFormField>
                        <UFormField
                            :label="$t('components.settings.property_definitions.constraint_latest_date')"
                            :error="constraintErrors.date"
                        >
                            <InputDatePicker v-model="constraintsDraft.dateMax" class="w-full" clearable />
                        </UFormField>
                    </div>
                </template>

                <template v-else-if="editor.valueType === PropertyValueType.TIMESTAMP">
                    <div class="grid gap-3 sm:grid-cols-2">
                        <UFormField
                            :label="$t('components.settings.property_definitions.constraint_earliest_instant')"
                            :error="constraintErrors.timestamp"
                        >
                            <InputDatePicker v-model="constraintsDraft.timestampMin" class="w-full" clearable time />
                        </UFormField>
                        <UFormField
                            :label="$t('components.settings.property_definitions.constraint_latest_instant')"
                            :error="constraintErrors.timestamp"
                        >
                            <InputDatePicker v-model="constraintsDraft.timestampMax" class="w-full" clearable time />
                        </UFormField>
                    </div>
                </template>

                <UFormField
                    name="valueFormat"
                    :label="$t('components.settings.property_definitions.value_format')"
                    :description="$t('components.settings.property_definitions.value_format_description')"
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
                    v-if="editor.valueFormat === PropertyValueFormat.MONEY"
                    :label="$t('components.settings.property_definitions.currency')"
                    :description="$t('components.settings.property_definitions.currency_description')"
                >
                    <UInput v-model="editor.currency" class="w-full uppercase" maxlength="3" />
                </UFormField>

                <UFormField
                    v-if="editor.valueFormat === PropertyValueFormat.DISTANCE"
                    :label="$t('components.settings.property_definitions.unit')"
                    :description="$t('components.settings.property_definitions.unit_description')"
                >
                    <UInput v-model="editor.unit" class="w-full" maxlength="32" />
                </UFormField>

                <UFormField
                    v-if="editor.valueType === PropertyValueType.NUMBER || editor.valueFormat === PropertyValueFormat.MONEY"
                    :label="$t('components.settings.property_definitions.decimal_places')"
                    :description="$t('components.settings.property_definitions.decimal_places_description')"
                >
                    <UInputNumber v-model="editor.decimalPlaces" class="w-full" :min="0" :max="9" />
                </UFormField>

                <USeparator :label="$t('components.settings.property_definitions.section_presentation')" />
                <UFormField
                    name="renderer"
                    :label="$t('components.settings.property_definitions.renderer')"
                    :description="$t('components.settings.property_definitions.renderer_description')"
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
                        :label="$t('components.settings.property_definitions.renderer_prefix')"
                        :description="$t('components.settings.property_definitions.renderer_prefix_description')"
                    >
                        <UInput v-model="rendererConfigDraft.prefix" class="w-full" />
                    </UFormField>
                    <UFormField
                        :label="$t('components.settings.property_definitions.renderer_suffix')"
                        :description="$t('components.settings.property_definitions.renderer_suffix_description')"
                    >
                        <UInput v-model="rendererConfigDraft.suffix" class="w-full" />
                    </UFormField>
                    <UFormField
                        v-if="hasNumberRendererConfig"
                        :label="$t('components.settings.property_definitions.renderer_unit')"
                        :description="$t('components.settings.property_definitions.renderer_unit_description')"
                    >
                        <UInput v-model="rendererConfigDraft.unit" class="w-full" />
                    </UFormField>
                    <UFormField
                        v-if="hasNumberRendererConfig"
                        :label="$t('components.settings.property_definitions.renderer_decimal_places')"
                        :description="$t('components.settings.property_definitions.renderer_decimal_places_description')"
                    >
                        <UInputNumber v-model="rendererConfigDraft.decimalPlaces" class="w-full" :min="0" :max="9" />
                    </UFormField>
                    <UFormField
                        v-if="editor.renderer === PropertyRenderer.STRING"
                        :label="$t('components.settings.property_definitions.renderer_truncate')"
                        :description="$t('components.settings.property_definitions.renderer_truncate_description')"
                    >
                        <USwitch v-model="rendererConfigDraft.truncate" />
                    </UFormField>
                </template>

                <USeparator :label="$t('components.settings.property_definitions.section_lifecycle')" />
                <div class="flex flex-col gap-2">
                    <UFormField
                        :label="$t('components.settings.property_definitions.state_policy')"
                        :description="$t('components.settings.property_definitions.state_policy_description')"
                    >
                        <USwitch
                            :model-value="editor.statePolicy?.enabled ?? false"
                            @update:model-value="
                                (enabled) =>
                                    (editor.statePolicy = enabled
                                        ? {
                                              enabled: true,
                                              supportsExpiry: false,
                                              expiryAction: PropertyExpiryAction.UNSPECIFIED,
                                          }
                                        : undefined)
                            "
                        />
                    </UFormField>
                    <template v-if="editor.statePolicy?.enabled">
                        <UFormField
                            :label="$t('components.settings.property_definitions.state_policy_expiry')"
                            :description="$t('components.settings.property_definitions.state_policy_expiry_description')"
                        >
                            <USwitch
                                v-model="editor.statePolicy.supportsExpiry"
                                @update:model-value="
                                    (enabled) =>
                                        (editor.statePolicy!.expiryAction = enabled
                                            ? PropertyExpiryAction.CLEAR_VALUE
                                            : PropertyExpiryAction.UNSPECIFIED)
                                "
                            />
                        </UFormField>
                    </template>
                    <UFormField
                        :label="$t('components.settings.property_definitions.reason_mode')"
                        :description="$t('components.settings.property_definitions.reason_mode_description')"
                    >
                        <USelectMenu
                            v-model="editor.reasonMode"
                            class="w-full"
                            :items="reasonModeOptions"
                            label-key="label"
                            value-key="value"
                        />
                    </UFormField>

                    <USeparator class="pt-2" :label="$t('components.settings.property_definitions.section_availability')" />
                    <UFormField
                        name="enabled"
                        :label="$t('components.settings.property_definitions.enabled')"
                        :description="$t('components.settings.property_definitions.enabled_description')"
                    >
                        <USwitch v-model="editor.enabled" />
                    </UFormField>
                    <UFormField
                        name="searchable"
                        :label="$t('components.settings.property_definitions.searchable')"
                        :description="$t('components.settings.property_definitions.searchable_description')"
                    >
                        <USwitch v-model="editor.searchable" />
                    </UFormField>
                    <UFormField
                        name="pinned"
                        :label="$t('components.settings.property_definitions.pinned')"
                        :description="$t('components.settings.property_definitions.pinned_description')"
                    >
                        <USwitch v-model="editor.pinned" />
                    </UFormField>
                    <UFormField
                        name="implicit"
                        :label="$t('components.settings.property_definitions.implicit')"
                        :description="$t('components.settings.property_definitions.implicit_description')"
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
                    :disabled="!canSubmit || hasConstraintErrors"
                    :loading="isSubmitting"
                    :label="$t('common.save')"
                    @click="formRef?.submit()"
                />
            </UFieldGroup>
        </template>
    </ResponsiveDrawerModal>
</template>
