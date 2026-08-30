export function isTemplateBlockActionValid(value: string): boolean {
    const action = value.trim();

    if (action === 'else' || action === 'break' || action === 'continue') return true;
    if (/^else\s+if\s+\S/.test(action)) return true;
    if (/^(range|if|with)\s+\S/.test(action)) return true;

    return false;
}
