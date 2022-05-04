import { Box, Button, Typography } from "@mui/material";
import PackageIcon from "./PackageIcon";

export interface Package {
    title: string;
    description: string;
    icon?: string;
    createAction?: () => void;
    deleteAction?: () => void;
    openAction?: () => void;
    showCreate?: boolean;
    showDelete?: boolean;
    showOpen?: boolean;
    helpLink?: string;
}

export default function PackageCard(props: Package) {
    const { title, description, icon, createAction, deleteAction, openAction, showCreate, showDelete, showOpen, helpLink, ...other } = props;
    const handleOpenLink = async () => {
        window.ddClient.host.openExternal(helpLink);
    };

    return (
        <Box sx={{
            minHeight: 150,
            maxHeight: 200,
            borderColor: (theme) => theme.palette.mode === 'dark' ? 'grey.300' : 'grey.800',
            border: 0,
            boxShadow: 3,
            mt: 3,
            display: 'flex',
            width: '100%'
        }}>
            <Box sx={{ margin: 3 }}>
                <PackageIcon icon={icon} width={50} />
            </Box>
            <Box sx={{ flexGrow: 1, mt: 4 }}>
                <Typography gutterBottom variant="h5" component="div">
                    {title}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    {description}
                </Typography>
            </Box>
            <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'stretch', justifyContent: 'space-evenly', ml: 2, mr: 5 }}>
                {(createAction !== undefined) &&
                    ((showCreate) &&
                        (<Button size="small" variant="contained" onClick={createAction} sx={{ width: '100px' }}>Create</Button>) ||
                        (<Button size="small" variant="outlined" sx={{ width: '100px' }} disabled>Create</Button>))
                }
                {(deleteAction !== undefined) && 
                    ((showDelete) && 
                        (<Button size="small" variant="outlined" onClick={deleteAction} sx={{ width: '100px' }}>Delete</Button>) ||
                        (<Button size="small" variant="outlined" sx={{ width: '100px' }} disabled>Delete</Button>))
                }
                {(openAction !== undefined) && 
                    ((showOpen) && 
                        (<Button size="small" variant="outlined" onClick={openAction} sx={{ width: '100px' }}>Open</Button>) ||
                        (<Button size="small" variant="outlined" sx={{ width: '100px' }} disabled>Open</Button>))
                }
                {(helpLink) && (
                    <Button size="small" variant="outlined" onClick={handleOpenLink} sx={{ width: '100px' }}>Learn more</Button>)}
            </Box>
        </Box>
    );
}