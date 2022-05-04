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

export enum ClusterClassVariableType {
    INTEGER,
    NUMBER,
    STRING,
    BOOLEAN,
}
