<script lang="ts" setup>
import { isToday, parse } from 'date-fns';
import { emojiBlasts } from 'emoji-blast';
import SelfServiceAbsenceDateModal from '~/components/jobs/colleagues/SelfServiceAbsenceDateModal.vue';
import SelfServiceAvatarModal from '~/components/jobs/colleagues/SelfServiceAvatarModal.vue';
import { getJobsJobsClient } from '~~/gen/ts/clients';

defineProps<{
    userId: number;
}>();

const overlay = useOverlay();

const { can, activeChar } = useAuth();

const jobsJobsClient = await getJobsJobsClient();

const { data: colleagueSelf } = useLazyAsyncData('jobs-selfcolleague', async () => {
    try {
        const call = jobsJobsClient.getSelf({});
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
});

onBeforeMount(() => {
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

const selfServicePropsAbsenceDateModal = overlay.create(SelfServiceAbsenceDateModal);
const selfServicePropsAvatarModal = overlay.create(SelfServiceAvatarModal);
</script>

<template>
    <UCard class="flex-1">
        <template #header>
            <h3 class="text-lg font-semibold">{{ $t('components.jobs.self_service.title') }}</h3>
        </template>

        <div class="flex flex-initial flex-col items-center gap-1 md:flex-row">
            <UButton
                v-if="
                    colleagueSelf?.colleague &&
                    activeChar?.userId === colleagueSelf.colleague?.userId &&
                    can('jobs.JobsService/SetColleagueProps').value
                "
                class="flex-1"
                block
                icon="i-mdi-island"
                @click="
                    selfServicePropsAbsenceDateModal.open({
                        userId: colleagueSelf.colleague.userId,
                        userProps: colleagueSelf.colleague.props,
                    })
                "
            >
                <span>{{ $t('components.jobs.self_service.set_absence_date') }}</span>
            </UButton>
            <UButton class="flex-1" block icon="i-mdi-camera" @click="selfServicePropsAvatarModal.open({})">
                <span>{{ $t('components.jobs.self_service.set_profile_picture') }}</span>
            </UButton>
        </div>
    </UCard>
</template>
