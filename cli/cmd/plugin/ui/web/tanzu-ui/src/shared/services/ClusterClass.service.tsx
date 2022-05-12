import { ClusterClass, ClusterClassVariable, ManagementService } from '../../swagger-api';
import { CCDefinition, CCUiClassification, CCVariable, ClusterClassVariableType } from '../models/ClusterClass';

// Retrieves the cluster class associated with the given MC+CC name, then calls the callback with the resulting data
export function retrieveClusterClass(clusterName: string, clusterClassName: string, callback: (ccDef: CCDefinition) => void) {
    ManagementService.getClusterClass(clusterName, clusterClassName).then(cc => {
        callback(createCCDefinition(cc))
    })
}

function inCategoryRequired(ccVar: CCVariable): boolean {
    return ccVar.required || false
}

function inCategoryBasic(ccVar: CCVariable): boolean {
    return !ccVar.required && ccVar.category === CCUiClassification.BASIC
}

function inCategoryIntermediate(ccVar: CCVariable): boolean {
    return !ccVar.required && ccVar.category === CCUiClassification.INTERMEDIATE
}

function inCategoryAdvanced(ccVar: CCVariable): boolean {
    return !ccVar.required && (ccVar.category === CCUiClassification.ADVANCED ||
        ccVar.category === CCUiClassification.NONE)
}

// Takes a cluster class definition from the backend and creates a frontend version of it
function createCCDefinition(cc: ClusterClass): CCDefinition {
    const variables = createCCVariables(cc.variables)
    return {
        name: cc.name,
        variables,
        requiredVariables:      () => { return variables?.filter(ccVar => inCategoryRequired(ccVar)) || [] },
        basicVariables:         () => { return variables?.filter(ccVar => inCategoryBasic(ccVar)) || [] },
        intermediateVariables:  () => { return variables?.filter(ccVar => inCategoryIntermediate(ccVar) ) || [] },
        advancedVariables:      () => { return variables?.filter(ccVar => inCategoryAdvanced(ccVar)) || [] },
    } as CCDefinition
}

function createCCVariables(ccVars: ClusterClassVariable[] | undefined): CCVariable[] | undefined {
    return ccVars?.map<CCVariable>( clusterClassVariable => createCCVar(clusterClassVariable))
}

function createCCVar(ccVar: ClusterClassVariable): CCVariable {
    // TODO: get the "real" children from the backend data and populate accordingly
    const children = ccVar?.children ? createCCVariables(ccVar.children) : undefined
    return {
        default: ccVar?.default,     // TODO: use taxonomy to create default, cuz might be complex object ???
        description: ccVar?.description,
        name: ccVar.name || '',
        possibleValues: [],
        required: ccVar.required,
        taxonomy: getCcVarTaxonomyFromBackendValue(ccVar?.taxonomy || ''),
        category: getCcVarUiClassificationFromBackendValue(ccVar?.category || ''),
        children
    }
}

// TODO: ClusterClassVariableType[backendValue] should work here, but hasn't
function getCcVarTaxonomyFromBackendValue(backendValue: string): ClusterClassVariableType {
    switch(backendValue) {
    case '':
        return ClusterClassVariableType.UNKNOWN
    case 'boolean':
        return ClusterClassVariableType.BOOLEAN
    case 'booleanEnabled':
        return ClusterClassVariableType.BOOLEAN_ENABLED
    case 'cidr':
        return ClusterClassVariableType.CIDR
    case 'int':
        return ClusterClassVariableType.INTEGER
    case 'ip':
        return ClusterClassVariableType.IP
    case 'ipList':
        return ClusterClassVariableType.IP_LIST
    case 'number':
        return ClusterClassVariableType.NUMBER
    case 'string':
        return ClusterClassVariableType.STRING
    case 'stringK8sCompliant':
        return ClusterClassVariableType.STRING_K8S_COMPLIANT
    case 'stringParagraph':
        return ClusterClassVariableType.STRING_PARAGRAPH
    case 'proxy':
        return ClusterClassVariableType.PROXY
    }
    return ClusterClassVariableType.UNKNOWN
}

function getCcVarUiClassificationFromBackendValue(backendValue: string): CCUiClassification {
    // return CCUiClassification[backendValue as keyof typeof CCUiClassification] || CCUiClassification.NONE
    switch (backendValue) {
        case CCUiClassification.BASIC:
            return CCUiClassification.BASIC
        case CCUiClassification.INTERMEDIATE:
            return CCUiClassification.INTERMEDIATE
        case CCUiClassification.ADVANCED:
            return CCUiClassification.ADVANCED
    }
    return CCUiClassification.NONE
}
