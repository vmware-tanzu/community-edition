import { CCVariable, ClusterClassVariableType } from '../../../shared/models/ClusterClass';

export function createProxyComponentCCVar(defaults: any): CCVariable {
    return populateProxyDefaults(defaults, ProxyComponentCCVar());
}

// NOTE: this fxn has the intentional side effect of populating default values for the children of the given ccVar
//       To change this to a "pure" fxn, we would need to do a deep clone
function populateProxyDefaults(defaults: any, ccVar: CCVariable): CCVariable {
    if (defaults) {
        Object.keys(defaults).forEach((key) => {
            const child = ccVar.children?.find((child) => child.name === key);
            if (child) {
                child.default = defaults[key];
            } else {
                console.warn(`A proxy type was found with a default key of ${key}, but no such child exists in our PROXY children`);
            }
        });
    }
    return ccVar;
}

function ProxyComponentCCVar(): CCVariable {
    return {
        name: 'proxy',
        label: 'Proxy',
        info: 'Use this panel to configure the proxy information of your proxy server and what IPs should be routed to it.',
        // taxonomy: ClusterClassVariableType.GROUP_OPTIONAL,
        taxonomy: ClusterClassVariableType.GROUP,
        required: false,
        children: ProxyComponentChildren(),
    };
}

function ProxyComponentChildren(): CCVariable[] {
    return [
        {
            name: 'httpProxy',
            prompt: 'Use proxy for HTTP calls to these IPs (list):',
            taxonomy: ClusterClassVariableType.IP_LIST,
        },
        {
            name: 'httpsProxy',
            prompt: 'Use proxy for HTTPS calls to these IPs (list):',
            taxonomy: ClusterClassVariableType.IP_LIST,
        },
        {
            name: 'noProxy',
            prompt: 'Use NO proxy for calls to these IPs (list):',
            taxonomy: ClusterClassVariableType.IP_LIST,
        },
        {
            name: 'proxyCA',
            prompt: 'Certificate for proxy',
            taxonomy: ClusterClassVariableType.STRING_PARAGRAPH,
        },
    ];
}
