interface Config {
  apiUrl: string;
  environment: string;
}

const configs: Record<string, Config> = {
  development: {
    apiUrl: 'http://localhost:8000',
    environment: 'development'
  },
  production: {
    apiUrl: 'https://api.canadahires.info',
    environment: 'production'
  },
  staging: {
    apiUrl: 'https://staging-api.canadahires.info',
    environment: 'staging'
  }
};

const getConfig = (): Config => {
  const env = import.meta.env.MODE || 'development';
  return configs[env] || configs.development;
};

export const config = getConfig();
export default config;