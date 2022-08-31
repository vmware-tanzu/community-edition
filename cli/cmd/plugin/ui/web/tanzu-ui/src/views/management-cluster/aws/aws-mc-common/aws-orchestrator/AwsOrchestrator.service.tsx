// App imports
import { AWS_FIELDS, AWS_NODE_PROFILE_NAMES } from '../../aws-mc-basic/AwsManagementClusterBasic.constants';
import { AwsDefaults } from '../default-service/AwsDefaults.service';
import { AWSKeyPair } from '../../../../../swagger-api/models/AWSKeyPair';
import { AwsService, AWSVirtualMachine } from '../../../../../swagger-api';
import { AWSAvailabilityZone } from '../../../../../swagger-api/models/AWSAvailabilityZone';
import { DefaultOrchestrator } from '../../../default-orchestrator/DefaultOrchestrator';
import { getDefaultNodeTypes } from '../../../../../shared/constants/defaults/aws.defaults';
import { INPUT_CHANGE } from '../../../../../state-management/actions/Form.actions';
import { NodeProfileType } from '../../../../../shared/components/FormInputComponents/NodeProfile/NodeProfile';
import { StoreDispatch, FormAction } from '../../../../../shared/types/types';
import { STORE_SECTION_FORM } from '../../../../../state-management/reducers/Form.reducer';

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
        DefaultOrchestrator.initResources<AWSKeyPair>({
            resourceName: AWS_FIELDS.EC2_KEY_PAIR,
            dispatch: awsDispatch,
            errorObject,
            setErrorObject,
            fetcher: () => AwsService.getAwsKeyPairs(),
            fxnSelectDefault: AwsDefaults.selectDefaultEC2KeyPairs,
        });
    }

    static initControlPlaneNodeType(props: AwsOrchestratorProps, nodeProfile: string) {
        const { awsDispatch, setErrorObject, errorObject } = props;
        DefaultOrchestrator.initResources<string>({
            resourceName: AWS_FIELDS.NODE_TYPE,
            dispatch: awsDispatch,
            errorObject,
            setErrorObject,
            fetcher: () => AwsService.getAwsNodeTypes(),
            fxnSelectDefault: (nodeTypes: string[]) => AwsDefaults.selectDefaultNodeType(nodeTypes, getDefaultNodeTypes(nodeProfile)),
        });
    }

    static initAvailabilityZones(props: AwsOrchestratorProps) {
        const { awsDispatch, setErrorObject, errorObject } = props;
        DefaultOrchestrator.initResources<AWSAvailabilityZone>({
            resourceName: AWS_FIELDS.AVAILABILITY_ZONES,
            dispatch: awsDispatch,
            errorObject,
            setErrorObject,
            fetcher: () => AwsService.getAwsAvailabilityZones(),
        });
    }

    static getAZFieldsForNodeProfile(nodeProfile: string): { [key: string]: AWS_FIELDS }[] {
        if (nodeProfile === AWS_NODE_PROFILE_NAMES.SINGLE_NODE) {
            return [{ name: AWS_FIELDS.AVAILABILITY_ZONE_1, workerNodeType: AWS_FIELDS.AVAILABILITY_ZONE_1_NODE_TYPE }];
        }
        return [
            { name: AWS_FIELDS.AVAILABILITY_ZONE_1, workerNodeType: AWS_FIELDS.AVAILABILITY_ZONE_1_NODE_TYPE },
            { name: AWS_FIELDS.AVAILABILITY_ZONE_2, workerNodeType: AWS_FIELDS.AVAILABILITY_ZONE_2_NODE_TYPE },
            { name: AWS_FIELDS.AVAILABILITY_ZONE_3, workerNodeType: AWS_FIELDS.AVAILABILITY_ZONE_3_NODE_TYPE },
        ];
    }

    static initNodeTypesForAz(props: AwsOrchestratorProps, azs: AWSAvailabilityZone[], nodeProfile: string) {
        const { awsDispatch, setErrorObject, errorObject } = props;
        const azFields = AwsOrchestrator.getAZFieldsForNodeProfile(nodeProfile);
        for (let i = 0; i < azs.length; i++) {
            AwsOrchestrator.initNodeTypeForAZ(
                azs[i].name ?? '',
                awsDispatch,
                errorObject,
                setErrorObject,
                nodeProfile,
                azFields[i].workerNodeType
            );
            awsDispatch({
                type: INPUT_CHANGE,
                field: azFields[i].name,
                payload: azs[i].name,
            } as FormAction);
        }
    }

    static initNodeTypeForAZ(
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
            fxnSelectDefault: (nodeTypes: string[]) => AwsDefaults.selectDefaultNodeType(nodeTypes, getDefaultNodeTypes(nodeProfile)),
            fieldName: field,
        });
    }
}
