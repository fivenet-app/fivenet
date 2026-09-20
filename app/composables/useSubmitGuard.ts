type AsyncSubmitAction<Args extends readonly unknown[], Result> = (...args: Args) => Promise<Result>;

type UseSubmitGuardOptions = {
    cooldown?: number;
};

/**
 * Guards asynchronous form actions against duplicate submissions.
 *
 * Unlike a throttle, this does not queue a trailing invocation. A second
 * submission is ignored while the action is running and during the optional
 * cooldown afterwards.
 */
export function useSubmitGuard<Args extends readonly unknown[], Result>(
    action: AsyncSubmitAction<Args, Result>,
    options: UseSubmitGuardOptions = {},
) {
    const isSubmitting = ref(false);
    const cooldown = options.cooldown ?? 400;

    async function submit(...args: Args): Promise<Result | undefined> {
        if (isSubmitting.value) return undefined;

        isSubmitting.value = true;

        try {
            return await action(...args);
        } finally {
            if (cooldown > 0) {
                await new Promise<void>((resolve) => setTimeout(resolve, cooldown));
            }

            isSubmitting.value = false;
        }
    }

    return {
        submit,
        isSubmitting: readonly(isSubmitting),
        canSubmit: computed(() => !isSubmitting.value),
    };
}
