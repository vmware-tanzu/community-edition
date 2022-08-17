// App imports
import { AwsDefaults } from '../default-service/AwsDefaults.service';
import { AwsService, AWSVirtualMachine } from '../../../../../swagger-api';
import { AWSKeyPair } from '../../../../../swagger-api/models/AWSKeyPair';
import { AWS_FIELDS } from '../../aws-mc-basic/AwsManagementClusterBasic.constants';
import { INPUT_CHANGE } from '../../../../../state-management/actions/Form.actions';
import { FormAction, StoreDispatch } from '../../../../../shared/types/types';
import { STORE_SECTION_FORM } from '../../../../../state-management/reducers/Form.reducer';
import { RESOURCE } from '../../../../../state-management/actions/Resources.actions';
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
