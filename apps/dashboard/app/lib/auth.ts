import { Client, createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { AuthService } from "proto-gen-ts/backend/Auth_pb";
import { getConfig } from "~/lib/config";

let authClient: Client<typeof AuthService>;

export async function getAuthClient() {
  if (authClient) return authClient;

  const config = await getConfig();

  const transport = createConnectTransport({
    baseUrl: config.BACKEND_HOST,
  });

  authClient = createClient(AuthService, transport);

  return authClient;
}
