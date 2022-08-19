// App imports
import { AwsDefaults } from '../default-service/AwsDefaults.service';
import { AwsService, AWSVirtualMachine } from '../../../../../swagger-api';
import { AWSKeyPair } from '../../../../../swagger-api/models/AWSKeyPair';
import { AWS_FIELDS, AWS_NODE_PROFILE_NAMES } from '../../aws-mc-basic/AwsManagementClusterBasic.constants';
import { INPUT_CHANGE } from '../../../../../state-management/actions/Form.actions';
import { FormAction, StoreDispatch } from '../../../../../shared/types/types';
import { STORE_SECTION_FORM } from '../../../../../state-management/reducers/Form.reducer';
import { RESOURCE } from '../../../../../state-management/actions/Resources.actions';
import { NodeInstanceType } from '../../../../../shared/components/FormInputComponents/NodeProfile/NodeProfile';
import {
    clearPreviousResourceData,
    saveCurrentResourceData,
    removeErrorInfo,
    addErrorInfo,
} from '../../../default-orchestrator/DefaultOrchestrator';
interface AwsOrchestratorProps {
    awsState: { [key: string]: any };
    awsDispatch: StoreDispatch;
    errorObject: { [key: string]: any };
    setErrorObject: (newErrorObject: { [key: string]: any }) => void;
}

export const nodeInstanceTypes: NodeInstanceType[] = [
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
        try {
            const osImages = await AwsService.getAwsosImages(awsState[STORE_SECTION_FORM].REGION);
            saveCurrentResourceData(awsDispatch, RESOURCE.AWS_ADD_RESOURCES, AWS_FIELDS.OS_IMAGE, osImages);
            setDefaultOsImage(awsDispatch, osImages);
            setErrorObject(removeErrorInfo(errorObject, AWS_FIELDS.OS_IMAGE));
        } catch (e) {
            clearPreviousResourceData(awsDispatch, RESOURCE.AWS_ADD_RESOURCES, AWS_FIELDS.OS_IMAGE);
            setErrorObject(addErrorInfo(errorObject, e, AWS_FIELDS.OS_IMAGE));
        }
    }

    static async initEC2KeyPairs(props: AwsOrchestratorProps, setKeyPairs: (keyPairs: AWSKeyPair[]) => void) {
        const { awsDispatch, setErrorObject, errorObject } = props;
        try {
            const keyPairs = await AwsService.getAwsKeyPairs();
            saveCurrentResourceData(awsDispatch, RESOURCE.AWS_ADD_RESOURCES, AWS_FIELDS.EC2_KEY_PAIR, keyPairs);
            setDefaultEC2KeyPair(awsDispatch, keyPairs);
            setErrorObject(removeErrorInfo(errorObject, AWS_FIELDS.OS_IMAGE));
            setKeyPairs(keyPairs);
        } catch (e) {
            clearPreviousResourceData(awsDispatch, RESOURCE.AWS_ADD_RESOURCES, AWS_FIELDS.EC2_KEY_PAIR);
            setErrorObject(addErrorInfo(errorObject, e, AWS_FIELDS.EC2_KEY_PAIR));
        }
    }

    static async initNodeProfile(props: AwsOrchestratorProps) {
        const { awsDispatch, setErrorObject, errorObject } = props;
        try {
            const nodeInstance = await AwsService.getAwsNodeTypes();
            const nodeProfileList: { [key: string]: string } = {
                [AWS_NODE_PROFILE_NAMES.SINGLE_NODE]: '',
                [AWS_NODE_PROFILE_NAMES.HIGH_AVAILABILITY]: '',
                [AWS_NODE_PROFILE_NAMES.PRODUCTION_READY]: '',
            };
            Object.keys(nodeProfileList).map((nodeProfile) => {
                nodeProfileList[nodeProfile] = AwsDefaults.setDefaultNodeType(nodeInstance, nodeProfile);
            });
            saveCurrentResourceData(awsDispatch, RESOURCE.AWS_ADD_RESOURCES, AWS_FIELDS.NODE_TYPE, nodeProfileList);
            setDefaultNodeType(awsDispatch, nodeProfileList[nodeInstanceTypes[0].id]);
            setErrorObject(removeErrorInfo(errorObject, AWS_FIELDS.NODE_TYPE));
        } catch (e) {
            clearPreviousResourceData(awsDispatch, RESOURCE.AWS_ADD_RESOURCES, AWS_FIELDS.NODE_TYPE);
            setErrorObject(addErrorInfo(errorObject, e, AWS_FIELDS.NODE_TYPE));
        }
    }
}

function setDefaultOsImage(awsDispatch: StoreDispatch, osImages: AWSVirtualMachine[]) {
    awsDispatch({
        type: INPUT_CHANGE,
        field: AWS_FIELDS.OS_IMAGE,
        payload: AwsDefaults.selectDefalutOsImage(osImages),
    } as FormAction);
}

function setDefaultEC2KeyPair(awsDispatch: StoreDispatch, keyPairs: AWSKeyPair[]) {
    awsDispatch({
        type: INPUT_CHANGE,
        field: AWS_FIELDS.EC2_KEY_PAIR,
        payload: AwsDefaults.selectDefalutEC2KeyPairs(keyPairs),
    } as FormAction);
}

function setDefaultNodeType(awsDispatch: StoreDispatch, nodeProfile: string) {
    awsDispatch({
        type: INPUT_CHANGE,
        field: AWS_FIELDS.NODE_TYPE,
        payload: nodeProfile,
    } as FormAction);
}
