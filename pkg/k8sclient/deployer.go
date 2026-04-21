package k8sclient

import (
	"context"
	"fmt"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"multicloud-paas-platform/internal/core/models"
)

type K8sDeployer interface {
	PushDeployment(ctx context.Context, kubeconfig string, req models.DeploymentRequest, customHost string) (string, error)
}

type k8sDeployer struct{}

func NewK8sDeployer() K8sDeployer {
	return &k8sDeployer{}
}

func (k *k8sDeployer) PushDeployment(ctx context.Context, kubeconfig string, req models.DeploymentRequest, customHost string) (string, error) {
	// Initialize the Kubernetes client (using your existing helper method)
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
	if err != nil {
		return "", fmt.Errorf("failed to parse kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", fmt.Errorf("failed to create k8s clientset: %v", err)
	}
	replicas := int32(req.Replicas)

	// ==========================================
	// 1. CREATE DEPLOYMENT
	// ==========================================
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.AppName,
			Namespace: "default",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": req.AppName},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": req.AppName},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  req.AppName,
							Image: req.Image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80, // Assuming standard web apps listen on 80
								},
							},
						},
					},
				},
			},
		},
	}

	_, err = clientset.AppsV1().Deployments("default").Create(ctx, deployment, metav1.CreateOptions{})
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return "", fmt.Errorf("failed to create deployment: %v", err)
	}

	// ==========================================
	// 2. CREATE SERVICE (ClusterIP instead of NodePort)
	// ==========================================
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.AppName + "-service",
			Namespace: "default",
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP, // Keeps traffic inside the cluster!
			Selector: map[string]string{
				"app": req.AppName, // Must match the Deployment labels
			},
			Ports: []corev1.ServicePort{
				{
					Port:       80,                 // Port exposed to Traefik
					TargetPort: intstr.FromInt(80), // Port exposed by your container
				},
			},
		},
	}

	_, err = clientset.CoreV1().Services("default").Create(ctx, service, metav1.CreateOptions{})
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return "", fmt.Errorf("failed to create service: %v", err)
	}

	// ==========================================
	// 3. CREATE INGRESS (Traefik Routing)
	// ==========================================
	pathType := networkingv1.PathTypePrefix
	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.AppName + "-ingress",
			Namespace: "default",
			Annotations: map[string]string{
				// Tell K3s to use its built-in Traefik router
				"kubernetes.io/ingress.class": "traefik",
			},
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: customHost, // This is the dynamic URL from service.go!
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathType,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: req.AppName + "-service", // Matches Service above
											Port: networkingv1.ServiceBackendPort{
												Number: 80, // Matches Service port above
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	_, err = clientset.NetworkingV1().Ingresses("default").Create(ctx, ingress, metav1.CreateOptions{})
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return "", fmt.Errorf("failed to create ingress: %v", err)
	}

	// ==========================================
	// 4. FETCH AND RETURN THE NODE IP
	// ==========================================
	// Retrieve the IP address of the first node in the cluster
	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil || len(nodes.Items) == 0 {
		return "", fmt.Errorf("failed to get cluster nodes to extract IP: %v", err)
	}

	var nodeIP string
	// Try to get ExternalIP first (if configured)
	for _, addr := range nodes.Items[0].Status.Addresses {
		if addr.Type == corev1.NodeExternalIP {
			nodeIP = addr.Address
			break
		}
	}
	// Fallback to InternalIP if no ExternalIP is found
	if nodeIP == "" {
		for _, addr := range nodes.Items[0].Status.Addresses {
			if addr.Type == corev1.NodeInternalIP {
				nodeIP = addr.Address
				break
			}
		}
	}

	// Return just the clean IP (e.g., "80.225.211.127"). No ports!
	return nodeIP, nil
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
			Type:     corev1.ServiceTypeNodePort,
			Ports: []corev1.ServicePort{
				{
					Port:       80,
					TargetPort: intstr.FromInt(int(targetPort)),
				},
			},
		},
	}

	_, err := clientset.CoreV1().Services("default").Create(ctx, service, metav1.CreateOptions{})
	return err
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
