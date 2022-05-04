import { Box, Button, Stack, Typography } from "@mui/material";
import { useLocation, useNavigate } from "react-router-dom";

export default function NotFoundView() {
  const navigate = useNavigate();
  const location = useLocation();
  return (
    <>
      <Box py={5} textAlign="center">
        <Typography variant="h2">Not found</Typography>
        <Typography> {location.pathname}</Typography>
        <Stack sx={{ mt: 4 }} direction="row" spacing={2} justifyContent="center">
          <Button onClick={() => navigate(-1)} size="large" variant="contained">
            Go back
          </Button>
        </Stack>
      </Box>
    </>
  );
}