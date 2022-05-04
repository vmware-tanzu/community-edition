// Maybe replace with https://github.com/mozilla-frontend-infra/react-lazylog
import React, { useState, useEffect, useRef } from 'react';
import { Box } from '@mui/material';
import AutoScroll from './AutoScroll';


export interface ILogTailerBoxProps {
  // Children to render in the scroll container.
  log: string[];
  // Height value of the scroll container.
  height?: number;
  // Ability to disable the smooth scrolling behavior.
  scrollBehavior?: 'smooth' | 'auto';
}

export default function LogTailerBox(props: ILogTailerBoxProps) {
  const [height, setHeight] = useState(0);
  const parentRef = useRef<any>();

  useEffect(() => {
    if (parentRef.current) {
      setHeight(parentRef.current.clientHeight);
    }
  });

  return (
    <Box id="parentautoscroll" sx={{ width: '100%', height: '100%' }} ref={parentRef}>
      <AutoScroll showOption={false} height={height} preventInteraction={false}>
        {props.log.map((line,i) => {
          return (
            <Box key={i} sx={{ whiteSpace: 'pre' }}>
              {line}
            </Box>
          );
        })}
      </AutoScroll>
    </Box>
  );
}