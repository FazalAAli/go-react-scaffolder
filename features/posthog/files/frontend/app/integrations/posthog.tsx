import { PostHogProvider as PHProvider } from "@posthog/react";
import type { ReactNode } from "react";

const key = import.meta.env.VITE_POSTHOG_KEY;
const host = import.meta.env.VITE_POSTHOG_HOST ?? "https://us.i.posthog.com";

export function PostHogProvider({ children }: { children: ReactNode }) {
  if (!key) return <>{children}</>;
  return (
    <PHProvider apiKey={key} options={{ api_host: host }}>
      {children}
    </PHProvider>
  );
}
