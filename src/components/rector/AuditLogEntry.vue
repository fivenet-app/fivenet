<script lang="ts" setup>
import { AuditEntry } from '@fivenet/gen/resources/rector/audit_pb';

defineProps({
    log: {
        type: AuditEntry,
        required: true,
    }
});
</script>

<template>
    <tr>
        <td>
            ID
        </td>
        <td>
            Time
        </td>
        <td>
            Type
        </td>
        <td>
            Actions
        </td>
    </tr>
</template>
