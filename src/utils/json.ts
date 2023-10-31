import { type StateTree } from 'pinia';

class JSONSerializer {
    serialize(value: StateTree): string {
        return jsonStringify(value);
    }

    deserialize(value: string): StateTree {
        return jsonParse(value) as StateTree;
    }
}

export const jsonSerializer = new JSONSerializer();

function jsonMarshal(_: string, value: any): any {
    if (typeof value === 'bigint') {
        return {
            type: 'bigint',
            value: value.toString(),
        };
    } else {
        return value;
    }
}

export function jsonStringify(obj: any, space?: string | number): string {
    return JSON.stringify(obj, jsonMarshal, space);
}

function jsonUnmarshal(_: string, value: any) {
    if (value && value.type === 'bigint') {
        return BigInt(value.value);
    }
    return value;
}

export function jsonParse(text: string): any {
    return JSON.parse(text, jsonUnmarshal);
}
