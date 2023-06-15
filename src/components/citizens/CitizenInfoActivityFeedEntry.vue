<script setup lang="ts">
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiAt, mdiBellAlert, mdiBellSleep, mdiBriefcase, mdiHandcuffs, mdiHelpCircle, mdiLicense, mdiReceipt } from '@mdi/js';
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

switch (props.activity.key) {
    case 'UserProps.Wanted': {
        actionText.value = t('components.citizens.citizen_info_activity_feed_entry.set_citizen_as');
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
        actionText.value = t('components.citizens.citizen_info_activity_feed_entry.set_ciizen_job');
        actionValue.value = props.activity.newValue;
        icon.value = mdiBriefcase;
        iconColor.value = 'text-secondary-400';

        break;
    }

    case 'DocStore.Relation': {
        actionText.value = t('components.citizens.citizen_info_activity_feed_entry.document_relation');
        icon.value = mdiAt;
        actionLink.value = { name: 'documents-id', params: { id: 0 } };
        actionLinkText.value = t('common.document', 1);
        actionReason.value = t(`enums.docstore.DOC_RELATION.${props.activity.reason}`);

        if (props.activity.newValue !== '') {
            iconColor.value = 'text-info-600';
            actionLink.value.params.id = props.activity.newValue;
        } else if (props.activity.oldValue !== '') {
            actionText.value = t('components.citizens.citizen_info_activity_feed_entry.document_relation_removed');
            iconColor.value = 'text-base-600';
            actionLink.value.params.id = props.activity.oldValue;
        }

        break;
    }

    case 'Plugin.Licenses': {
        // TODO
        icon.value = mdiLicense;
        break;
    }

    case 'Plugin.Jail': {
        // TODO
        icon.value = mdiHandcuffs;
        break;
    }

    case 'Plugin.Fines': {
        // TODO
        icon.value = mdiReceipt;
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
                    {{ activity.sourceUser?.firstname }}
                    {{ activity.sourceUser?.lastname }}
                </h3>
                <p class="text-sm text-gray-400">
                    {{ useLocaleTimeAgo(toDate(activity.createdAt)!).value }}
                </p>
            </div>
            <p class="text-sm text-gray-300">
                {{ actionText }}
                <span class="font-bold">
                    <NuxtLink v-if="actionLink" :to="actionLink">
                        {{ actionLinkText }}
                    </NuxtLink>
                    <span v-else v-html="actionValue"></span>
                </span>
                <span v-if="actionReason">
                    {{ ' ' + $t('components.citizens.citizen_info_activity_feed_entry.with_reason') }}
                    <span class="font-bold">
                        {{ actionReason }}
                    </span>
                </span>
            </p>
        </div>
    </div>
</template>
