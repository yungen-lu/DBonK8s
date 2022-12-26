package client

import (
	"context"
	"fmt"
	"strings"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
)

type K8sClient struct {
	clientset *kubernetes.Clientset
}

func New() *K8sClient {
	k := &K8sClient{}
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	k.clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return k
}

func (k *K8sClient) GetEndPoint(ctx context.Context, namespace, dbname string) (string, error) {
	namespace = strings.ToLower(namespace)
	nodes, err := k.clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return "", err
	}
	var ips []string
	for i := range nodes.Items {
		node := &nodes.Items[i]
		ips = append(ips, getExternalIP(node.Status.Addresses))

	}
	serviceClient := k.clientset.CoreV1().Services(namespace)
	service, err := serviceClient.Get(ctx, dbname, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	nodeport := service.Spec.Ports[0].NodePort
	iplist := ""
	for _, ip := range ips {
		iplist += fmt.Sprintf("%s:%d\n", ip, nodeport)
	}
	strings.TrimRight(iplist, "\n")
	return iplist, nil

}
func (k *K8sClient) upsertNamespace(ctx context.Context, b64namespace string) error {
	_, err := k.clientset.CoreV1().Namespaces().Get(ctx, b64namespace, metav1.GetOptions{})
	if err != nil {
		_, err := k.clientset.CoreV1().Namespaces().Create(ctx, &apiv1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: b64namespace}}, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}
	return nil

}
func getExternalIP(addrs []apiv1.NodeAddress) string {
	for _, v := range addrs {
		if v.Type == apiv1.NodeExternalIP {
			return v.Address
		}
	}
	return ""
}
