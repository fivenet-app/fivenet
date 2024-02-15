<script lang="ts" setup>
import { ClipboardPlusIcon } from 'mdi-vue3';
import VueJsonPretty from 'vue-json-pretty';
import 'vue-json-pretty/lib/styles.css';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useNotificatorStore } from '~/store/notificator';
import { AuditEntry, EventType } from '~~/gen/ts/resources/rector/audit';

const { d } = useI18n();

const props = defineProps<{
    log: AuditEntry;
}>();

const notifications = useNotificatorStore();

async function addToClipboard(): Promise<void> {
    const user = props.log.user;
    let text = `**Audit Log Entry ${props.log.id} - ${d(toDate(props.log.createdAt)!, 'short')}**

`;
    if (user) {
        text += `User: ${user?.firstname} ${user?.lastname} (${user?.userId}; ${user?.identifier})
`;
    }
    text += `Action: \`${props.log.service}/${props.log.method}\`
Event: \`${EventType[props.log.state]}\`
`;
    if (props.log.data) {
        text += `Data:
\`\`\`json
${jsonStringify(jsonParse(props.log.data!), 2)}
\`\`\`
`;
    } else {
        text += `Data: N/A
`;
    }

    notifications.dispatchNotification({
        title: { key: 'notifications.rector.audit_log.title', parameters: {} },
        content: { key: 'notifications.rector.audit_log.content', parameters: {} },
        type: 'info',
    });

    return copyToClipboardWrapper(text);
}
</script>

<template>
    <tr class="transition-colors even:bg-base-800 hover:bg-neutral/5">
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-1">
            {{ log.id }}
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-1">
            <GenericTime :value="log.createdAt" type="long" />
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-1">
            <CitizenInfoPopover :user="log.user" />
        </td>
        <td class="break-all py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-1">{{ log.service }}/{{ log.method }}</td>
        <td class="py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-1">
            {{ EventType[log.state] }}
        </td>
        <td class="max-w-3xl py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-1">
            <span v-if="!log.data">N/A</span>
            <span v-else>
                <VueJsonPretty
                    :data="jsonParse(props.log.data!) as any"
                    :show-icon="true"
                    :show-length="true"
                    :virtual="true"
                    :height="160"
                />
            </span>
        </td>
        <td class="break-all py-2 pl-3 pr-4 text-right text-sm font-medium sm:pr-0">
            <button
                class="flex-initial text-primary-500 hover:text-primary-400"
                :title="$t('components.clipboard.clipboard_button.add')"
            >
                <ClipboardPlusIcon class="ml-auto mr-2.5 w-5 h-auto" @click="addToClipboard" />
            </button>
        </td>
    </tr>
</template>
