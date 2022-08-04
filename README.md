# kubectl-low-evict-pod
This tool is used for all versions of kubernetes to evict pod, especially for low versions of kubernetes that do not support evict interface. The main mechanism of the scheme is cordon delete uncordon. The specific parameters are "--grace-orders --cordon-orders", which can be viewed through "kubectl low evict pod --help"
