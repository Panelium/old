interface Config {
  BACKEND_HOST: string;
}

let configCache: Config;

export async function getConfig(): Promise<Config> {
  if (configCache) return configCache;

  const response = await fetch("/config.json");
  if (!response.ok) {
    throw new Error("Failed to fetch configuration");
  }

  configCache = await response.json();

  return configCache;
}
