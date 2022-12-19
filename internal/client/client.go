package client

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/yungen-lu/TOC-Project-2022/internal/models"
	appsv1 "k8s.io/api/apps/v1"
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
func (k *K8sClient) GetEndPoint(ctx context.Context, namespace, dbname string) (string, error) {
	// b64namespace := b64.StdEncoding.EncodeToString([]byte(namespace))
	namespace = strings.ToLower(namespace)
	nodes, err := k.clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return "", err
	}
	var ips []string
	for _, node := range nodes.Items {
		ips = append(ips, getExternalIP(node.Status.Addresses))
	}
	serviceClient := k.clientset.CoreV1().Services(namespace)
	service, err := serviceClient.Get(ctx, dbname, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	nodeport := service.Spec.Ports[0].NodePort
	return fmt.Sprintf("%s:%d", ips[0], nodeport), nil

}

func (k *K8sClient) DeleteService(ctx context.Context, namespace, dbname string) error {
	// b64namespace := b64.StdEncoding.EncodeToString([]byte(namespace))
	namespace = strings.ToLower(namespace)
	serviceClient := k.clientset.CoreV1().Services(namespace)
	// _, err := serviceClient.Create(ctx, service, metav1.CreateOptions{})
	deletePolicy := metav1.DeletePropagationForeground
	return serviceClient.Delete(ctx, dbname, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}

func (k *K8sClient) Create(ctx context.Context, namespace, dbtype, dbname, username, password string) error {
	// b64namespace := b64.StdEncoding.EncodeToString([]byte(namespace))
	namespace = strings.ToLower(namespace)
	err := k.upsertNamespace(ctx, namespace)
	if err != nil {
		return err
	}
	deploymentsClient := k.clientset.AppsV1().Deployments(namespace)
	serviceClient := k.clientset.CoreV1().Services(namespace)
	var deployment *appsv1.Deployment
	var service *apiv1.Service
	switch dbtype {
	case "postgres":
		deployment = buildPostgresDeployment(dbname, username, password)
		service = buildPostgresService(dbname)
	case "mysql":
		deployment = buildMysqlDeployment(dbname, username, password)
		service = buildMysqlService(dbname)
	case "redis":
		deployment = buildRedisDeployment(dbname, username, password)
		service = buildRedisService(dbname)
	case "mongodb":
		deployment = buildMongodbDeployment(dbname, username, password)
		service = buildMongodbService(dbname)
	default:
		deployment = buildPostgresDeployment(dbname, username, password)
		service = buildPostgresService(dbname)
	}

	_, err = deploymentsClient.Create(ctx, deployment, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	_, err = serviceClient.Create(ctx, service, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil

}
func (k *K8sClient) Delete(ctx context.Context, namespace, dbname string) error {
	// b64namespace := b64.StdEncoding.EncodeToString([]byte(namespace))
	namespace = strings.ToLower(namespace)
	deploymentsClient := k.clientset.AppsV1().Deployments(namespace)
	deletePolicy := metav1.DeletePropagationForeground
	err1 := deploymentsClient.Delete(ctx, dbname, metav1.DeleteOptions{PropagationPolicy: &deletePolicy}) // TODO
	err2 := k.DeleteService(ctx, namespace, dbname)
	if err1 != nil && err2 != nil {
		return fmt.Errorf("%s\n%s", err1.Error(), err2.Error())
	} else if err1 != nil {
		return err1
	} else if err2 != nil {
		return err2
	}
	return nil
}

func (k *K8sClient) ListInNamespace(ctx context.Context, namespace string) ([]models.Instance, error) {
	// b64namespace := b64.StdEncoding.EncodeToString([]byte(namespace))
	namespace = strings.ToLower(namespace)
	deploymentsClient := k.clientset.AppsV1().Deployments(namespace)
	filter := &metav1.LabelSelector{}
	filter = metav1.AddLabelToSelector(filter, "tag", "linebot")
	list, err := deploymentsClient.List(ctx, metav1.ListOptions{Limit: 12, LabelSelector: metav1.FormatLabelSelector(filter)})
	if err != nil {
		return nil, err
	}
	if len(list.Items) == 0 {
		return nil, nil
	}
	instances := make([]models.Instance, len(list.Items))
	for i, d := range list.Items {
		instances[i].Type = d.Labels["dbtype"]
		instances[i].Name = d.Name
		instances[i].Namespace = d.Namespace
	}
	return instances, nil
}

func (k *K8sClient) GetPodInNamespace(ctx context.Context, namespace, dbname string) (*models.Instance, error) {
	// b64namespace := b64.StdEncoding.EncodeToString([]byte(namespace))
	namespace = strings.ToLower(namespace)
	podClient := k.clientset.CoreV1().Pods(namespace)
	filter := &metav1.LabelSelector{}
	filter = metav1.AddLabelToSelector(filter, "app", dbname)
	filter = metav1.AddLabelToSelector(filter, "tag", "linebot")
	// println(filter.String())
	pod, err := podClient.List(ctx, metav1.ListOptions{Limit: 12, LabelSelector: metav1.FormatLabelSelector(filter)})
	if err != nil {
		return nil, err
	}
	if len(pod.Items) < 1 {
		return nil, errors.New("no instances running")
	}
	tmp := pod.Items[0]
	dbtype := tmp.Labels["dbtype"]
	// username, password := findUserPassword(dbtype, tmp.Spec.Containers[0].Env)
	username, password := tmp.Spec.Containers[0].Env[0].Value, tmp.Spec.Containers[0].Env[1].Value
	instance := &models.Instance{
		Type:      dbtype,
		Name:      tmp.Name,
		User:      username,
		Password:  password,
		Namespace: tmp.Namespace,
	}
	end, err := k.GetEndPoint(ctx, namespace, dbname)
	if err != nil {
		return nil, err
	}
	instance.Endpoint = end
	return instance, nil
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
