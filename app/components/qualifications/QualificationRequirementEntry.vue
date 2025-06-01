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

const { $grpc } = useNuxtApp();

const qualificationsLoading = ref(false);
async function listQualifications(search?: string): Promise<Qualification[]> {
    qualificationsLoading.value = true;
    try {
        const call = $grpc.qualifications.qualifications.listQualifications({
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
        <UFormGroup class="flex-1" name="selectedQualification">
            <ClientOnly>
                <USelectMenu
                    v-model="selectedQualification"
                    option-attribute="title"
                    :search-attributes="['title', 'abbreviation']"
                    block
                    searchable-lazy
                    :searchable="(query: string) => listQualifications(query)"
                    :searchable-placeholder="$t('common.search_field')"
                    :loading="qualificationsLoading"
                >
                    <template #label>
                        <span v-if="selectedQualification" class="truncate">
                            {{ selectedQualification.abbreviation }}: {{ selectedQualification.title }}
                        </span>
                    </template>

                    <template #option="{ option: quali }">
                        <span class="truncate">
                            <template v-if="quali.abbreviation">{{ quali.abbreviation }}: </template
                            >{{ !quali.title ? $t('common.untitled') : quali.title }}
                        </span>
                    </template>

                    <template #option-empty="{ query: search }">
                        <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                    </template>

                    <template #empty> {{ $t('common.not_found', [$t('common.qualification', 2)]) }} </template>
                </USelectMenu>
            </ClientOnly>
        </UFormGroup>

        <UTooltip :text="$t('components.qualifications.remove_requirement')">
            <UButton class="ml-2" :ui="{ rounded: 'rounded-full' }" icon="i-mdi-close" @click="$emit('remove')" />
        </UTooltip>
    </div>
</template>
