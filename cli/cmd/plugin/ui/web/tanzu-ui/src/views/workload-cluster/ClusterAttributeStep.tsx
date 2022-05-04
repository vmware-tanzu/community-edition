import React, { useContext } from 'react';

import { StepProps } from '../../shared/components/wizard/Wizard';
import { ClusterClassDefinition, ClusterClassVariable, ClusterClassVariableType } from '../../shared/models/ClusterClass';
import { WcStore } from '../../state-management/stores/Store.wc';
import { SubmitHandler, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup/dist/yup';
import { ClusterClassMultipleVariablesDisplay } from './ClusterClassVariableDisplay';
import { CdsButton } from '@cds/react/button';
import * as yup from 'yup';
import { TOGGLE_NAV, TOGGLE_WC_CC_OPTIONAL, TOGGLE_WC_CC_REQUIRED } from '../../state-management/actions/Ui.actions';

interface ClusterAttributeStepProps extends StepProps {
    retrieveClusterClassDefinition: (mc: string) => ClusterClassDefinition | undefined
}

function ClusterAttributeStep(props: Partial<ClusterAttributeStepProps>) {
    const {handleValueChange, currentStep, goToStep, submitForm, retrieveClusterClassDefinition} = props;
    const {state, dispatch} = useContext(WcStore);

    const mcName = state.data.SELECTED_MANAGEMENT_CLUSTER?.name
    const cc = mcName && retrieveClusterClassDefinition ? retrieveClusterClassDefinition(mcName) : undefined

    const formSchema = createFormSchema(cc)
    const methods = useForm({
        resolver: formSchema ? yupResolver(formSchema) : undefined,
    });
    const {
        register,
        handleSubmit,
        setValue,
        formState: {errors},
    } = methods;

    if (!retrieveClusterClassDefinition) {
        return <div>Programmer error: ClusterAttributeStep did not receive retrieveClusterClassDefinition!</div>
    }
    if (!mcName) {
        return <div>No management cluster has been selected (how did you get to this step?!)</div>
    }
    if (!cc) {
        return <div>We were unable to retrieve a ClusterClass object for management cluster {mcName}</div>
    }

    const onSubmit: SubmitHandler<any> = (data) => {
        if (Object.keys(errors).length === 0) {
            if (goToStep && currentStep && submitForm) {
                // TODO: figure out where to go!
                // goToStep(currentStep + 1);
                submitForm(currentStep)
                alert('form submitted - where do I go?!')
            }
        }
    };
    const toggleRequired = () => { dispatch({type: TOGGLE_WC_CC_REQUIRED});}
    const toggleOptional = () => { dispatch({type: TOGGLE_WC_CC_OPTIONAL});}

    const requiredVars = cc.requiredVariables ? cc.requiredVariables : []
    const optionalVars = cc.optionalVariables ? cc.optionalVariables : []
    return <div>
        {ClusterAttributeStepInstructions(cc)}
        <br/>
        {ClusterClassMultipleVariablesDisplay(requiredVars, 'Required Variables',
            { register, errors, expanded: state.ui.wcCcRequiredExpanded, toggleExpanded: toggleRequired }) }
        <br/>
        {ClusterClassMultipleVariablesDisplay(optionalVars, 'Optional Variables',
            { register, errors, expanded: state.ui.wcCcOptionalExpanded, toggleExpanded: toggleOptional }) }
        <br/>
        <br/>
        <CdsButton onClick={handleSubmit(onSubmit)}>CREATE WORKLOAD CLUSTER</CdsButton>
    </div>
}

function createFormSchema(cc: ClusterClassDefinition | undefined) {
    if (!cc) {
        return undefined
    }

    let schemaObject = {}
    schemaObject = addRequiredFieldsToSchemaObject(schemaObject, cc)
    schemaObject = addOptionalFieldsToSchemaObject(schemaObject, cc)

    return yup.object(schemaObject);
}

function addRequiredFieldsToSchemaObject(obj: any, cc: ClusterClassDefinition) {
    const result = {...obj}
    cc.requiredVariables?.forEach(ccVar => addSingleFieldToSchemaObject(result, ccVar, true))
    return result
}

function addOptionalFieldsToSchemaObject(obj: any, cc: ClusterClassDefinition) {
    const result = {...obj}
    cc.optionalVariables?.forEach(ccVar => addSingleFieldToSchemaObject(result, ccVar, false))
    return result
}

function addSingleFieldToSchemaObject(obj: any, ccVar: ClusterClassVariable, required: boolean) {
    if (ccVar.valueType === ClusterClassVariableType.STRING) {
        obj[ccVar.name] = yup.string().nullable()
    } else if (ccVar.valueType === ClusterClassVariableType.BOOLEAN) {
        obj[ccVar.name] = yup.boolean().nullable()
    } else {
        obj[ccVar.name] = yup.number().nullable()
    }
    if (required) {
        const prompt = promptFromClusterClassType(ccVar)
        obj[ccVar.name] = obj[ccVar.name].required(prompt)
    }
}

function promptFromClusterClassType(ccVar: ClusterClassVariable): string {
    switch (ccVar.valueType) {
        case ClusterClassVariableType.STRING:
            if (ccVar.possibleValues && ccVar.possibleValues.length > 0) {
                return 'Please select a value'
            }
            return 'Please enter a value'
        case ClusterClassVariableType.INTEGER:
        case ClusterClassVariableType.NUMBER:
            return 'Please enter a value'
    }
    return 'Value required'
}

function ClusterAttributeStepInstructions(cc: ClusterClassDefinition | undefined) {
    if (!cc) {
        return <div>There is no cluster class definition, so you cannot do this step! So sorry.</div>
    }
    const nRequiredVars = cc.requiredVariables?.length
    const nOptionalVars = cc.optionalVariables?.length
    return <div>So you have a cluster class with {nRequiredVars ? nRequiredVars : 'no'} required variables
        and {nOptionalVars ? nOptionalVars : 'no'} optional variables. Deal with it.</div>
}

export default ClusterAttributeStep;
