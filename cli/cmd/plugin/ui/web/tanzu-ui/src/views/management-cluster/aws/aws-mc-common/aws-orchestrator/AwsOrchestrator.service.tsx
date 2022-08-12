import { AwsService, AWSVirtualMachine } from '../../../../../swagger-api';
import { AWS_FIELDS } from '../../aws-mc-basic/AwsManagementClusterBasic.constants';
import { INPUT_CHANGE } from '../../../../../state-management/actions/Form.actions';
import { AwsDefaults } from '../default-service/AwsDefaults.service';
import { STORE_SECTION_FORM } from '../../../../../state-management/reducers/Form.reducer';
import { AWS_ADD_RESOURCES } from '../../../../../state-management/actions/Resources.actions';
import { AwsResourceAction, FormAction, StoreDispatch } from '../../../../../shared/types/types';
import { AWSKeyPair } from '../../../../../swagger-api/models/AWSKeyPair';
interface AwsOrchestratorProps {
    awsState: { [key: string]: any };
    awsDispatch: StoreDispatch;
}
export class AwsOrchestrator {
    static async initOsImages(props: AwsOrchestratorProps) {
        clearPreviousValues(props, AWS_FIELDS.OS_IMAGE);
        const osImages = await AwsService.getAwsosImages(props.awsState[STORE_SECTION_FORM].REGION);
        saveCurrentValues(props, AWS_FIELDS.OS_IMAGE, osImages);
        setDefaultOsImage(props, osImages);
    }

    static async initEC2KeyPairs(props: AwsOrchestratorProps) {
        clearPreviousValues(props, AWS_FIELDS.EC2_KEY_PAIR);
        const keyPairs = await AwsService.getAwsKeyPairs();
        saveCurrentValues(props, AWS_FIELDS.EC2_KEY_PAIR, keyPairs);
        setDefaultEC2KeyPair(props, keyPairs);
        return keyPairs;
    }
}

function clearPreviousValues(props: AwsOrchestratorProps, resourceName: AWS_FIELDS) {
    props.awsDispatch({
        type: AWS_ADD_RESOURCES,
        resourceName: resourceName,
        payload: [],
    } as AwsResourceAction);
}

function saveCurrentValues(props: AwsOrchestratorProps, resourceName: AWS_FIELDS, currentValues: any[]) {
    props.awsDispatch({
        type: AWS_ADD_RESOURCES,
        resourceName: resourceName,
        payload: currentValues,
    } as AwsResourceAction);
}

function setDefaultOsImage(props: AwsOrchestratorProps, osImages: AWSVirtualMachine[]) {
    props.awsDispatch({
        type: INPUT_CHANGE,
        field: AWS_FIELDS.OS_IMAGE,
        payload: AwsDefaults.selectDefalutOsImage(osImages),
    } as FormAction);
}

function setDefaultEC2KeyPair(props: AwsOrchestratorProps, keyPairs: AWSKeyPair[]) {
    props.awsDispatch({
        type: INPUT_CHANGE,
        field: AWS_FIELDS.EC2_KEY_PAIR,
        payload: AwsDefaults.selectDefalutEC2KeyPairs(keyPairs),
    } as FormAction);
}
