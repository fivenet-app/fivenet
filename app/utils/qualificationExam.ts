export function examTextLength(value: string): number {
    return [...value].length;
}

export function areExamChoicesUnique(choices: readonly string[]): boolean {
    return new Set(choices).size === choices.length;
}

export function areExamChoicesAllowed(choices: readonly string[], allowed: readonly string[]): boolean {
    return choices.every((choice) => allowed.includes(choice));
}

// The exam form initializes single-choice responses so its controls always
// have a v-model target. A value is only an answer once it is an actual
// configured choice; UI libraries can also transiently write non-strings.
export function isAnsweredExamSingleChoice(choice: unknown, allowed: readonly string[]): choice is string {
    return typeof choice === 'string' && allowed.includes(choice);
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
