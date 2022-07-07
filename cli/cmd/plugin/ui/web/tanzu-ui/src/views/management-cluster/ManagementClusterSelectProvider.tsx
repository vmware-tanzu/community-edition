// React imports
import React from 'react';
import { useNavigate } from 'react-router-dom';

// Library imports
import { CdsCard } from '@cds/react/card';
import { ClarityIcons, cloudIcon } from '@cds/core/icon';

// App imports
import './ManagementClusterSelectProvider.scss';

ClarityIcons.addIcons(cloudIcon);

const ManagementClusterSelectProvider: React.FC = () => {
    const navigate = useNavigate();
    const cards = [
        {
            name: 'Docker',
            path: '/docker',
            imgClass: 'docker-logo-sm',
        },
        {
            name: 'Amazon Web Services',
            path: '/aws',
            imgClass: 'aws-logo-sm',
        },
        {
            name: 'VMware vSphere',
            path: '/vsphere',
            imgClass: 'vsphere-logo-sm',
        },
        {
            name: 'Microsoft Azure',
            path: '/azure',
            imgClass: 'azure-logo-sm',
        },
    ];
    return (
        <div className="management-cluster-landing-container" cds-layout="vertical gap:md col@sm:12 grid">
            <div cds-layout="vertical col:8 gap:lg">
                <div cds-text="title" cds-layout="m-t:md">
                    Management Cluster
                </div>
                <div cds-text="body" cds-layout="p-y:md">
                    Managed Clusters is a deployment model that features one management cluster and multiple workload clusters. The
                    management cluster provides management and operations for Tanzu. It runs Cluster-API which is used to manage workload
                    clusters and multi-cluster services. The workload clusters are where developer’s workloads run.
                </div>
                <div cds-text="section">Select a supported cloud provider</div>
                <div cds-layout="grid cols@md:6 cols@lg:6 gap:lg">
                    {cards.map((card, index) => {
                        return (
                            <CdsCard
                                className="card-container"
                                aria-labelledby="containerOfCards1"
                                tabIndex={0}
                                key={index}
                                cds-layout="vertical"
                                onClick={() => navigate(card.path)}
                                onKeyPress={(e) => (e.key === 'Enter' ? navigate(card.path) : '')}
                            >
                                <div className={card.imgClass} cds-layout="vertical p:xl p-b:lg">
                                    <div cds-layout="align:horizontal-center p-t:xl m-t:md" className="logo-name" cds-text="section">
                                        {card.name}
                                    </div>
                                </div>
                            </CdsCard>
                        );
                    })}
                </div>
            </div>
            <div cds-layout="col:4" className="mgmt-cluster-admins-img"></div>
        </div>
    );
};

export default ManagementClusterSelectProvider;
