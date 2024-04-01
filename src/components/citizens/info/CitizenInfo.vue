<script lang="ts" setup>
import AddToButton from '~/components/clipboard/AddToButton.vue';
import { useClipboardStore } from '~/store/clipboard';
import { useNotificatorStore } from '~/store/notificator';
import { User } from '~~/gen/ts/resources/users/users';
import CitizenActivityFeed from '~/components/citizens/info/CitizenActivityFeed.vue';
import CitizenDocuments from '~/components/citizens/info/CitizenDocuments.vue';
import CitizenProfile from '~/components/citizens/info/CitizenProfile.vue';
import CitizenVehicles from '~/components/citizens/info/CitizenVehicles.vue';
import ClipboardButton from '~/components/clipboard/ClipboardButton.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import type { Perms } from '~~/gen/ts/perms';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';

const props = defineProps<{
    id: number;
}>();

const { $grpc } = useNuxtApp();

const clipboardStore = useClipboardStore();
const notifications = useNotificatorStore();

const { t } = useI18n();

const tabs: { slot: string; label: string; icon: string; permission: Perms }[] = [
    {
        slot: 'profile',
        label: t('common.profile'),
        icon: 'i-mdi-account',
        permission: 'CitizenStoreService.ListCitizens' as Perms,
    },
    {
        slot: 'vehicles',
        label: t('common.vehicle', 2),
        icon: 'i-mdi-car',
        permission: 'DMVService.ListVehicles' as Perms,
    },
    {
        slot: 'documents',
        label: t('common.document', 2),
        icon: 'i-mdi-file-document-multiple',
        permission: 'DocStoreService.ListUserDocuments' as Perms,
    },
    {
        slot: 'activity',
        label: t('common.activity'),
        icon: 'i-mdi-bulletin-board',
        permission: 'CitizenStoreService.ListUserActivity' as Perms,
    },
].filter((tab) => can(tab.permission));

const { data: user, pending, refresh, error } = useLazyAsyncData(`citizen-${props.id}`, () => getUser(props.id));

async function getUser(userId: number): Promise<User> {
    try {
        const call = $grpc.getCitizenStoreClient().getUser({ userId });
        const { response } = await call;

        if (response.user?.props === undefined) {
            response.user!.props = {
                userId: response.user!.userId,
            };
        }

        return response.user!;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

function addToClipboard(): void {
    if (user.value === null) {
        return;
    }

    clipboardStore.addUser(user.value);

    notifications.dispatchNotification({
        title: { key: 'notifications.clipboard.citizen_add.title', parameters: {} },
        content: { key: 'notifications.clipboard.citizen_add.content', parameters: {} },
        duration: 3250,
        type: 'info',
    });
}
</script>

<template>
    <div class="mx-2">
        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.citizen', 1)])" />
        <DataErrorBlock
            v-else-if="error"
            :title="$t('common.unable_to_load', [$t('common.citizen', 1)])"
            :message="$t(error.message)"
            :retry="refresh"
        />
        <DataNoDataBlock v-else-if="user === null" />

        <template v-else>
            <div class="mb-14">
                <div class="my-4 flex gap-4 px-4">
                    <ProfilePictureImg
                        :url="user.props?.mugShot?.url"
                        :name="`${user.firstname} ${user.lastname}`"
                        size="xl"
                        :rounded="false"
                        :enable-popup="true"
                        :alt-text="$t('common.mug_shot')"
                    />
                    <div class="w-full">
                        <div class="flex snap-x flex-row flex-wrap justify-between gap-2 overflow-x-auto">
                            <h1 class="flex-1 break-words px-0.5 py-1 text-4xl font-bold text-neutral sm:pl-1">
                                {{ user?.firstname }} {{ user?.lastname }}
                            </h1>
                            <IDCopyBadge
                                :id="user.userId"
                                prefix="CIT"
                                :title="{ key: 'notifications.citizen_info.copy_citizen_id.title', parameters: {} }"
                                :content="{ key: 'notifications.citizen_info.copy_citizen_id.content', parameters: {} }"
                                class="min-h-9 place-self-end"
                            />
                        </div>
                        <div class="inline-flex gap-2">
                            <UBadge>
                                {{ user.jobLabel }}
                                <template v-if="user.jobGrade > 0">
                                    ({{ $t('common.rank') }}: {{ user.jobGradeLabel }})</template
                                >
                            </UBadge>
                            <UBadge v-if="user.props?.wanted" color="red">
                                {{ $t('common.wanted').toUpperCase() }}
                            </UBadge>
                        </div>
                    </div>
                </div>

                <UTabs :items="tabs" class="w-full" :unmount="true">
                    <template #default="{ item, selected }">
                        <div class="flex items-center gap-2 relative truncate">
                            <UIcon :name="item.icon" class="w-4 h-4 flex-shrink-0" />

                            <span class="truncate">{{ item.label }}</span>

                            <span
                                v-if="selected"
                                class="absolute -right-4 w-2 h-2 rounded-full bg-primary-500 dark:bg-primary-400"
                            />
                        </div>
                    </template>

                    <template #profile>
                        <CitizenProfile
                            :user="user"
                            @update:wanted-status="user.props!.wanted = $event"
                            @update:job="
                                user.job = $event.job.name;
                                user.jobLabel = $event.job.label;
                                user.jobGrade = $event.grade.grade;
                                user.jobGradeLabel = $event.grade.label;
                            "
                            @update:traffic-infraction-points="user.props!.trafficInfractionPoints = $event"
                            @update:mug-shot="user.props!.mugShot = $event"
                        />
                    </template>
                    <template #vehicles v-if="can('DMVService.ListVehicles')">
                        <CitizenVehicles :user-id="user.userId" />
                    </template>
                    <template #documents v-if="can('DocStoreService.ListUserDocuments')">
                        <CitizenDocuments :user-id="user.userId" />
                    </template>
                    <template #activity v-if="can('CitizenStoreService.ListUserActivity')">
                        <CitizenActivityFeed :user-id="user.userId" />
                    </template>
                </UTabs>
            </div>

            <ClipboardButton />
        </template>
    </div>

    <AddToButton :callback="addToClipboard" :title="$t('components.clipboard.clipboard_button.add')" />
</template>
