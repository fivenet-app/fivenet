<script lang="ts" setup>
import CitizenSetJobModal from '~/components/citizens/info/props/CitizenSetJobModal.vue';
import CitizenSetMugshotModal from '~/components/citizens/info/props/CitizenSetMugshotModal.vue';
import CitizenSetTrafficPointsModal from '~/components/citizens/info/props/CitizenSetTrafficPointsModal.vue';
import CitizenSetWantedModal from '~/components/citizens/info/props/CitizenSetWantedModal.vue';
import TemplatesModal from '~/components/documents/templates/TemplatesModal.vue';
import { checkIfCanAccessColleague } from '~/components/jobs/colleagues/helpers';
import { useClipboardStore } from '~/stores/clipboard';
import type { File } from '~~/gen/ts/resources/file/file';
import type { Job, JobGrade } from '~~/gen/ts/resources/jobs/jobs';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { User } from '~~/gen/ts/resources/users/users';

const props = withDefaults(
    defineProps<{
        user: User;
        registerShortcuts?: boolean;
    }>(),
    {
        registerShortcuts: false,
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

const w = window;

const modal = useModal();

function openTemplates(): void {
    if (!props.user) {
        return;
    }

    clipboardStore.addUser(props.user, true);

    modal.open(TemplatesModal, {});
}

function copyLinkToClipboard(): void {
    copyToClipboardWrapper(w.location.href);

    notifications.add({
        title: { key: 'notifications.clipboard.link_copied.title', parameters: {} },
        description: { key: 'notifications.clipboard.link_copied.content', parameters: {} },
        timeout: 3250,
        type: NotificationType.INFO,
    });
}

if (props.registerShortcuts) {
    defineShortcuts({
        'c-w': () => {
            if (!attr('citizens.CitizensService/SetUserProps', 'Fields', 'Wanted').value) {
                return;
            }

            modal.open(CitizenSetWantedModal, {
                user: props.user,
                'onUpdate:wantedStatus': ($event) => emit('update:wantedStatus', $event),
            });
        },
        'c-j': () => {
            if (!attr('citizens.CitizensService/SetUserProps', 'Fields', 'Job').value) {
                return;
            }

            modal.open(CitizenSetJobModal, {
                user: props.user,
                'onUpdate:job': ($event) => emit('update:job', $event),
            });
        },
        'c-p': () => {
            if (!attr('citizens.CitizensService/SetUserProps', 'Fields', 'TrafficInfractionPoints').value) {
                return;
            }

            modal.open(CitizenSetTrafficPointsModal, {
                user: props.user,
                'onUpdate:trafficInfractionPoints': ($event) => emit('update:trafficInfractionPoints', $event),
            });
        },
        'c-m': () => {
            if (!attr('citizens.CitizensService/SetUserProps', 'Fields', 'Mugshot').value) {
                return;
            }

            modal.open(CitizenSetMugshotModal, {
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
            :shortcuts="['C', 'W']"
        >
            <UButton
                :color="user?.props?.wanted ? 'error' : 'primary'"
                block
                truncate
                :icon="user?.props?.wanted ? 'i-mdi-account-alert' : 'i-mdi-account-cancel'"
                @click="
                    modal.open(CitizenSetWantedModal, {
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
            :shortcuts="['C', 'J']"
        >
            <UButton
                block
                icon="i-mdi-briefcase"
                @click="
                    modal.open(CitizenSetJobModal, {
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
            :shortcuts="['C', 'P']"
        >
            <UButton
                block
                icon="i-mdi-counter"
                @click="
                    modal.open(CitizenSetTrafficPointsModal, {
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
            :shortcuts="['C', 'M']"
        >
            <UButton
                block
                icon="i-mdi-camera"
                @click="
                    modal.open(CitizenSetMugshotModal, {
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
            :shortcuts="['C', 'D']"
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

        <UDivider />

        <UButton block icon="i-mdi-link-variant" @click="copyLinkToClipboard()">
            {{ $t('components.citizens.CitizenInfoProfile.copy_profile_link') }}
        </UButton>
    </div>
</template>
