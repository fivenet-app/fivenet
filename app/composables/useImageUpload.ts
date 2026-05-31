import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

type I18nKeyMessage = {
    key: string;
    parameters: I18NItemParamters;
};

type I18nNotification = {
    title: I18nKeyMessage;
    description: I18nKeyMessage;
};

type UploadImagesOptions<TResponse> = {
    files: File[];
    uploadOne: (file: File) => Promise<TResponse>;
    allowedMimeTypes?: string[];
    currentFileCount?: number;
    fileLimit?: number;
    onUploaded?: (response: TResponse, file: File) => void;
    onUploadError?: (error: unknown, file: File) => void;
    invalidTypeNotification?: I18nNotification;
    fileLimitNotification?: I18nNotification;
    uploadFailedNotification?: (error: unknown) => I18nNotification;
    successNotification?: I18nNotification;
};

type UploadImagesResult<TResponse> =
    | { ok: true; uploaded: TResponse[] }
    | { ok: false; reason: 'file_limit' | 'invalid_file_type' };

function matchesAcceptedType(file: File, acceptedType: string): boolean {
    if (!acceptedType) return false;

    if (acceptedType.startsWith('.')) {
        return file.name.toLowerCase().endsWith(acceptedType.toLowerCase());
    }

    if (acceptedType.endsWith('/*')) {
        const prefix = acceptedType.slice(0, -1);
        return file.type.startsWith(prefix);
    }

    return file.type === acceptedType;
}

export function useImageUpload() {
    const notifications = useNotificationsStore();
    const { fileUpload } = useAppConfig();

    const notify = (payload: I18nNotification, type: NotificationType = NotificationType.ERROR): void => {
        notifications.add({
            title: payload.title,
            description: payload.description,
            type,
        });
    };

    const uploadImages = async <TResponse>(options: UploadImagesOptions<TResponse>): Promise<UploadImagesResult<TResponse>> => {
        const allowedMimeTypes = options.allowedMimeTypes ?? fileUpload.types.images;
        const currentFileCount = options.currentFileCount ?? 0;

        if (options.fileLimit !== undefined && currentFileCount + options.files.length > options.fileLimit) {
            if (options.fileLimitNotification) notify(options.fileLimitNotification);
            return { ok: false, reason: 'file_limit' };
        }

        for (const file of options.files) {
            const isAllowed = allowedMimeTypes.some((acceptedType) => matchesAcceptedType(file, acceptedType));
            if (!isAllowed) {
                if (options.invalidTypeNotification) notify(options.invalidTypeNotification);
                return { ok: false, reason: 'invalid_file_type' };
            }
        }

        const uploaded: TResponse[] = [];

        for (const file of options.files) {
            try {
                const response = await options.uploadOne(file);
                uploaded.push(response);

                options.onUploaded?.(response, file);
                if (options.successNotification) notify(options.successNotification, NotificationType.SUCCESS);
            } catch (error) {
                options.onUploadError?.(error, file);
                if (options.uploadFailedNotification) {
                    notify(options.uploadFailedNotification(error), NotificationType.ERROR);
                }
                throw error;
            }
        }

        return { ok: true, uploaded };
    };

    return { uploadImages };
}
