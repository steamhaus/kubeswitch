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
```
$ kubeswitch
  Latest stable release is: v1.18.5
  Do you want to install this version? [yes/no]  // Yes installs stable, no gives you the list of available options

$ kubeswitch -v v1.x.z  // this installs your selected version, currently, if you get permission denied, the version you requested doesn't exist

```




[![asciicast](https://asciinema.org/a/qHmIcoVAScse9o0sGn2BFDZ11.svg)](https://asciinema.org/a/qHmIcoVAScse9o0sGn2BFDZ11)

##### *Linux support added, freebsd support can be added if needs be*
