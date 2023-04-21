// Code generated by protoc-gen-customizerweb. DO NOT EDIT.
// source: resources/rector/audit.proto

import { EVENT_TYPE } from './audit_pb';

// EVENT_TYPE
export class EVENT_TYPE_Util {
    public static toEnumKey(input: EVENT_TYPE): string | undefined {
        const index = Object.values(EVENT_TYPE).indexOf(input);
        if (index <= -1) {
            return "N/A";
        }
        return Object.keys(EVENT_TYPE)[index];
    }

    public static fromInt(input: Number): EVENT_TYPE {
        switch (input) {
            case 0:
                return EVENT_TYPE.UNKNOWN;
            
            case 1:
                return EVENT_TYPE.ERRORED;
            
            case 2:
                return EVENT_TYPE.VIEWED;
            
            case 3:
                return EVENT_TYPE.CREATED;
            
            case 4:
                return EVENT_TYPE.UPDATED;
            
            case 5:
                return EVENT_TYPE.DELETED;
            }
        return EVENT_TYPE.UNKNOWN;
    }

    public static fromString(input: String): EVENT_TYPE {
        switch (input) {
            case 'UNKNOWN':
                return EVENT_TYPE.UNKNOWN;
            
            case 'ERRORED':
                return EVENT_TYPE.ERRORED;
            
            case 'VIEWED':
                return EVENT_TYPE.VIEWED;
            
            case 'CREATED':
                return EVENT_TYPE.CREATED;
            
            case 'UPDATED':
                return EVENT_TYPE.UPDATED;
            
            case 'DELETED':
                return EVENT_TYPE.DELETED;
            }
        return EVENT_TYPE.UNKNOWN;
    }
}
