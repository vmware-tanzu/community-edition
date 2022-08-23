// App imports
import { AwsDefaults } from '../default-service/AwsDefaults.service';
import { AwsService, AWSVirtualMachine } from '../../../../../swagger-api';
import { AWSKeyPair } from '../../../../../swagger-api/models/AWSKeyPair';

import { AWS_FIELDS, AWS_NODE_PROFILE_NAMES } from '../../aws-mc-basic/AwsManagementClusterBasic.constants';
import { DefaultOrchestrator } from '../../../default-orchestrator/DefaultOrchestrator';
import { StoreDispatch } from '../../../../../shared/types/types';
import { STORE_SECTION_FORM } from '../../../../../state-management/reducers/Form.reducer';
import { NodeProfileType } from '../../../../../shared/components/FormInputComponents/NodeProfile/NodeProfile';
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
    static async initOsImages(props: AwsOrchestratorProps) {
        const { awsState, awsDispatch, setErrorObject, errorObject } = props;
        await DefaultOrchestrator.initResources<AWSVirtualMachine>({
            resourceName: AWS_FIELDS.OS_IMAGE,
            dispatch: awsDispatch,
            errorObject,
            setErrorObject,
            fetcher: () => AwsService.getAwsosImages(awsState[STORE_SECTION_FORM].REGION),
            fxnSelectDefault: AwsDefaults.selectDefaultOsImage,
        });
    }

    static async initEC2KeyPairs(props: AwsOrchestratorProps) {
        const { awsDispatch, setErrorObject, errorObject } = props;
        await DefaultOrchestrator.initResources<AWSKeyPair>({
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
            Object.keys(nodeProfileList).map((nodeProfile) => {
                nodeProfileList[nodeProfile] = AwsDefaults.setDefaultNodeType(nodeTypeList, nodeProfile);
            });
            saveCurrentResourceData(awsDispatch, RESOURCE.AWS_ADD_RESOURCES, AWS_FIELDS.NODE_TYPE, nodeProfileList);
            setErrorObject(removeErrorInfo(errorObject, AWS_FIELDS.NODE_TYPE));
        } catch (e) {
            clearPreviousResourceData(awsDispatch, RESOURCE.AWS_ADD_RESOURCES, AWS_FIELDS.NODE_TYPE);
            setErrorObject(addErrorInfo(errorObject, e, AWS_FIELDS.NODE_TYPE));
        }
    }
}
