# Kubernetes Version Switcher

### Easily switch kubectl binary versions.
### Currently only works on OSX however I'm working on the Linux implementation



### To run

```
$ go build main.go

$ mv main /usr/local/bin/kubeswitch

$ kubeswitch
  Do you want to install this version?
  yes/no

OR

$ mv kubeswitch /usr/local/bin/

```

### Enhancements:

1. Get OS version and decide on which URL to download from
2. Install Helm version 
3. For helm install, add logic to check server version to match
