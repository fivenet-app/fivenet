<script lang="ts" setup>
import { useClipboardStore } from '~/store/clipboard';
import { useNotificatorStore } from '~/store/notificator';
import TemplatesModal from '~/components/documents/templates/TemplatesModal.vue';
import CitizenSetJobModal from '~/components/citizens/info/props/CitizenSetJobModal.vue';
import CitizenSetTrafficPointsModal from '~/components/citizens/info/props/CitizenSetTrafficPointsModal.vue';
import CitizenSetWantedModal from '~/components/citizens/info/props/CitizenSetWantedModal.vue';
import CitizenSetMugShotModal from '~/components/citizens/info/props/CitizenSetMugShotModal.vue';
import type { Job, JobGrade } from '~~/gen/ts/resources/users/jobs';
import type { User } from '~~/gen/ts/resources/users/users';
import type { File } from '~~/gen/ts/resources/filestore/file';

const props = withDefaults(
    defineProps<{
        user: User;
        registerShortcuts?: boolean;
    }>(),
    {
        registerShortcuts: false,
    },
);

const emits = defineEmits<{
    (e: 'update:wantedStatus', value: boolean): void;
    (e: 'update:job', value: { job: Job; grade: JobGrade }): void;
    (e: 'update:trafficInfractionPoints', value: number): void;
    (e: 'update:mugShot', value?: File): void;
}>();

const clipboardStore = useClipboardStore();
const notifications = useNotificatorStore();

const w = window;

const modal = useModal();

function openTemplates(): void {
    if (!props.user) {
        return;
    }

    clipboardStore.addUser(props.user);

    modal.open(TemplatesModal, {});
}

function copyLinkToClipboard(): void {
    copyToClipboardWrapper(w.location.href);

    notifications.add({
        title: { key: 'notifications.clipboard.link_copied.title', parameters: {} },
        description: { key: 'notifications.clipboard.link_copied.content', parameters: {} },
        timeout: 3250,
        type: 'info',
    });
}

if (props.registerShortcuts) {
    defineShortcuts({
        'c-w': () => {
            if (!attr('CitizenStoreService.SetUserProps', 'Fields', 'Wanted')) {
                return;
            }

            modal.open(CitizenSetWantedModal, {
                user: props.user,
                'onUpdate:wantedStatus': ($event) => emits('update:wantedStatus', $event),
            });
        },
        'c-j': () => {
            if (!attr('CitizenStoreService.SetUserProps', 'Fields', 'Job')) {
                return;
            }

            modal.open(CitizenSetJobModal, {
                user: props.user,
                'onUpdate:job': ($event) => emits('update:job', $event),
            });
        },
        'c-p': () => {
            if (!attr('CitizenStoreService.SetUserProps', 'Fields', 'TrafficInfractionPoints')) {
                return;
            }

            modal.open(CitizenSetTrafficPointsModal, {
                user: props.user,
                'onUpdate:trafficInfractionPoints': ($event) => emits('update:trafficInfractionPoints', $event),
            });
        },
        'c-m': () => {
            if (!attr('CitizenStoreService.SetUserProps', 'Fields', 'MugShot')) {
                return;
            }

            modal.open(CitizenSetMugShotModal, {
                user: props.user,
                'onUpdate:mugShot': ($event) => emits('update:mugShot', $event),
            });
        },
        'c-d': () => {
            if (!can('DocStoreService.CreateDocument')) {
                return;
            }

            openTemplates();
        },
    });
}
</script>

<template>
    <div class="flex w-full min-w-80 flex-col gap-2">
        <UTooltip
            v-if="attr('CitizenStoreService.SetUserProps', 'Fields', 'Wanted')"
            :text="
                user?.props?.wanted
                    ? $t('components.citizens.CitizenInfoProfile.revoke_wanted')
                    : $t('components.citizens.CitizenInfoProfile.set_wanted')
            "
            :shortcuts="['C', 'W']"
        >
            <UButton
                color="red"
                block
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
            v-if="attr('CitizenStoreService.SetUserProps', 'Fields', 'Job')"
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
            v-if="attr('CitizenStoreService.SetUserProps', 'Fields', 'TrafficInfractionPoints')"
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
            v-if="attr('CitizenStoreService.SetUserProps', 'Fields', 'MugShot')"
            :text="
                user?.props?.wanted
                    ? $t('components.citizens.CitizenInfoProfile.revoke_wanted')
                    : $t('components.citizens.CitizenInfoProfile.set_wanted')
            "
            :shortcuts="['C', 'M']"
        >
            <UButton
                block
                icon="i-mdi-camera"
                @click="
                    modal.open(CitizenSetMugShotModal, {
                        user: user,
                        'onUpdate:mugShot': ($event) => $emit('update:mugShot', $event),
                    })
                "
            >
                {{ $t('components.citizens.CitizenInfoProfile.set_mug_shot') }}
            </UButton>
        </UTooltip>

        <UTooltip
            v-if="can('DocStoreService.CreateDocument')"
            :text="
                user?.props?.wanted
                    ? $t('components.citizens.CitizenInfoProfile.revoke_wanted')
                    : $t('components.citizens.CitizenInfoProfile.set_wanted')
            "
            :shortcuts="['C', 'D']"
        >
            <UButton block icon="i-mdi-file-document-plus" @click="openTemplates()">
                {{ $t('components.citizens.CitizenInfoProfile.create_new_document') }}
            </UButton>
        </UTooltip>

        <UDivider />

        <UButton block icon="i-mdi-link" @click="copyLinkToClipboard()">
            {{ $t('components.citizens.CitizenInfoProfile.copy_profile_link') }}
        </UButton>
    </div>
</template>