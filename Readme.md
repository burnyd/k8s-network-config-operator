# Kubernetes Network operator
Will add more to the readme later :D

### Operations

The idea for this repo is that a user would enter a full configuration within the runningconfig: key within the CR that covers each switch. For example, checking manifests/ceos/ceos1.yaml

Logging the operator.

```
2022-01-26T02:52:38.298Z        INFO    Checking on Device 192.168.4.244
2022-01-26T02:52:38.372Z        INFO    Switch config matches CR
2022-01-26T02:52:38.383Z        INFO    Starting to check on the device
2022-01-26T02:52:38.383Z        INFO    Device is EOS
2022-01-26T02:52:38.383Z        INFO    Checking on Device 192.168.4.244
2022-01-26T02:52:38.452Z        INFO    Switch config matches CR
```