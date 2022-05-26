import { CCVariable, ClusterClassVariableType } from '../../../shared/models/ClusterClass';
import { populateDefaults } from './CCUtil';

export function createImageRepositoryCCVar(defaults: any): CCVariable {
    return populateDefaults(defaults, ImageRepositoryComponentCCVar());
}

function ImageRepositoryComponentCCVar(): CCVariable {
    return {
        name: 'imageRepository',
        label: 'Image Repository',
        info: 'Use this panel to configure the proxy information of your proxy server and what IPs should be routed to it.',
        // taxonomy: ClusterClassVariableType.GROUP_OPTIONAL,
        taxonomy: ClusterClassVariableType.GROUP,
        required: false,
        children: ImageRepositoryComponentChildren(),
    };
}

function ImageRepositoryComponentChildren(): CCVariable[] {
    return [
        {
            name: 'host',
            prompt: 'Repository',
            taxonomy: ClusterClassVariableType.STRING,
        },
        {
            name: 'tlsCertificateValidation',
            prompt: 'Validate TLS certificate',
            taxonomy: ClusterClassVariableType.BOOLEAN,
            // default: true,
        },
    ];
}
