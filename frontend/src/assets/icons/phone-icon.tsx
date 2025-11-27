import type { SVGProps } from "react";

export function PhoneIcon(props: SVGProps<SVGSVGElement>) {
  const { color = "#888888" } = props;
  return (
    <svg xmlns="http://www.w3.org/2000/svg" width="1em" height="1em" viewBox="0 0 24 24" {...props}>{/* Icon from Solar by 480 Design - https://creativecommons.org/licenses/by/4.0/ */}<path fill={color} d="M12.052 2c1.68 0 3.01 0 4.052.142c1.072.147 1.94.456 2.624 1.152s.988 1.58 1.132 2.67c.14 1.06.14 2.414.14 4.123v3.826c0 1.71 0 3.064-.14 4.123c-.144 1.09-.448 1.974-1.132 2.67s-1.552 1.005-2.624 1.152C15.063 22 13.732 22 12.052 22h-.104c-1.68 0-3.01 0-4.052-.142c-1.072-.147-1.94-.456-2.624-1.152s-.988-1.58-1.132-2.67C4 16.976 4 15.622 4 13.913v-3.826c0-1.71 0-3.064.14-4.123c.144-1.09.448-1.974 1.132-2.67S6.824 2.29 7.896 2.142C8.937 2 10.268 2 11.948 2z" opacity=".3" /><path fill={color} d="M8.571 18.511c0-.385.307-.697.686-.697h5.486c.378 0 .685.312.685.697a.69.69 0 0 1-.685.698H9.257a.69.69 0 0 1-.686-.698" /></svg>
  )
}