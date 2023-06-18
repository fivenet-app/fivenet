<script setup lang="ts">
import SvgIcon from '@jamescoyle/vue-icon';
import {
    mdiAt,
    mdiBellAlert,
    mdiBellSleep,
    mdiBriefcase,
    mdiDoorOpen,
    mdiHandcuffs,
    mdiHelpCircle,
    mdiLicense,
    mdiReceiptTextCheck,
    mdiReceiptTextPlus,
    mdiReceiptTextRemove,
    mdiRunFast,
    mdiTrafficCone,
} from '@mdi/js';
import { RoutesNamedLocations } from '~~/.nuxt/typed-router/__routes';
import { UserActivity } from '~~/gen/ts/resources/users/users';

const props = defineProps<{
    activity: UserActivity;
}>();

const { t } = useI18n();

const icon = ref<string>(mdiHelpCircle);
const iconColor = ref<string>('text-neutral');
const actionText = ref<string>(props.activity.key);
const actionValue = ref<string>(`${props.activity.oldValue} -> ${props.activity.newValue}`);
const actionReason = ref<string>(props.activity.reason);
const actionLink = ref<RoutesNamedLocations>();
const actionLinkText = ref<string>('');
const reasonText = ref<string>('components.citizens.citizen_info_activity_feed_entry.with_reason');

switch (props.activity.key) {
    case 'DocStore.Relation': {
        actionText.value = t('components.citizens.citizen_info_activity_feed_entry.document_relation.added');
        icon.value = mdiAt;
        actionLink.value = { name: 'documents-id', params: { id: 0 } };
        actionLinkText.value = t('common.document', 1);
        actionReason.value = t(`enums.docstore.DOC_RELATION.${props.activity.reason}`);

        if (props.activity.newValue !== '') {
            iconColor.value = 'text-info-600';
            actionLink.value.params.id = props.activity.newValue;
        } else if (props.activity.oldValue !== '') {
            actionText.value = t('components.citizens.citizen_info_activity_feed_entry.document_relation.removed');
            iconColor.value = 'text-base-600';
            actionLink.value.params.id = props.activity.oldValue;
        }
        break;
    }

    case 'UserProps.Wanted': {
        actionText.value = t('components.citizens.citizen_info_activity_feed_entry.bool_set_citizen');
        actionValue.value =
            props.activity.newValue === 'true' ? t('common.wanted') : `${t('common.not').toLowerCase()} ${t('common.wanted')}`;

        if (props.activity.newValue === 'true') {
            icon.value = mdiBellAlert;
            iconColor.value = 'text-error-400';
        } else {
            icon.value = mdiBellSleep;
            iconColor.value = 'text-success-400';
        }
        break;
    }

    case 'UserProps.Job': {
        actionText.value = t('components.citizens.citizen_info_activity_feed_entry.userprops_job_set');
        actionValue.value = props.activity.newValue;
        icon.value = mdiBriefcase;
        iconColor.value = 'text-secondary-400';
        break;
    }

    case 'UserProps.TrafficInfractionPoints': {
        actionText.value = t('components.citizens.citizen_info_activity_feed_entry.traffic_infraction_points.action_text');
        actionValue.value = `${props.activity.oldValue} ${t('common.to').toLocaleLowerCase()} ${props.activity.newValue}`;
        icon.value = mdiTrafficCone;
        iconColor.value = 'text-secondary-400';
        break;
    }

    case 'Plugin.Licenses': {
        reasonText.value = '';
        icon.value = mdiLicense;
        if (props.activity.newValue !== '') {
            actionText.value = t('components.citizens.citizen_info_activity_feed_entry.plugin_license.added');
            actionValue.value = '';
            iconColor.value = 'text-info-600';
        } else {
            actionText.value = t('components.citizens.citizen_info_activity_feed_entry.plugin_license.removed');
            actionValue.value = '';
            iconColor.value = 'text-warn-600';
        }
        break;
    }

    case 'Plugin.Jail': {
        reasonText.value = '';
        actionValue.value = '';
        if (props.activity.oldValue === '' && props.activity.newValue !== '0') {
            icon.value = mdiHandcuffs;
            actionText.value = t('components.citizens.citizen_info_activity_feed_entry.plugin_jail.jailed');
            actionValue.value = fromSecondsToFormattedDuration(parseInt(props.activity.newValue));
        } else if (props.activity.newValue === '0') {
            icon.value = mdiDoorOpen;
            actionText.value = t('components.citizens.citizen_info_activity_feed_entry.plugin_jail.unjailed');
        } else {
            icon.value = mdiRunFast;
            actionText.value = t('components.citizens.citizen_info_activity_feed_entry.plugin_jail.escaped');
        }
        break;
    }

    case 'Plugin.Billing.Fines': {
        if (props.activity.newValue === '0') {
            icon.value = mdiReceiptTextCheck;
            actionText.value = t('components.citizens.citizen_info_activity_feed_entry.plugin_billing_fines.paid');
            actionValue.value = '$' + props.activity.oldValue;
            iconColor.value = 'text-green-400';
        } else if (props.activity.newValue === props.activity.oldValue) {
            icon.value = mdiReceiptTextRemove;
            actionText.value = t('components.citizens.citizen_info_activity_feed_entry.plugin_billing_fines.removed');
            actionValue.value = '$' + props.activity.oldValue;
            iconColor.value = 'text-secondary-400';
        } else {
            icon.value = mdiReceiptTextPlus;
            actionText.value = t('components.citizens.citizen_info_activity_feed_entry.plugin_billing_fines.created');
            actionValue.value = '$' + props.activity.newValue;
            iconColor.value = 'text-info-400';
        }
        break;
    }
}
</script>

<template>
    <div class="flex space-x-3">
        <div class="h-10 w-10 rounded-full flex items-center justify-center my-auto">
            <SvgIcon :class="[iconColor, 'w-full h-full']" type="mdi" :path="icon" />
        </div>
        <div class="flex-1 space-y-1">
            <div class="flex items-center justify-between">
                <h3 class="text-sm font-medium text-neutral">
                    {{ actionText }}
                    <span class="font-bold">
                        <NuxtLink v-if="actionLink" :to="actionLink">
                            {{ actionLinkText }}
                        </NuxtLink>
                        <span v-else v-html="actionValue"></span>
                    </span>
                    <span v-if="reasonText">
                        {{ ' ' + $t(reasonText) }}
                    </span>
                    <span v-if="actionReason">
                        <span class="font-bold">
                            {{ ' ' + actionReason }}
                        </span>
                    </span>
                </h3>
                <p class="text-sm text-gray-400">
                    {{ useLocaleTimeAgo(toDate(activity.createdAt)!).value }}
                </p>
            </div>
            <p class="text-sm text-gray-300">
                {{ $t('common.created_by') + ' ' }}
                <NuxtLink
                    :to="{ name: 'citizens-id', params: { id: activity.sourceUser?.userId! } }"
                    class="underline decoration-solid"
                >
                    {{ activity.sourceUser?.firstname }}
                    {{ activity.sourceUser?.lastname }}
                </NuxtLink>
            </p>
        </div>
    </div>
</template>
