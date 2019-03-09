# Simple Pod Watcher
This is simple example using the client-go kubernetes client.  The service will watch as pods are added or deleted and will print the pod name along with it's assigned IP Address

## Running Pod Watcher
You can run the pod watcher either in-cluster or out-of-cluster.

### Running Pod Watcher Out-of-Cluster
When running Pod Watcher out-of-cluster, it uses the kubeconfig set in your KUBECONFIG environment variable.  If that is not set it will use $HOME/.kube/config

At a minimum, the authenticated user must have cluster-wide read-access to pods

The Pod Watcher doesn't take any arguments.  Simply run the following:

```
# First build the podwatcher
cd src/
go build -o ../bin/podwatcher-controller

# Now run it
../bin/podwatcher-controller
```

### Running Pod Watcher In-Cluster
When running Pod Watcher in-cluster, it uses the associated service account to authenticate with the API Server.  At a minimum, the service account must be given a ClusterRoleBinding that has read-access to pods in all namespaces

To deploy Pod Watcher into your cluster, create a project or namespace named pod-watcher and run the following:
```
oc apply -n pod-watcher -f k8s/
```
**Note:** 
If using a different namespace be sure to change it in k8s/cluster-role-binding.yaml

If you build your own pod watcher image you must update the k8s/pod-watcher-deployment.yaml manifest to point to your image.
