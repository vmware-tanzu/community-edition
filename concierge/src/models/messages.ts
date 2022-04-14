export enum Messages {
    REQUEST_INSTALL = 'app:install-tanzu',
    REQUEST_INSTALLED_PLUGIN_LIST = 'app:plugin-list-request',
    REQUEST_LAUNCH_KICKSTART = 'app:launch-kickstart',
    REQUEST_LAUNCH_UI = 'app:launch-ui',
    REQUEST_PRE_INSTALL = 'app:pre-install-tanzu',

    RESPONSE_INSTALLED_PLUGIN_LIST = 'app:plugin-list-response',
    RESPONSE_PROGRESS = 'app:install-progress',
    RESPONSE_PRE_INSTALL = 'app:pre-install-tanzu-response',
}
