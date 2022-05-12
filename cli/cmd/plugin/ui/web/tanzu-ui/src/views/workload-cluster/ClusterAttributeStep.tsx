// React imports
import React, { ChangeEvent, useContext, useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { SubmitHandler, useForm } from 'react-hook-form';

// Library imports
import { CdsButton } from '@cds/react/button';
import { yupResolver } from '@hookform/resolvers/yup/dist/yup';

// App imports
import { CCVAR_CHANGE } from '../../state-management/actions/Form.actions';
import { CCDefinition, ClusterClassDefinition } from '../../shared/models/ClusterClass';
import { CCMultipleVariablesDisplay, createFormSchemaCC } from './ClusterClassVariableDisplay';
import { getSelectedManagementCluster, getValueFromChangeEvent, } from './WorkloadClusterUtility';
import ManagementClusterInfoBanner from './ManagementClusterInfoBanner';
import { NavRoutes } from '../../shared/constants/NavRoutes.constants';
import { StepProps } from '../../shared/components/wizard/Wizard';
import {
    TOGGLE_WC_CC_ADVANCED,
    TOGGLE_WC_CC_BASIC,
    TOGGLE_WC_CC_INTERMEDIATE,
    TOGGLE_WC_CC_REQUIRED
} from '../../state-management/actions/Ui.actions';
import { WcStore } from '../../state-management/stores/Store.wc';
import { retrieveClusterClass } from '../../shared/services/ClusterClass.service';
import { AnyObjectSchema } from 'yup';

interface ClusterAttributeStepProps extends StepProps {
    retrieveClusterClassDefinition: (mc: string | undefined) => ClusterClassDefinition | undefined
}

function ClusterAttributeStep(props: Partial<ClusterAttributeStepProps>) {
    const { handleValueChange, currentStep, goToStep, submitForm, retrieveClusterClassDefinition } = props;
    const { state, dispatch } = useContext(WcStore);
    const [ ccDefinition, setCcDefinition ] = useState<CCDefinition>()
    const [ formSchema, setFormSchema ] = useState<AnyObjectSchema>()

    const cluster = getSelectedManagementCluster(state)

    const methods = useForm({
        resolver: formSchema ? yupResolver(formSchema) : undefined,
    });
    const { register, handleSubmit, formState: { errors }, setValue } = methods;

    const navigate = useNavigate();

    useEffect(() => {
        if (cluster.name) {
            // TODO: actually get the cluster class list (instead of hard-coded GET), and allow user to select if multiple
            retrieveClusterClass(cluster.name, `tkg-${cluster.provider}-default`, (ccDef) => {
                setCcDefinition(ccDef)
                setFormSchema(createFormSchemaCC(ccDef))
                console.log(`TODO: let's do something with the ccDef: ${JSON.stringify(ccDef)}`)
            })
        }
    }, [cluster])

    // TODO: we will likely need to navigate to a WC-specific progress route, but for now, just to be able to demo...
    const navigateToProgress = (): void => {
        navigate('/' + NavRoutes.DEPLOY_PROGRESS);
    };

    if (!retrieveClusterClassDefinition) {
        return <div>Programmer error: ClusterAttributeStep did not receive retrieveClusterClassDefinition!</div>
    }
    if (!cluster) {
        return <div>No management cluster has been selected (how did you get to this step?!)</div>
    }
    if (!ccDefinition) {
        return <div>We were unable to retrieve a ClusterClass object for management cluster {cluster.name}</div>
    }

    const onSubmit: SubmitHandler<any> = (data) => {
        const nErrors = Object.keys(errors).length
        if (nErrors === 0) {
            if (goToStep && currentStep && submitForm) {
                // TODO: we'll need to call the backend service to actually do the deploying
                // goToStep(currentStep + 1);
                submitForm(currentStep)
                navigateToProgress()
            }
        } else {
            console.log(`ClusterAttributeStep has an invalid form submission (${nErrors} errors)`)
        }
    };
    const toggleRequired = () => { dispatch({ type: TOGGLE_WC_CC_REQUIRED }) }
    const toggleBasic = () => { dispatch({ type: TOGGLE_WC_CC_BASIC }) }
    const toggleIntermediate = () => { dispatch({ type: TOGGLE_WC_CC_INTERMEDIATE }) }
    const toggleAdvanced = () => { dispatch({ type: TOGGLE_WC_CC_ADVANCED }) }

    const onValueChange = (evt: ChangeEvent<HTMLSelectElement>) => {
        const value = getValueFromChangeEvent(evt)
        const varName = evt.target.name
        setValue(varName, value, { shouldValidate: true })
        if (handleValueChange) {
            handleValueChange(CCVAR_CHANGE, varName, value, currentStep, errors,
                { clusterName: cluster.name })
        } else {
            console.error('ClusterAttributeStep unable to find a handleValueChange handler!')
        }
    }

    const requiredVars = ccDefinition ? ccDefinition.requiredVariables() : []
    const basicVars = ccDefinition ? ccDefinition.basicVariables() : []
    const intermediateVars = ccDefinition ? ccDefinition.intermediateVariables() : []
    const advancedVars = ccDefinition ? ccDefinition.advancedVariables() : []
    return <div>
        { ManagementClusterInfoBanner(cluster) }
        <br/>
        { CCStepInstructions(ccDefinition) }
        <br/>
        { CCMultipleVariablesDisplay(requiredVars, 'Required Variables',
            { register, errors, expanded: state.ui.wcCcRequiredExpanded, toggleExpanded: toggleRequired, onValueChange }) }
        <br/>
        { CCMultipleVariablesDisplay(basicVars, `Basic Variables (${basicVars.length})`,
            { register, errors, expanded: state.ui.wcCcBasicExpanded, toggleExpanded: toggleBasic, onValueChange }) }
        <br/>
        { CCMultipleVariablesDisplay(intermediateVars, `Intermediate Variables (${intermediateVars.length})`,
            { register, errors, expanded: state.ui.wcCcIntermediateExpanded, toggleExpanded: toggleIntermediate, onValueChange }) }
        <br/>
        { CCMultipleVariablesDisplay(advancedVars, `Advanced Variables (${advancedVars.length})`,
            { register, errors, expanded: state.ui.wcCcAdvancedExpanded, toggleExpanded: toggleAdvanced, onValueChange }) }
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

function CCStepInstructions(cc: CCDefinition | undefined) {
    if (!cc) {
        return <div>There is no cluster class definition, so you cannot do this step! So sorry.</div>
    }
    const nRequiredVars = cc.requiredVariables().length
    const nBasicVars = cc.basicVariables().length
    const nIntermediateVars = cc.intermediateVariables().length
    const nAdvancedVars = cc.advancedVariables().length
    return <div>So you have a cluster class with {nRequiredVars ? nRequiredVars : 'no'} required
        variables, {nBasicVars ? nBasicVars : 'no'} basic
        variables, {nIntermediateVars ? nIntermediateVars : 'no'} intermediate
        variables and {nAdvancedVars ? nAdvancedVars : 'no'} advanced variables. Deal with it.
        {
            nRequiredVars === 0 && <div><br/>
                Because there are no <b>required</b> variables, you can just click the &quot;Create Workload Cluster&quot; button below
                to create your workload cluster; all the cluster class variables on this page are optional.
            </div>
        }
    </div>
}

export default ClusterAttributeStep;
