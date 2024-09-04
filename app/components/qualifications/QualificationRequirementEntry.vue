<script lang="ts" setup>
import type {
    Qualification,
    QualificationRequirement,
    QualificationShort,
} from '~~/gen/ts/resources/qualifications/qualifications';

const props = withDefaults(
    defineProps<{
        requirement: QualificationRequirement;
        readOnly?: boolean;
    }>(),
    {
        readOnly: false,
    },
);

const emits = defineEmits<{
    (e: 'update-qualification', qualification?: QualificationShort): void;
    (e: 'remove'): void;
}>();

const qualificationsLoading = ref(false);
async function listQualifications(search?: string): Promise<Qualification[]> {
    qualificationsLoading.value = true;
    try {
        const call = getGRPCQualificationsClient().listQualifications({
            pagination: {
                offset: 0,
            },
            search: search,
        });
        const { response } = await call;

        qualificationsLoading.value = false;
        return response.qualifications;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    } finally {
        qualificationsLoading.value = false;
    }
}
const selectedQualification = ref<QualificationShort | undefined>(props.requirement.targetQualification);

watch(selectedQualification, () => emits('update-qualification', selectedQualification.value));
</script>

<template>
    <div class="my-2 flex flex-row items-center">
        <UFormGroup name="selectedQualification" class="flex-1">
            <UInputMenu
                v-model="selectedQualification"
                option-attribute="title"
                :search-attributes="['title']"
                block
                :search="(query: string) => listQualifications(query)"
                search-lazy
                :search-placeholder="$t('common.search_field')"
                :loading="qualificationsLoading"
            >
                <template #option-empty="{ query: search }">
                    <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                </template>
                <template #empty> {{ $t('common.not_found', [$t('common.category', 2)]) }} </template>
            </UInputMenu>
        </UFormGroup>

        <UButton :ui="{ rounded: 'rounded-full' }" class="ml-2" icon="i-mdi-close" @click="$emit('remove')" />
    </div>
</template>
