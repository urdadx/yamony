import type { SVGProps } from "react";

export function BrushIcon(props: SVGProps<SVGSVGElement>) {
  const { color = "#888888" } = props;
  return (
    <svg xmlns="http://www.w3.org/2000/svg" width="1em" height="1em" viewBox="0 0 24 24" {...props}>{/* Icon from Huge Icons by Hugeicons - undefined */}<g fill="none" stroke={color} strokeLinejoin="round" strokeWidth="1.5"><path strokeLinecap="round" d="m12 13l3 2m-3-2c-4.48 2.746-7.118 1.78-9 .85c0 2.08.681 3.82 1.696 5.15M12 13l3-4.586M15 15c-.219 2.037-1.573 6.185-4.844 7c-1.815 0-3.988-1.07-5.46-3M15 15l2.598-5m0 0l3.278-6.31a1.166 1.166 0 0 0-.524-1.567a1.174 1.174 0 0 0-1.544.47L15 8.414M17.598 10L15 8.414M4.696 19c1.038.167 3.584.2 5.46-1" /><path d="m6 4l.221.597c.29.784.435 1.176.72 1.461c.286.286.678.431 1.462.72L9 7l-.597.221c-.784.29-1.176.435-1.461.72c-.286.286-.431.678-.72 1.462L6 10l-.221-.597c-.29-.784-.435-1.176-.72-1.461c-.286-.286-.678-.431-1.462-.72L3 7l.597-.221c.784-.29 1.176-.435 1.461-.72c.286-.286.431-.678.72-1.462z" /></g></svg>
  )
}