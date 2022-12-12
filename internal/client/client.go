package client

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

type K8sClient struct {
	clientset *kubernetes.Clientset
}
type Instance struct {
	Type     string
	Name     string
	User     string
	Password string
	// Owner     string
	Namespace string
	Endpoint  string
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

func (k *K8sClient) Create(ctx context.Context, namespace, dbtype, containername, image, dbname, username, password string) error {
	deploymentsClient := k.clientset.AppsV1().Deployments(namespace)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: dbname,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":    dbname,
					"dbtype": dbtype,
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RecreateDeploymentStrategyType,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":    dbname,
						"dbtype": dbtype,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  containername,
							Image: image,
							Ports: []apiv1.ContainerPort{
								{
									// Name:          "mysql",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 3306,
								},
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      "mysql-emptydir",
									MountPath: "/var/lib/mysql",
								},
							},
						},
					},
					Volumes: []apiv1.Volume{
						{
							Name: "mysql-emptydir",
							VolumeSource: apiv1.VolumeSource{
								EmptyDir: &apiv1.EmptyDirVolumeSource{},
							},
						},
					},
				},
			},
		},
	}
	_, err := deploymentsClient.Create(ctx, deployment, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil

}
func (k *K8sClient) ListInNamespace(ctx context.Context, namespace string) ([]Instance, error) {
	deploymentsClient := k.clientset.AppsV1().Deployments(namespace)
	list, err := deploymentsClient.List(ctx, metav1.ListOptions{Limit: 12})
	if err != nil {
		return nil, err
	}
	if len(list.Items) == 0 {
		return nil, nil
	}
	instances := make([]Instance, len(list.Items))
	for i, d := range list.Items {
		instances[i].Name = d.Name
		instances[i].Type = d.Labels["dbtype"]
		instances[i].Namespace = namespace
		// d.Labels
		// d.Name
	}
	return instances, nil
}
