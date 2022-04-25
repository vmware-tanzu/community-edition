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

ClarityIcons.addIcons(clusterIcon, cloudScaleIcon, applicationsIcon);

export interface RolloverConfigItem {
    logo: string,
    icon: string,
    backgroundImage?: string,
    title: string,
    description: string,
    actionText?: string,
    actionUrl?: string
}

const RolloverBanner = () => {

    const [currentBannerItem, setCurrentBannerItem] = useState<number>(0);

    // Rollover banner config entries populate a line item and corresponding detail for display
    const rolloverBannerConfig:Array<RolloverConfigItem> = [
        {
            logo: TceLogo,
            icon: '',
            // backgroundImage: TceExperienceBg,
            title: 'Community-supported experience',
            description: 'Tanzu Community Edition is an open source distribution of Tanzu that can be installed and' +
                'configured in minutes on your local workstation.',
            actionText: 'Visit the project on Github',
            actionUrl: 'https://github.com/vmware-tanzu/community-edition',
        },
        {
            logo: '',
            icon: 'cluster',
            title: 'How do I decide which type of cluster to create?',
            description: 'There are two different types of Tanzu clusters that can be deployed; managed and unmanaged clusters.',
            actionText: 'Learn more about types of clusters',
            actionUrl: '',
        },
        {
            logo: '',
            icon: 'cloud-scale',
            title: 'How do I manage my clusters at scale?',
            description: 'VMware Tanzu Mission Control is a is a centralized hub for simplified, multi-cloud, ' +
                'multi-cluster Kubernetes management.',
            actionText: 'Learn more about Tanzu Mission Control',
            actionUrl: '',
        },
        {
            logo: '',
            icon: 'applications',
            // backgroundImage: TapExperience,
            title: 'How do I manage my application platform?',
            description: 'VMware Tanzu is a complete portfolio of products and services enabling developers and ' +
                'operators to run and manage Kubernetes across multiple cloud providers.',
            actionText: 'Learn more about Tanzu Application Platform',
            actionUrl: 'https://tanzu.vmware.com/application-platform',
        }
    ];

    const setCurrentBannerCallback = (id: number) => {
        setCurrentBannerItem(id);
    };

    return (
        <>
            <div cds-layout="col:5">
                <div cds-layout="vertical gap:md">
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
                                selected={currentBannerItem===index}/>
                        ))
                    }
                </div>
            </div>
            <div className="banner-content" cds-layout="col:7" style={{
                backgroundImage: `url(${rolloverBannerConfig[currentBannerItem].backgroundImage})`,
                backgroundPosition: 'right',
                backgroundRepeat: 'no-repeat'
            }}>
                <div cds-text="h3" className="banner-content-title text-blue">
                    {rolloverBannerConfig[currentBannerItem].title}
                </div>
                <div className="banner-content-description">
                    {rolloverBannerConfig[currentBannerItem].description}
                </div>
                {rolloverBannerConfig[currentBannerItem].actionText && rolloverBannerConfig[currentBannerItem].actionUrl &&
                    <div className="banner-content-action text-blue" onClick={() => {
                        window.open(rolloverBannerConfig[currentBannerItem].actionUrl, '_blank');
                    }}>
                        {rolloverBannerConfig[currentBannerItem].actionText}
                    </div>
                }

            </div>
        </>
    );
};

export default RolloverBanner;
