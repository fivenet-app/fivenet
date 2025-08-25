<script lang="ts" setup>
import { getQualificationsQualificationsClient } from '~~/gen/ts/clients';
import type {
    Qualification,
    QualificationRequirement,
    QualificationShort,
} from '~~/gen/ts/resources/qualifications/qualifications';

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

const qualificationsLoading = ref(false);

async function listQualifications(search?: string): Promise<Qualification[]> {
    qualificationsLoading.value = true;
    try {
        const call = qualificationsQualificationsClient.listQualifications({
            pagination: {
                offset: 0,
            },
            search: search,
        });
        const { response } = await call;

        qualificationsLoading.value = false;
        if (props.qualificationId === undefined) {
            return response.qualifications;
        }

        return response.qualifications.filter((q) => q.id !== props.qualificationId);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    } finally {
        qualificationsLoading.value = false;
    }
}

const selectedQualification = ref<QualificationShort | undefined>(props.requirement.targetQualification);

watch(selectedQualification, () => emit('update-qualification', selectedQualification.value));
</script>

<template>
    <div class="my-2 flex flex-row items-center">
        <UFormField class="flex-1" name="selectedQualification">
            <ClientOnly>
                <USelectMenu
                    v-model="selectedQualification"
                    option-attribute="title"
                    :search-attributes="['title', 'abbreviation']"
                    block
                    searchable-lazy
                    :searchable="(q: string) => listQualifications(q)"
                    :searchable-placeholder="$t('common.search_field')"
                    :loading="qualificationsLoading"
                >
                    <template #item-label>
                        <span v-if="selectedQualification" class="truncate">
                            {{ selectedQualification.abbreviation }}: {{ selectedQualification.title }}
                        </span>
                    </template>

                    <template #item="{ option: quali }">
                        <span class="truncate">
                            <template v-if="quali.abbreviation">{{ quali.abbreviation }}: </template
                            >{{ !quali.title ? $t('common.untitled') : quali.title }}
                        </span>
                    </template>

                    <template #empty> {{ $t('common.not_found', [$t('common.qualification', 2)]) }} </template>
                </USelectMenu>
            </ClientOnly>
        </UFormField>

        <UTooltip :text="$t('components.qualifications.remove_requirement')">
            <UButton class="ml-2" icon="i-mdi-close" @click="$emit('remove')" />
        </UTooltip>
    </div>
</template>
