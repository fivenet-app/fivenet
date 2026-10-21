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
    key: z.string().min(1).max(64),
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

const groupOptions = computed(() => [
    {
        label: t('components.settings.property_definitions.ungrouped'),
        value: '',
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
                <UFormField name="key" :label="t('components.settings.property_definitions.key')" required>
                    <UInput v-model="editor.key" class="w-full" :disabled="isEditing" />
                </UFormField>

                <UFormField name="label" :label="t('components.settings.property_definitions.label')" required>
                    <UInput v-model="editor.label" class="w-full" />
                </UFormField>

                <UFormField name="description" :label="t('components.settings.property_definitions.description')">
                    <UTextarea v-model="editor.description" class="w-full" />
                </UFormField>

                <UFormField name="groupKey" :label="t('components.settings.property_definitions.group')">
                    <USelectMenu
                        v-model="editor.groupKey"
                        class="w-full"
                        :items="groupOptions"
                        label-key="label"
                        value-key="value"
                    />
                </UFormField>

                <UFormField name="valueType" :label="t('components.settings.property_definitions.value_type')">
                    <USelectMenu
                        v-model="editor.valueType"
                        class="w-full"
                        :items="valueTypeOptions"
                        label-key="label"
                        value-key="value"
                        :disabled="isEditing"
                    />
                </UFormField>

                <UFormField name="valueFormat" :label="t('components.settings.property_definitions.value_format')">
                    <USelectMenu
                        v-model="editor.valueFormat"
                        class="w-full"
                        :items="valueFormatOptions"
                        label-key="label"
                        value-key="value"
                    />
                </UFormField>

                <UFormField name="renderer" :label="t('components.settings.property_definitions.renderer')">
                    <USelectMenu
                        v-model="editor.renderer"
                        class="w-full"
                        :items="rendererOptions"
                        label-key="label"
                        value-key="value"
                    />
                </UFormField>

                <template v-if="hasRendererConfig">
                    <UFormField :label="t('components.settings.property_definitions.renderer_prefix')">
                        <UInput v-model="rendererConfigDraft.prefix" class="w-full" />
                    </UFormField>
                    <UFormField :label="t('components.settings.property_definitions.renderer_suffix')">
                        <UInput v-model="rendererConfigDraft.suffix" class="w-full" />
                    </UFormField>
                    <UFormField
                        v-if="hasNumberRendererConfig"
                        :label="t('components.settings.property_definitions.renderer_unit')"
                    >
                        <UInput v-model="rendererConfigDraft.unit" class="w-full" />
                    </UFormField>
                    <UFormField
                        v-if="hasNumberRendererConfig"
                        :label="t('components.settings.property_definitions.renderer_decimal_places')"
                    >
                        <UInput
                            v-model.number="rendererConfigDraft.decimalPlaces"
                            class="w-full"
                            type="number"
                            min="0"
                            max="9"
                        />
                    </UFormField>
                    <UCheckbox
                        v-if="editor.renderer === PropertyRenderer.STRING"
                        v-model="rendererConfigDraft.truncate"
                        :label="t('components.settings.property_definitions.renderer_truncate')"
                    />
                </template>

                <UFormField
                    v-if="editor.valueFormat === PropertyValueFormat.MONEY"
                    :label="t('components.settings.property_definitions.currency')"
                >
                    <UInput v-model="editor.currency" class="w-full uppercase" maxlength="3" />
                </UFormField>

                <UFormField
                    v-if="editor.valueFormat === PropertyValueFormat.DISTANCE"
                    :label="t('components.settings.property_definitions.unit')"
                >
                    <UInput v-model="editor.unit" class="w-full" maxlength="32" />
                </UFormField>

                <UFormField
                    v-if="editor.valueType === PropertyValueType.NUMBER || editor.valueFormat === PropertyValueFormat.MONEY"
                    :label="t('components.settings.property_definitions.decimal_places')"
                >
                    <UInput v-model.number="editor.decimalPlaces" class="w-full" type="number" min="0" max="9" />
                </UFormField>

                <div class="flex flex-col gap-2">
                    <UCheckbox v-model="editor.enabled" :label="t('components.settings.property_definitions.enabled')" />
                    <UCheckbox v-model="editor.searchable" :label="t('components.settings.property_definitions.searchable')" />
                    <UCheckbox v-model="editor.pinned" :label="t('components.settings.property_definitions.pinned')" />
                    <UCheckbox v-model="editor.implicit" :label="t('components.settings.property_definitions.implicit')" />
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
