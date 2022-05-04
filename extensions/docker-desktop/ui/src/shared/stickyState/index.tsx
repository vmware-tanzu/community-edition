// Taken from https://github.com/kubeapps/kubeapps-dd/blob/main/client/src/shared/stickyState/index.tsx
import { useEffect, useState } from "react";

// Taken from https://www.joshwcomeau.com/react/persisting-react-state-in-localstorage
export default function useStickyState(defaultValue: any, key: string) {
  const [value, setValue] = useState(() => {
    const stickyValue: string | null = window.localStorage.getItem(key);
    return stickyValue !== null ? JSON.parse(stickyValue).value : defaultValue;
  });
  useEffect(() => {
    window.localStorage.setItem(key, JSON.stringify({ value: value }));
  }, [key, value]);
  return [value, setValue];
}