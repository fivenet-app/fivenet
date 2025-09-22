<script lang="ts" setup>
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import { useCentrumStore } from '~/stores/centrum';
import { CentrumMode } from '~~/gen/ts/resources/centrum/settings';

defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const centrumStore = useCentrumStore();
const { dispatchers, anyDispatchersActive, getCurrentMode } = storeToRefs(centrumStore);
</script>

<template>
    <UDrawer
        :title="$t('common.dispatchers', 2)"
        :overlay="false"
        :close="{ onClick: () => $emit('close', false) }"
        :ui="{ title: 'flex flex-row gap-2' }"
    >
        <template #title>
            <span class="flex-1">{{ $t('common.dispatchers', 2) }}</span>
            <UButton icon="i-mdi-close" color="neutral" variant="link" size="sm" @click="$emit('close', false)" />
        </template>

        <template #actions>
            <UBadge color="neutral">
                {{ $t('common.mode') }}: {{ $t(`enums.centrum.CentrumMode.${CentrumMode[getCurrentMode ?? 0]}`) }}
            </UBadge>
        </template>

        <template #body>
            <UContainer>
                <DataNoDataBlock
                    v-if="!dispatchers || !anyDispatchersActive"
                    icon="i-mdi-monitor"
                    :type="$t('common.dispatcher')"
                />
                <template v-else>
                    <div v-for="dispas in dispatchers" :key="dispas.job" class="gap-4 p-4" :cols="2">
                        <h3 class="mb-4 text-lg font-semibold">
                            {{ dispas.jobLabel ?? dispas.job }}
                        </h3>

                        <UPageGrid class="lg:grid-cols-2 xl:grid-cols-2">
                            <UPageCard
                                v-for="dispatcher in dispas.dispatchers"
                                :key="dispatcher.userId"
                                :title="`${dispatcher.firstname} ${dispatcher.lastname}`"
                                icon="i-mdi-account"
                                :ui="{
                                    title: 'text-highlighted text-base font-semibold flex items-center gap-1.5 line-clamp-2 whitespace-break-spaces',
                                }"
                            >
                                <template #default>
                                    <PhoneNumberBlock :number="dispatcher.phoneNumber" />
                                </template>

                                <template v-if="dispatcher.profilePicture" #leading>
                                    <ProfilePictureImg
                                        class="mr-2"
                                        :src="dispatcher?.profilePicture"
                                        :name="`${dispatcher.firstname} ${dispatcher.lastname}`"
                                        size="sm"
                                        :enable-popup="false"
                                        :alt="$t('common.profile_picture')"
                                    />
                                </template>
                            </UPageCard>
                        </UPageGrid>
                    </div>
                </template>
            </UContainer>
        </template>

        <template #footer>
            <UButton class="flex-1" color="neutral" block @click="$emit('close', false)">
                {{ $t('common.close', 1) }}
            </UButton>
        </template>
    </UDrawer>
</template>
