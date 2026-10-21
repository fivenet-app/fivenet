import { describe, expect, it, vi } from 'vitest';
import type { JobProps } from '~~/gen/ts/resources/jobs/props/props';
import { removeJobLogo } from './jobPropsLogo';

const jobProps = (): JobProps => ({
    job: 'police',
    livemapMarkerColor: '#ffffff',
    logoFileId: 42,
    logoFile: {
        id: 42,
        filePath: 'jobprops/logo.png',
    },
});

describe('removeJobLogo', () => {
    it('clears the logo state and reports success after deletion', async () => {
        const currentJobProps = jobProps();
        const setAuthJobProps = vi.fn();
        const clearPendingLogo = vi.fn();
        const notifySuccess = vi.fn();

        await removeJobLogo({
            deleteLogo: vi.fn().mockResolvedValue(undefined),
            jobProps: currentJobProps,
            setAuthJobProps,
            clearPendingLogo,
            notifySuccess,
        });

        expect(currentJobProps.logoFile).toBeUndefined();
        expect(currentJobProps.logoFileId).toBeUndefined();
        expect(setAuthJobProps).toHaveBeenCalledWith(currentJobProps);
        expect(clearPendingLogo).toHaveBeenCalledOnce();
        expect(notifySuccess).toHaveBeenCalledOnce();
    });

    it('does not clear state or notify when deletion fails', async () => {
        const currentJobProps = jobProps();
        const setAuthJobProps = vi.fn();
        const clearPendingLogo = vi.fn();
        const notifySuccess = vi.fn();
        const error = new Error('delete failed');

        await expect(
            removeJobLogo({
                deleteLogo: vi.fn().mockRejectedValue(error),
                jobProps: currentJobProps,
                setAuthJobProps,
                clearPendingLogo,
                notifySuccess,
            }),
        ).rejects.toThrow(error);

        expect(currentJobProps.logoFileId).toBe(42);
        expect(setAuthJobProps).not.toHaveBeenCalled();
        expect(clearPendingLogo).not.toHaveBeenCalled();
        expect(notifySuccess).not.toHaveBeenCalled();
    });
});
