import { Logger, type ILogger } from '~/utils/logger';

const logger = new Logger();

export function useLogger(prefix?: string): ILogger {
    if (prefix === '') {
        return logger;
    }
    return new Logger(prefix);
}
