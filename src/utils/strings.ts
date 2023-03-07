export function toTitleCase(value: string) {
    return value.replace(/(?:^|\s|-)\S/g, x => x.toUpperCase());
}
