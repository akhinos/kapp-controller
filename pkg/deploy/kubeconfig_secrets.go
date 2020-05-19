package deploy

import (
	"fmt"
	"strings"

	"github.com/k14s/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type KubeconfigSecrets struct {
	coreClient kubernetes.Interface
}

func NewKubeconfigSecrets(coreClient kubernetes.Interface) *KubeconfigSecrets {
	return &KubeconfigSecrets{coreClient}
}

func (s *KubeconfigSecrets) Find(nsName string, clusterOpts *v1alpha1.AppCluster) (string, error) {
	if clusterOpts == nil {
		return "", nil
	}

	if clusterOpts.KubeconfigSecretRef == nil {
		return "", fmt.Errorf("Expected kubeconfig secret reference to be specified")
	}

	secret, err := s.coreClient.CoreV1().Secrets(nsName).Get(
		clusterOpts.KubeconfigSecretRef.Name, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("Getting kubeconfig secret: %s", err)
	}

	key := clusterOpts.KubeconfigSecretRef.Key
	if len(key) == 0 {
		key = "value"
	}

	val, found := secret.Data[key]
	if !found {
		var otherKeys []string
		for otherKey, _ := range secret.Data {
			otherKeys = append(otherKeys, otherKey)
		}

		return "", fmt.Errorf("Expected to find key '%s' in secret (keys: %s)",
			key, strings.Join(otherKeys, ", "))
	}

	return string(val), nil
}
