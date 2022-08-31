// App imports
import { addErrorInfo, removeErrorInfo } from '../../../../../shared/utilities/Error.util';
import { AWS_FIELDS, AWS_NODE_PROFILE_NAMES } from '../../aws-mc-basic/AwsManagementClusterBasic.constants';
import { AwsDefaults, SelectedAvailabiltyZoneData } from '../default-service/AwsDefaults.service';
import { AWSKeyPair } from '../../../../../swagger-api/models/AWSKeyPair';
import { AWSAvailabilityZone } from '../../../../../swagger-api/models/AWSAvailabilityZone';
import { AwsService, AWSVirtualMachine, CancelablePromise } from '../../../../../swagger-api';
import {
    clearPreviousResourceData,
    DefaultOrchestrator,
    saveSingleResourceObject,
} from '../../../default-orchestrator/DefaultOrchestrator';
import { NodeProfileType } from '../../../../../shared/components/FormInputComponents/NodeProfile/NodeProfile';
import { RESOURCE } from '../../../../../state-management/actions/Resources.actions';
import { StoreDispatch, FormAction } from '../../../../../shared/types/types';
import { STORE_SECTION_FORM } from '../../../../../state-management/reducers/Form.reducer';
import { STORE_SECTION_RESOURCES } from '../../../../../state-management/reducers/Resources.reducer';
import { INPUT_CHANGE } from '../../../../../state-management/actions/Form.actions';

interface AwsOrchestratorProps {
    awsState: { [key: string]: any };
    awsDispatch: StoreDispatch;
    errorObject: { [key: string]: any };
    setErrorObject: (newErrorObject: { [key: string]: any }) => void;
}

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

export class AwsOrchestrator {
    static initOsImages(props: AwsOrchestratorProps) {
        const { awsState, awsDispatch, setErrorObject, errorObject } = props;
        DefaultOrchestrator.initResources<AWSVirtualMachine>({
            resourceName: AWS_FIELDS.OS_IMAGE,
            dispatch: awsDispatch,
            errorObject,
            setErrorObject,
            fetcher: () => AwsService.getAwsosImages(awsState[STORE_SECTION_FORM].REGION),
            fxnSelectDefault: AwsDefaults.selectDefaultOsImage,
        });
    }

    static initEC2KeyPairs(props: AwsOrchestratorProps) {
        const { awsDispatch, setErrorObject, errorObject } = props;
        DefaultOrchestrator.initResources<AWSAvailabilityZone>({
            resourceName: AWS_FIELDS.EC2_KEY_PAIR,
            dispatch: awsDispatch,
            errorObject,
            setErrorObject,
            fetcher: () => AwsService.getAwsKeyPairs(),
            fxnSelectDefault: AwsDefaults.selectDefaultEC2KeyPairs,
        });
    }

    static async initNodeProfile(props: AwsOrchestratorProps) {
        const { awsDispatch, setErrorObject, errorObject } = props;
        try {
            const nodeTypeList = await AwsService.getAwsNodeTypes();
            const nodeProfileList: { [key: string]: string } = {
                [AWS_NODE_PROFILE_NAMES.SINGLE_NODE]: '',
                [AWS_NODE_PROFILE_NAMES.HIGH_AVAILABILITY]: '',
                [AWS_NODE_PROFILE_NAMES.PRODUCTION_READY]: '',
            };
            Object.keys(nodeProfileList).forEach((nodeProfile) => {
                nodeProfileList[nodeProfile] = AwsDefaults.setDefaultNodeType(nodeTypeList, nodeProfile);
            });
            saveSingleResourceObject(awsDispatch, RESOURCE.ADD_RESOURCES, AWS_FIELDS.NODE_TYPE, nodeProfileList);
            setErrorObject(removeErrorInfo(errorObject, AWS_FIELDS.NODE_TYPE));
        } catch (e) {
            clearPreviousResourceData(awsDispatch, RESOURCE.ADD_RESOURCES, AWS_FIELDS.NODE_TYPE);
            setErrorObject(addErrorInfo(errorObject, e, AWS_FIELDS.NODE_TYPE));
        }
    }

    static initAvailabilityZones(props: AwsOrchestratorProps) {
        const { awsDispatch, setErrorObject, errorObject } = props;
        DefaultOrchestrator.initResources<AWSKeyPair>({
            resourceName: AWS_FIELDS.AVAILABILITY_ZONES,
            dispatch: awsDispatch,
            errorObject,
            setErrorObject,
            fetcher: () => AwsService.getAwsAvailabilityZones(),
        });
    }

    static initNodeTypesForAz(props: AwsOrchestratorProps, nodeProfile: string) {
        const { awsState, awsDispatch, setErrorObject, errorObject } = props;
        const azList = awsState[STORE_SECTION_RESOURCES][AWS_FIELDS.AVAILABILITY_ZONES];
        const defaultAZList: { [key: string]: string }[] = AwsDefaults.defaulAvailabilityZoneNameStrategy(azList, nodeProfile);
        const selectedAZObjects: { [key: string]: any } = {};
        for (let i = 0; i < defaultAZList.length; i++) {
            //  (a) store all the instance types in a segmented resource under THAT AZ
            AwsOrchestrator.initNodeTypeForAZ(
                defaultAZList[i],
                awsDispatch,
                errorObject,
                setErrorObject,
                selectedAZObjects,
                nodeProfile,
                i + 1
            );
        }
        // (b) store the "selected" instance type in the form section.
        awsDispatch({
            type: INPUT_CHANGE,
            field: AWS_FIELDS.SELECTED_AZ_OBJECTS,
            payload: selectedAZObjects,
        } as FormAction);
    }

    static initNodeTypeForAZ(
        az: { [key: string]: string },
        awsDispatch: StoreDispatch,
        errorObject: { [fieldName: string]: any },
        setErrorObject: (errorObject: { [fieldName: string]: any }) => void,
        selectedAZObjects: { [key: string]: any },
        nodeProfile: string,
        index: number
    ) {
        DefaultOrchestrator.initResources<SelectedAvailabiltyZoneData>({
            resourceName: AWS_FIELDS.AVAILIABILITY_ZONE_NODE_TYPES,
            segment: az.name,
            dispatch: awsDispatch,
            errorObject,
            setErrorObject,
            fetcher: () => {
                return new CancelablePromise(async (res, rej) => {
                    const nodeTypesByAZ = await AwsDefaults.createAZNodeType(az);
                    res(nodeTypesByAZ);
                });
            },

            fxnSelectDefault: (resources: SelectedAvailabiltyZoneData[]) => {
                const defaultAZ = AwsDefaults.defaulAvailabilityZoneNodeTypeStrategy(resources, nodeProfile, az.name);
                if (defaultAZ !== undefined) {
                    AwsOrchestrator.createStoredAZObjects(defaultAZ, selectedAZObjects, index);
                }
                return defaultAZ;
            },
        });
    }

    static createStoredAZObjects(defaultAZ: SelectedAvailabiltyZoneData, selectedAZObjects: { [key: string]: any }, index: number) {
        selectedAZObjects['availability-zone-' + index] = { id: defaultAZ.id, name: defaultAZ.name };
        selectedAZObjects['availability-zone-' + index + '-work-node-type'] = defaultAZ.workerNodeType;
        selectedAZObjects['availability-zone-' + index + '-public-subnet-id'] = defaultAZ.publicSubnetID;
        selectedAZObjects['availability-zone-' + index + '-private-subnet-id'] = defaultAZ.privateSubnetID;
    }
}
