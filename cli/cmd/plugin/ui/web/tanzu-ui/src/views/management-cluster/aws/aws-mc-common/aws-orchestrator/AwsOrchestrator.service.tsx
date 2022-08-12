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
    errorMessage: { [key: string]: any };
    setErrorMessage: (newErrorMessage: { [key: string]: any }) => void;
}
export class AwsOrchestrator {
    static async initOsImages(props: AwsOrchestratorProps) {
        const { awsState, setErrorMessage, errorMessage } = props;
        try {
            clearPreviousDiscoveryData(props, AWS_FIELDS.OS_IMAGE);
            const osImages = await AwsService.getAwsosImages(awsState[STORE_SECTION_FORM].REGION);
            saveCurrentDiscoveryData(props, AWS_FIELDS.OS_IMAGE, osImages);
            setDefaultOsImage(props, osImages);
            setErrorMessage(removeErrorInfo(errorMessage, AWS_FIELDS.OS_IMAGE));
            // throw '404';
        } catch (e) {
            setErrorMessage(addErrorInfo(errorMessage, e, AWS_FIELDS.OS_IMAGE));
        }
    }

    static async initEC2KeyPairs(props: AwsOrchestratorProps, setKeyPairs: (keyPairs: AWSKeyPair[]) => void) {
        const { awsState, setErrorMessage, errorMessage } = props;
        try {
            clearPreviousDiscoveryData(props, AWS_FIELDS.EC2_KEY_PAIR);
            const keyPairs = await AwsService.getAwsKeyPairs();
            saveCurrentDiscoveryData(props, AWS_FIELDS.EC2_KEY_PAIR, keyPairs);
            setDefaultEC2KeyPair(props, keyPairs);
            setErrorMessage(removeErrorInfo(errorMessage, AWS_FIELDS.OS_IMAGE));
            setKeyPairs(keyPairs);
            // throw '40411';
        } catch (e) {
            setErrorMessage(addErrorInfo(errorMessage, e, AWS_FIELDS.EC2_KEY_PAIR));
        }
    }
}

function clearPreviousDiscoveryData(props: AwsOrchestratorProps, resourceName: AWS_FIELDS) {
    props.awsDispatch({
        type: AWS_ADD_RESOURCES,
        resourceName: resourceName,
        payload: [],
    } as AwsResourceAction);
}

function saveCurrentDiscoveryData(props: AwsOrchestratorProps, resourceName: AWS_FIELDS, currentValues: any[]) {
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

function removeErrorInfo(errorMessage: { [key: string]: any }, field: AWS_FIELDS) {
    const copy = { ...errorMessage };
    delete copy[field];
    return copy;
}

function addErrorInfo(errorMessage: { [key: string]: any }, error: any, field: AWS_FIELDS) {
    return {
        ...errorMessage,
        [field]: error,
    };
}
