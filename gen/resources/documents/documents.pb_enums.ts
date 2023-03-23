// Code generated by protoc-gen-customizerweb. DO NOT EDIT.
// source: resources/documents/documents.proto

import * as enums from './documents_pb';


// DOC_CONTENT_TYPE
export class DOC_CONTENT_TYPE_Util {
    public static toEnumKey(input: enums.DOC_CONTENT_TYPE): string | undefined {
        const index = Object.values(enums.DOC_CONTENT_TYPE).indexOf(input);
        if (index <= -1) {
            return "N/A";
        }
        return Object.keys(enums.DOC_CONTENT_TYPE)[index];
    }

    public static fromInt(input: Number): enums.DOC_CONTENT_TYPE {
        switch (input) {
            case 0:
                return enums.DOC_CONTENT_TYPE.HTML;
            
            case 1:
                return enums.DOC_CONTENT_TYPE.PLAIN;
            }
        return;
    }

    public static fromString(input: String): enums.DOC_CONTENT_TYPE {
        switch (input) {
            case 'HTML':
                return enums.DOC_CONTENT_TYPE.HTML;
            
            case 'PLAIN':
                return enums.DOC_CONTENT_TYPE.PLAIN;
            }
        return enums.DOC_CONTENT_TYPE.HTML;
    }
}

// DOC_ACCESS
export class DOC_ACCESS_Util {
    public static toEnumKey(input: enums.DOC_ACCESS): string | undefined {
        const index = Object.values(enums.DOC_ACCESS).indexOf(input);
        if (index <= -1) {
            return "N/A";
        }
        return Object.keys(enums.DOC_ACCESS)[index];
    }

    public static fromInt(input: Number): enums.DOC_ACCESS {
        switch (input) {
            case 0:
                return enums.DOC_ACCESS.BLOCKED;
            
            case 1:
                return enums.DOC_ACCESS.VIEW;
            
            case 2:
                return enums.DOC_ACCESS.COMMENT;
            
            case 3:
                return enums.DOC_ACCESS.ACCESS;
            
            case 4:
                return enums.DOC_ACCESS.EDIT;
            }
        return;
    }

    public static fromString(input: String): enums.DOC_ACCESS {
        switch (input) {
            case 'BLOCKED':
                return enums.DOC_ACCESS.BLOCKED;
            
            case 'VIEW':
                return enums.DOC_ACCESS.VIEW;
            
            case 'COMMENT':
                return enums.DOC_ACCESS.COMMENT;
            
            case 'ACCESS':
                return enums.DOC_ACCESS.ACCESS;
            
            case 'EDIT':
                return enums.DOC_ACCESS.EDIT;
            }
        return enums.DOC_ACCESS.BLOCKED;
    }
}

// DOC_REFERENCE
export class DOC_REFERENCE_Util {
    public static toEnumKey(input: enums.DOC_REFERENCE): string | undefined {
        const index = Object.values(enums.DOC_REFERENCE).indexOf(input);
        if (index <= -1) {
            return "N/A";
        }
        return Object.keys(enums.DOC_REFERENCE)[index];
    }

    public static fromInt(input: Number): enums.DOC_REFERENCE {
        switch (input) {
            case 0:
                return enums.DOC_REFERENCE.LINKED;
            
            case 1:
                return enums.DOC_REFERENCE.SOLVES;
            
            case 2:
                return enums.DOC_REFERENCE.CLOSES;
            
            case 3:
                return enums.DOC_REFERENCE.DEPRECATES;
            }
        return;
    }

    public static fromString(input: String): enums.DOC_REFERENCE {
        switch (input) {
            case 'LINKED':
                return enums.DOC_REFERENCE.LINKED;
            
            case 'SOLVES':
                return enums.DOC_REFERENCE.SOLVES;
            
            case 'CLOSES':
                return enums.DOC_REFERENCE.CLOSES;
            
            case 'DEPRECATES':
                return enums.DOC_REFERENCE.DEPRECATES;
            }
        return enums.DOC_REFERENCE.LINKED;
    }
}

// DOC_RELATION
export class DOC_RELATION_Util {
    public static toEnumKey(input: enums.DOC_RELATION): string | undefined {
        const index = Object.values(enums.DOC_RELATION).indexOf(input);
        if (index <= -1) {
            return "N/A";
        }
        return Object.keys(enums.DOC_RELATION)[index];
    }

    public static fromInt(input: Number): enums.DOC_RELATION {
        switch (input) {
            case 0:
                return enums.DOC_RELATION.MENTIONED;
            
            case 1:
                return enums.DOC_RELATION.TARGETS;
            
            case 2:
                return enums.DOC_RELATION.CAUSED;
            }
        return;
    }

    public static fromString(input: String): enums.DOC_RELATION {
        switch (input) {
            case 'MENTIONED':
                return enums.DOC_RELATION.MENTIONED;
            
            case 'TARGETS':
                return enums.DOC_RELATION.TARGETS;
            
            case 'CAUSED':
                return enums.DOC_RELATION.CAUSED;
            }
        return enums.DOC_RELATION.MENTIONED;
    }
}
