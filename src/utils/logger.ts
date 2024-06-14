export interface ILogger {
    log(message?: any, ...optionalParams: any[]): void;
    debug(message?: any, ...optionalParams: any[]): void;
    info(message?: any, ...optionalParams: any[]): void;
    warn(message?: any, ...optionalParams: any[]): void;
    error(message?: any, ...optionalParams: any[]): void;
}

export class Logger {
    readonly prefix?: string;

    constructor(prefix?: string) {
        this.prefix = prefix?.trim();
    }

    getPrefix(): string {
        return this.prefix ? this.prefix + ': ' : '';
    }

    log(message?: any, ...optionalParams: any[]): void {
        console.log(this.getPrefix() + message, ...optionalParams);
    }

    debug(message?: any, ...optionalParams: any[]): void {
        console.debug(this.getPrefix() + message, ...optionalParams);
    }

    info(message?: any, ...optionalParams: any[]): void {
        console.info(this.getPrefix() + message, ...optionalParams);
    }

    warn(message?: any, ...optionalParams: any[]): void {
        console.warn(this.getPrefix() + message, ...optionalParams);
    }

    error(message?: any, ...optionalParams: any[]): void {
        console.error(this.getPrefix() + message, ...optionalParams);
    }
}
