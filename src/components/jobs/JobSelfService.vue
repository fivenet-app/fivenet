<script lang="ts" setup>
import SelfServicePropsAbsenceDateModal from '~/components/jobs/colleagues/SelfServicePropsAbsenceDateModal.vue';
import SelfServicePropsProfilePictureModal from '~/components/jobs/colleagues/SelfServicePropsProfilePictureModal.vue';
import { checkIfCanAccessColleague } from '~/components/jobs/colleagues/helpers';
import { useAuthStore } from '~/store/auth';

defineProps<{
    userId: number;
}>();

const { $grpc } = useNuxtApp();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const { data: colleagueSelf } = useLazyAsyncData('jobs-selfcolleague', async () => {
    try {
        const call = $grpc.getJobsClient().getSelf({});
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
});

const modal = useModal();
</script>

<template>
    <UCard class="flex-1">
        <template #header>
            <h3 class="text-lg font-semibold">{{ $t('components.jobs.self_service.title') }}</h3>
        </template>

        <div class="flex flex-initial flex-col items-center gap-1 md:flex-row">
            <UButtonGroup class="flex w-full">
                <UButton
                    v-if="
                        colleagueSelf?.colleague &&
                        can('JobsService.SetJobsUserProps') &&
                        checkIfCanAccessColleague(activeChar!, colleagueSelf.colleague, 'JobsService.SetJobsUserProps')
                    "
                    block
                    class="flex-1"
                    icon="i-mdi-island"
                    @click="
                        modal.open(SelfServicePropsAbsenceDateModal, {
                            userId: colleagueSelf.colleague.userId,
                        })
                    "
                >
                    <span>{{ $t('components.jobs.self_service.set_absence_date') }}</span>
                </UButton>
                <UButton block class="flex-1" icon="i-mdi-camera" @click="modal.open(SelfServicePropsProfilePictureModal, {})">
                    <span>{{ $t('components.jobs.self_service.set_profile_picture') }}</span>
                </UButton>
            </UButtonGroup>
        </div>
    </UCard>
</template>
