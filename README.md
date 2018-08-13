Website Operator
================
Inspiration:
* `Extending Kubernetes` chapter from `Kubernetes in Action` book by [Marko Luksa](https://twitter.com/markoluksa) talks about controller and with example of website controller which you can find [here](https://github.com/luksa/k8s-website-controller). 


This example is created using [operator framework sdk](https://github.com/operator-framework/operator-sdk)


Installation
---------------

* Configure RBAC
  
```
kubectl create -f deploy/rbac.yaml
```

* Install CRD

```
kubectl create -f deploy/crd.yaml
```

* Install operator

```
kubectl create -f deploy/operator.yaml
```

Example
-------

* Check sample example at `deploy/cr.yaml`

```
apiVersion: "website.example.com/v1alpha1"
kind: "Website"
metadata:
  name: "example-1"
spec:
  GitRepo: "https://github.com/surajnarwade/samplewebsite-1"
```

* Create website

```
kubectl create -f deploy/cr.yaml
```

* you can see website resource,
  
```
$ kubectl get website
NAME        AGE
example-1   1h
```

* it creates deployment and service,

```
$ kubectl get deployments
NAME               DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
example-1          1         1         1            1           1h
website-operator   1         1         1            1           2h

$ kubectl get services
NAME               TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)        AGE
example-1          NodePort    10.100.168.97    <none>        80:30157/TCP   1h
kubernetes         ClusterIP   10.96.0.1        <none>        443/TCP        2h
website-operator   ClusterIP   10.96.124.114    <none>        60000/TCP      2h
```

* Access website,

```
curl $(minikube ip):30157
```
