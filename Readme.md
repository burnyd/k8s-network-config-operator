# K8s network operator

![Alt text](/images/operator.jpg?raw=true "overall")

This is a network configuration operator very similar to https://github.com/burnyd/k8sNetworkgNMIOperator that it is not exactly the same where that operator uses openconfig to do all of the pushing of configuration.  This operator uses eAPI to push the configuration.

This should not be seen as production ready but rather as a means of demo of what is possible with network automation and kubernetes.

[goEAPI](https://github.com/aristanetworks/goeapi) is pretty simplistic and easy to use.  Plenty of general examples can be found [here within this gist](https://gist.github.com/burnyd/4d742ecbd2010d5e725f3649954f7370)

This is build on the [Kubernetes operator framework](https://operatorframework.io/) which is really amazing for providing scaffolding to create kubernetes operators.

Here is a [CRD](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) as an example state for a description of a network device.  Today this is built using EOS but anything could be substituted with enough work and a Go module related to a device.

```shell
apiVersion: network.networkconfigoperator/v1
kind: NetDevs
metadata:
  name: ceos1
  namespace: networkconfig
spec:
  username: admin
  password: admin
  host: 172.20.20.2
  transport: http
  port: 80
  networkos: eos
  runningconfig:  |+
    transceiver qsfp default-mode 4x10G
    !
    service routing protocols model multi-agent
    !
    hostname ceos1
    !
    spanning-tree mode mstp
    !
    no aaa root

```

The general idea is that a user can use their current tooling(salt, ansible etc) render a config a long with the rest of the kubernetes yaml manifest or even talk directly to the kube api and create this cr based off of tooling.

I truncated this so it is east to view.  Just an example of keys within the data structure.
```shell
apiVersion: network.networkconfigoperator/v1 # The API for the CRD
kind: NetDevs # The CR
metadata:
  name: ceos1 # Field for the name of the device.
  namespace: networkconfig # This is necessary as this needs to run within this name space
spec:
  username: admin # Username this can be either plain text or as a k8s secret.
  password: admin # Password this can be either plain text or as a k8s secret.
  host: 172.20.20.2 # IP address of the device
  transport: http # Transport http,https or ssh.
  port: 80 # Port of the device.
  networkos: eos # Network OS eos is supported today but can add functionality for anything else.
  runningconfig: # String of the network configuration.

```


# Demo

Create a kind cluster unless you already have access to a kubernetes cluster:
```
kind create cluster
```

Apply everything needed for the operator to run inside of the cluster with kustomize this will install the CRD, service account and rbac information.
```
kubectl apply -k manifests/.
```

Edit the configuration files or create your own. For this demo they are located inside of manifests/ceos
```
kubectl apply -f manifests/ceos/.
```

You should see two devices if you are following the demo
```
➜  network-config-operator git:(master) ✗ kubectl get netdevs -n networkconfig
NAME    AGE
ceos1   7m7s
ceos2   7m7s
```

```
2022-01-26T02:52:38.298Z        INFO    Checking on Device 192.168.4.244
2022-01-26T02:52:38.372Z        INFO    Switch config matches CR
2022-01-26T02:52:38.383Z        INFO    Starting to check on the device
2022-01-26T02:52:38.383Z        INFO    Device is EOS
2022-01-26T02:52:38.383Z        INFO    Checking on Device 192.168.4.244
2022-01-26T02:52:38.452Z        INFO    Switch config matches CR
```

## Tested versions
- Kubernetes 1.22.2
- cEOS 4.26.1F
- Operator framework v3
- go 1.16