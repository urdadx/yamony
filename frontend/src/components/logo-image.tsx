import LOGO from "@/assets/icons/nucleopass-icon.png";
import Alfred from "@/assets/yamony-logo.png";
import { cn } from "@/lib/utils";

type ImageProps = {
  width?: number;
  height?: number;
  className?: string;
  [key: string]: any;
};

export const Logo = ({
  width = 40,
  height = 40,
  className,
  ...props
}: ImageProps) => {
  return (
    // biome-ignore lint/a11y/useAltText: <explanation>
    <img
      src={LOGO}
      alt="logo"
      width={width}
      height={height}
      className={cn(`w-[${width}px] h-[${height}px]`, className)}
      {...props}
    />
  );
};

export const UserImage = ({
  width = 40,
  height = 40,
  className,
  ...props
}: ImageProps) => {
  return (
    // biome-ignore lint/a11y/useAltText: <explanation>
    <img
      src={Alfred}
      alt="user"
      width={width}
      height={height}
      className={cn(`w-[${width}px] h-[${height}px]`, className)}
      {...props}
    />
  );
};