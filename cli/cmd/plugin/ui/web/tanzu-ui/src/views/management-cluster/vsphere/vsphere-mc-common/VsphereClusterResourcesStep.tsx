// React imports
import React, { ChangeEvent, useContext, useEffect, useState } from 'react';

// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsControlMessage } from '@cds/react/forms';
import { CdsSelect } from '@cds/react/select';
import { SubmitHandler, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';

// App imports
import { FormAction, StoreDispatch } from '../../../../shared/types/types';
import { getResource } from '../../../../state-management/reducers/Resources.reducer';
import { initDatastores, initFolders, initNetworks, initResources } from './VsphereOrchestrator.service';
import { INPUT_CHANGE } from '../../../../state-management/actions/Form.actions';
import { SelectionType, TreeSelectItem } from '../../../../shared/components/TreeSelect/TreeSelect.interface';
import { StepProps } from '../../../../shared/components/wizard/Wizard';
import { STORE_SECTION_FORM } from '../../../../state-management/reducers/Form.reducer';
import TreeSelect from '../../../../shared/components/TreeSelect/TreeSelect';
import UseUpdateTabStatus from '../../../../shared/components/wizard/UseUpdateTabStatus.hooks';
import { VSphereDatastore, VSphereFolder, VSphereManagementObject, VSphereNetwork } from '../../../../swagger-api';
import { VSPHERE_FIELDS } from '../VsphereManagementCluster.constants';
import { VsphereStore } from '../Store.vsphere.mc';

export interface VSphereClusterResourcesStepInputs {
    [VSPHERE_FIELDS.VMFolder]: string;
    [VSPHERE_FIELDS.DataStore]: string;
    [VSPHERE_FIELDS.Network]: string;
    [VSPHERE_FIELDS.Pool]: string;
}

const schema = yup
    .object()
    .shape({
        [VSPHERE_FIELDS.VMFolder]: yup.string().required('Please select a VM folder.'),
        [VSPHERE_FIELDS.DataStore]: yup.string().required('Please select a Datastore.'),
        [VSPHERE_FIELDS.Network]: yup.string().required('Please select a vSphere network name.'),
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

    const treeData: TreeSelectItem[] = inputData?.map((r) => ({
        id: r.moid ?? '',
        label: r.name ?? '',
        value: parsePath(r) ?? '',
        type: r.resourceType ?? '',
        parentId: r.parentMoid ?? '',
    }));

    const nodeMap: Map<string, TreeSelectItem[]> = new Map();
    const resourceTree: TreeSelectItem[] = [];

    treeData?.forEach((resource) => {
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
    const [errorObject, setErrorObject] = useState({});
    const datacenter = vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.DATACENTER];

    const { currentStep, updateTabStatus, goToStep, submitForm } = props;

    const methods = useForm<VSphereClusterResourcesStepInputs>({
        resolver: yupResolver(schema),
        mode: 'all',
    });

    const {
        register,
        formState: { errors },
        handleSubmit,
        control,
        setValue,
    } = methods;

    if (updateTabStatus) {
        UseUpdateTabStatus(errors, currentStep, updateTabStatus);
    }

    useEffect(() => {
        initFolders(errorObject, setErrorObject, vsphereDispatch, datacenter).then(() =>
            setValue(VSPHERE_FIELDS.VMFolder, vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.VMFolder]?.name)
        );
        initDatastores(errorObject, setErrorObject, vsphereDispatch, datacenter).then(() =>
            setValue(VSPHERE_FIELDS.DataStore, vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.DataStore]?.name)
        );
        initNetworks(errorObject, setErrorObject, vsphereDispatch, datacenter).then(() =>
            setValue(VSPHERE_FIELDS.Network, vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.Network]?.name)
        );
        initResources(errorObject, setErrorObject, vsphereDispatch, datacenter).then(() =>
            setValue(VSPHERE_FIELDS.Pool, vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.Pool])
        );
    }, [vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.DATACENTER]]);

    const vSphereFolders = getResource<VSphereFolder[]>(VSPHERE_FIELDS.VMFolder, vsphereState) || [];
    const vSphereDatastores = getResource<VSphereDatastore[]>(VSPHERE_FIELDS.DataStore, vsphereState) || [];
    const vSphereNetworks = getResource<VSphereNetwork[]>(VSPHERE_FIELDS.Network, vsphereState) || [];
    const vSphereResourcePools = getResource<VSphereManagementObject[]>(VSPHERE_FIELDS.Pool, vsphereState) || [];
    const mappedVSphereComputeResources = treeDataMapper(vSphereResourcePools);

    const canContinue = () => {
        return Object.keys(errors).length === 0;
    };

    const onSubmit: SubmitHandler<VSphereClusterResourcesStepInputs> = (data) => {
        goToStep && goToStep((currentStep ?? 0) + 1);
        submitForm && submitForm(currentStep);
    };

    const onSelectPool = (fieldName: string, selectedId: string) =>
        vsphereDispatch({
            type: INPUT_CHANGE,
            field: fieldName,
            payload: selectedId,
        } as FormAction);

    const onSelectVmFolder = fxnOnSelectArrayObject<VSphereFolder>(
        VSPHERE_FIELDS.VMFolder,
        vsphereDispatch,
        vSphereFolders,
        (obj, selectedValue) => obj.moid === selectedValue
    );

    const onSelectDatastore = fxnOnSelectArrayObject<VSphereDatastore>(
        VSPHERE_FIELDS.DataStore,
        vsphereDispatch,
        vSphereDatastores,
        (datastore, selectedValue) => datastore.moid === selectedValue
    );

    const onSelectNetwork = fxnOnSelectArrayObject<VSphereNetwork>(
        VSPHERE_FIELDS.Network,
        vsphereDispatch,
        vSphereNetworks,
        (network, selectedValue) => network.moid === selectedValue
    );

    return (
        <div className="wizard-content-container" cds-layout="m:lg">
            <h3 cds-layout="m-t:md m-b:xl" cds-text="title">
                vSphere Cluster Resources
            </h3>
            <div cds-layout="vertical gap:md align:stretch">
                <p>The following are settings for VMWare vSphere that could be changed from their default values or enabled.</p>

                <h3>Resources</h3>
                <div cds-layout="horizontal gap:md align:stretch">
                    <VMFolder
                        register={register}
                        onSelect={onSelectVmFolder}
                        errors={errors}
                        vmFolders={vSphereFolders}
                        selected={vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.VMFolder]}
                    />
                    <Datastore
                        register={register}
                        onSelect={onSelectDatastore}
                        errors={errors}
                        datastores={vSphereDatastores}
                        selected={vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.DataStore]}
                    />
                    <VSphereNetworkSelect
                        register={register}
                        onSelect={onSelectNetwork}
                        errors={errors}
                        vSphereNetworks={vSphereNetworks}
                        selected={vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.Network]}
                    />
                </div>

                <div cds-layout="m-y:md">
                    <h4 cds-layout="m-y:none"> Clusters, hosts, and resource pools</h4>
                    <TreeSelect
                        data={mappedVSphereComputeResources}
                        control={control}
                        name={VSPHERE_FIELDS.Pool}
                        selectionType={SelectionType.Single}
                        onChange={onSelectPool}
                        selectedValue={vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.Pool]}
                    />
                </div>
            </div>
            <CdsButton onClick={handleSubmit(onSubmit)} disabled={!canContinue()}>
                NEXT
            </CdsButton>
        </div>
    );
}

interface VMFolderParams {
    register: any;
    onSelect: (e: ChangeEvent<HTMLSelectElement>) => void;
    errors: any;
    vmFolders: VSphereFolder[];
    selected?: VSphereFolder;
}
function VMFolder({ register, onSelect, errors, vmFolders, selected }: VMFolderParams) {
    return (
        <div>
            <CdsSelect layout="vertical" controlWidth="shrink">
                <label>VM folder</label>
                <select {...register(VSPHERE_FIELDS.VMFolder)} onChange={onSelect} value={selected?.moid}>
                    <option />
                    {vmFolders?.map((folder: any) => (
                        <option key={folder.moid} value={folder.moid}>
                            {folder.name}
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

interface DatastoreParams {
    register: any;
    onSelect: (e: ChangeEvent<HTMLSelectElement>) => void;
    errors: any;
    datastores: VSphereDatastore[];
    selected?: VSphereDatastore;
}
function Datastore({ register, onSelect, errors, datastores, selected }: DatastoreParams) {
    return (
        <div>
            <CdsSelect layout="vertical" controlWidth="shrink">
                <label>Datastore</label>
                <select {...register(VSPHERE_FIELDS.DataStore)} onChange={onSelect} value={selected?.moid}>
                    <option />
                    {datastores?.map((details: any) => (
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

interface VSphereNetworkSelectParams {
    register: any;
    onSelect: (e: ChangeEvent<HTMLSelectElement>) => void;
    errors: any;
    vSphereNetworks: VSphereNetwork[];
    selected?: VSphereNetwork;
}
function VSphereNetworkSelect({ register, onSelect, errors, vSphereNetworks, selected }: VSphereNetworkSelectParams) {
    return (
        <div>
            <CdsSelect layout="vertical" controlWidth="shrink">
                <label>VSphere network name</label>
                <select {...register(VSPHERE_FIELDS.Network)} onChange={onSelect} value={selected?.moid}>
                    <option />
                    {vSphereNetworks?.map((details: any) => (
                        <option key={details.moid} value={details.moid}>
                            {details.name}
                        </option>
                    ))}
                </select>
                {errors[VSPHERE_FIELDS.Network] && (
                    <CdsControlMessage status="error">{errors[VSPHERE_FIELDS.Network].message}</CdsControlMessage>
                )}
            </CdsSelect>
        </div>
    );
}

// TODO: move this function to a utility class for reuse
// returns a function that can be used as an event handler for a select box,
// which will find the target object (after the user selects a value) and then dispatch an event which stores the corresponding OBJECT
// into the data store
function fxnOnSelectArrayObject<OBJ>(
    fieldName: string,
    dispatch: StoreDispatch,
    source: OBJ[],
    matcher: (obj: OBJ, selectedId: string) => boolean
) {
    return (event: ChangeEvent<HTMLSelectElement>) => {
        const selectedId = event.target.value;
        recordSelectedArrayObject<OBJ>(selectedId, fieldName, dispatch, source, matcher);
    };
}

function recordSelectedArrayObject<OBJ>(
    selectedId: string,
    fieldName: string,
    dispatch: StoreDispatch,
    source: OBJ[],
    matcher: (obj: OBJ, selectedId: string) => boolean
) {
    const selectedObjectIndex = source.findIndex((obj) => matcher(obj, selectedId));
    const selectedObject = selectedObjectIndex >= 0 ? source[selectedObjectIndex] : undefined;
    if (selectedId && !selectedObject) {
        console.error(`recordSelectedArrayObject() is unable to find selected id "${selectedId}" in array ${JSON.stringify(source)}`);
    }
    dispatch({
        type: INPUT_CHANGE,
        field: fieldName,
        payload: selectedObject,
    } as FormAction);
}
