package stub

import (
	"context"

	"github.com/surajnarwade/website-operator/pkg/apis/website/v1alpha1"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"

)

func NewHandler() sdk.Handler {
	return &Handler{}
}

type Handler struct {
	//deploySpec *appsv1.Deployment
	//serviceSpec *corev1.Service
	// Fill me
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *v1alpha1.Website:
		deploySpec, serviceSpec := newbusyBoxPod(o)
		err := sdk.Create(&deploySpec)
		if err != nil && !errors.IsAlreadyExists(err) {
			logrus.Errorf("failed to create busybox pod : %v", err)
			return err
		}
		err = sdk.Create(&serviceSpec)
		if err != nil && !errors.IsAlreadyExists(err) {
			logrus.Errorf("failed to create busybox pod : %v", err)
			return err
		}
	}
	return nil
}

// newbusyBoxPod demonstrates how to create a busybox pod
func newbusyBoxPod(cr *v1alpha1.Website) (appsv1.Deployment, corev1.Service) {
	labels := map[string]string{
		"app": "busy-box",
	}

	webDeployment := appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    "Website",
				}),
			},
			Labels: labels,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: cr.Name,
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						{
							Name: "git-clone",
							Image: "alpine/git",
							Args: []string{
								 "clone",
								 "--single-branch",
								 "--",
								  cr.Spec.GitRepo,
								 "/repo",
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name: "workdir",
									MountPath: "/repo",
								},
							},
						},
					},
					Containers: []corev1.Container{
					{
						Name: "website",
						Image: "nginx:alpine",
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 80,
								Protocol: corev1.ProtocolTCP,
							},
						},
						VolumeMounts: []corev1.VolumeMount{
						{
							Name: "workdir",
							MountPath: "/usr/share/nginx/html",
							ReadOnly: true,
						},
						},
					},
					},
					Volumes: []corev1.Volume{
						{
							Name: "workdir",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
				},

			},

		},

	}

	webService := corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    "Website",
				}),
			},
			Labels: labels,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeNodePort,
			Ports: []corev1.ServicePort{
				{
				Port: 80,
				Protocol: corev1.ProtocolTCP,
				TargetPort: intstr.FromInt(80),
				},
			},
			Selector: labels,
		},
	}

	return webDeployment, webService
}
