<script lang="ts" setup>
type RsvpScope = 'series' | 'occurrence';

const props = withDefaults(
    defineProps<{
        open: boolean;
        remove?: boolean;
    }>(),
    {
        remove: false,
    },
);

const emit = defineEmits<{
    (e: 'update:open', value: boolean): void;
    (e: 'confirm', scope: RsvpScope): void;
}>();

const { t } = useI18n();

const scope = ref<RsvpScope>('occurrence');

const isOpen = computed({
    get: () => props.open,
    set: (value: boolean) => emit('update:open', value),
});

watch(
    () => props.open,
    (open) => {
        if (open) {
            scope.value = 'occurrence';
        }
    },
);

function confirm(): void {
    emit('confirm', scope.value);
    isOpen.value = false;
}
</script>

<template>
    <UModal
        v-model:open="isOpen"
        :title="$t('components.calendar.rsvp_scope.title')"
        :description="
            props.remove
                ? $t('components.calendar.rsvp_scope.remove_description')
                : $t('components.calendar.rsvp_scope.description')
        "
    >
        <template #body>
            <div class="space-y-4">
                <URadioGroup
                    v-model="scope"
                    class="w-full"
                    name="rsvp-scope"
                    variant="table"
                    orientation="horizontal"
                    :items="[
                        {
                            value: 'occurrence',
                            label: t('components.calendar.rsvp_scope.occurrence_only'),
                        },
                        {
                            value: 'series',
                            label: t('components.calendar.rsvp_scope.series'),
                        },
                    ]"
                    value-key="value"
                    :ui="{ wrapper: 'space-y-2', fieldset: 'w-full', item: 'flex-1' }"
                />
            </div>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton class="flex-1" block color="neutral" :label="$t('common.cancel')" @click="isOpen = false" />

                <UButton class="flex-1" block color="primary" :label="$t('common.confirm')" @click="confirm()" />
            </UFieldGroup>
        </template>
    </UModal>
</template>
