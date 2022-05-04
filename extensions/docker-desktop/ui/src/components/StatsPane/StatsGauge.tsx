import CircularProgress, {
    CircularProgressProps,
} from '@mui/material/CircularProgress';
import Typography from '@mui/material/Typography';
import Box from '@mui/material/Box';

function CircularProgressWithLabel(
    props: CircularProgressProps & { value: number, label: string },
) {
    return (
        <Box sx={{ top: 50, left: 50, position: 'relative', display: 'inline-flex' }}>
            <CircularProgress variant="determinate" color="success" {...props} />
            <Box
                sx={{
                    top: 0,
                    left: 0,
                    bottom: 0,
                    right: 0,
                    position: 'absolute',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                }}
            >
                <Typography
                    variant="caption"
                    component="div"
                    color="text.secondary"
                >{`${props.value.toFixed(2)}%`}</Typography>
            </Box>
            <Box
                sx={{
                    top: 150,
                    left: 0,
                    bottom: 0,
                    right: 0,
                    position: 'absolute',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                }}
            >
                <Typography
                    variant="caption"
                    component="div"
                    color="text.secondary"
                >{props.label}</Typography>
            </Box>
        </Box>
    );
}


export interface IStatsGaugeProps {
    value: number;
    scale: number;
    label: string
}

export default function StatsGauge(props: IStatsGaugeProps) {
    const avgval = (props.value * 100) / props.scale;

    return <CircularProgressWithLabel size='120px' value={avgval} label={props.label} />;
}