import type { JobProps } from '~~/gen/ts/resources/jobs/props/props';

type RemoveJobLogoOptions = {
    deleteLogo: () => Promise<unknown>;
    jobProps: JobProps | null;
    setAuthJobProps: (jobProps: JobProps) => void;
    clearPendingLogo: () => void;
    notifySuccess: () => void;
};

export async function removeJobLogo(options: RemoveJobLogoOptions): Promise<void> {
    await options.deleteLogo();

    if (options.jobProps) {
        options.jobProps.logoFile = undefined;
        options.jobProps.logoFileId = undefined;
        options.setAuthJobProps(options.jobProps);
    }

    options.clearPendingLogo();
    options.notifySuccess();
}
