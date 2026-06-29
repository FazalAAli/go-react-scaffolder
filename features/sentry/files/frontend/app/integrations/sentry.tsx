import * as Sentry from "@sentry/react";
import type { ReactNode } from "react";

const dsn = import.meta.env.VITE_SENTRY_DSN;

if (dsn) {
  Sentry.init({
    dsn,
    environment: import.meta.env.MODE,
    integrations: [
      Sentry.browserTracingIntegration(),
      Sentry.replayIntegration(),
    ],
    tracesSampleRate: import.meta.env.PROD ? 0.1 : 1.0,
    tracePropagationTargets: ["localhost", import.meta.env.VITE_API_URL],
    replaysSessionSampleRate: 0.1,
    replaysOnErrorSampleRate: 1.0,
  });
}

export function SentryProvider({ children }: { children: ReactNode }) {
  if (!dsn) return <>{children}</>;
  return (
    <Sentry.ErrorBoundary fallback={<p>An error has occurred</p>}>
      {children}
    </Sentry.ErrorBoundary>
  );
}
