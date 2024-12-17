# OpenIM Chat Deployment

## Preconditions

- Successfully deployed OpenIM Server and its dependencies(mongo, kafka, redis, minio).
- Expose the corresponding SVCs and ports of OpenIM Server.

## Deploy OpenIM Chat

Chat depends on OpenIM Server, so you need to deploy OpenIM Server first.

### Modify ConfigMap

You need to modify the `chat-config.yml` file to match your environment. Focus on the following fields:
**discovery.yml**

- kubernetes.namespace
- rpcService

**log.yml**

- storageLocation
- isStdout

**mongodb.yml**

- address (modify to the address or mongodb service)
- database
- username
- password
- authSource

**redis.yml**

- address (modify to the address or redis service)
- password

**share.yml**

- openIM.apiURL (modify to the address or `openim-api` service)
- openIM.adminUserID (same to IM Server `imAdminUserID` field config)
- chatAdmin

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
