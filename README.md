# Kubernetes Version Switcher 

![Go](https://github.com/steamhaus/kubeswitch/workflows/Go/badge.svg?branch=master)

## Easily switch kubectl binary versions.

### Build:

```
$ go build -o kubeswitch .

$ chmod +x kubeswitch 

$ mv kubeswitch $PATH

```

### Run:
#### Yes installs stable, no gives you the list of available options
```
$ kubeswitch
  Latest stable release is: v1.18.5
  Do you want to install this version? [yes/no]
```

#### This installs your selected version, currently, if you get permission denied, the version you requested doesn't exist
```
$ kubeswitch -v/ --version v1.x.z  
```

#### Match to your client to your EKS version
```
$ kubeswitch -a / --aws 
```




[![asciicast](https://asciinema.org/a/qHmIcoVAScse9o0sGn2BFDZ11.svg)](https://asciinema.org/a/qHmIcoVAScse9o0sGn2BFDZ11)

