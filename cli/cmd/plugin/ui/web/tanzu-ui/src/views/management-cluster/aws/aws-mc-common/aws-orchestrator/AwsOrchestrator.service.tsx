// React imports
import { useContext, useEffect, useState } from 'react';
// App imports
import { AWS_FIELDS, AWS_NODE_PROFILE_NAMES } from '../../aws-mc-basic/AwsManagementClusterBasic.constants';
import { AwsStore } from '../../store/Aws.store.mc';
import { AwsDefaults } from '../default-service/AwsDefaults.service';
import { AWSKeyPair } from '../../../../../swagger-api/models/AWSKeyPair';
import { AwsService, AWSVirtualMachine } from '../../../../../swagger-api';
import { AWSAvailabilityZone } from '../../../../../swagger-api/models/AWSAvailabilityZone';
import { CONNECTION_STATUS } from '../../../../../shared/components/ConnectionNotification/ConnectionNotification';
import { DefaultOrchestrator } from '../../../default-orchestrator/DefaultOrchestrator';
import { getDefaultNodeTypes } from '../../../../../shared/constants/defaults/aws.defaults';
import { STORE_SECTION_RESOURCES } from '../../../../../state-management/reducers/Resources.reducer';
import { INPUT_CHANGE } from '../../../../../state-management/actions/Form.actions';
import { NodeProfileType } from '../../../../../shared/components/FormInputComponents/NodeProfile/NodeProfile';
import { StoreDispatch, FormAction } from '../../../../../shared/types/types';
import { STORE_SECTION_FORM } from '../../../../../state-management/reducers/Form.reducer';
import { findFirstMatchingOption } from '../../../../../shared/utilities/Array.util';

export const nodeProfiles: NodeProfileType[] = [
    {
        id: 'SINGLE_NODE',
        label: 'Single node',
        icon: 'block',
        description: 'Create a single control plane node with a medium instance type',
    },
    {
        id: 'HIGH_AVAILABILITY',
        label: 'High availability',
        icon: 'blocks-group',
        description: 'Create a multi-node control plane with a medium instance type',
    },
    {
        id: 'PRODUCTION_READY',
        label: 'Production-ready (High availability)',
        icon: 'blocks-group',
        isSolidIcon: true,
        description: 'Create a multi-node control plane with a large instance type',
    },
];

function usePrerequisite() {
    const [errorObject, setErrorObject] = useState({});
    const { awsState, awsDispatch } = useContext(AwsStore);

    return [{ awsState, awsDispatch, errorObject, setErrorObject }];
}

function useInitOsImages(connectionStatus: CONNECTION_STATUS) {
    const [{ awsState, awsDispatch, errorObject, setErrorObject }] = usePrerequisite();

    useEffect(() => {
        const fetchOsImages = async () => {
            await DefaultOrchestrator.initResources<AWSVirtualMachine>({
                resourceName: AWS_FIELDS.OS_IMAGE,
                errorObject,
                setErrorObject,
                dispatch: awsDispatch,
                fetcher: () => AwsService.getAwsosImages(awsState[STORE_SECTION_FORM].REGION),
                fxnSelectDefault: AwsDefaults.selectDefaultOsImage,
            });
        };
        if (connectionStatus === CONNECTION_STATUS.CONNECTED) {
            fetchOsImages();
        }
    }, [connectionStatus]);

    return [errorObject, setErrorObject];
}

function useInitEC2KeyPairs(connectionStatus: CONNECTION_STATUS, setKeyPairLoading: React.Dispatch<React.SetStateAction<boolean>>) {
    const [{ awsDispatch, errorObject, setErrorObject }] = usePrerequisite();

    useEffect(() => {
        const fetchEC2KeyPairs = async () => {
            setKeyPairLoading(true);
            await DefaultOrchestrator.initResources<AWSKeyPair>({
                resourceName: AWS_FIELDS.EC2_KEY_PAIR,
                dispatch: awsDispatch,
                errorObject,
                setErrorObject,
                fetcher: () => AwsService.getAwsKeyPairs(),
                fxnSelectDefault: AwsDefaults.selectDefaultEC2KeyPairs,
            });
            setKeyPairLoading(false);
        };
        if (connectionStatus === CONNECTION_STATUS.CONNECTED) {
            fetchEC2KeyPairs();
        }
    }, [connectionStatus]);

    return [errorObject, setErrorObject];
}

function useInitControlPlaneNodeType(connectionStatus: CONNECTION_STATUS) {
    const [{ awsState, awsDispatch, errorObject, setErrorObject }] = usePrerequisite();

    const selectedNodeProfile = awsState[STORE_SECTION_FORM][AWS_FIELDS.NODE_PROFILE];

    useEffect(() => {
        const fetchControlPlaneNodeType = async () => {
            DefaultOrchestrator.initResources<string>({
                resourceName: AWS_FIELDS.NODE_TYPE,
                dispatch: awsDispatch,
                errorObject,
                setErrorObject,
                fetcher: () => AwsService.getAwsNodeTypes(),
                fxnSelectDefault: (nodeTypes: string[]) =>
                    findFirstMatchingOption(nodeTypes, getDefaultNodeTypes(selectedNodeProfile) ?? []),
            });
        };
        if (connectionStatus === CONNECTION_STATUS.CONNECTED) {
            fetchControlPlaneNodeType();
        }
    }, [selectedNodeProfile, connectionStatus]);
    return [errorObject, setErrorObject];
}

function useInitAvailabilityZones(connectionStatus: CONNECTION_STATUS) {
    const [{ awsDispatch, errorObject, setErrorObject }] = usePrerequisite();

    useEffect(() => {
        const fetchAvailabilityZones = async () => {
            DefaultOrchestrator.initResources<AWSAvailabilityZone>({
                resourceName: AWS_FIELDS.AVAILABILITY_ZONES,
                dispatch: awsDispatch,
                errorObject,
                setErrorObject,
                fetcher: () => AwsService.getAwsAvailabilityZones(),
            });
        };
        if (connectionStatus === CONNECTION_STATUS.CONNECTED) {
            fetchAvailabilityZones();
        }
    }, [connectionStatus]);

    return [errorObject, setErrorObject];
}

function getAZFieldsForNodeProfile(nodeProfile: string): { [key: string]: AWS_FIELDS }[] {
    if (nodeProfile === AWS_NODE_PROFILE_NAMES.SINGLE_NODE) {
        return [{ name: AWS_FIELDS.AVAILABILITY_ZONE_1, workerNodeType: AWS_FIELDS.AVAILABILITY_ZONE_1_NODE_TYPE }];
    }
    return [
        { name: AWS_FIELDS.AVAILABILITY_ZONE_1, workerNodeType: AWS_FIELDS.AVAILABILITY_ZONE_1_NODE_TYPE },
        { name: AWS_FIELDS.AVAILABILITY_ZONE_2, workerNodeType: AWS_FIELDS.AVAILABILITY_ZONE_2_NODE_TYPE },
        { name: AWS_FIELDS.AVAILABILITY_ZONE_3, workerNodeType: AWS_FIELDS.AVAILABILITY_ZONE_3_NODE_TYPE },
    ];
}

function useInitNodeTypesForAz() {
    const [{ awsState, awsDispatch, errorObject, setErrorObject }] = usePrerequisite();

    const nodeProfile = awsState[STORE_SECTION_FORM][AWS_FIELDS.NODE_PROFILE];

    useEffect(() => {
        if (awsState[STORE_SECTION_RESOURCES][AWS_FIELDS.AVAILABILITY_ZONES]) {
            const azs = AwsDefaults.getDefaulAvailabilityZones(
                awsState[STORE_SECTION_RESOURCES][AWS_FIELDS.AVAILABILITY_ZONES],
                awsState[STORE_SECTION_FORM][AWS_FIELDS.NODE_PROFILE]
            );
            const fetchOsImages = async () => {
                const azFields = AwsOrchestrator.getAZFieldsForNodeProfile(nodeProfile);
                for (let i = 0; i < azs.length; i++) {
                    initNodeTypeForAZ(azs[i].name ?? '', awsDispatch, errorObject, setErrorObject, nodeProfile, azFields[i].workerNodeType);
                    awsDispatch({
                        type: INPUT_CHANGE,
                        field: azFields[i].name,
                        payload: azs[i].name,
                    } as FormAction);
                }
            };
            fetchOsImages();
        }
    }, [nodeProfile, awsState[STORE_SECTION_RESOURCES][AWS_FIELDS.AVAILABILITY_ZONES]]);

    return [errorObject, setErrorObject];
}

function initNodeTypeForAZ(
    azName: string,
    awsDispatch: StoreDispatch,
    errorObject: { [fieldName: string]: any },
    setErrorObject: (errorObject: { [fieldName: string]: any }) => void,
    nodeProfile: string,
    field: AWS_FIELDS
) {
    DefaultOrchestrator.initResources<string>({
        resourceName: AWS_FIELDS.NODE_TYPES_BY_AZ,
        segment: azName,
        dispatch: awsDispatch,
        errorObject,
        setErrorObject,
        fetcher: () => AwsService.getAwsNodeTypes(azName),
        fxnSelectDefault: (nodeTypes: string[]) => findFirstMatchingOption(nodeTypes, getDefaultNodeTypes(nodeProfile) ?? []),
        fieldName: field,
    });
}

const AwsOrchestrator = {
    useInitOsImages,
    useInitControlPlaneNodeType,
    useInitEC2KeyPairs,
    useInitAvailabilityZones,
    useInitNodeTypesForAz,
    getAZFieldsForNodeProfile,
};
export default AwsOrchestrator;
