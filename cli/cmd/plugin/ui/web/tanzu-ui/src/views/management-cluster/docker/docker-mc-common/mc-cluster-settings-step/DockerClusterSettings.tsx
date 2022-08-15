// React imports
import React, { useContext } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';

// Library imports
import { CdsAlert, CdsAlertGroup } from '@cds/react/alert';
import { CdsButton } from '@cds/react/button';
import { CdsIcon } from '@cds/react/icon';
import * as yup from 'yup';
import { yupResolver } from '@hookform/resolvers/yup/dist/yup';

// App imports
import { clusterNameValidation, ClusterName } from '../../../../../shared/components/FormInputComponents/ClusterName/ClusterName';
import { DockerStore } from '../../../../../state-management/stores/Docker.store';
import { DOCKER_FIELDS } from '../../docker-mc-basic/DockerManagementClusterBasic.constants';
import { FormAction } from '../../../../../shared/types/types';
import { INPUT_CHANGE } from '../../../../../state-management/actions/Form.actions';
import { StepProps } from '../../../../../shared/components/wizard/Wizard';
import { STORE_SECTION_FORM } from '../../../../../state-management/reducers/Form.reducer';
import UseUpdateTabStatus from '../../../../../shared/components/wizard/UseUpdateTabStatus.hooks';

interface DockerClusterSettingFormInputs {
    [DOCKER_FIELDS.CLUSTER_NAME]: string;
}

function DockerClusterSettings(props: Partial<StepProps>) {
    const { updateTabStatus, currentStep, deploy } = props;
    const { dockerState, dockerDispatch } = useContext(DockerStore);

    const formSchema = yup.object({ [DOCKER_FIELDS.CLUSTER_NAME]: clusterNameValidation() }).required();
    const methods = useForm<DockerClusterSettingFormInputs>({
        resolver: yupResolver(formSchema),
        mode: 'all',
    });

    const {
        handleSubmit,
        formState: { errors },
    } = methods;
    // update tab status bar
    if (updateTabStatus) {
        UseUpdateTabStatus(errors, currentStep, updateTabStatus);
    }

    const canContinue = (): boolean => {
        return Object.keys(errors).length === 0 && dockerState[STORE_SECTION_FORM][DOCKER_FIELDS.CLUSTER_NAME];
    };

    const onSubmit: SubmitHandler<DockerClusterSettingFormInputs> = (data) => {
        if (canContinue() && deploy) {
            deploy();
        }
    };

    const onClusterNameChange = (clusterName: string) => {
        dockerDispatch({
            type: INPUT_CHANGE,
            field: DOCKER_FIELDS.CLUSTER_NAME,
            payload: clusterName,
        } as FormAction);
    };
    return (
        <FormProvider {...methods}>
            <div className="wizard-content-container">
                <h2 cds-layout="m-t:md m-b:xl" cds-text="title">
                    Docker Management Cluster Settings
                </h2>
                <div cds-layout="grid gap:m p-b:lg" key="section-holder">
                    <div cds-layout="col:4" key="cluster-name-section">
                        <ClusterName
                            field={DOCKER_FIELDS.CLUSTER_NAME}
                            clusterNameChange={onClusterNameChange}
                            placeholderClusterName="my-docker-cluster"
                            defaultClusterName={dockerState[STORE_SECTION_FORM][DOCKER_FIELDS.CLUSTER_NAME]}
                        />
                    </div>
                    <div cds-layout="col:8 p-l:xl" key="instance-type-section">
                        <CdsAlertGroup status="info" cds-layout="m-b:md">
                            <CdsAlert>A single node Management Cluster will be created on your local workstation.</CdsAlert>
                        </CdsAlertGroup>
                    </div>
                </div>
                <div cds-layout="grid col:12 p-t:lg">
                    <CdsButton cds-layout="col:start-1" status="success" onClick={handleSubmit(onSubmit)} disabled={!canContinue()}>
                        <CdsIcon shape="cluster" size="sm"></CdsIcon>
                        Create Management cluster
                    </CdsButton>
                </div>
            </div>
        </FormProvider>
    );
}

export default DockerClusterSettings;
