// React imports
import React from 'react';

// Library imports
import { CdsIcon } from '@cds/react/icon';

// App imports
import './RolloverBannerItem.scss';

export interface RolloverProps {
    onMouseEnter: () => void,
    index: number,
    logo: string,
    icon: string,
    title: string,
    selected: boolean,
    mouseEnterCallback: (index: number) => void
}

const RolloverBannerItem = (props:RolloverProps) => {
    const {
        index,
        logo,
        icon,
        title,
        selected,
        mouseEnterCallback
    }: RolloverProps = props;

    return (
        <>
            <div
                cds-text="message"
                className="rollover-item"
                onMouseEnter={() => {
                    mouseEnterCallback(index);
                }}>
                <div className="rollover-item-icon">
                    { logo ?
                        <img src={logo} className="logo-42" alt="tce logo"/> :
                        <CdsIcon shape={icon} size="lg" className="icon-blue"></CdsIcon>
                    }
                </div>
                <div className="rollover-item-title">
                    <span className={(selected ? '' : 'text-blurred')}>{title}</span>
                </div>
                <div className="rollover-item-arrow">
                    <CdsIcon className={'icon-blue ' + (selected ? '' : 'hidden')} shape="angle" direction="right" size="lg"></CdsIcon>
                </div>
            </div>
        </>
    );
};

export default RolloverBannerItem;