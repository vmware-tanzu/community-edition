export interface ClusterClassDefinition {
    name: string,
    requiredVariables?: ClusterClassVariable[],
    optionalVariables?: ClusterClassVariable[],
}

export interface ClusterClassVariable {
    name: string,
    valueType: ClusterClassVariableType,
    description?: string,
    defaultValue?: string,
    possibleValues?: string[],
}

// for some reason, eslint is reporting these enum values as unused
/* eslint-disable no-unused-vars */
export enum ClusterClassVariableType {
    INTEGER,
    NUMBER,
    STRING,
    BOOLEAN,
}
