// React imports
import React, { useContext } from 'react';

// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsControlMessage } from '@cds/react/forms';
import { CdsIcon } from '@cds/react/icon';
import { CdsSelect } from '@cds/react/select';
import { yupResolver } from '@hookform/resolvers/yup';
import { SubmitHandler, useForm } from 'react-hook-form';
import * as yup from 'yup';

// App imports
import { FormAction } from '../../../../shared/types/types';
import { INPUT_CHANGE } from '../../../../state-management/actions/Form.actions';
import { SelectionType, TreeSelectItem } from '../../../../shared/components/TreeSelect/TreeSelect.interface';
import { StepProps } from '../../../../shared/components/wizard/Wizard';
import { STORE_SECTION_FORM } from '../../../../state-management/reducers/Form.reducer';
import TreeSelect from '../../../../shared/components/TreeSelect/TreeSelect';
import useVSphereComputeResources from '../../../../shared/hooks/VSphere/UseVSphereComputeResources';
import UseUpdateTabStatus from '../../../../shared/components/wizard/UseUpdateTabStatus.hooks';
import useVSphereDatastores from '../../../../shared/hooks/VSphere/UseVSphereDatastores';
import useVSphereFolders from '../../../../shared/hooks/VSphere/UseVSphereFolders';
import useVSphereNetworkNames from '../../../../shared/hooks/VSphere/UseVSphereNetworkNames';
import { VSphereDatastore, VSphereFolder, VSphereManagementObject, VSphereNetwork } from '../../../../swagger-api';
import { VSPHERE_FIELDS } from '../VsphereManagementCluster.constants';
import { VsphereStore } from '../Store.vsphere.mc';

export interface VSphereClusterResourcesStepInputs {
    [VSPHERE_FIELDS.VMFolder]: string;
    [VSPHERE_FIELDS.DataStore]: string;
    [VSPHERE_FIELDS.VSphereNetworkName]: string;
    [VSPHERE_FIELDS.Pool]: string;
}

const schema = yup
    .object()
    .shape({
        [VSPHERE_FIELDS.VMFolder]: yup.string().required('Please select a VM folder.'),
        [VSPHERE_FIELDS.DataStore]: yup.string().required('Please select a Datastore.'),
        [VSPHERE_FIELDS.VSphereNetworkName]: yup.string().required('Please select a vSphere network name.'),
        [VSPHERE_FIELDS.Pool]: yup.string().required('Please select a resource pool.'),
    })
    .required();

const treeDataMapper = (inputData: VSphereManagementObject[]) => {
    const parsePath = (node: VSphereManagementObject) => {
        const hostOrCluster = [VSphereManagementObject.resourceType.HOST, VSphereManagementObject.resourceType.CLUSTER];
        if (node.resourceType) {
            return hostOrCluster.includes(node.resourceType) ? `${node.path}/Resources` : node.path;
        }
        return node.path;
    };

    const constructTree = (treeNodes: Array<TreeSelectItem>, map: Map<string, Array<TreeSelectItem>>): void => {
        if (!treeNodes || treeNodes.length <= 0) {
            return;
        }

        treeNodes.forEach((node) => {
            const childNodes = map.get(node.id) || [];
            node.children = childNodes;
            constructTree(childNodes, map);
        });
    };

    const removeDatacenter = (resourceTree: TreeSelectItem[]): TreeSelectItem[] => {
        let rootNodes: TreeSelectItem[] = [];
        resourceTree.forEach((resource) => {
            if (resource.type === VSphereManagementObject.resourceType.DATACENTER) {
                if (resource?.children?.length) {
                    rootNodes = [...rootNodes, ...resource.children];
                }
            } else {
                rootNodes.push(resource);
            }
        });
        return rootNodes;
    };

    const treeData: TreeSelectItem[] = inputData.map((r) => ({
        id: r.moid ?? '',
        label: r.name ?? '',
        value: parsePath(r) ?? '',
        type: r.resourceType ?? '',
        parentId: r.parentMoid ?? '',
    }));

    const nodeMap: Map<string, TreeSelectItem[]> = new Map();
    const resourceTree: TreeSelectItem[] = [];

    treeData.forEach((resource) => {
        const parentId = resource.parentId;
        if (parentId) {
            const a: TreeSelectItem[] = nodeMap.get(parentId) as TreeSelectItem[];
            nodeMap.set(parentId, a ? [...a, resource] : [resource]);
        } else {
            resourceTree.push(resource); // it contains root nodes
        }
    });
    constructTree(resourceTree, nodeMap);
    return removeDatacenter(resourceTree);
};

export function VsphereClusterResourcesStep(props: Partial<StepProps>) {
    const { vsphereState, vsphereDispatch } = useContext(VsphereStore);
    const datacenter = vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.DATACENTER];

    const { currentStep, deploy, updateTabStatus, goToStep } = props;

    const methods = useForm<VSphereClusterResourcesStepInputs>({
        resolver: yupResolver(schema),
        mode: 'all',
    });

    const {
        register,
        formState: { errors },
        handleSubmit,
        control,
    } = methods;

    // update tab status bar
    if (updateTabStatus) {
        UseUpdateTabStatus(errors, currentStep, updateTabStatus);
    }

    const { data: vSphereFolders } = useVSphereFolders(datacenter);
    const { data: vSphereDatastores } = useVSphereDatastores(datacenter);
    const { data: vSphereNetworkNames } = useVSphereNetworkNames(datacenter);
    const { data: vSphereComputeResources } = useVSphereComputeResources(datacenter);
    const mappedVSphereComputeResources = treeDataMapper(vSphereComputeResources);

    const canContinue = () => {
        return Object.keys(errors).length === 0;
    };

    const onSubmit: SubmitHandler<VSphereClusterResourcesStepInputs> = (data) => {
        goToStep && goToStep((currentStep ?? 0) + 1);
    };

    const onChange = (field: any, value: string) => {
        vsphereDispatch({
            type: INPUT_CHANGE,
            field,
            payload: value,
        } as FormAction);
    };

    return (
        <div className="cluster-settings-container" cds-layout="m:lg">
            <h3>vSphere Cluster Resources</h3>
            <div cds-layout="vertical gap:md align:stretch">
                <p>The following are settings for VMWare vSphere that could be changed from their default values or enabled.</p>

                <h3>Resources</h3>
                <div cds-layout="horizontal gap:md align:stretch">
                    <VMFolder register={register} onChange={onChange} errors={errors} vmFolders={vSphereFolders} />
                    <Datastore register={register} onChange={onChange} errors={errors} datastores={vSphereDatastores} />
                    <VSphereNetworkName register={register} onChange={onChange} errors={errors} vSphereNetworkNames={vSphereNetworkNames} />
                </div>

                <div cds-layout="m-y:md">
                    <h4 cds-layout="m-y:none"> Clusters, hosts, and resource pools</h4>
                    <TreeSelect
                        data={mappedVSphereComputeResources}
                        control={control}
                        name={VSPHERE_FIELDS.Pool}
                        selectionType={SelectionType.Single}
                        onChange={onChange}
                    />
                </div>
            </div>
            <CdsButton onClick={handleSubmit(onSubmit)} disabled={!canContinue()}>
                NEXT
            </CdsButton>
        </div>
    );
}

function VMFolder({
    register,
    onChange,
    errors,
    vmFolders,
}: {
    register: any;
    onChange: (field: string, value: string) => void;
    errors: any;
    vmFolders: VSphereFolder[];
}) {
    return (
        <div>
            <CdsSelect layout="vertical" controlWidth="shrink">
                <label>VM folder</label>
                <select {...register(VSPHERE_FIELDS.VMFolder)} onChange={(e) => onChange(VSPHERE_FIELDS.VMFolder, e?.target?.value ?? '')}>
                    <option />
                    {vmFolders.map((details: any) => (
                        <option key={details.moid} value={details.moid}>
                            {details.name}
                        </option>
                    ))}
                </select>
                {errors[VSPHERE_FIELDS.VMFolder] && (
                    <CdsControlMessage status="error">{errors[VSPHERE_FIELDS.VMFolder].message}</CdsControlMessage>
                )}
            </CdsSelect>
        </div>
    );
}

function Datastore({
    register,
    onChange,
    errors,
    datastores,
}: {
    register: any;
    onChange: (field: string, value: string) => void;
    errors: any;
    datastores: VSphereDatastore[];
}) {
    return (
        <div>
            <CdsSelect layout="vertical" controlWidth="shrink">
                <label>Datastore</label>
                <select
                    {...register(VSPHERE_FIELDS.DataStore, {
                        onChange: (e: any) => onChange(VSPHERE_FIELDS.DataStore, e?.target?.value ?? ''),
                    })}
                >
                    <option />
                    {datastores.map((details: any) => (
                        <option key={details.moid} value={details.moid}>
                            {details.name}
                        </option>
                    ))}
                </select>
                {errors[VSPHERE_FIELDS.DataStore] && (
                    <CdsControlMessage status="error">{errors[VSPHERE_FIELDS.DataStore].message}</CdsControlMessage>
                )}
            </CdsSelect>
        </div>
    );
}

function VSphereNetworkName({
    register,
    onChange,
    errors,
    vSphereNetworkNames,
}: {
    register: any;
    onChange: (field: string, value: string) => void;
    errors: any;
    vSphereNetworkNames: VSphereNetwork[];
}) {
    return (
        <div>
            <CdsSelect layout="vertical" controlWidth="shrink">
                <label>VSphere network name</label>
                <select
                    {...register(VSPHERE_FIELDS.VSphereNetworkName, {
                        onChange: (e: any) => onChange(VSPHERE_FIELDS.VSphereNetworkName, e?.target?.value ?? ''),
                    })}
                >
                    <option />
                    {vSphereNetworkNames.map((details: any) => (
                        <option key={details.moid} value={details.moid}>
                            {details.name}
                        </option>
                    ))}
                </select>
                {errors[VSPHERE_FIELDS.VSphereNetworkName] && (
                    <CdsControlMessage status="error">{errors[VSPHERE_FIELDS.VSphereNetworkName].message}</CdsControlMessage>
                )}
            </CdsSelect>
        </div>
    );
}
