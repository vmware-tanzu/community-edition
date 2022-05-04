// Maybe replace with https://github.com/mozilla-frontend-infra/react-lazylog
import { Box, IconButton, Typography } from '@mui/material';
import ContentCopyIcon from '@mui/icons-material/ContentCopy';
import useCopyToClipboard from 'shared/copyToClipboard';


export interface IKubeconfigBoxProps {
  kubeconfig: string;
}

export default function KubeconfigBox(props: IKubeconfigBoxProps) {
  const [kubeconfig, copy] = useCopyToClipboard();

  return (
    <Box sx={{ width: '100%', height: '100%', position: 'relative' }}>
      <Box sx={{ width: '100%', height: '100%', position: 'absolute', overflow: 'auto' }}>
        <Typography sx={{ whiteSpace: 'pre-wrap', wordWrap: 'break-word' }}>
          {props.kubeconfig}
        </Typography>
      </Box>
      {(props.kubeconfig !== undefined && props.kubeconfig != '') &&
        <IconButton sx={{ position: 'absolute', top: 0, right: 0, }} onClick={() => copy(props.kubeconfig)}>
          <ContentCopyIcon />
        </IconButton>
      }
    </Box>
  );
}