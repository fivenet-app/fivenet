<script lang="ts" setup>
import { AuditEntry } from '@fivenet/gen/resources/rector/audit_pb';
import { toDateRelativeString } from '~/utils/time';
import { EVENT_TYPE_Util } from '@fivenet/gen/resources/rector/audit.pb_enums';

defineProps({
    log: {
        type: AuditEntry,
        required: true,
    }
});
</script>

<template>
    <tr>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ log.getId() }}
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ toDateRelativeString(log.getCreatedAt()) }}
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ log.hasUser() ? (log.getUser()?.getFirstname() + ' ' + log.getUser()?.getLastname()) : 'N/A' }}
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ log.getService() }}/{{ log.getMethod() }}
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ EVENT_TYPE_Util.toEnumKey(log.getState()) }}
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ log.getData() ? log.getData() : 'N/A' }}
        </td>
        <td class="whitespace-nowrap py-2 pl-3 pr-4 text-right text-sm font-medium sm:pr-0">
            {{ $t('common.copy').toUpperCase() }}
        </td>
    </tr>
</template>
