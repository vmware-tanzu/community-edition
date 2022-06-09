// React imports
import React, { useState } from 'react';

// Library imports
import { ClarityIcons, clusterIcon, cloudScaleIcon, applicationsIcon } from '@cds/core/icon';

// App imports
import './RolloverBanner.scss';
import RolloverBannerItem from './RolloverBannerItem/RolloverBannerItem';
import TceLogo from '../../../assets/tce-logo.svg';
import TceExperienceBg from '../../../assets/tce-experience-bg.svg';
import TapExperience from '../../../assets/tap-experience-bg.svg';
import TmcExperience from '../../../assets/tmc-experience-bg.svg';
import ClusterExperience from '../../../assets/cluster-experience-bg.svg';

ClarityIcons.addIcons(clusterIcon, cloudScaleIcon, applicationsIcon);

export interface RolloverConfigItem {
    logo: string;
    icon: string;
    backgroundImage?: string;
    title: string;
    description: string;
    actionText?: string;
    actionUrl?: string;
}

const RolloverBanner = () => {
    const [currentBannerItem, setCurrentBannerItem] = useState<number>(0);

    // Rollover banner config entries populate a line item and corresponding detail for display
    const rolloverBannerConfig: Array<RolloverConfigItem> = [
        {
            logo: TceLogo,
            icon: '',
            backgroundImage: TceExperienceBg,
            title: 'Community-supported experience',
            description:
                'Tanzu Community Edition is an open source distribution of Tanzu that can be installed and ' +
                'configured in minutes on your local workstation.',
            actionText: 'Visit the project on GitHub',
            actionUrl: 'https://github.com/vmware-tanzu/community-edition',
        },
        {
            logo: '',
            icon: 'cluster',
            backgroundImage: ClusterExperience,
            title: 'How do I decide which type of cluster to create?',
            description:
                'There are two different types of Tanzu clusters that can be deployed; managed and unmanaged clusters' +
                'Managed clusters are for producton-ready enviroments that features a Management Cluster and Workload Clusters' +
                'Unmanaged Clusters offer Tanzu environments for development and experimentation.',
            actionText: 'Learn more about types of clusters',
            actionUrl: 'https://tanzucommunityedition.io/docs/main/architecture/#tanzu-clusters',
        },
        {
            logo: '',
            icon: 'cloud-scale',
            backgroundImage: TmcExperience,
            title: 'How do I manage my clusters at scale?',
            description:
                'Easily manage clusters on multiple-clouds with Tanzu Mission Control, a central hub for operators. ' +
                'Bringing consistancy to your platform by connecting your clusters.',
            actionText: 'Learn more about Tanzu Mission Control',
            actionUrl: 'https://tanzu.vmware.com/mission-control',
        },
        {
            logo: '',
            icon: 'applications',
            backgroundImage: TapExperience,
            title: 'How do I manage my application platform?',
            description:
                'VMware Tanzu is a complete portfolio of products and services enabling developers and ' +
                'operators to run and manage Kubernetes across multiple cloud providers.',
            actionText: 'Learn more about Tanzu Application Platform',
            actionUrl: 'https://tanzu.vmware.com/application-platform',
        },
    ];

    const setCurrentBannerCallback = (id: number) => {
        setCurrentBannerItem(id);
    };

    return (
        <>
            <div cds-layout="col:5">
                <div cds-layout="vertical">
                    {rolloverBannerConfig.length &&
                        rolloverBannerConfig.map(({ logo, icon, title }, index) => (
                            <RolloverBannerItem
                                onMouseEnter={() => {
                                    setCurrentBannerItem(index);
                                }}
                                key={index}
                                index={index}
                                logo={logo}
                                icon={icon}
                                title={title}
                                mouseEnterCallback={setCurrentBannerCallback}
                                selected={currentBannerItem === index}
                            />
                        ))}
                </div>
            </div>
            <div
                className="banner-content"
                cds-layout="col:7 p:md"
                style={{
                    backgroundImage: `url(${rolloverBannerConfig[currentBannerItem].backgroundImage})`,
                }}
            >
                <div cds-text="title" className="banner-content-title text-aqua">
                    {rolloverBannerConfig[currentBannerItem].title}
                </div>
                <div cds-layout="m-y:sm" className="banner-content-description">
                    {rolloverBannerConfig[currentBannerItem].description}
                </div>
                {rolloverBannerConfig[currentBannerItem].actionText && rolloverBannerConfig[currentBannerItem].actionUrl && (
                    <div
                        cds-text="semibold"
                        className="banner-content-action text-aqua inline"
                        onClick={() => {
                            window.open(rolloverBannerConfig[currentBannerItem].actionUrl, '_blank');
                        }}
                    >
                        {rolloverBannerConfig[currentBannerItem].actionText}
                    </div>
                )}
            </div>
        </>
    );
};

export default RolloverBanner;
