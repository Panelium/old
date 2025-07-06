import { Client, createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { GenService } from "@bufbuild/protobuf/codegenv2";
import { getConfig } from "~/lib/config";
import { AuthService } from "proto-gen-ts/backend/Auth_pb";
import { ClientService } from "proto-gen-ts/backend/Client_pb";
import { ServerService } from "proto-gen-ts/daemon/Server_pb";

const clientCache = new Map<string, Client<any>>();

async function createTransport(baseUrl?: string) {
  if (baseUrl) {
    return createConnectTransport({ baseUrl });
  }
  const config = await getConfig();
  return createConnectTransport({ baseUrl: config.BACKEND_HOST });
}

export async function getClient<T extends GenService<any>>(
  service: T,
  baseUrl?: string
): Promise<Client<T>> {
  const key = baseUrl ? `${service.typeName}|${baseUrl}` : service.typeName;
  if (clientCache.has(key)) {
    return clientCache.get(key) as Client<T>;
  }
  const client = createClient(service, await createTransport(baseUrl));
  clientCache.set(key, client);
  return client;
}

export const getAuthClient = () => getClient(AuthService);
export const getClientClient = () => getClient(ClientService);
export const getDaemonServerClient = (baseUrl: string) =>
  getClient(ServerService, baseUrl);
