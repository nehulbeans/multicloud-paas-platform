package k8sclient

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"multicloud-paas-platform/internal/core/models"
)

type K8sDeployer interface {
	PushDeployment(ctx context.Context, kubeconfig string, req models.DeploymentRequest) (string, error)
}

type k8sDeployer struct{}

func NewK8sDeployer() K8sDeployer {
	return &k8sDeployer{}
}

func (k *k8sDeployer) PushDeployment(ctx context.Context, kubeconfig string, req models.DeploymentRequest) (string, error) {
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
	if err != nil {
		return "", fmt.Errorf("failed to parse kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", fmt.Errorf("failed to create k8s client: %v", err)
	}

	// Create Deployment
	if err := k.createDeployment(ctx, clientset, req); err != nil {
		return "", err
	}

	// Create Service
	if err := k.createService(ctx, clientset, req); err != nil {
		return "", err
	}

	// poll for the LoadBalancer IP
	serviceName := req.AppName + "-svc"
	return k.waitForLoadBalancerIP(ctx, clientset, serviceName)
}

func (k *k8sDeployer) createDeployment(ctx context.Context, clientset *kubernetes.Clientset, req models.DeploymentRequest) error {
	labels := map[string]string{"app": req.AppName}

	// Build Environment Variables
	var envVars []corev1.EnvVar
	for key, val := range req.EnvVars {
		envVars = append(envVars, corev1.EnvVar{
			Name:  key,
			Value: val,
		})
	}

	// Build Container Ports
	var containerPorts []corev1.ContainerPort
	for _, port := range req.Ports {
		containerPorts = append(containerPorts, corev1.ContainerPort{
			ContainerPort: port,
		})
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: req.AppName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &req.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  req.AppName,
							Image: req.Image,
							Ports: containerPorts,
							Env:   envVars,
						},
					},
				},
			},
		},
	}

	_, err := clientset.AppsV1().Deployments("default").Create(ctx, deployment, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create k8s deployment: %v", err)
	}
	return nil
}

func (k *k8sDeployer) createService(ctx context.Context, clientset *kubernetes.Clientset, req models.DeploymentRequest) error {
	labels := map[string]string{"app": req.AppName}

	// Map the first requested port to external port 80. Fallback to 8080 if none provided.
	targetPort := int32(8080)
	if len(req.Ports) > 0 {
		targetPort = req.Ports[0]
	}

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: req.AppName + "-svc",
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     corev1.ServiceTypeLoadBalancer,
			Ports: []corev1.ServicePort{
				{
					Port:       80,
					TargetPort: intstr.FromInt(int(targetPort)),
				},
			},
		},
	}

	_, err := clientset.CoreV1().Services("default").Create(ctx, service, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create k8s service: %v", err)
	}
	return nil
}

func (k *k8sDeployer) waitForLoadBalancerIP(ctx context.Context, clientset *kubernetes.Clientset, serviceName string) (string, error) {
	// Timeout after 2 minutes to prevent hanging API requests
	timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeoutCtx.Done():
			return "", fmt.Errorf("timed out waiting for LoadBalancer IP for service %s", serviceName)
		case <-ticker.C:
			svc, err := clientset.CoreV1().Services("default").Get(timeoutCtx, serviceName, metav1.GetOptions{})
			if err != nil {
				continue
			}

			if len(svc.Status.LoadBalancer.Ingress) > 0 {
				ingress := svc.Status.LoadBalancer.Ingress[0]
				if ingress.Hostname != "" {
					return ingress.Hostname, nil // AWS EKS returns a Hostname
				}
				if ingress.IP != "" {
					return ingress.IP, nil // GCP/Hetzner K3s returns an IP
				}
			}
		}
	}
}
