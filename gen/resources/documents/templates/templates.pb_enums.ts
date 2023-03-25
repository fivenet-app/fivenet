// Code generated by protoc-gen-customizerweb. DO NOT EDIT.
// source: resources/documents/templates/templates.proto

import * as enums from './templates_pb';


// DATA_ITEM_TYPE
export class DATA_ITEM_TYPE_Util {
    public static toEnumKey(input: enums.DATA_ITEM_TYPE): string | undefined {
        const index = Object.values(enums.DATA_ITEM_TYPE).indexOf(input);
        if (index <= -1) {
            return "N/A";
        }
        return Object.keys(enums.DATA_ITEM_TYPE)[index];
    }

    public static fromInt(input: Number): enums.DATA_ITEM_TYPE {
        switch (input) {
            case 0:
                return enums.DATA_ITEM_TYPE.USER;
            }
        return enums.DATA_ITEM_TYPE.USER;
    }

    public static fromString(input: String): enums.DATA_ITEM_TYPE {
        switch (input) {
            case 'USER':
                return enums.DATA_ITEM_TYPE.USER;
            }
        return enums.DATA_ITEM_TYPE.USER;
    }
}
