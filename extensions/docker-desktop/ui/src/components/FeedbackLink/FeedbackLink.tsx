import ChatTwoToneIcon from '@mui/icons-material/ChatTwoTone';
import ArticleTwoToneIcon from '@mui/icons-material/ArticleTwoTone';
import { Button, ButtonGroup } from '@mui/material';
import {
  EXTENSION_REPO_ISSUES,
  EXTENSION_REPO_LICENSE,
} from 'shared/constants';

const handleFeedback = async () => {
  window.ddClient.host.openExternal(EXTENSION_REPO_ISSUES);
};

const handleLicense = async () => {
  window.ddClient.host.openExternal(EXTENSION_REPO_LICENSE);
};

export default function FeedbackLink() {
  return (
    <ButtonGroup size="small" sx={{ position: 'absolute', top: 10, right: 10 }}>
      <Button variant="text" aria-label="Feedback" onClick={handleFeedback}>
        <ChatTwoToneIcon fontSize="small" color="primary" />
        Feedback
      </Button>
      <Button variant="text" aria-label="License" onClick={handleLicense}>
        <ArticleTwoToneIcon fontSize="small" color="primary" />
        License
      </Button>
    </ButtonGroup>
  );
}
