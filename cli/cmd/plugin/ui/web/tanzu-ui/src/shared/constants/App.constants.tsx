export const AppEnvironment = {
    DEV: 'development',
    PROD: 'production',
};

export const WebsocketAddress = {
    DEV_LOCATION: 'localhost:8008',
    DEFAULT_PROTOCOL: 'ws',
    SECURE_PROTOCOL: 'wss'
};

export enum STATUS {
    VALID = 'valid',
    DISABLED = 'disabled',
    CURRENT = 'current',
    INVALID = 'invalid',
    TOUCHED = 'touched'
}
