# kssh

A utility for generating ssh config files for every node of a cluster.

## Usage

1. Set `KUBECONFIG` so that you get the desired nodes when you run `$ kubectl get nodes`.

2. Add `Include config.d/*` to the top of your `~/.ssh/config`.

3. `$ mkdir ~/.ssh/config.d`

4. `$ kssh --bastion bastionIpAddress > ~/.ssh/config.d/clustername`

5. `$ chmod 600 ~/.ssh/config.d/*`

To get the list of host names to connect to:

`$ kssh --nodes`

Then try sshing into the node:

```
$ ssh ip-10-16-111-77.us-west-1.compute.internal
Container Linux by CoreOS stable (1520.9.0)
Update Strategy: No Reboots
core@ip-10-16-111-77 ~ $ 
```

## Usage with `cssh`

`kssh` can output valid `cssh` config files for easy cluster ssh access.

`cssh` is in `yum/dnf` and `apt` repos as `clusterssh`.

1. `$ cssh -u # set up cssh`
2. `$ kssh -cssh clustername > ~/.clusterssh/clustername`
3. `$ vi ~/.clusterssh/config # edit extra_cluster_file to add ~/.clusterssh/clustername`

An example `extra_cluster_file` line:

`extra_cluster_file=~/.clusterssh/prod-west,~/.clusterssh/prod-east`
