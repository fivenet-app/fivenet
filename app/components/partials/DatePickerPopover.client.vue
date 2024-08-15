<script lang="ts" setup>
import { format } from 'date-fns';
import DatePickerClient from './DatePicker.client.vue';

const props = defineProps<{
    modelValue?: Date | undefined;
    popover?: any;
    button?: any;
    datePicker?: any;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', modelValue: Date | undefined): void;
}>();

const date = useVModel(props, 'modelValue', emit);

const open = ref(false);

defineOptions({
    inheritAttrs: false,
});
</script>

<template>
    <UPopover v-bind="props.popover" v-model:open="open">
        <UButton
            v-bind="button"
            variant="outline"
            color="gray"
            block
            icon="i-mdi-calendar-month"
            :label="modelValue ? format(modelValue, 'dd.MM.yyyy') : 'dd.mm.yyyy'"
            @touchstart="open = true"
        />

        <template #panel="{ close }">
            <DatePickerClient v-bind="datePicker" v-model="date" @close="close" />
        </template>
    </UPopover>
</template>
