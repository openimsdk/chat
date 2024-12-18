# OpenIM Chat Deployment

## Preconditions

- Ensure deployed OpenIM Server and its dependencies.
  - Redis
  - MongoDB
  - Kafka
  - MinIO
- Expose the corresponding Services and ports of OpenIM Server.

## Deploy OpenIM Chat

**Chat depends on OpenIM Server, so you need to deploy OpenIM Server first.**

enter the target directory

```shell
cd deployments/deploy
```

### Modify ConfigMap

You need to modify the `chat-config.yml` file to match your environment. Focus on the following fields:
**discovery.yml**

- `kubernetes.namespace`: default is `default`, you can change it to your namespace.

**mongodb.yml**

- `address`: set to your already mongodb address or mongo Service name and port in your deployed.
- `database`: set to your mongodb database name.
- `authSource`: et to your mongodb authSource. (authSource is specify the database name associated with the user's credentials, user need create in this database.)

**redis.yml**

- `address`: set to your already redis address or redis Service name and port in your deployed.

**share.yml**

- `openIM.apiURL`: modify to your already API address or use your `openim-api` service name and port
- `openIM.adminUserID`: same to IM Server `imAdminUserID` field value.

### Set the secret

A Secret is an object that contains a small amount of sensitive data. Such as password and secret. Secret is similar to ConfigMaps.

#### Example:

create a secret for redis password. You can update `redis-secret.yml`.

you need update `redis-password` value to your redis password in base64.

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: openim-redis-secret
type: Opaque
data:
  redis-password: b3BlbklNMTIz # you need update to your redis password in base64
```

#### Usage:

use secret in deployment file. If you apply the secret to IM Server, you need adapt the Env Name to config file and all toupper.

OpenIM Chat use prefix `CHATENV_`. Next adapt is the config file name. Like `redis.yml`. Such as `CHATENV_REDIS_PASSWORD` is mapped to `redis.yml` password filed in OpenIM Server.

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
            - name: CHATENV_REDIS_PASSWORD # adapt to redis.yml password field in OpenIM Server config, Don't modify it.
              valueFrom:
                secretKeyRef:
                  name: openim-redis-secret # You deployed secret name
                  key: redis-password # You deployed secret key name
```

So, you need following configurations to set secret:

- `MONGODB_USERNAME`
- `MONGODB_PASSWORD`
- `REDIS_PASSWORD`

Apply the secret.

```shell
kubectl apply -f redis-secret.yml -f mongo-secret.yml
```

### Apply Config and Services

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
