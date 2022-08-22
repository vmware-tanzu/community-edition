// App imports
import { AwsDefaults } from '../default-service/AwsDefaults.service';
import { AwsService, AWSVirtualMachine } from '../../../../../swagger-api';
import { AWSKeyPair } from '../../../../../swagger-api/models/AWSKeyPair';
import { AWS_FIELDS } from '../../aws-mc-basic/AwsManagementClusterBasic.constants';
import { DefaultOrchestrator } from '../../../default-orchestrator/DefaultOrchestrator';
import { StoreDispatch } from '../../../../../shared/types/types';
import { STORE_SECTION_FORM } from '../../../../../state-management/reducers/Form.reducer';

interface AwsOrchestratorProps {
    awsState: { [key: string]: any };
    awsDispatch: StoreDispatch;
    errorObject: { [key: string]: any };
    setErrorObject: (newErrorObject: { [key: string]: any }) => void;
}

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
}
