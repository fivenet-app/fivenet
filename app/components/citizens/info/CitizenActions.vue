<script lang="ts" setup>
import SetJobModal from '~/components/citizens/info/props/SetJobModal.vue';
import SetMugshotModal from '~/components/citizens/info/props/SetMugshotModal.vue';
import SetTrafficPointsModal from '~/components/citizens/info/props/SetTrafficPointsModal.vue';
import SetWantedModal from '~/components/citizens/info/props/SetWantedModal.vue';
import TemplateModal from '~/components/documents/templates/TemplateModal.vue';
import { checkIfCanAccessColleague } from '~/components/jobs/colleagues/helpers';
import { useClipboardStore } from '~/stores/clipboard';
import type { File } from '~~/gen/ts/resources/file/file';
import type { Job, JobGrade } from '~~/gen/ts/resources/jobs/jobs';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { User } from '~~/gen/ts/resources/users/users';

const props = withDefaults(
    defineProps<{
        user: User;
        registerKBDs?: boolean;
    }>(),
    {
        registerKBDs: false,
    },
);

const emit = defineEmits<{
    (e: 'update:wantedStatus', value: boolean): void;
    (e: 'update:job', value: { job: Job; grade: JobGrade }): void;
    (e: 'update:trafficInfractionPoints', value: number): void;
    (e: 'update:mugshot', value?: File): void;
}>();

const { attr, can, activeChar } = useAuth();

const clipboardStore = useClipboardStore();

const notifications = useNotificationsStore();

const overlay = useOverlay();

const templateModal = overlay.create(TemplateModal);
const setWantedModal = overlay.create(SetWantedModal);
const setJobModal = overlay.create(SetJobModal);
const setTrafficPointsModal = overlay.create(SetTrafficPointsModal);
const setMugshotModal = overlay.create(SetMugshotModal);

function openTemplates(): void {
    if (!props.user) {
        return;
    }

    clipboardStore.addUser(props.user, true);

    templateModal.open({});
}

function copyLinkToClipboard(): void {
    copyToClipboardWrapper(window.location.href);

    notifications.add({
        title: { key: 'notifications.clipboard.link_copied.title', parameters: {} },
        description: { key: 'notifications.clipboard.link_copied.content', parameters: {} },
        duration: 3250,
        type: NotificationType.INFO,
    });
}

if (props.registerKBDs) {
    defineShortcuts({
        'c-w': () => {
            if (!attr('citizens.CitizensService/SetUserProps', 'Fields', 'Wanted').value) {
                return;
            }

            setWantedModal.open({
                user: props.user,
                'onUpdate:wantedStatus': ($event) => emit('update:wantedStatus', $event),
            });
        },
        'c-j': () => {
            if (!attr('citizens.CitizensService/SetUserProps', 'Fields', 'Job').value) {
                return;
            }

            setJobModal.open({
                user: props.user,
                'onUpdate:job': ($event) => emit('update:job', $event),
            });
        },
        'c-p': () => {
            if (!attr('citizens.CitizensService/SetUserProps', 'Fields', 'TrafficInfractionPoints').value) {
                return;
            }

            setTrafficPointsModal.open({
                user: props.user,
                'onUpdate:trafficInfractionPoints': ($event) => emit('update:trafficInfractionPoints', $event),
            });
        },
        'c-m': () => {
            if (!attr('citizens.CitizensService/SetUserProps', 'Fields', 'Mugshot').value) {
                return;
            }

            setMugshotModal.open({
                user: props.user,
                'onUpdate:mugshot': ($event) => emit('update:mugshot', $event),
            });
        },
        'c-d': () => {
            if (!can('documents.DocumentsService/UpdateDocument').value) {
                return;
            }

            openTemplates();
        },
    });
}
</script>

<template>
    <div class="flex w-full flex-col gap-2">
        <UTooltip
            v-if="attr('citizens.CitizensService/SetUserProps', 'Fields', 'Wanted').value"
            :text="
                user?.props?.wanted
                    ? $t('components.citizens.CitizenInfoProfile.revoke_wanted')
                    : $t('components.citizens.CitizenInfoProfile.set_wanted')
            "
            :kbds="['C', 'W']"
        >
            <UButton
                :color="user?.props?.wanted ? 'error' : 'primary'"
                block
                truncate
                :icon="user?.props?.wanted ? 'i-mdi-account-alert' : 'i-mdi-account-cancel'"
                @click="
                    setWantedModal.open({
                        user: user,
                        'onUpdate:wantedStatus': ($event) => $emit('update:wantedStatus', $event),
                    })
                "
            >
                {{
                    user?.props?.wanted
                        ? $t('components.citizens.CitizenInfoProfile.revoke_wanted')
                        : $t('components.citizens.CitizenInfoProfile.set_wanted')
                }}
            </UButton>
        </UTooltip>

        <UTooltip
            v-if="attr('citizens.CitizensService/SetUserProps', 'Fields', 'Job').value"
            :text="$t('components.citizens.CitizenInfoProfile.set_job')"
            :kbds="['C', 'J']"
        >
            <UButton
                block
                icon="i-mdi-briefcase"
                @click="
                    setJobModal.open({
                        user: user,
                        'onUpdate:job': ($event) => $emit('update:job', $event),
                    })
                "
            >
                {{ $t('components.citizens.CitizenInfoProfile.set_job') }}
            </UButton>
        </UTooltip>

        <UTooltip
            v-if="attr('citizens.CitizensService/SetUserProps', 'Fields', 'TrafficInfractionPoints').value"
            :text="$t('components.citizens.CitizenInfoProfile.set_traffic_points')"
            :kbds="['C', 'P']"
        >
            <UButton
                block
                icon="i-mdi-counter"
                @click="
                    setTrafficPointsModal.open({
                        user: user,
                        'onUpdate:trafficInfractionPoints': ($event) => $emit('update:trafficInfractionPoints', $event),
                    })
                "
            >
                {{ $t('components.citizens.CitizenInfoProfile.set_traffic_points') }}
            </UButton>
        </UTooltip>

        <UTooltip
            v-if="attr('citizens.CitizensService/SetUserProps', 'Fields', 'Mugshot').value"
            :text="$t('components.citizens.CitizenInfoProfile.set_mugshot')"
            :kbds="['C', 'M']"
        >
            <UButton
                block
                icon="i-mdi-camera"
                @click="
                    setMugshotModal.open({
                        user: user,
                        'onUpdate:mugshot': ($event) => $emit('update:mugshot', $event),
                    })
                "
            >
                {{ $t('components.citizens.CitizenInfoProfile.set_mugshot') }}
            </UButton>
        </UTooltip>

        <UTooltip
            v-if="can('documents.DocumentsService/UpdateDocument').value"
            :text="$t('components.citizens.CitizenInfoProfile.create_new_document')"
            :kbds="['C', 'D']"
        >
            <UButton block icon="i-mdi-file-document-plus" @click="openTemplates()">
                {{ $t('components.citizens.CitizenInfoProfile.create_new_document') }}
            </UButton>
        </UTooltip>

        <UButton
            v-if="
                activeChar?.job === user.job &&
                can('jobs.JobsService/GetColleague').value &&
                checkIfCanAccessColleague(user, 'jobs.JobsService/GetColleague')
            "
            block
            icon="i-mdi-account-circle"
            :to="`/jobs/colleagues/${user.userId}/info`"
        >
            {{ $t('components.citizens.CitizenInfoProfile.go_to_colleague_info') }}
        </UButton>

        <USeparator />

        <UButton block icon="i-mdi-link-variant" @click="copyLinkToClipboard()">
            {{ $t('components.citizens.CitizenInfoProfile.copy_profile_link') }}
        </UButton>
    </div>
</template>
