import React, { useContext } from 'react';
import { SubmitHandler, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { StepProps } from '../../shared/components/wizard/Wizard';
import * as yup from 'yup';
import { CdsButton } from '@cds/react/button';
import { WcStore } from '../../state-management/stores/Store.wc';
import { SUBMIT_FORM } from '../../state-management/actions/Form.actions';

interface ClusterTopologyStepProps extends StepProps {
    selectedManagementCluster: string,
}

interface ClusterTopologyStepFormInputs {
    WORKLOAD_CLUSTER_NAME: string;
}

const clusterTopologyStepFormSchema = yup.object({
    WORKLOAD_CLUSTER_NAME: yup.string().nullable().required('Please enter a name for your workload cluster')
}).required();

function ClusterTopologyStep(props: Partial<ClusterTopologyStepProps>) {
    const { handleValueChange, currentStep, goToStep, submitForm } = props;
    const { state, dispatch } = useContext(WcStore);
    const methods = useForm<ClusterTopologyStepFormInputs>({
        resolver: yupResolver(clusterTopologyStepFormSchema),
    });
    const {
        register,
        handleSubmit,
        setValue,
        formState: { errors },
    } = methods;

    const onSubmit: SubmitHandler<ClusterTopologyStepFormInputs> = (data) => {
        if (Object.keys(errors).length === 0) {
            if (goToStep && currentStep && submitForm) {
                console.log(`submitting form is dispatching a SUBMIT_FORM with ${JSON.stringify(data)}`);
                dispatch({
                    type: SUBMIT_FORM,
                    payload: data
                });
                goToStep(currentStep + 1);
                submitForm(currentStep);
            }
        }
    };

    return <div>
        I guess we are talking about provisioning a workload cluster using management cluster: {state.data.WORKLOAD_CLUSTER_NAME} <br/>
        state={JSON.stringify(state)}
        <CdsButton onClick={handleSubmit(onSubmit)}>NEXT</CdsButton>
    </div>;
}

export default ClusterTopologyStep;
