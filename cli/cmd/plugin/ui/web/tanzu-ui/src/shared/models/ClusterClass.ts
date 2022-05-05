export interface ClusterClassDefinition {
    name: string,
    requiredVariables?: ClusterClassVariable[],
    optionalVariables?: ClusterClassVariable[],
    advancedVariables?: ClusterClassVariable[],
}

export interface ClusterClassVariable {
    name: string,
    valueType: ClusterClassVariableType,
    description?: string,
    defaultValue?: string,
    possibleValues?: string[],
    required?: boolean,
}

// for some reason, eslint is reporting these enum values as unused
/* eslint-disable no-unused-vars */
export enum ClusterClassVariableType {
    BOOLEAN,
    CIDR,
    INTEGER,
    IP,
    IP_LIST,
    NUMBER,
    STRING,
    STRING_PARAGRAPH,
}
