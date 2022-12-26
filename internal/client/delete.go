package client

import (
	"context"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8sClient) DeleteService(ctx context.Context, namespace, dbname string) error {
	namespace = strings.ToLower(namespace)
	serviceClient := k.clientset.CoreV1().Services(namespace)
	deletePolicy := metav1.DeletePropagationForeground
	return serviceClient.Delete(ctx, dbname, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}

func (k *K8sClient) Delete(ctx context.Context, namespace, dbname string) error {
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
