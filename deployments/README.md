# OpenIM Chat Deployment

## Preconditions

- Successfully deployed OpenIM Server and its dependencies(mongo, kafka, redis, minio).
- Expose the corresponding Services and ports of OpenIM Server.

## Deploy OpenIM Chat

Chat depends on OpenIM Server, so you need to deploy OpenIM Server first.

### Modify ConfigMap

You need to modify the `chat-config.yml` file to match your environment. Focus on the following fields:
**discovery.yml**

- `kubernetes.namespace`: default is `default`, you can change it to your namespace.
- `enable`: set to `kubernetes`
- `rpcService`: Every field value need to same to the corresponding service name. Such as `chat` value in same to `openim-chat-rpc-service.yml` service name.

**log.yml**

- `storageLocation`: log save path in container.
- `isStdout`: output in kubectl log.
- `isJson`: log format to JSON.

**mongodb.yml**

- `address`: set to your already mongodb address or mongo Service name and port in your deployed.
- `username`: set to your mongodb username.
- `database`: set to your mongodb database name.
- `password`: **need to set to secret use base64 encode.**
- `authSource`: set to your mongodb authSource, default is `openim_v3`.

**redis.yml**

- `address`: set to your already redis address or redis Service name and port in your deployed.
- `password`: **need to set to secret use base64 encode.**

**share.yml**

- `openIM.apiURL`: modify to your already API address or use your `openim-api` service name and port
- `openIM.adminUserID`: same to IM Server `imAdminUserID` field value.
- `chatAdmin`: default is `chatAdmin`.

### Set the secret

A Secret is an object that contains a small amount of sensitive data. Such as password and secret. Secret is similar to ConfigMaps.

#### Example:

create a secret for redis password. You can create new file is `redis-secret.yml` or append contents to `chat-config.yml` use `---` split it.

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: redis-secret
type: Opaque
data:
  redis-password: b3BlbklNMTIz # "openIM123" in base64
```

#### Usage:

use secret in deployment file. If you apply the secret to IM Server, you need adapt the Env Name to config file and all toupper.

OpenIM Server use prefix `IMENV_`, OpenIM Chat use prefix `CHATENV_`. Next adapt is the config file name. Like `redis.yml`. Such as `CHATENV_REDIS_PASSWORD` is mapped to `redis.yml` password filed in OpenIM Server.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: chat-rpc-server
spec:
  template:
    spec:
      containers:
        - name: chat-rpc-server
          env:
            - name: CHATENV_REDIS_PASSWORD # adapt to redis.yml password field
              valueFrom:
                secretKeyRef:
                  name: redis-secret
                  key: redis-password
```

So, you need following configurations to set secret:

- `MONGODB_USERNAME`
- `MONGODB_PASSWORD`
- `REDIS_PASSWORD`
- `SHARE_OPENIM_SECRET`

### Apply Config and Services

enter the target directory

```shell
cd deployments/deploy
```

deploy the config and services

```shell
kubectl apply -f chat-config.yml -f openim-admin-api-service.yml -f openim-chat-api-service.yml -f openim-admin-rpc-service.yml -f openim-chat-rpc-service.yml
```

### Start Chat Deployments

```shell
kubectl apply -f openim-chat-api-deployment.yml -f openim-admin-api-deployment.yml -f openim-chat-rpc-deployment.yml -f openim-admin-rpc-deployment.yml
```

## Verify

After the deployment is complete, you can verify the deployment status.

```shell
# Check the status of all pods
kubectl get pods

# Check the status of services
kubectl get svc

# Check the status of deployments
kubectl get deployments

# View all resources
kubectl get all

```

## clean all

`kubectl delete -f ./`

## Notes:

- If you use a specific namespace for your deployment, be sure to append the -n <namespace> flag to your kubectl commands.
