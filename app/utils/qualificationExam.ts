export function examTextLength(value: string): number {
    return [...value].length;
}

export function areExamChoicesUnique(choices: readonly string[]): boolean {
    return new Set(choices).size === choices.length;
}

export function areExamChoicesAllowed(choices: readonly string[], allowed: readonly string[]): boolean {
    return choices.every((choice) => allowed.includes(choice));
}

export function isExamChoiceLimitExceeded(choices: readonly string[], limit?: number): boolean {
    return limit !== undefined && limit > 0 && choices.length > limit;
}

// Image questions intentionally use separator responses because they do not
// accept an answer.
export function isValidExamResponseKind(questionKind: string | undefined, responseKind: string | undefined): boolean {
    if (questionKind === undefined || responseKind === undefined) return false;

    return responseKind === (questionKind === 'image' ? 'separator' : questionKind);
}
