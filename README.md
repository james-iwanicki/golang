# golang

rest: 
	Has a client/server project to demonstrate a golang kubernetes client making REST API calls to golang kubernetes server over mTLS.

demo-istio: 
	Has some deployment files to demonstrate use of ISTIO service mesh

shell: 
	Has some two clients which are used to ping from one to another in same POD. 
	The shell-inspect container, then inspects inter-container traffic in foreign POD using namespace to get the loopback device and you can run tcpdump on it.
	This is to illustrate two containers in same POD talk over loopback and you can intercept that traffic using namespace.
