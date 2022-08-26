// App imports
import { AWS_FIELDS, AWS_NODE_PROFILE_NAMES } from '../../aws-mc-basic/AwsManagementClusterBasic.constants';
import { AwsDefaults } from '../default-service/AwsDefaults.service';
import { AWSKeyPair } from '../../../../../swagger-api/models/AWSKeyPair';
import { AwsService, AWSVirtualMachine, CancelablePromise } from '../../../../../swagger-api';
import { DefaultOrchestrator } from '../../../default-orchestrator/DefaultOrchestrator';
import { NodeProfileType } from '../../../../../shared/components/FormInputComponents/NodeProfile/NodeProfile';
import { StoreDispatch } from '../../../../../shared/types/types';
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

    static async initNodeProfile(props: AwsOrchestratorProps) {
        const { awsDispatch, setErrorObject, errorObject } = props;
        DefaultOrchestrator.initResources<{ [key: string]: { [key: string]: any } }>({
            resourceName: AWS_FIELDS.NODE_TYPE,
            dispatch: awsDispatch,
            errorObject,
            setErrorObject,
            fetcher: () => {
                return new CancelablePromise(async (res, rej) => {
                    const defaultNodeProfileValues: { [key: string]: { [key: string]: any } } = {
                        [AWS_NODE_PROFILE_NAMES.SINGLE_NODE]: {},
                        [AWS_NODE_PROFILE_NAMES.HIGH_AVAILABILITY]: {},
                        [AWS_NODE_PROFILE_NAMES.PRODUCTION_READY]: {},
                    };
                    const keyList = Object.keys(defaultNodeProfileValues);
                    for (const nodeProfile of keyList) {
                        defaultNodeProfileValues[nodeProfile] = {
                            nodeType: await AwsDefaults.setDefaultNodeType(nodeProfile).catch((e) => {
                                rej(e);
                                return '';
                            }),
                            azs: await AwsDefaults.setDefaultAvailabilityZones(nodeProfile).catch((e) => {
                                rej(e);
                                return '';
                            }),
                        };
                    }
                    res([defaultNodeProfileValues]);
                });
            },
        });
    }
}
