import { AwsService, AWSVirtualMachine } from '../../../../../swagger-api';
import { AWS_FIELDS } from '../../aws-mc-basic/AwsManagementClusterBasic.constants';
import { INPUT_CHANGE } from '../../../../../state-management/actions/Form.actions';
import { AwsDefaults } from '../default-service/AwsDefaults.service';
import { STORE_SECTION_FORM } from '../../../../../state-management/reducers/Form.reducer';
import { AWS_ADD_RESOURCES } from '../../../../../state-management/actions/Resources.actions';
import { AwsResourceAction, FormAction, StoreDispatch } from '../../../../../shared/types/types';
interface AwsOrchestratorProps {
    awsState: { [key: string]: any };
    awsDispatch: StoreDispatch;
}
export class AwsOrchestrator {
    static async initOsImages(props: AwsOrchestratorProps) {
        clearOsImages(props);
        const osImages = await AwsService.getAwsosImages(props.awsState[STORE_SECTION_FORM].REGION);
        saveOsImages(props, osImages);
        setDefaultOsImage(props, osImages);
    }
}

function clearOsImages(props: AwsOrchestratorProps) {
    props.awsDispatch({
        type: AWS_ADD_RESOURCES,
        resourceName: AWS_FIELDS.OS_IMAGE,
        payload: [],
    } as AwsResourceAction);
}

function saveOsImages(props: AwsOrchestratorProps, osImages: AWSVirtualMachine[]) {
    props.awsDispatch({
        type: AWS_ADD_RESOURCES,
        resourceName: AWS_FIELDS.OS_IMAGE,
        payload: osImages,
    } as AwsResourceAction);
}

function setDefaultOsImage(props: AwsOrchestratorProps, osImages: AWSVirtualMachine[]) {
    props.awsDispatch({
        type: INPUT_CHANGE,
        field: AWS_FIELDS.OS_IMAGE,
        payload: AwsDefaults.selectDefalutOsImage(osImages),
    } as FormAction);
}
