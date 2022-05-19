import { ClusterClass, ClusterClassVariable, ClusterClassVariableCategory, ManagementService } from '../../swagger-api';
import { CCCategory, CCDefinition, CCVariable, ClusterClassVariableType } from '../models/ClusterClass';
import { createProxyComponentCCVar } from '../../views/workload-cluster/cluster-class-display-widgets/CCProxyComponent';

// NOTE: style guide: backend object variables use a FULLY spelled out name, front end objects abbreviate to CC.
//       For example "clusterClass" is an instance of ClusterClass (backend), whereas "cc" is an instance of CCDefinition (front end)
//       This style helps us understand exactly what kind of object we're dealing with

// Retrieves the cluster class associated with the given MC+CC name, then calls the callback with the resulting data
export function retrieveClusterClass(clusterName: string, clusterClassName: string, callback: (ccDef: CCDefinition) => void) {
    ManagementService.getClusterClass(clusterName, clusterClassName).then((cc) => {
        callback(createCCDefinition(cc));
    });
}

// Takes a cluster class definition from the backend and creates a frontend version of it
function createCCDefinition(clusterClass: ClusterClass): CCDefinition {
    return {
        name: clusterClass.name,
        categories: createCCCategories(clusterClass),
    } as CCDefinition;
}

function createCCCategories(cc: ClusterClass): CCCategory[] {
    return (
        cc.categories?.map<CCCategory>((clusterClassVariableCategory: ClusterClassVariableCategory) => {
            return {
                displayOpen: clusterClassVariableCategory.displayOpen,
                label: clusterClassVariableCategory.label,
                name: clusterClassVariableCategory.name,
                variables: createCCVariables(clusterClassVariableCategory.variables),
            } as CCCategory;
        }) || []
    );
}

function createCCVariables(clusterClassVariables: ClusterClassVariable[] | undefined): CCVariable[] | undefined {
    return clusterClassVariables?.map<CCVariable>((clusterClassVariable) => createCCVar(clusterClassVariable));
}

function createCCVar(clusterClassVariable: ClusterClassVariable): CCVariable {
    switch (clusterClassVariable.taxonomy) {
        case ClusterClassVariableType.PROXY:
            return createProxyComponentCCVar(clusterClassVariable.default);
    }

    return {
        default: clusterClassVariable?.default, // TODO: use taxonomy to create default, cuz might be complex object ???
        prompt: clusterClassVariable?.prompt,
        label: clusterClassVariable?.label,
        info: clusterClassVariable?.info,
        name: clusterClassVariable.name || '',
        possibleValues: clusterClassVariable.possibleValues ? clusterClassVariable.possibleValues : [],
        required: clusterClassVariable.required,
        taxonomy: getCcVarTaxonomyFromBackendValue(clusterClassVariable?.taxonomy || ''),
        children: createCCVarChildren(clusterClassVariable),
    };
}

// createCCVarChildren() takes a ClusterClassVariable object (backend) and creates an array of CCVariable objects (front end) which
// are the children for the given ClusterClassVariable.
function createCCVarChildren(clusterClassVariable: ClusterClassVariable): CCVariable[] | undefined {
    return clusterClassVariable?.children ? createCCVariables(clusterClassVariable.children) : undefined;
}

// TODO: ClusterClassVariableType[backendValue] should work here, but hasn't
function getCcVarTaxonomyFromBackendValue(backendValue: string): ClusterClassVariableType {
    switch (backendValue) {
        case '':
            return ClusterClassVariableType.UNKNOWN;
        case 'boolean':
            return ClusterClassVariableType.BOOLEAN;
        case 'cidr':
            return ClusterClassVariableType.CIDR;
        case 'int':
            return ClusterClassVariableType.INTEGER;
        case 'ip':
            return ClusterClassVariableType.IP;
        case 'ipList':
            return ClusterClassVariableType.IP_LIST;
        case 'number':
            return ClusterClassVariableType.NUMBER;
        case 'string':
            return ClusterClassVariableType.STRING;
        case 'stringK8sCompliant':
            return ClusterClassVariableType.STRING_K8S_COMPLIANT;
        case 'stringParagraph':
            return ClusterClassVariableType.STRING_PARAGRAPH;
        case 'proxy':
            return ClusterClassVariableType.PROXY;
    }
    return ClusterClassVariableType.UNKNOWN;
}
