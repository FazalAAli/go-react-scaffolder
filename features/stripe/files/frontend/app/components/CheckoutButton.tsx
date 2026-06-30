import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { StripeService } from "gen/ts/app/v1/app_pb";

const transport = createConnectTransport({
  baseUrl: import.meta.env.VITE_API_URL ?? "http://localhost:8000",
});

const client = createClient(StripeService, transport);

export function CheckoutButton() {
  async function handleClick() {
    // TODO: pass the price ID for your product.
    const { url } = await client.createCheckoutSession({ priceId: "" });
    if (url) window.location.href = url;
  }

  return <button onClick={handleClick}>Checkout</button>;
}
