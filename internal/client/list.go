package client

import (
	"context"
	"errors"
	"strings"

	"github.com/yungen-lu/TOC-Project-2022/internal/models"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
