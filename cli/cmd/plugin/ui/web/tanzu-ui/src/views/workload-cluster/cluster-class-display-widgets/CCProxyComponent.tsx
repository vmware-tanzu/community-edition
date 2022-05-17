import { CCVariable, ClusterClassVariableType } from '../../../shared/models/ClusterClass';
import { CCParentVariableDisplay, ClusterClassVariableDisplayOptions } from '../ClusterClassVariableDisplay';

export function ProxyComponent(options: ClusterClassVariableDisplayOptions) {
    return CCParentVariableDisplay(ProxyComponentVars(), options)
}

export function ProxyComponentVars(): CCVariable {
    return {
        name: 'proxy',
        description: 'Proxy',
        taxonomy: ClusterClassVariableType.GROUP_OPTIONAL,
        required: false,
        children: [
            {
                name: 'httpProxy',
                description: 'Use HTTP proxy for (list):',
                taxonomy: ClusterClassVariableType.IP_LIST,
            },
            {
                name: 'httpsProxy',
                description: 'Use HTTPS proxy for (list):',
                taxonomy: ClusterClassVariableType.IP_LIST,
            },
            {
                name: 'noProxy',
                description: 'Use NO proxy for (list):',
                taxonomy: ClusterClassVariableType.IP_LIST,
            },
            {
                name: 'proxyCA',
                description: 'Certificate for proxy',
                taxonomy: ClusterClassVariableType.STRING_PARAGRAPH,
            },
        ]
    }
}
