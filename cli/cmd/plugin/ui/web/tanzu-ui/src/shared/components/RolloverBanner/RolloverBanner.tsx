// React imports
import React, { useState } from 'react';

// Library imports
import { ClarityIcons, clusterIcon, cloudScaleIcon, applicationsIcon } from '@cds/core/icon';
import styled from 'styled-components';

// App imports
import './RolloverBanner.scss';
import RolloverBannerItem from './RolloverBannerItem/RolloverBannerItem';
import TanzuLogo from '../../../assets/tanzu-logo.svg';

ClarityIcons.addIcons(clusterIcon, cloudScaleIcon, applicationsIcon);

const RolloverList = styled.div`
    width: 100%;
`;

export interface RolloverConfigItem {
    logo: string,
    icon: string,
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
            logo: TanzuLogo,
            icon: '',
            title: 'What is VMware Tanzu?',
            description: 'VMware Tanzu is a complete portfolio of products and services enabling developers and ' +
                'operators to run and manage Kubernetes across multiple cloud providers.',
            actionText: 'Visit the project on Github',
            actionUrl: '',
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
            title: 'How do I manage my application platform?',
            description: 'VMware Tanzu is a complete portfolio of products and services enabling developers and ' +
                'operators to run and manage Kubernetes across multiple cloud providers.',
            actionText: 'Learn more about Tanzu Application Platform',
            actionUrl: '',
        }
    ];

    const setCurrentBannerCallback = (id: number) => {
        setCurrentBannerItem(id);
    };

    return (
        <>
            <div cds-layout="col@sm:5">
                <div cds-layout="vertical gap:md">
                    <RolloverList>
                        {rolloverBannerConfig.length &&
                            rolloverBannerConfig.map(({ logo, icon, title, description }, index) => (
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
                            ))}
                    </RolloverList>
                </div>
            </div>
            <div className="banner-content" cds-layout="col@sm:7">
                <div cds-text="h3" className="banner-content-title text-blue">
                    {rolloverBannerConfig[currentBannerItem].title}
                </div>
                <div className="banner-content-description">
                    {rolloverBannerConfig[currentBannerItem].description}
                </div>
                <div className="banner-content-action text-blue">
                    {rolloverBannerConfig[currentBannerItem].actionText}
                </div>
            </div>
        </>
    );
};

export default RolloverBanner;