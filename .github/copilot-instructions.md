1. Every dialog component should behave like `responsive=dialog.tsx`. and add this class to the DialogContent `rounded-xl border-none bg-clip-padding shadow-2xl ring-4 ring-neutral-200/80 outline-none md:max-w-2xl dark:bg-neutral-800 dark:ring-neutral-900`
2. Use the `useIsMobile` hook from `hooks/use-mobile.ts` to determine if the device is mobile.
3. If the device is mobile, render a `Drawer` component.
4. If the device is not mobile, render a `Dialog` component.
5. For inline dynamic tailwind classes, use the `clsx` library to conditionally apply classes.