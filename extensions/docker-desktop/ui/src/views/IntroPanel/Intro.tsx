import React, { useState, useRef, useEffect } from 'react';

import { Box, Button, Stack, Typography, Link } from '@mui/material';
import { FeedbackLink } from 'components';

// https://mui.com/guides/routing/#button
import { Link as RouterLink } from 'react-router-dom';
import { TCE_DEMO_SITE, DOCS_TCE_SITE, TCE_SLACK_SITE } from 'shared/constants';

const handleGoToTCEDemo = async () => {
  window.ddClient.host.openExternal(TCE_DEMO_SITE);
};

const handleGoToTCEDocs = async () => {
  window.ddClient.host.openExternal(DOCS_TCE_SITE);
};

const handleGoToTCEInSlack = async () => {
  window.ddClient.host.openExternal(TCE_SLACK_SITE);
};

export default function Intro() {
  return (
    <Stack spacing={2}>
      <FeedbackLink />
      <Box py={5} textAlign="center" alignItems="center">
        <Box sx={{ mt: 2 }}>
          <img
            src="TCE-logo.svg"
            width="300"
            className="Header-logo"
            alt="logo"
          />
        </Box>
        <Box sx={{ mt: 2, mx: 15 }}>
          <Typography>
            Tanzu Community Edition is a full-featured, easy-to-manage
            Kubernetes platform for all technical practitioners, from developers
            to platform operators and architects. It is a community supported,
            open source distribution including VMware Tanzu’s ecosystem and
            beyond, that can be installed and configured in minutes on your
            local workstation or your favorite cloud.
          </Typography>
        </Box>
        <Stack
          sx={{ mt: 4 }}
          direction="row"
          spacing={2}
          justifyContent="center"
        >
          <Button
            component={RouterLink}
            variant="contained"
            to="/create-cluster"
          >
            Create cluster
          </Button>
        </Stack>
      </Box>
      <Box
        position="absolute"
        bottom="50px"
        width="90%"
        sx={{ py: 5, px: 3, borderRadius: 4, boxShadow: 3 }}
      >
        <Stack
          direction="row"
          sx={{ display: 'flex', justifyContent: 'space-between' }}
        >
          <Box>
            <Typography style={{ fontWeight: 600 }} textAlign="left">
              New to VMware Tanzu Community Edition?
            </Typography>
            <Typography textAlign="left">
              Not sure if Tanzu Community Edition is right for you? Check out
              the{' '}
              <Link href="#" onClick={handleGoToTCEDocs}>
                docs
              </Link>
              , join us on the{' '}
              <Link href="#" onClick={handleGoToTCEInSlack}>
                Kubernetes slack
              </Link>{' '}
              or watch this walkthrough demo!
            </Typography>
          </Box>
          <Box
            sx={{
              display: 'flex',
              alignItems: 'center',
              width: '20%',
              justifyContent: 'center',
            }}
          >
            <Button variant="outlined" onClick={handleGoToTCEDemo}>
              Watch demo
            </Button>
          </Box>
        </Stack>
      </Box>
    </Stack>
  );
}
