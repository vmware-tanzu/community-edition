// React imports
import React from 'react';

// Library imports
import { CdsControlMessage } from '@cds/react/forms';

export interface ThumbprintDisplayInputs {
    serverName: string;
    thumbprint: string;
    errorMessage: string;
}
export function ThumbprintDisplay(props: ThumbprintDisplayInputs) {
    return displayThumbprint(props.serverName, props.thumbprint, props.errorMessage);
}

function displayThumbprint(servername: string, print: string, errMsg: string) {
    if (errMsg) {
        return displayErrorThumbprint(servername, errMsg);
    }
    if (!print) {
        return emptyThumbprint();
    }
    const parts = print.split(':');
    if (parts.length === 1) {
        // This code handles the case where there is no colon in the thumbprint; it should never happen
        const halfway = print.length / 2;
        const firstHalf = print.substring(0, halfway);
        const secondHalf = print.substring(halfway);
        return displayThreePartThumbprint(servername, firstHalf, secondHalf);
    }
    const halfway = parts.length / 2;
    const firstHalf = parts.slice(0, halfway).join(':');
    const secondHalf = parts.slice(halfway).join(':');
    return displayThreePartThumbprint(servername, firstHalf, secondHalf);
}

// The point here is to keep the vertical spacing the same (as when there is a thumbprint to display) so
// that the display doesn't jiggle between empty and non-empty
function emptyThumbprint() {
    return (
        <>
            <div cds-layout="vertical gap:sm">
                <div>
                    <CdsControlMessage status="neutral">&nbsp;</CdsControlMessage>
                </div>
                <div className="thumbprint">
                    &nbsp;
                    <br />
                    &nbsp;
                </div>
            </div>
        </>
    );
}

function displayThreePartThumbprint(servername: string, part1: string, part2: string) {
    return (
        <>
            <div cds-layout="vertical gap:sm">
                <div>
                    <CdsControlMessage status="neutral">
                        SSL thumbprint for{' '}
                        <b>
                            <span>{servername}</span>
                        </b>
                    </CdsControlMessage>
                </div>
                <div className="thumbprint">
                    <span>{part1}</span>
                    <br />
                    <span>{part2}</span>
                </div>
            </div>
        </>
    );
}

function displayErrorThumbprint(servername: string, errMsg: string) {
    return (
        <>
            <div>
                <CdsControlMessage status="error">
                    Error retrieving SSL thumbprint of{' '}
                    <b>
                        <span>{servername}</span>
                    </b>
                    <br />
                    <span>{errMsg}</span>
                </CdsControlMessage>
            </div>
            <div className="thumbprint">
                &nbsp;
                <br />
                &nbsp;
            </div>
        </>
    );
}
