// React imports
import React, { ChangeEvent, useContext } from 'react';
import { useNavigate } from 'react-router-dom';
import { SubmitHandler, useForm } from 'react-hook-form';

// Library imports
import { CdsButton } from '@cds/react/button';
import * as yup from 'yup';
import { yupResolver } from '@hookform/resolvers/yup/dist/yup';

// App imports
import { StepProps } from '../../shared/components/wizard/Wizard';
import { ClusterClassDefinition, ClusterClassVariable, ClusterClassVariableType } from '../../shared/models/ClusterClass';
import { WcStore } from '../../state-management/stores/Store.wc';
import { ClusterClassMultipleVariablesDisplay } from './ClusterClassVariableDisplay';
import { TOGGLE_WC_CC_OPTIONAL, TOGGLE_WC_CC_REQUIRED } from '../../state-management/actions/Ui.actions';
import { NavRoutes } from '../../shared/constants/NavRoutes.constants';

interface ClusterAttributeStepProps extends StepProps {
    retrieveClusterClassDefinition: (mc: string) => ClusterClassDefinition | undefined
}

function ClusterAttributeStep(props: Partial<ClusterAttributeStepProps>) {
    const { handleValueChange, currentStep, goToStep, submitForm, retrieveClusterClassDefinition } = props;
    const { state, dispatch } = useContext(WcStore);

    const mcName = state.data.SELECTED_MANAGEMENT_CLUSTER?.name
    const cc = mcName && retrieveClusterClassDefinition ? retrieveClusterClassDefinition(mcName) : undefined

    const formSchema = createFormSchema(cc)
    const methods = useForm({
        resolver: formSchema ? yupResolver(formSchema) : undefined,
    });
    const { register, handleSubmit, formState: { errors } } = methods;

    const navigate = useNavigate();

    // TODO: we will likely need to navigate to a WC-specific progress route, but for now, just to be able to demo...
    const navigateToProgress = (): void => {
        navigate('/' + NavRoutes.DEPLOY_PROGRESS);
    };

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
                // TODO: we'll need to call the backend service to actually do the deploying
                // goToStep(currentStep + 1);
                submitForm(currentStep)
                navigateToProgress()
            }
        }
    };
    const toggleRequired = () => { dispatch({ type: TOGGLE_WC_CC_REQUIRED }) }
    const toggleOptional = () => { dispatch({ type: TOGGLE_WC_CC_OPTIONAL }) }

    const onValueChange = (evt: ChangeEvent<HTMLSelectElement>) => {
        if (handleValueChange) {
            // NOTE: Rather than store the CC variables at the top level of our data store,
            //       we accumulate them in a single object. This makes using them much easier (later).
            // NOTE: we assume that state.data.CLUSTER_CLASS_VARIABLE_VALUES is never null or undefined
            const accumulatedAttributes = state.data.CLUSTER_CLASS_VARIABLE_VALUES
            accumulatedAttributes[evt.target.name] = evt.target.value
            handleValueChange('CLUSTER_CLASS_VARIABLE_VALUES', accumulatedAttributes, currentStep, errors)
        } else {
            console.error('ClusterAttributeStep unable to find a handleValueChange handler!')
        }
    }

    const requiredVars = cc.requiredVariables ? cc.requiredVariables : []
    const optionalVars = cc.optionalVariables ? cc.optionalVariables : []
    return <div>
        {ClusterAttributeStepInstructions(cc)}
        <br/>
        {ClusterClassMultipleVariablesDisplay(requiredVars, 'Required Variables',
            { register, errors, expanded: state.ui.wcCcRequiredExpanded, toggleExpanded: toggleRequired, onValueChange }) }
        <br/>
        {ClusterClassMultipleVariablesDisplay(optionalVars, 'Optional Variables',
            { register, errors, expanded: state.ui.wcCcOptionalExpanded, toggleExpanded: toggleOptional, onValueChange }) }
        <br/>
        <br/>
        <CdsButton
            className="cluster-action-btn"
            status="primary"
            onClick={handleSubmit(onSubmit)}>
            Create Workload Cluster
        </CdsButton>
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
    const result = { ...obj }
    // TODO: SHIMON: feels like reduce() would handle this better
    cc.requiredVariables?.forEach(ccVar => addSingleFieldToSchemaObject(result, ccVar, true))
    return result
}

function addOptionalFieldsToSchemaObject(obj: any, cc: ClusterClassDefinition) {
    const result = { ...obj }
    // TODO: SHIMON: feels like reduce() would handle this better
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