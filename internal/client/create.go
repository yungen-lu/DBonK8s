package client

import (
	"context"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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

func buildMysqlDeployment(dbname, username, password string) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: dbname,
			Labels: map[string]string{
				"app":    dbname,
				"dbtype": "mysql",
				"tag":    "linebot",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":    dbname,
					"dbtype": "mysql",
					"tag":    "linebot",
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RecreateDeploymentStrategyType,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":    dbname,
						"dbtype": "mysql",
						"tag":    "linebot",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  dbname,
							Image: "mysql:8",
							Ports: []apiv1.ContainerPort{
								{
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 3306,
								},
							},
							Env: []apiv1.EnvVar{
								{
									Name:  "MYSQL_USER",
									Value: username,
								},
								{
									Name:  "MYSQL_PASSWORD",
									Value: password,
								},
								{
									Name:  "MYSQL_ROOT_PASSWORD",
									Value: password,
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
}

func buildMysqlService(dbname string) *apiv1.Service {
	return &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: dbname,
			Labels: map[string]string{
				"app":    dbname,
				"dbtype": "mysql",
				"tag":    "linebot",
			},
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"app":    dbname,
				"dbtype": "mysql",
				"tag":    "linebot",
			},
			// Type: apiv1.ServiceTypeLoadBalancer,
			Type: apiv1.ServiceTypeNodePort,
			Ports: []apiv1.ServicePort{
				{
					Protocol: apiv1.ProtocolTCP,
					Port:     3306,
					// TargetPort: intstr.FromInt(5432),
				},
			},
		},
	}

}

func buildPostgresDeployment(dbname, username, password string) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: dbname,
			Labels: map[string]string{
				"app":    dbname,
				"dbtype": "postgres",
				"tag":    "linebot",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":    dbname,
					"dbtype": "postgres",
					"tag":    "linebot",
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RecreateDeploymentStrategyType,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":    dbname,
						"dbtype": "postgres",
						"tag":    "linebot",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  dbname,
							Image: "postgres:alpine",
							Ports: []apiv1.ContainerPort{
								{
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 5432,
								},
							},
							Env: []apiv1.EnvVar{
								{
									Name:  "POSTGRES_USER",
									Value: username,
								},
								{
									Name:  "POSTGRES_PASSWORD",
									Value: password,
								},
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      "postgres-emptydir",
									MountPath: "/var/lib/postgresql/data",
								},
							},
						},
					},
					Volumes: []apiv1.Volume{
						{
							Name: "postgres-emptydir",
							VolumeSource: apiv1.VolumeSource{
								EmptyDir: &apiv1.EmptyDirVolumeSource{},
							},
						},
					},
				},
			},
		},
	}
}

func buildPostgresService(dbname string) *apiv1.Service {
	return &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: dbname,
			Labels: map[string]string{
				"app":    dbname,
				"dbtype": "postgres",
				"tag":    "linebot",
			},
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"app":    dbname,
				"dbtype": "postgres",
				"tag":    "linebot",
			},
			// Type: apiv1.ServiceTypeLoadBalancer,
			Type: apiv1.ServiceTypeNodePort,
			Ports: []apiv1.ServicePort{
				{
					Protocol: apiv1.ProtocolTCP,
					Port:     5432,
					// TargetPort: intstr.FromInt(5432),
				},
			},
		},
	}
}

func buildMongodbDeployment(dbname, username, password string) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: dbname,
			Labels: map[string]string{
				"app":    dbname,
				"dbtype": "mongodb",
				"tag":    "linebot",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":    dbname,
					"dbtype": "mongodb",
					"tag":    "linebot",
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RecreateDeploymentStrategyType,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":    dbname,
						"dbtype": "mongodb",
						"tag":    "linebot",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  dbname,
							Image: "mongo:6",
							Ports: []apiv1.ContainerPort{
								{
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 27017,
								},
							},
							Env: []apiv1.EnvVar{
								{
									Name:  "MONGO_INITDB_ROOT_USERNAME",
									Value: username,
								},
								{
									Name:  "MONGO_INITDB_ROOT_PASSWORD",
									Value: password,
								},
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      "mongodb-emptydir",
									MountPath: "/data/db",
								},
							},
						},
					},
					Volumes: []apiv1.Volume{
						{
							Name: "mongodb-emptydir",
							VolumeSource: apiv1.VolumeSource{
								EmptyDir: &apiv1.EmptyDirVolumeSource{},
							},
						},
					},
				},
			},
		},
	}
}

func buildMongodbService(dbname string) *apiv1.Service {
	return &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: dbname,
			Labels: map[string]string{
				"app":    dbname,
				"dbtype": "mongodb",
				"tag":    "linebot",
			},
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"app":    dbname,
				"dbtype": "mongodb",
				"tag":    "linebot",
			},
			// Type: apiv1.ServiceTypeLoadBalancer,
			Type: apiv1.ServiceTypeNodePort,
			Ports: []apiv1.ServicePort{
				{
					Protocol: apiv1.ProtocolTCP,
					Port:     27017,
					// TargetPort: intstr.FromInt(5432),
				},
			},
		},
	}
}

func buildRedisDeployment(dbname, username, password string) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: dbname,
			Labels: map[string]string{
				"app":    dbname,
				"dbtype": "redis",
				"tag":    "linebot",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":    dbname,
					"dbtype": "redis",
					"tag":    "linebot",
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RecreateDeploymentStrategyType,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":    dbname,
						"dbtype": "redis",
						"tag":    "linebot",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  dbname,
							Image: "bitnami/redis:7.0",
							Ports: []apiv1.ContainerPort{
								{
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 6379,
								},
							},
							Env: []apiv1.EnvVar{
								{
									Name:  "REDIS_USERNAME",
									Value: username,
								},
								{
									Name:  "REDIS_PASSWORD",
									Value: password,
								},
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      "redis-emptydir",
									MountPath: "/bitnami/redis/data",
								},
							},
						},
					},
					Volumes: []apiv1.Volume{
						{
							Name: "redis-emptydir",
							VolumeSource: apiv1.VolumeSource{
								EmptyDir: &apiv1.EmptyDirVolumeSource{},
							},
						},
					},
				},
			},
		},
	}
}

func buildRedisService(dbname string) *apiv1.Service {
	return &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: dbname,
			Labels: map[string]string{
				"app":    dbname,
				"dbtype": "redis",
				"tag":    "linebot",
			},
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"app":    dbname,
				"dbtype": "redis",
				"tag":    "linebot",
			},
			// Type: apiv1.ServiceTypeLoadBalancer,
			Type: apiv1.ServiceTypeNodePort,
			Ports: []apiv1.ServicePort{
				{
					Protocol: apiv1.ProtocolTCP,
					Port:     6379,
					// TargetPort: intstr.FromInt(5432),
				},
			},
		},
	}
}
