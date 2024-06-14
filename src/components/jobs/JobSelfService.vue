<script lang="ts" setup>
import { isToday, parse } from 'date-fns';
import { emojiBlasts } from 'emoji-blast';
import SelfServicePropsAbsenceDateModal from '~/components/jobs/colleagues/SelfServicePropsAbsenceDateModal.vue';
import SelfServicePropsProfilePictureModal from '~/components/jobs/colleagues/SelfServicePropsProfilePictureModal.vue';
import { checkIfCanAccessColleague } from '~/components/jobs/colleagues/helpers';
import { useAuthStore } from '~/store/auth';

defineProps<{
    userId: number;
}>();

const modal = useModal();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const { data: colleagueSelf } = useLazyAsyncData('jobs-selfcolleague', async () => {
    try {
        const call = getGRPCJobsClient().getSelf({});
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
});

onMounted(() => {
    if (!colleagueSelf.value?.colleague?.dateofbirth) {
        return;
    }

    const birthday = parse(colleagueSelf.value?.colleague?.dateofbirth, 'dd.MM.yyyy', new Date());
    birthday.setFullYear(new Date().getFullYear());

    if (isToday(birthday)) {
        const { cancel } = emojiBlasts({
            emojis: ['🎂', '🎁', '🍰', '🎈', '🎉', '🥳', '🎊', '✨'],
        });
        useTimeoutFn(cancel, 5000);
    }
});
</script>

<template>
    <UCard class="flex-1">
        <template #header>
            <h3 class="text-lg font-semibold">{{ $t('components.jobs.self_service.title') }}</h3>
        </template>

        <div class="flex flex-initial flex-col items-center gap-1 md:flex-row">
            <UButtonGroup class="inline-flex w-full">
                <UButton
                    v-if="
                        colleagueSelf?.colleague &&
                        can('JobsService.SetJobsUserProps').value &&
                        checkIfCanAccessColleague(activeChar!, colleagueSelf.colleague, 'JobsService.SetJobsUserProps')
                    "
                    block
                    class="flex-1"
                    icon="i-mdi-island"
                    @click="
                        modal.open(SelfServicePropsAbsenceDateModal, {
                            userId: colleagueSelf.colleague.userId,
                            userProps: colleagueSelf.colleague.props,
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
