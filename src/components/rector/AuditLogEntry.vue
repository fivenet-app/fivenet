<script lang="ts" setup>
import { AuditEntry } from '@fivenet/gen/resources/rector/audit_pb';
import { EVENT_TYPE_Util } from '@fivenet/gen/resources/rector/audit.pb_enums';
import { ClipboardDocumentIcon } from '@heroicons/vue/24/solid';

const { d } = useI18n();

const props = defineProps({
    log: {
        type: AuditEntry,
        required: true,
    }
});

async function addToClipboard(): Promise<void> {
    const user = props.log.getUser();
    const text = `**Audit Log Entry ${props.log.getId()} - ${d(props.log.getCreatedAt()?.getTimestamp()?.toDate()!, 'short')}**
User: ${user?.getFirstname()}, ${user?.getLastname()} (${user?.getUserId()}; ${user?.getIdentifier()})
Action: ${props.log.getMethod()}/${props.log.getService()}
Event: ${EVENT_TYPE_Util.toEnumKey(props.log.getState())}
Data:
\`\`\`
${props.log.getData()}
\`\`\`
`;

    return navigator.clipboard.writeText(text);
}
</script>

<template>
    <tr>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ log.getId() }}
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ $d(log.getCreatedAt()?.getTimestamp()?.toDate()!, 'short') }}
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ log.hasUser() ? (log.getUser()?.getFirstname() + ' ' + log.getUser()?.getLastname()) : 'N/A' }}
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ log.getService() }}: {{ log.getMethod() }}
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ EVENT_TYPE_Util.toEnumKey(log.getState()) }}
        </td>
        <td class="py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ log.getData() ? log.getData() : 'N/A' }}
        </td>
        <td class="whitespace-nowrap py-2 pl-3 pr-4 text-right text-sm font-medium sm:pr-0">
            <button class="flex-initial text-primary-500 hover:text-primary-400"
                :title="$t('components.clipboard.clipboard_button.add')">
                <ClipboardDocumentIcon class="w-6 h-auto ml-auto mr-2.5" @click="addToClipboard" />
            </button>
        </td>
    </tr>
</template>
