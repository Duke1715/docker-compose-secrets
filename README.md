# docker-compose-secrets

Manage secrets as environment variables in Docker Compose with HashiCorp Vault

## Info

This setup represents a self-hosted version of the [Doppler](https://www.doppler.com) SecretOps Platform, using HashiCorp [Vault](https://www.vaultproject.io) as the secret management back-end.
The secrets are stored and managed securely in Vault and are injected into Docker Compose stacks by the DCS CLI application.

With this setup, no more secrets can be leaked through insufficiently protected `.env` or `docker-compose.yml` files.

## Demo

In Vault, the secrets (environment variables) of an application (in this case: [Logto](https://logto.io)) are stored in a [KV Secrets Engine - Version 2](https://developer.hashicorp.com/vault/docs/secrets/kv/kv-v2), located at the default "secret" path.
The secrets can be imported to Vault either by using the Web UI or the Vault CLI.

![image](https://user-images.githubusercontent.com/36514045/209883550-1f2851e1-ecdc-4d72-96a3-adb1d87c3844.png)

Then, the values of the secrets are removed from the Docker Compose file, only leaving the empty environment variables (marked in yellow) behind.

![image](https://user-images.githubusercontent.com/36514045/209885615-bd330192-09c6-4d73-9ce0-dc4a3fdc1c7e.png)

Now, the Docker Compose stack is started via the DCS application, automatically injecting the secrets provided by Vault as environment variables into the containers.

![image](https://user-images.githubusercontent.com/36514045/209884922-fc129f45-bbc4-4919-888a-cb303fa2b2ff.png)

The DCS application also supports stopping, restarting and updating the Docker Compose stack.

![image](https://user-images.githubusercontent.com/36514045/209885715-e9a5139a-0d3a-4e67-be53-9c3cbb3e3c54.png)

Instead of using Docker Compose commands, the DCS CLI is used for interacting with the containers.

## Guide

Please follow this guide to recreate the setup used in the demo above.

### Deploy Vault server

First of all, you need to set up a Vault server.
I recommend following the official [HashiCorp Vault tutorials](https://developer.hashicorp.com/vault/tutorials/getting-started) or using Docker to deploy the [Vault server](https://hub.docker.com/_/vault).

You may use these commands to create the Vault server container.

```bash
docker volume create vault_file
docker volume create vault_logs
docker volume create vault_config
docker run -d --name=vault --restart=unless-stopped -v vault_file:/vault/file -v vault_logs:/vault/logs -v vault_config:/vault/config --cap-add=IPC_LOCK -e 'VAULT_LOCAL_CONFIG={"storage": {"file": {"path": "/vault/file"}}, "listener": [{"tcp": { "address": "0.0.0.0:8200", "tls_disable": true}}], "default_lease_ttl": "168h", "max_lease_ttl": "720h", "ui": true}' -p 8200:8200 vault server
```

Make sure you use a reverse proxy such as [Nginx Proxy Manager](https://nginxproxymanager.com) or [Traefik](https://traefik.io/traefik) to protect the container from unauthorized access and issue valid SSL certificates for it.

### Configuring Vault server

You can either configure the Vault server from its own Web UI or use the Vault CLI.

In this guide, we are going to use the Vault CLI.

Make sure you have the Vault CLI installed on your server.
If you should not, please follow [this guide](https://developer.hashicorp.com/vault/tutorials/getting-started/getting-started-install) to install it.

First, we need to set the `VAULT_ADDR` environment variable to the url of the Vault server:

```bash
export VAULT_ADDR='http://127.0.0.1:8200'
```

Then we need to initialize the Vault:

```bash
vault operator init
```

Make sure you save the Unseal Keys and the Initial Root Token. You can not loose them!

Now log into the Vault by setting the `VAULT_TOKEN` environment variable to the Initial Root Token from the step above:

```bash
export VAULT_TOKEN='initial-root-token-here'
```

We now need to unseal the Vault with the Unseal Keys:

```bash
vault operator unseal
```

After that, ensure that the settings are correct by running following command:

```bash
vault secrets list
```

This should show you all available secrets engines.

Then we create our own secrets engine where we will store all our environment variables:

```bash
vault secrets enable -version=2 -path=secret kv
```

Now we can write our secrets to the secrets engine we just created:

```bash
vault kv put secret/logto DB_URL="postgres://postgres:p0stgr3s@postgres:5432/logto" ENDPOINT="https://logto.yourdomain.com" POSTGRES_DB="logto" POSTGRES_PASSWORD="p0stgr3s" POSTGRES_USER="postgres" TRUST_PROXY_HEADER="1"
```

Note: In production, you should use files or the Web UI instead of providing the secrets to the Vault CLI, as they will probably be logged to shell history, potentially exposing them.

To see, if the secrets were stored successfully, we use following command to print the secrets to the terminal:

```bash
vault kv get secret/logto
```

### Installing the DCS CLI

After successfully setting up Vault, we now install the DCS CLI.

First, download the latest binary from GitHub and make it executable:

```bash
sudo mkdir /usr/local/docker-compose-secrets && cd /usr/local/docker-compose-secrets
sudo wget https://github.com/KNIF/docker-compose-secrets/releases/download/1.0/dcs
sudo chmod +x dcs
```

Now add the DCS CLI to the `PATH` environment variable. You can also add the `VAULT_ADDR` and `VAULT_TOKEN` so you don't have to provide them at each login.
You can do this by adding the following line to your $HOME/.profile or /etc/profile (for a system-wide installation):

```bash
export PATH=$PATH:/usr/local/docker-compose-secrets
export VAULT_ADDR='http://127.0.0.1:8200'
export VAULT_TOKEN='initial-root-token-here'
```

After that, re-login to your terminal and run DCS to verify the installation was successful:

```bash
dcs
```

### Running a Docker Compose stack with the DCS CLI

Now that you have a fully working installation, navigate to the folder containing your `docker-compose.yml` for your stack and run following command:

```bash
dcs start
```

This will now fetch the secrets from Vault and inject them into the Docker Compose process.
