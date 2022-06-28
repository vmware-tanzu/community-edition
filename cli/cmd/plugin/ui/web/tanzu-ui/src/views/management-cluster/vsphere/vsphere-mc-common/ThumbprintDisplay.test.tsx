import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { ThumbprintDisplay } from './ThumbprintDisplay';

const testServerName = 'TEST_SERVER';
const firstHalfThumbprint = '11:22:33:44:55:66:77:88:99';
const secondHalfThumbprint = '00:AA:BB:CC:DD:EE:FF:GG:HH';
const testThumbprint = firstHalfThumbprint + ':' + secondHalfThumbprint;
const testErrorMessage = 'Totally unable to get thumbprint';
describe('ThumbprintDisplay component', () => {
    it('should render', async () => {
        const view = render(<ThumbprintDisplay serverName={testServerName} thumbprint={testThumbprint} errorMessage={testErrorMessage} />);
        await waitFor(() => {
            expect(view).toBeDefined();
        });
    });
    it('should display non-empty colon thumbprint in two parts with server name', async () => {
        render(<ThumbprintDisplay serverName={testServerName} thumbprint={testThumbprint} errorMessage={''} />);
        expect(await screen.findByText(testServerName)).toBeInTheDocument();
        expect(await screen.findByText(firstHalfThumbprint)).toBeInTheDocument();
        expect(await screen.findByText(secondHalfThumbprint)).toBeInTheDocument();
        expect(screen.queryByText(testThumbprint)).toBeNull();
    });
    it('on error, should display only error message with server name', async () => {
        render(<ThumbprintDisplay serverName={testServerName} thumbprint={testThumbprint} errorMessage={testErrorMessage} />);
        expect(await screen.findByText(testServerName)).toBeInTheDocument();
        expect(await screen.findByText(testErrorMessage)).toBeInTheDocument();
        expect(screen.queryByText(firstHalfThumbprint)).toBeNull();
        expect(screen.queryByText(secondHalfThumbprint)).toBeNull();
        expect(screen.queryByText(testThumbprint)).toBeNull();
    });
    // the non-colon thumbprint is for completeness; we never expect to encounter it
    it('should display non-colon thumbprint in two parts with server name', async () => {
        const firstPartThumbprint = 'abcdef';
        const secondPartThumbprint = 'GHIJKL';
        const thumbprint = firstPartThumbprint + secondPartThumbprint;
        render(<ThumbprintDisplay serverName={testServerName} thumbprint={thumbprint} errorMessage={''} />);
        expect(await screen.findByText(testServerName)).toBeInTheDocument();
        expect(await screen.findByText(firstPartThumbprint)).toBeInTheDocument();
        expect(await screen.findByText(secondPartThumbprint)).toBeInTheDocument();
        expect(screen.queryByText(thumbprint)).toBeNull();
    });
    it('should display empty thumbprint without any server name', async () => {
        render(<ThumbprintDisplay serverName={testServerName} thumbprint={''} errorMessage={''} />);
        expect(screen.queryByText(testServerName)).toBeNull();
    });
});
