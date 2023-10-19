export function jsonStringify(obj: any, space?: string | number): string {
    return JSON.stringify(
        obj,
        (_, value) => {
            if (typeof value === 'bigint') {
                return {
                    type: 'bigint',
                    value: value.toString(),
                };
            } else {
                return value;
            }
        },
        space,
    );
}

export function jsonParse(text: string): any {
    JSON.parse(text, (_, value) => {
        if (value && value.type == 'bigint') {
            return BigInt(value.value);
        }
        return value;
    });
}
