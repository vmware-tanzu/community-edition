// From https://github.com/kubeapps/kubeapps-dd/blob/3c46864485d9eb85a7d84e6df1b6023084ad7185/client/src/components/Icon/Icon.tsx
import { useEffect, useState } from "react";
import { icons } from "shared/images";

export interface IIconProps {
  width?: number;
  icon?: any;
}

export default function PackageIcon({ icon, width=150 }: IIconProps) {
  const [srcIcon, setSrcIcon] = useState(icons.tce);
  const [iconErrored, setIconErrored] = useState(false);
  useEffect(() => {
    if (srcIcon !== icon && icon && !iconErrored) {
      setSrcIcon(icon);
    }
  }, [srcIcon, icon, iconErrored]);

  const onError = () => {
    setIconErrored(true);
    setSrcIcon(icons.tce);
  };

  return <img src={srcIcon} alt="icon" width={width} height={width} onError={onError} />;
}