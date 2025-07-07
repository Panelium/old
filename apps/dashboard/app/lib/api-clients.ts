import {
  Client,
  Code,
  ConnectError,
  createClient,
  StreamRequest,
  StreamResponse,
  UnaryRequest,
  UnaryResponse,
} from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { GenService } from "@bufbuild/protobuf/codegenv2";
import { Config, getConfig } from "~/lib/config";
import { AuthService } from "proto-gen-ts/backend/Auth_pb";
import { ClientService } from "proto-gen-ts/backend/Client_pb";
import { ServerService } from "proto-gen-ts/daemon/Server_pb";
import { DescMessage } from "@bufbuild/protobuf";
import { setSessionAuthenticated } from "~/providers/SessionProvider";

const clientCache = new Map<string, Client<any>>();

async function createTransport(baseUrl?: string) {
  const config = baseUrl ? ({} as Config) : await getConfig();
  return createConnectTransport({
    baseUrl: baseUrl || config.BACKEND_HOST,
    fetch: (input, init) => fetch(input, { ...init, credentials: "include" }),
    interceptors: [
      (next) =>
        async (
          req: UnaryRequest | StreamRequest
        ): Promise<UnaryResponse<DescMessage, DescMessage> | StreamResponse<DescMessage, DescMessage>> => {
          // Do not intercept refreshToken requests
          if (req.service.typeName === "backend.AuthService" && req.method.name === "RefreshToken") {
            return await next(req);
          }
          try {
            return await next(req);
          } catch (err) {
            if (
              err instanceof ConnectError &&
              (err.code === Code.Unauthenticated || err.rawMessage?.toLowerCase().includes("unauthenticated"))
            ) {
              try {
                // Try to refresh the token
                const authClient = await getAuthClient();
                const refreshRes = await authClient.refreshToken({});
                if (refreshRes.success) {
                  // Retry the original request
                  return await next(req);
                } else {
                  setSessionAuthenticated(false);
                  console.error(err);
                  throw err;
                }
              } catch (refreshErr) {
                setSessionAuthenticated(false);
                throw err;
              }
            }
            throw err;
          }
        },
    ],
  });
}

export async function getClient<T extends GenService<any>>(service: T, baseUrl?: string): Promise<Client<T>> {
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
export const getDaemonServerClient = (baseUrl: string) => getClient(ServerService, baseUrl);
