rm testbase64.ign
rm testremote.ign
./igntool a testbase64.ign igntool
./igntool ar testremote.ign https://dl.k8s.io/v1.13.7/kubernetes-node-windows-amd64.tar.gz /k/kube1.13.7.tar.gz

