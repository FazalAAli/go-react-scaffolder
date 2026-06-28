import type { ReactNode } from "react";
// scaffold:region:providers-imports:start
// scaffold:region:providers-imports:end

type ProviderComponent = ({ children }: { children: ReactNode }) => ReactNode;

const providers: ProviderComponent[] = [
  // scaffold:region:providers:start
  // scaffold:region:providers:end
];

export function Providers({ children }: { children: ReactNode }) {
  return providers.reduceRight<ReactNode>(
    (acc, Provider) => <Provider>{acc}</Provider>,
    children,
  );
}
