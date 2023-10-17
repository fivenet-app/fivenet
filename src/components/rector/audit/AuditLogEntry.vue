<script lang="ts" setup>
import { useClipboard } from '@vueuse/core';
import { ClipboardPlusIcon } from 'mdi-vue3';
import VueJsonPretty from 'vue-json-pretty';
import 'vue-json-pretty/lib/styles.css';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import Time from '~/components/partials/elements/Time.vue';
import { useNotificatorStore } from '~/store/notificator';
import { AuditEntry, EventType } from '~~/gen/ts/resources/rector/audit';

const clipboard = useClipboard();

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
        text += `User: ${user?.firstname}, ${user?.lastname} (${user?.userId}; ${user?.identifier})
`;
    }
    text += `Action: \`${props.log.service}/${props.log.method}\`
Event: \`${EventType[props.log.state]}\`
`;
    if (props.log.data) {
        text += `Data:
\`\`\`json
${JSON.stringify(JSON.parse(props.log.data!), null, 2)}
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

    return clipboard.copy(text);
}
</script>

<template>
    <tr>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ log.id }}
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            <Time :value="log.createdAt" type="long" />
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            <CitizenInfoPopover :user="log.user" />
        </td>
        <td class="break-all py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">{{ log.service }}/{{ log.method }}</td>
        <td class="py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ EventType[log.state] }}
        </td>
        <td class="py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0 max-w-3xl">
            <span v-if="!log.data">N/A</span>
            <span v-else>
                <VueJsonPretty
                    :data="JSON.parse(props.log.data!) as any"
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
                <ClipboardPlusIcon class="w-6 h-auto ml-auto mr-2.5" @click="addToClipboard" />
            </button>
        </td>
    </tr>
</template>
