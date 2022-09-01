# kubectl-low-evict-pod
  
  This tool is used for all versions of kubernetes to evict pod, especially that do not successfully evict fron its node. 
  
  ## Build

  ```
  cd  kubectl-low-evict-pod
  make clear
  make build
  ```

  ## Usage

  ### 

  ```
  $ kubectl low-evict-pod --help

      evict pod for low-version kubernetes, example 1.15.x

      Usage:
        evict pod  [flags]

      Flags:
        -c, --cordon-secords int             cordon seconds (default 1)
        -g, --grace-secords int              grace period seconds (default 120)
        -n, --namespace string               If present, the namespace scope for this CLI request
            --request-timeout string         The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests. (default "0")

```

  ###  pod evicted scenario 
  
  The running process : *cordon - delete - uncordon*

  #### before running

```
$ kubectl get pod -owide
    
      NAME                          READY   STATUS        RESTARTS   AGE   IP              NODE    NOMINATED NODE   READINESS GATES
      nginx-demo-68fc8c5cd5-s2g86   1/1     Running       0          47d   172.30.204.6    kube2   <none>           <none>

```

#### evicting

```
$ kubectl low-evict-pod nginx-demo-68fc8c5cd5-s2g86 -n default

      2022/08/17 16:01:33 cordon kube2
      2022/08/17 16:01:33 deleting pod nginx-demo-68fc8c5cd5-s2g86
      2022/08/17 16:01:34 uncordon kube2

```
#### observe the pod

```
$ kubectl get pod -owide

      NAME                          READY   STATUS        RESTARTS   AGE   IP             NODE    NOMINATED NODE   READINESS GATES
      nginx-demo-68fc8c5cd5-g4grv   1/1     Running       0          90s   172.30.204.8   kube2   <none>           <none>
      nginx-demo-68fc8c5cd5-s2g86   1/1     Terminating   0          47d   172.30.204.6   kube2   <none>           <none>

```

if you want the pod 'nginx-demo-68fc8c5cd5-s2g86 '  terminate immediately, using parameter "--grace-secords=[n second which default 120]" 
