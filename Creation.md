# Steps to produce with the operator sdk

### Create the api and repo
operator-sdk init --skip-go-version-check --domain networkconfigoperator --repo github.com/burnyd/networkconfigoperator --verbose

### Create the controller and scaffold the rest.
operator-sdk create api --group network --version v1 --kind NetDevs --resource --controller

## Edit the spec value for the cr

Can be found within the api/netdevs_types.go

type NetDevsSpec struct {
	Username  string `json:"username"`  // Username for the network device.
	Password  string `json:"password"`  // Password for the device.
	Host      string `json:"host"`      // Address for the device.
	Port      int    `json:"port"`      // Port the switch is running on
	Transport string `json:"transport"` // http, https , socket , http_local for goeapi
	NetworkOs string `json:"networkos"` // eos
}

## Apply the CRD
kubectl apply -f config/crd/bases/network.networkconfigoperator_netdevs.yaml
customresourcedefinition.apiextensions.k8s.io/netdevs.network.networkconfigoperator configured

## Apply the CR for leaf1

kubectl apply -f config/samples/ceos1.yaml
netdevs.network.networkconfigoperator/netdevs-sample created