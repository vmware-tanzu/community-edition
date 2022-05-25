// React imports
import React, { ChangeEvent, useContext, useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { SubmitHandler, useForm } from 'react-hook-form';

// Library imports
import { AnyObjectSchema } from 'yup';
import { CdsButton } from '@cds/react/button';
import { CdsProgressCircle } from '@cds/react/progress-circle';
import { CdsSelect } from '@cds/react/select';
import { yupResolver } from '@hookform/resolvers/yup/dist/yup';

// App imports
import { CancelablePromise } from '../../swagger-api';
import { CCVAR_CHANGE, INPUT_CHANGE } from '../../state-management/actions/Form.actions';
import { CCCategory, CCDefinition } from '../../shared/models/ClusterClass';
import { CCMultipleVariablesDisplay, createFormSchemaCC } from './ClusterClassVariableDisplay';
import { getFieldData } from '../../state-management/reducers/Form.reducer';
import { getSelectedManagementCluster, getValueFromChangeEvent } from './WorkloadClusterUtility';
import ManagementClusterInfoBanner from './ManagementClusterInfoBanner';
import { NavRoutes } from '../../shared/constants/NavRoutes.constants';
import { retrieveClusterClass } from '../../shared/services/ClusterClass.service';
import { StepProps } from '../../shared/components/wizard/Wizard';
import { TOGGLE_WC_CC_CATEGORY } from '../../state-management/actions/Ui.actions';
import { WcStore } from '../../state-management/stores/Store.wc';

interface ClusterAttributeStepProps extends StepProps {
    retrieveAvailableClusterClasses: (mcName: string) => CancelablePromise<Array<string>>;
}

const INSTRUCTION_CC_DEFAULT = 'Fill out the variables you wish to set as you create your workload cluster.';
const INSTRUCTION_CC_NO_DEFINITION = 'We were unable to retrieve the cluster class definition for this cluster class! So sorry.';

function ClusterAttributeStep(props: Partial<ClusterAttributeStepProps>) {
    const { retrieveAvailableClusterClasses, handleValueChange, currentStep, goToStep, submitForm } = props;
    const { state, dispatch } = useContext(WcStore);
    const [ccDefinition, setCcDefinition] = useState<CCDefinition>();
    const [formSchema, setFormSchema] = useState<AnyObjectSchema>();
    const [ccNames, setCcNames] = useState<string[]>(state.data.AVAILABLE_CLUSTER_CLASSES);
    const [selectedCc, setSelectedCc] = useState<string>(state.data.SELECTED_CLUSTER_CLASS);
    const [loadingClusterClass, setLoadingClusterClass] = useState<boolean>(false);
    // associates a category name with a fxn that will toggle the expanded flag in the data store for that category
    const [categoryToggleFxns] = useState<{ [category: string]: () => void }>({});

    const cluster = getSelectedManagementCluster(state);
    const getFieldValue = (fieldName: string) => {
        return cluster && cluster.name ? getFieldData(fieldName, cluster.name, state) : undefined;
    };

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
        if (cluster && cluster.name) {
            // get the cluster classes on the selected management cluster
            if (retrieveAvailableClusterClasses && cluster.name && handleValueChange) {
                retrieveAvailableClusterClasses(cluster.name).then((availableClusterClasses: string[]) => {
                    handleValueChange(INPUT_CHANGE, 'AVAILABLE_CLUSTER_CLASSES', availableClusterClasses, currentStep, errors);
                    setCcNames(availableClusterClasses);
                });
            }
        }
    }, [cluster]); // eslint-disable-line react-hooks/exhaustive-deps

    useEffect(() => {
        // createToggleCategoryExpandedFxn returns a fxn that will toggle the expanded flag in the data store for that category
        // (The point is: the accordion requires a method that doesn't take a parameter, and we need the
        // category, so we create a custom fxn that already knows the category and doesn't need a parameter)
        const createToggleCategoryExpandedFxn = (category: string): (() => void) => {
            return () => {
                dispatch({ type: TOGGLE_WC_CC_CATEGORY, locationData: category });
            };
        };

        if (selectedCc) {
            // TODO: remove setTimeout(), which is just here to simulate a backend call delay
            setTimeout(() => {
                if (cluster.name) {
                    retrieveClusterClass(cluster.name, selectedCc, (ccDef) => {
                        setLoadingClusterClass(false);
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
            }, 500);
        }
    }, [selectedCc]); // eslint-disable-line react-hooks/exhaustive-deps

    // TODO: we will likely need to navigate to a WC-specific progress route, but for now, just to be able to demo...
    const navigateToProgress = (): void => {
        navigate('/' + NavRoutes.DEPLOY_PROGRESS);
    };

    if (!cluster) {
        return <div>No management cluster has been selected (how did you get to this step?!)</div>;
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

    const onSelectCc = (evt: ChangeEvent<HTMLSelectElement>) => {
        const selCc = evt.target.value;
        if (handleValueChange) {
            handleValueChange(INPUT_CHANGE, 'SELECTED_CLUSTER_CLASS', selCc, currentStep, errors);
        }
        setSelectedCc(selCc);
        if (selCc) {
            setLoadingClusterClass(true);
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
            {SelectClusterClass(ccNames, onSelectCc)}
            <br />
            {loadingClusterClass && ClusterClassLoading()}
            {!loadingClusterClass && selectedCc && ccDefinition && CCStepInstructions(ccDefinition)}
            {!loadingClusterClass &&
                selectedCc &&
                ccDefinition?.categories?.map((ccCategory: CCCategory) => {
                    const expanded = state.ui.wcCcCategoryExpanded[ccCategory.name];
                    const toggleCategoryExpanded = categoryToggleFxns[ccCategory.name];
                    const options = { register, errors, expanded, onValueChange, toggleCategoryExpanded, getFieldValue };
                    return CCMultipleVariablesDisplay(ccCategory.variables, ccCategory, options);
                })}
            <br />
            <CdsButton
                className="cluster-action-btn"
                status="primary"
                onClick={handleSubmit(onSubmit)}
                disabled={!ccDefinition || !selectedCc || loadingClusterClass}
            >
                Create Workload Cluster
            </CdsButton>
        </div>
    );
}

function ClusterClassLoading() {
    return (
        <>
            <div className="text-white" cds-layout="col:1"></div>
            <div cds-layout="horizontal gap:sm col:11" cds-theme="dark">
                <CdsProgressCircle size="xl" status="info"></CdsProgressCircle>
            </div>
        </>
    );
}

function CCStepInstructions(cc: CCDefinition | undefined) {
    let instructions = INSTRUCTION_CC_DEFAULT;
    if (!cc) {
        instructions = INSTRUCTION_CC_NO_DEFINITION;
    } else if (cc.instructions) {
        instructions = cc.instructions;
    }
    return (
        <div>
            {instructions}
            <br />
            &nbsp;
        </div>
    );
}

function SelectClusterClass(availableCCs: string[], onValueChange: (evt: ChangeEvent<HTMLSelectElement>) => void) {
    return (
        <CdsSelect layout="compact" controlWidth="shrink">
            <label>Use cluster class:</label>
            <select className="select-md-width" onChange={onValueChange}>
                {availableCCs.length > 1 && <option></option>}
                {availableCCs &&
                    availableCCs.map((ccName) => (
                        <option key={ccName} value={ccName}>
                            {ccName}
                        </option>
                    ))}
            </select>
        </CdsSelect>
    );
}

export default ClusterAttributeStep;
