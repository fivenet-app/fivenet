<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import ResponsiveDrawerModal from '~/components/partials/ResponsiveDrawerModal.vue';
import { getCitizensCitizensClient, getVehiclesVehiclesClient } from '~~/gen/ts/clients';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import InputDatePicker from '~/components/partials/InputDatePicker.vue';
import {
    PropertyScope,
    PropertyReasonMode,
    PropertyValueType,
    type Properties,
    type Property,
    type PropertyMutation,
    type PropertyValue,
} from '~~/gen/ts/resources/properties/properties';
import { toDate, toTimestamp } from '~/utils/time';
import { stringPropertyInputMode } from '~/utils/propertyValueInput';

const props = defineProps<{
    property: Property;
    scope: PropertyScope;
    userId?: number;
    plate?: string;
}>();

const emit = defineEmits<{
    (e: 'close', value: boolean): void;
    (e: 'update:properties', value: Properties | undefined): void;
}>();

const { t } = useI18n();
const { propertyLabel } = usePropertyRenderer();
const notifications = useNotificationsStore();
const citizensClient = await getCitizensCitizensClient();
const vehiclesClient = await getVehiclesVehiclesClient();
const overlay = useOverlay();
const resetConfirmModal = overlay.create(ConfirmModal);

const reasonMode = computed(() => props.property.definition?.reasonMode ?? PropertyReasonMode.REQUIRED);
const schema = z.object({ reason: z.coerce.string().max(255) }).superRefine((values, ctx) => {
    if (reasonMode.value === PropertyReasonMode.REQUIRED && values.reason.trim().length < 3) {
        ctx.addIssue({ code: z.ZodIssueCode.custom, path: ['reason'], message: t('components.properties.reason_required') });
    }
});
type Schema = z.output<typeof schema>;

const state = reactive({
    reason: '',
    stringValue: '',
    numberValue: 0,
    boolValue: false,
    dateValue: undefined as Date | undefined,
    timestampValue: undefined as Date | undefined,
    reset: false,
});

const valueType = computed(() => props.property.definition?.valueType ?? PropertyValueType.UNSPECIFIED);
const allowedStringValues = computed(() => props.property.definition?.valueConstraints?.stringAllowedValues ?? []);
const stringInputMode = computed(() => stringPropertyInputMode(props.property.definition));
const stringValueOptions = computed(() => allowedStringValues.value.map((value) => ({ label: value, value })));
const stringValueSuggestionsId = computed(() => `property-value-suggestions-${props.property.key}`);
const canReset = computed(() => canActOnDynamicProperty(props.property));
const title = computed(() => `${t(canReset.value ? 'common.edit' : 'common.set')} ${propertyLabel(props.property)}`);
const { hasUnsavedChanges, confirmLeave } = useSnapshotChanges(state);

function resetState(): void {
    const value = props.property.value?.value;
    state.reason = '';
    state.stringValue = value?.oneofKind === 'stringValue' ? value.stringValue : '';
    state.numberValue =
        value?.oneofKind === 'integerValue' ? value.integerValue : value?.oneofKind === 'numberValue' ? value.numberValue : 0;
    state.boolValue = value?.oneofKind === 'boolValue' ? value.boolValue : false;
    state.dateValue =
        value?.oneofKind === 'dateValue'
            ? new Date(value.dateValue.year, value.dateValue.month - 1, value.dateValue.day)
            : undefined;
    state.timestampValue = value?.oneofKind === 'timestampValue' ? toDate(value.timestampValue) : undefined;
    state.reset = false;
}

watch(() => props.property, resetState, { immediate: true, deep: true });

function propertyValue(): PropertyValue | undefined {
    if (state.reset) return undefined;
    switch (valueType.value) {
        case PropertyValueType.STRING:
            return { value: { oneofKind: 'stringValue', stringValue: state.stringValue } };
        case PropertyValueType.INTEGER:
            return { value: { oneofKind: 'integerValue', integerValue: Math.trunc(state.numberValue) } };
        case PropertyValueType.NUMBER:
            return { value: { oneofKind: 'numberValue', numberValue: state.numberValue } };
        case PropertyValueType.BOOL:
            return { value: { oneofKind: 'boolValue', boolValue: state.boolValue } };
        case PropertyValueType.DATE:
            if (!state.dateValue) return undefined;
            return {
                value: {
                    oneofKind: 'dateValue',
                    dateValue: {
                        year: state.dateValue.getFullYear(),
                        month: state.dateValue.getMonth() + 1,
                        day: state.dateValue.getDate(),
                    },
                },
            };
        case PropertyValueType.TIMESTAMP:
            if (!state.timestampValue) return undefined;
            return { value: { oneofKind: 'timestampValue', timestampValue: toTimestamp(state.timestampValue)! } };
        default:
            return undefined;
    }
}

async function save(values: Schema): Promise<void> {
    const value = propertyValue();
    if (!state.reset && value === undefined) return;

    try {
        const mutation: PropertyMutation = value
            ? { key: props.property.key, valueChange: { oneofKind: 'setValue', setValue: value }, state: undefined }
            : { key: props.property.key, valueChange: { oneofKind: 'clearValue', clearValue: {} }, state: undefined };
        if (reasonMode.value !== PropertyReasonMode.NONE && values.reason.trim() !== '') mutation.reason = values.reason.trim();
        const propertyMutations: PropertyMutation[] = [mutation];
        if (props.scope === PropertyScope.USER) {
            if (props.userId === undefined) return;
            const { response } = await citizensClient.setUserProps({
                props: { userId: props.userId },
                reason: values.reason.trim(),
                propertyMutations,
            });
            emit('update:properties', response.props?.properties);
        } else if (props.scope === PropertyScope.VEHICLE) {
            if (!props.plate) return;
            const { response } = await vehiclesClient.setVehicleProps({
                props: { plate: props.plate },
                reason: values.reason.trim(),
                propertyMutations,
            });
            emit('update:properties', response.props?.properties);
        }

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });
        emit('close', false);
    } catch (error) {
        handleGRPCError(error as RpcError);
        throw error;
    }
}

const { submit, isSubmitting, canSubmit } = useSubmitGuard(async (event: FormSubmitEvent<Schema>) => save(event.data));
const formRef = useTemplateRef('formRef');

function submitChanges(): void {
    state.reset = false;
    formRef.value?.submit();
}

function confirmReset(): void {
    resetConfirmModal.open({
        title: `${t('common.reset')} ${propertyLabel(props.property)}`,
        confirm: () => {
            state.reset = true;
            formRef.value?.submit();
        },
    });
}

async function closeModal(): Promise<void> {
    if (!canSubmit.value) return;
    if (hasUnsavedChanges.value && !(await confirmLeave())) return;
    emit('close', false);
}
</script>

<template>
    <ResponsiveDrawerModal :title="title" :close="false" :dismissible="!hasUnsavedChanges && canSubmit">
        <template #header>
            <div class="flex w-full items-center justify-between gap-2">
                <h3 class="font-semibold text-highlighted">{{ title }}</h3>
                <UButton color="neutral" variant="ghost" icon="i-mdi-close" :disabled="!canSubmit" @click="closeModal" />
            </div>
        </template>

        <template #body>
            <UForm ref="formRef" class="grid gap-2" :schema="schema" :state="state" @submit="submit">
                <UFormField
                    v-if="reasonMode !== PropertyReasonMode.NONE"
                    name="reason"
                    :label="$t('common.reason')"
                    :required="reasonMode === PropertyReasonMode.REQUIRED"
                >
                    <UInput v-model="state.reason" class="w-full" :placeholder="$t('common.reason')" />
                </UFormField>

                <UFormField :label="propertyLabel(property)" required>
                    <USelectMenu
                        v-if="valueType === PropertyValueType.STRING && stringInputMode === 'strict'"
                        v-model="state.stringValue"
                        class="w-full"
                        :items="stringValueOptions"
                        label-key="label"
                        value-key="value"
                    />
                    <template v-else-if="valueType === PropertyValueType.STRING">
                        <UInput
                            v-model="state.stringValue"
                            class="w-full"
                            :list="stringInputMode === 'suggestions' ? stringValueSuggestionsId : undefined"
                            :minlength="property.definition?.valueConstraints?.stringMinLength"
                            :maxlength="property.definition?.valueConstraints?.stringMaxLength"
                        />
                        <datalist v-if="stringInputMode === 'suggestions'" :id="stringValueSuggestionsId">
                            <option v-for="value in allowedStringValues" :key="value" :value="value" />
                        </datalist>
                    </template>
                    <UInputNumber
                        v-else-if="valueType === PropertyValueType.INTEGER || valueType === PropertyValueType.NUMBER"
                        v-model="state.numberValue"
                        class="w-full"
                        :step="valueType === PropertyValueType.INTEGER ? 1 : 0.01"
                    />
                    <USwitch v-else-if="valueType === PropertyValueType.BOOL" v-model="state.boolValue" />
                    <InputDatePicker
                        v-else-if="valueType === PropertyValueType.DATE"
                        v-model="state.dateValue"
                        class="w-full"
                        required
                    />
                    <InputDatePicker
                        v-else-if="valueType === PropertyValueType.TIMESTAMP"
                        v-model="state.timestampValue"
                        class="w-full"
                        time
                        required
                    />
                </UFormField>
            </UForm>
        </template>

        <template #footer>
            <div class="flex w-full flex-col-reverse gap-2 sm:flex-row">
                <UButton
                    class="flex-1"
                    color="neutral"
                    :disabled="!canSubmit"
                    :label="$t('common.close', 1)"
                    @click="closeModal"
                />
                <UButton
                    v-if="canReset"
                    class="flex-1"
                    color="error"
                    :disabled="!canSubmit"
                    :loading="isSubmitting"
                    :label="$t('common.reset')"
                    @click="confirmReset"
                />
                <UButton
                    class="flex-1"
                    :disabled="!canSubmit"
                    :loading="isSubmitting"
                    :label="$t('common.save')"
                    @click="submitChanges"
                />
            </div>
        </template>
    </ResponsiveDrawerModal>
</template>
