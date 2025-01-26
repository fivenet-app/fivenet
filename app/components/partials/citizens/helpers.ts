export function sexToColor(sex: string): string {
    if (sex === 'f') {
        return 'pink';
    } else if (sex === 'm') {
        return 'blue';
    } else if (sex === 'd') {
        return 'gray';
    }

    return 'white';
}

export function sexToTextColor(sex: string): string {
    if (sex === 'f') {
        return 'text-pink-300';
    } else if (sex === 'm') {
        return 'text-blue-300';
    } else if (sex === 'd') {
        return 'text-gray-300';
    }

    return 'text-white-300';
}

export function sexToIcon(sex: string): string {
    if (sex === 'f') {
        return 'i-mdi-gender-female';
    } else if (sex === 'm') {
        return 'i-mdi-gender-male';
    }

    return 'i-mdi-gender-non-binary';
}
