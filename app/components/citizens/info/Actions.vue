<script lang="ts" setup>
import SetJobDrawer from '~/components/citizens/info/props/SetJobDrawer.vue';
import SetMugshotDrawer from '~/components/citizens/info/props/SetMugshotDrawer.vue';
import SetTrafficPointsDrawer from '~/components/citizens/info/props/SetTrafficPointsDrawer.vue';
import SetWantedDrawer from '~/components/citizens/info/props/SetWantedDrawer.vue';
import TemplateDrawer from '~/components/documents/templates/TemplateDrawer.vue';
import SetDynamicPropertyDrawer from '~/components/partials/properties/SetDynamicPropertyDrawer.vue';
import { canActOnDynamicProperty, dynamicPropertiesForAction } from '~/composables/useDynamicPropertyAction';
import { checkIfCanAccessColleague } from '~/components/jobs/colleagues/helpers';
import { useClipboardStore } from '~/stores/clipboard';
import { PropertyScope, type Properties, type Property } from '~~/gen/ts/resources/properties/properties';
import type { File } from '~~/gen/ts/resources/file/file';
import type { Job, JobGrade } from '~~/gen/ts/resources/jobs/jobs';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { UserProps } from '~~/gen/ts/resources/users/props/props';
import type { User } from '~~/gen/ts/resources/users/user';

const props = withDefaults(
    defineProps<{
        user: User;
        registerKBDs?: boolean;
    }>(),
    {
        registerKBDs: false,
    },
);

const emits = defineEmits<{
    (e: 'update:wantedStatus', value: UserProps): void;
    (e: 'update:job', value: { job: Job; grade: JobGrade }): void;
    (e: 'update:trafficInfractionPoints', value: number): void;
    (e: 'update:mugshot', value?: File): void;
    (e: 'update:properties', value: Properties | undefined): void;
}>();

const { attr, can, activeChar } = useAuth();
const { propertyLabel } = usePropertyRenderer();

const clipboardStore = useClipboardStore();
const { open: openClipboardModal } = useClipboardModal();

const notifications = useNotificationsStore();

const overlay = useOverlay();

const templateDrawer = overlay.create(TemplateDrawer);
const setWantedDrawer = overlay.create(SetWantedDrawer);
const setJobDrawer = overlay.create(SetJobDrawer);
const setTrafficPointsDrawer = overlay.create(SetTrafficPointsDrawer);
const setMugshotDrawer = overlay.create(SetMugshotDrawer);
const setDynamicPropertyDrawer = overlay.create(SetDynamicPropertyDrawer);

const actionVisibility = computed(() => ({
    wanted: attr('citizens.CitizensService/SetUserProps', 'Properties', 'user.wanted').value,
    job: attr('citizens.CitizensService/SetUserProps', 'Properties', 'user.job').value,
    trafficPoints: attr('citizens.CitizensService/SetUserProps', 'Properties', 'user.traffic_infraction_points').value,
    mugshot: attr('citizens.CitizensService/SetUserProps', 'Properties', 'user.mugshot').value,
    document: can('documents.DocumentsService/UpdateDocument').value,
    colleague:
        activeChar.value?.job === props.user.job &&
        can('jobs.ColleaguesService/GetColleague').value &&
        checkIfCanAccessColleague(props.user, 'jobs.ColleaguesService/GetColleague'),
}));

const hasVisibleAction = computed(() => Object.values(actionVisibility.value).some(Boolean));
const dynamicProperties = computed(() =>
    dynamicPropertiesForAction(
        props.user.properties?.properties,
        (property) => attr('citizens.CitizensService/SetUserProps', 'Properties', property.key).value,
    ),
);

function openDynamicProperty(property: Property): void {
    setDynamicPropertyDrawer.open({
        property,
        scope: PropertyScope.USER,
        userId: props.user.userId,
        'onUpdate:properties': (value) => emits('update:properties', value),
    });
}

function openTemplates(): void {
    if (!props.user) return;

    const added = clipboardStore.addUser(props.user, true);

    if (!added) {
        notifications.add({
            title: { key: 'notifications.clipboard.limit_reached.title', parameters: {} },
            description: { key: 'notifications.clipboard.limit_reached.content', parameters: {} },
            duration: 3250,
            type: NotificationType.WARNING,
            actions: [
                {
                    label: { key: 'common.open', parameters: {} },
                    icon: 'i-mdi-clipboard-list-outline',
                    onClick: () => void openClipboardModal(),
                },
            ],
        });
    }

    templateDrawer.open({});
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

defineShortcuts({
    'c-w': () => {
        if (!attr('citizens.CitizensService/SetUserProps', 'Properties', 'user.wanted').value) return;

        setWantedDrawer.open({
            user: props.user,
            'onUpdate:wantedStatus': ($event) => emits('update:wantedStatus', $event),
        });
    },
    'c-j': () => {
        if (!attr('citizens.CitizensService/SetUserProps', 'Properties', 'user.job').value) return;

        setJobDrawer.open({
            user: props.user,
            'onUpdate:job': ($event) => emits('update:job', $event),
        });
    },
    'c-p': () => {
        if (!attr('citizens.CitizensService/SetUserProps', 'Properties', 'user.traffic_infraction_points').value) return;

        setTrafficPointsDrawer.open({
            user: props.user,
            'onUpdate:trafficInfractionPoints': ($event) => emits('update:trafficInfractionPoints', $event),
        });
    },
    'c-m': () => {
        if (!attr('citizens.CitizensService/SetUserProps', 'Properties', 'user.mugshot').value) return;

        setMugshotDrawer.open({
            user: props.user,
            'onUpdate:mugshot': ($event) => emits('update:mugshot', $event),
        });
    },
    'c-d': () => {
        if (!can('documents.DocumentsService/UpdateDocument').value) return;

        openTemplates();
    },
});
</script>

<template>
    <div class="flex w-full flex-col gap-2">
        <UTooltip
            v-if="actionVisibility.wanted"
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
                :label="
                    user?.props?.wanted
                        ? $t('components.citizens.CitizenInfoProfile.revoke_wanted')
                        : $t('components.citizens.CitizenInfoProfile.set_wanted')
                "
                @click="
                    setWantedDrawer.open({
                        user: user,
                        'onUpdate:wantedStatus': ($event) => $emit('update:wantedStatus', $event),
                    })
                "
            />
        </UTooltip>

        <USeparator
            v-if="
                actionVisibility.wanted && (actionVisibility.job || actionVisibility.trafficPoints || actionVisibility.mugshot)
            "
        />

        <UTooltip v-if="actionVisibility.job" :text="$t('components.citizens.CitizenInfoProfile.set_job')" :kbds="['C', 'J']">
            <UButton
                color="neutral"
                variant="outline"
                block
                icon="i-mdi-briefcase"
                :label="$t('components.citizens.CitizenInfoProfile.set_job')"
                @click="
                    setJobDrawer.open({
                        user: user,
                        'onUpdate:job': ($event) => $emit('update:job', $event),
                    })
                "
            />
        </UTooltip>

        <UTooltip
            v-if="actionVisibility.trafficPoints"
            :text="$t('components.citizens.CitizenInfoProfile.set_traffic_points')"
            :kbds="['C', 'P']"
        >
            <UButton
                color="neutral"
                variant="outline"
                block
                icon="i-mdi-counter"
                :label="$t('components.citizens.CitizenInfoProfile.set_traffic_points')"
                @click="
                    setTrafficPointsDrawer.open({
                        user: user,
                        'onUpdate:trafficInfractionPoints': ($event) => $emit('update:trafficInfractionPoints', $event),
                    })
                "
            />
        </UTooltip>

        <UTooltip
            v-if="actionVisibility.mugshot"
            :text="$t('components.citizens.CitizenInfoProfile.set_mugshot')"
            :kbds="['C', 'M']"
        >
            <UButton
                color="neutral"
                variant="outline"
                block
                icon="i-mdi-camera"
                :label="$t('components.citizens.CitizenInfoProfile.set_mugshot')"
                @click="
                    setMugshotDrawer.open({
                        user: user,
                        'onUpdate:mugshot': ($event) => $emit('update:mugshot', $event),
                    })
                "
            />
        </UTooltip>

        <UTooltip
            v-if="actionVisibility.document"
            :text="$t('components.citizens.CitizenInfoProfile.create_new_document')"
            :kbds="['C', 'D']"
        >
            <UButton
                block
                icon="i-mdi-file-document-plus"
                :label="$t('components.citizens.CitizenInfoProfile.create_new_document')"
                @click="openTemplates()"
            />
        </UTooltip>

        <UButton
            v-if="actionVisibility.colleague"
            block
            icon="i-mdi-account-circle"
            :to="`/jobs/colleagues/${user.userId}/info`"
            :label="$t('components.citizens.CitizenInfoProfile.go_to_colleague_info')"
        />

        <template v-if="dynamicProperties.length">
            <USeparator />

            <UButton
                v-for="property in dynamicProperties"
                :key="property.key"
                color="neutral"
                variant="outline"
                block
                icon="i-mdi-pencil"
                :label="`${$t(canActOnDynamicProperty(property) ? 'common.edit' : 'common.set')}: ${propertyLabel(property)}`"
                @click="openDynamicProperty(property)"
            />
        </template>

        <USeparator v-if="hasVisibleAction || dynamicProperties.length" />

        <UButton
            block
            color="neutral"
            variant="subtle"
            icon="i-mdi-link-variant"
            :label="$t('components.citizens.CitizenInfoProfile.copy_profile_link')"
            @click="copyLinkToClipboard()"
        />

        <div v-if="registerKBDs" class="flex flex-wrap items-center gap-x-2 gap-y-1 text-xs text-muted">
            <span>{{ $t('common.shortcuts') }}:</span>
            <span v-if="actionVisibility.wanted" class="inline-flex items-center gap-0.5">
                <UKbd value="C" /><UKbd value="W" />
            </span>
            <span v-if="actionVisibility.job" class="inline-flex items-center gap-0.5">
                <UKbd value="C" /><UKbd value="J" />
            </span>
            <span v-if="actionVisibility.trafficPoints" class="inline-flex items-center gap-0.5">
                <UKbd value="C" /><UKbd value="P" />
            </span>
            <span v-if="actionVisibility.mugshot" class="inline-flex items-center gap-0.5">
                <UKbd value="C" /><UKbd value="M" />
            </span>
            <span v-if="actionVisibility.document" class="inline-flex items-center gap-0.5">
                <UKbd value="C" /><UKbd value="D" />
            </span>
        </div>
    </div>
</template>
