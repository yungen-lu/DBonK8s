package client

import (
	"context"
	"fmt"
	"strings"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/apimachinery/pkg/util/intstr"
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
	// b64namespace := b64.StdEncoding.EncodeToString([]byte(namespace))
	namespace = strings.ToLower(namespace)
	nodes, err := k.clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return "", err
	}
	var ips []string
	for _, node := range nodes.Items {
		ips = append(append(ips, getExternalIP(node.Status.Addresses)), "\n")
	}
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
		iplist += fmt.Sprintf("%s:%d", ip, nodeport)
	}
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

// func findUserPassword(dbtype string, env []apiv1.EnvVar) (string, string) {
// 	switch dbtype {
// 	case "postgres":
// 		// env[0].Name = "POSTGRES_USER"
// 		// env[1].Name = "POSTGRES_PASSWORD"
// 		return env[0].Value, env[1].Value
// 	case "mysql":
// 		// env[0].Name = "MYSQL_ROOT_PASSWORD"
// 		// env[1].Name = "MYSQL_USER"
// 		// env[2].Name = "MYSQL_PASSWORD"
// 		return env[0].Value, env[1].Value
// 	case "redis":
// 		// env[0].Name = "REDIS_PASSWORD"
// 		return "default", env[0].Value
// 	case "mongodb":
// 		// env[0].Name = "MONGO_INITDB_ROOT_USERNAME"
// 		// env[1].Name = "MONGO_INITDB_ROOT_PASSWORD"
// 		return env[0].Value, env[1].Value
// 	default:
// 		return env[0].Value, env[1].Value
// 	}
// }

// func findPassword(dbtype string, env []apiv1.EnvVar) string {
// 	switch dbtype {
// 	case "postgres":
// 		env[0].Name = "POSTGRES_USER"
// 		env[1].Name = "POSTGRES_PASSWORD"
// 	case "mysql":
// 		env = make([]apiv1.EnvVar, 3)
// 		env[0].Name = "MYSQL_ROOT_PASSWORD"
// 		env[1].Name = "MYSQL_USER"
// 		env[2].Name = "MYSQL_PASSWORD"
// 	case "redis":
// 		env = make([]apiv1.EnvVar, 1)
// 		env[0].Name = "REDIS_PASSWORD"
// 	case "mongodb":
// 		env = make([]apiv1.EnvVar, 2)
// 		env[0].Name = "MONGO_INITDB_ROOT_USERNAME"
// 		env[1].Name = "MONGO_INITDB_ROOT_PASSWORD"
// 	default:
// 		env = make([]apiv1.EnvVar, 2)

//		}
//		for _, e := range env {
//			if e.Name == "POSTGRES_PASSWORD" {
//				return e.Value
//			}
//		}
//		return "default"
//	}
// func setDBEnv(dbtype, username, password string) []apiv1.EnvVar {
// 	var env []apiv1.EnvVar
// 	switch dbtype {
// 	case "postgres":
// 		env = make([]apiv1.EnvVar, 2)
// 		env[0].Name = "POSTGRES_USER"
// 		env[0].Value = username
// 		env[1].Name = "POSTGRES_PASSWORD"
// 		env[1].Value = password
// 	case "mysql":
// 		env = make([]apiv1.EnvVar, 3)
// 		env[0].Name = "MYSQL_ROOT_PASSWORD"
// 		env[0].Value = password
// 		env[1].Name = "MYSQL_USER"
// 		env[1].Value = username
// 		env[2].Name = "MYSQL_PASSWORD"
// 		env[2].Value = password
// 	case "redis":
// 		env = make([]apiv1.EnvVar, 1)
// 		env[0].Name = "REDIS_PASSWORD"
// 		env[0].Value = password
// 	case "mongodb":
// 		env = make([]apiv1.EnvVar, 2)
// 		env[0].Name = "MONGO_INITDB_ROOT_USERNAME"
// 		env[1].Name = "MONGO_INITDB_ROOT_PASSWORD"
// 	default:
// 		env = make([]apiv1.EnvVar, 2)
// 		env[0].Name = "POSTGRES_USER"
// 		env[0].Value = username
// 		env[1].Name = "POSTGRES_PASSWORD"
// 		env[1].Value = password

// 	}
// 	return env
// }

//	func (k *K8sClient) ListPodInNamespace(ctx context.Context, namespace string) ([]models.Instance, error) {
//		b64namespace := b64.StdEncoding.EncodeToString([]byte(namespace))
//		b64namespace = strings.ToLower(b64namespace)
//		podClient := k.clientset.CoreV1().Pods(b64namespace)
//		// filter := &metav1.LabelSelector{}
//		// filter = metav1.AddLabelToSelector(filter, "", "")
//		pod, err := podClient.List(ctx, metav1.ListOptions{})
//		if err != nil {
//			return nil, err
//		}
//		if len(pod.Items) == 0 {
//			return nil, nil
//		}
//		instances := make([]models.Instance, len(pod.Items))
//		for i, d := range pod.Items {
//			instances[i].Type = d.Labels["dbtype"]
//			instances[i].Name = d.Name
//			instances[i].User = "default"
//			instances[i].Password = findPassword(instances[i].Type, d.Spec.Containers[0].Env)
//			instances[i].Namespace = namespace
//			instances[i].Endpoint = "default"
//		}
//		return instances, nil
//	}
