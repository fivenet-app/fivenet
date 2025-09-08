<script lang="ts" setup>
import { getQualificationsQualificationsClient } from '~~/gen/ts/clients';
import type { QualificationRequirement, QualificationShort } from '~~/gen/ts/resources/qualifications/qualifications';
import SelectMenu from '../partials/SelectMenu.vue';

const props = withDefaults(
    defineProps<{
        requirement: QualificationRequirement;
        readOnly?: boolean;
        qualificationId?: number;
    }>(),
    {
        readOnly: false,
        qualificationId: undefined,
    },
);

const emit = defineEmits<{
    (e: 'update-qualification', qualification?: QualificationShort): void;
    (e: 'remove'): void;
}>();

const qualificationsQualificationsClient = await getQualificationsQualificationsClient();

async function listQualifications(search?: string): Promise<QualificationShort[]> {
    try {
        const call = qualificationsQualificationsClient.listQualifications({
            pagination: {
                offset: 0,
            },
            search: search,
        });
        const { response } = await call;

        if (props.qualificationId === undefined) {
            return response.qualifications;
        }

        return (response.qualifications as QualificationShort[]).filter((q) => q.id !== props.qualificationId);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const selectedQualification = ref<QualificationShort | undefined>(props.requirement.targetQualification);

watch(selectedQualification, () => emit('update-qualification', selectedQualification.value));
</script>

<template>
    <div class="flex flex-row items-center">
        <UFormField class="flex-1" name="selectedQualification">
            <SelectMenu
                v-model="selectedQualification"
                block
                :searchable="(q: string) => listQualifications(q)"
                :searchable-key="`qualification-${qualificationId}-requirement-entry`"
                :search-input="{ placeholder: $t('common.search_field') }"
                class="w-full"
            >
                <template v-if="selectedQualification" #default>
                    <span class="truncate"> {{ selectedQualification.abbreviation }}: {{ selectedQualification.title }} </span>
                </template>

                <template #item="{ item }">
                    <span class="truncate">
                        <template v-if="item?.abbreviation">{{ item.abbreviation }}: </template
                        >{{ !item.title ? $t('common.untitled') : item.title }}
                    </span>
                </template>

                <template #empty> {{ $t('common.not_found', [$t('common.qualification', 2)]) }} </template>
            </SelectMenu>
        </UFormField>

        <UTooltip :text="$t('components.qualifications.remove_requirement')">
            <UButton class="ml-2" icon="i-mdi-close" color="error" @click="$emit('remove')" />
        </UTooltip>
    </div>
</template>
