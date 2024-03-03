export const dateRequiredValidator = (value: unknown): boolean => {
    return value !== null && value !== undefined && typeof value === 'string';
};
