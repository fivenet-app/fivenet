<script setup lang="ts">
import {
    AtIcon,
    BellAlertIcon,
    BellSleepIcon,
    BriefcaseIcon,
    DoorOpenIcon,
    HandcuffsIcon,
    HelpCircleIcon,
    LicenseIcon,
    ReceiptTextCheckIcon,
    ReceiptTextPlusIcon,
    ReceiptTextRemoveIcon,
    RunFastIcon,
    TrafficConeIcon,
} from 'mdi-vue3';
import { DefineComponent } from 'vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import Time from '~/components/partials/elements/Time.vue';
import { RoutesNamedLocations } from '~~/.nuxt/typed-router/__routes';
import { UserActivity } from '~~/gen/ts/resources/users/users';

const props = defineProps<{
    activity: UserActivity;
}>();

const { t, n } = useI18n();

let icon: DefineComponent = markRaw(HelpCircleIcon);
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
        icon = markRaw(AtIcon);
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
            icon = BellAlertIcon;
            iconColor.value = 'text-error-400';
        } else {
            icon = BellSleepIcon;
            iconColor.value = 'text-success-400';
        }
        break;
    }

    case 'UserProps.Job': {
        actionText.value = t('components.citizens.citizen_info_activity_feed_entry.userprops_job_set');
        actionValue.value = props.activity.newValue;
        icon = BriefcaseIcon;
        iconColor.value = 'text-secondary-400';
        break;
    }

    case 'UserProps.TrafficInfractionPoints': {
        actionText.value = t('components.citizens.citizen_info_activity_feed_entry.traffic_infraction_points.action_text');
        actionValue.value = `${props.activity.oldValue} ${t('common.to').toLocaleLowerCase()} ${props.activity.newValue}`;
        icon = TrafficConeIcon;
        iconColor.value = 'text-secondary-400';
        break;
    }

    case 'Plugin.Licenses': {
        reasonText.value = '';
        icon = LicenseIcon;
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
            icon = HandcuffsIcon;
            actionText.value = t('components.citizens.citizen_info_activity_feed_entry.plugin_jail.jailed');
            actionValue.value = fromSecondsToFormattedDuration(parseInt(props.activity.newValue));
        } else if (props.activity.newValue === '0') {
            icon = DoorOpenIcon;
            actionText.value = t('components.citizens.citizen_info_activity_feed_entry.plugin_jail.unjailed');
        } else {
            icon = RunFastIcon;
            actionText.value = t('components.citizens.citizen_info_activity_feed_entry.plugin_jail.escaped');
        }
        break;
    }

    case 'Plugin.Billing.Fines': {
        if (props.activity.newValue === '0') {
            icon = ReceiptTextCheckIcon;
            actionText.value = t('components.citizens.citizen_info_activity_feed_entry.plugin_billing_fines.paid');
            actionValue.value = n(parseInt(props.activity.oldValue), 'currency');
            iconColor.value = 'text-green-400';
        } else if (props.activity.newValue === props.activity.oldValue) {
            icon = ReceiptTextRemoveIcon;
            actionText.value = t('components.citizens.citizen_info_activity_feed_entry.plugin_billing_fines.removed');
            actionValue.value = n(parseInt(props.activity.oldValue), 'currency');
            iconColor.value = 'text-secondary-400';
        } else {
            icon = ReceiptTextPlusIcon;
            actionText.value = t('components.citizens.citizen_info_activity_feed_entry.plugin_billing_fines.created');
            actionValue.value = n(parseInt(props.activity.newValue), 'currency');
            iconColor.value = 'text-info-400';
        }
        break;
    }
}
</script>

<template>
    <div class="flex space-x-3">
        <div class="h-10 w-10 rounded-full flex items-center justify-center my-auto">
            <component :is="icon" :class="[iconColor, 'w-full h-full']" />
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
                    <Time :value="activity.createdAt" type="long" />
                </p>
            </div>
            <p class="text-sm text-gray-300 inline-flex">
                {{ $t('common.created_by') }}
                &nbsp;
                <CitizenInfoPopover :user="activity.sourceUser" />
            </p>
        </div>
    </div>
</template>
