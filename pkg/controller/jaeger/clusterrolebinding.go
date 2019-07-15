package jaeger

import (
	"context"

	log "github.com/sirupsen/logrus"
	rbac "k8s.io/api/rbac/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "github.com/jaegertracing/jaeger-operator/pkg/apis/jaegertracing/v1"
	"github.com/jaegertracing/jaeger-operator/pkg/inventory"
)

func (r *ReconcileJaeger) applyClusterRoleBindingBindings(jaeger v1.Jaeger, desired []rbac.ClusterRoleBinding) error {
	opts := client.MatchingLabels(map[string]string{
		"app.kubernetes.io/instance":   jaeger.Name,
		"app.kubernetes.io/managed-by": "jaeger-operator",
	})
	list := &rbac.ClusterRoleBindingList{}
	if err := r.client.List(context.Background(), opts, list); err != nil {
		return err
	}

	inv := inventory.ForClusterRoleBindings(list.Items, desired)
	for _, d := range inv.Create {
		jaeger.Logger().WithFields(log.Fields{
			"clusteRoleBinding": d.Name,
			"namespace":         d.Namespace,
		}).Debug("creating cluster role binding")
		if err := r.client.Create(context.Background(), &d); err != nil {
			return err
		}
	}

	for _, d := range inv.Update {
		jaeger.Logger().WithFields(log.Fields{
			"clusteRoleBinding": d.Name,
			"namespace":         d.Namespace,
		}).Debug("updating cluster role binding")
		if err := r.client.Update(context.Background(), &d); err != nil {
			return err
		}
	}

	for _, d := range inv.Delete {
		jaeger.Logger().WithFields(log.Fields{
			"clusteRoleBinding": d.Name,
			"namespace":         d.Namespace,
		}).Debug("deleting cluster role binding")
		if err := r.client.Delete(context.Background(), &d); err != nil {
			return err
		}
	}

	return nil
}
