// React imports
import React, { ChangeEvent, useContext, useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { SubmitHandler, useForm } from 'react-hook-form';

// Library imports
import { CdsButton } from '@cds/react/button';
import { yupResolver } from '@hookform/resolvers/yup/dist/yup';

// App imports
import { CCVAR_CHANGE } from '../../state-management/actions/Form.actions';
import { CCCategory, CCDefinition, ClusterClassDefinition } from '../../shared/models/ClusterClass';
import { CCMultipleVariablesDisplay, createFormSchemaCC } from './ClusterClassVariableDisplay';
import { getSelectedManagementCluster, getValueFromChangeEvent } from './WorkloadClusterUtility';
import ManagementClusterInfoBanner from './ManagementClusterInfoBanner';
import { NavRoutes } from '../../shared/constants/NavRoutes.constants';
import { StepProps } from '../../shared/components/wizard/Wizard';
import { TOGGLE_WC_CC_CATEGORY } from '../../state-management/actions/Ui.actions';
import { WcStore } from '../../state-management/stores/Store.wc';
import { retrieveClusterClass } from '../../shared/services/ClusterClass.service';
import { AnyObjectSchema } from 'yup';

interface ClusterAttributeStepProps extends StepProps {
    retrieveClusterClassDefinition: (mc: string | undefined) => ClusterClassDefinition | undefined;
}

function ClusterAttributeStep(props: Partial<ClusterAttributeStepProps>) {
    const { handleValueChange, currentStep, goToStep, submitForm, retrieveClusterClassDefinition } = props;
    const { state, dispatch } = useContext(WcStore);
    const [ccDefinition, setCcDefinition] = useState<CCDefinition>();
    const [formSchema, setFormSchema] = useState<AnyObjectSchema>();
    // associates a category name with a fxn that will toggle the expanded flag in the data store for that category
    const [categoryToggleFxns] = useState<{ [category: string]: () => void }>({});

    const cluster = getSelectedManagementCluster(state);

    const methods = useForm({
        resolver: formSchema ? yupResolver(formSchema) : undefined,
    });
    const {
        register,
        handleSubmit,
        formState: { errors },
        setValue,
    } = methods;

    const navigate = useNavigate();

    useEffect(() => {
        // This fxn returns a fxn that will toggle the expanded flag in the data store for that category
        // (The point is: the accordion requires a method that doesn't take a parameter, and we need the
        // category, so we create a custom fxn that already knows the category and doesn't need a parameter)
        const createToggleCategoryExpandedFxn = (category: string): (() => void) => {
            return () => {
                dispatch({ type: TOGGLE_WC_CC_CATEGORY, locationData: category });
            };
        };

        if (cluster.name) {
            // TODO: get the cluster class list (instead of hard-coded GET), and allow user to select
            retrieveClusterClass(cluster.name, `tkg-${cluster.provider}-default`, (ccDef) => {
                setCcDefinition(ccDef);
                setFormSchema(createFormSchemaCC(ccDef));
                ccDef.categories?.forEach((category) => {
                    categoryToggleFxns[category.name] = createToggleCategoryExpandedFxn(category.name);
                    // if the category wants to default to display "open", toggle it now using the fxn we just created
                    if (category.displayOpen) {
                        categoryToggleFxns[category.name]();
                    }
                });
            });
        }
    }, [cluster, categoryToggleFxns, dispatch]); // only first needed, others included for linter only

    // TODO: we will likely need to navigate to a WC-specific progress route, but for now, just to be able to demo...
    const navigateToProgress = (): void => {
        navigate('/' + NavRoutes.DEPLOY_PROGRESS);
    };

    if (!retrieveClusterClassDefinition) {
        return <div>Programmer error: ClusterAttributeStep did not receive retrieveClusterClassDefinition!</div>;
    }
    if (!cluster) {
        return <div>No management cluster has been selected (how did you get to this step?!)</div>;
    }
    if (!ccDefinition) {
        return <div>We were unable to retrieve a ClusterClass object for management cluster {cluster.name}</div>;
    }

    const onSubmit: SubmitHandler<any> = (data) => {
        const nErrors = Object.keys(errors).length;
        if (nErrors === 0) {
            if (goToStep && currentStep && submitForm) {
                // TODO: we'll need to call the backend service to actually do the deploying
                // goToStep(currentStep + 1);
                submitForm(currentStep);
                navigateToProgress();
            }
        } else {
            console.log(`ClusterAttributeStep has an invalid form submission (${nErrors} errors)`);
        }
    };

    const onValueChange = (evt: ChangeEvent<HTMLSelectElement>) => {
        const value = getValueFromChangeEvent(evt);
        const varName = evt.target.name;
        setValue(varName, value, { shouldValidate: true });
        if (handleValueChange) {
            handleValueChange(CCVAR_CHANGE, varName, value, currentStep, errors, { clusterName: cluster.name });
        } else {
            console.error('ClusterAttributeStep unable to find a handleValueChange handler!');
        }
    };

    return (
        <div>
            {ManagementClusterInfoBanner(cluster)}
            <br />
            {CCStepInstructions(ccDefinition)}
            {ccDefinition?.categories?.map((ccCategory: CCCategory) => {
                const expanded = state.ui.wcCcCategoryExpanded[ccCategory.name];
                const toggleCategoryExpanded = categoryToggleFxns[ccCategory.name];
                const options = { register, errors, expanded, onValueChange, toggleCategoryExpanded };
                return CCMultipleVariablesDisplay(ccCategory.variables, ccCategory, options);
            })}
            <br />
            <CdsButton className="cluster-action-btn" status="primary" onClick={handleSubmit(onSubmit)}>
                Create Workload Cluster
            </CdsButton>
        </div>
    );
}

function CCStepInstructions(cc: CCDefinition | undefined) {
    if (!cc) {
        return <div>There is no cluster class definition, so you cannot do this step! So sorry.</div>;
    }
    return <div>Fill out the variables you wish to set as you create your workload cluster.</div>;
}

export default ClusterAttributeStep;
