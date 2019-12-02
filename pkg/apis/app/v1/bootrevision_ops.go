package v1

import (
	"fmt"
	"github.com/logancloud/logan-app-operator/pkg/logan"
	"github.com/logancloud/logan-app-operator/pkg/logan/util/hash"
	"github.com/logancloud/logan-app-operator/pkg/logan/util/keys"
	"hash/fnv"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"strconv"
)

// GetRevisionId return bootrevision's ID
func (in *BootRevision) GetRevisionId() int {
	if in.Annotations == nil {
		return -1
	}

	if idStr, found := in.Annotations[keys.BootRevisionIdAnnotationKey]; found {
		id, err := strconv.Atoi(idStr)
		if err == nil {
			return id
		}
	}
	return -1
}

// BootHash returns a hash value calculated from BootRevision
// from https://github.com/kubernetes/kubernetes/blob/28e800245e910b65b56548f36172ce525a554dc8/pkg/controller/controller_utils.go#L1027
func (in *BootRevision) BootHash() string {
	revisionCopy := in.DeepCopy()
	revisionCopy.Spec.Env = cleanEnv(revisionCopy.Spec.Env)

	bootTemplateSpecHasher := fnv.New32a()
	hash.DeepHashObject(bootTemplateSpecHasher, *revisionCopy)
	return rand.SafeEncodeString(fmt.Sprint(bootTemplateSpecHasher.Sum32()))
}

// cleanEnv will remove the env's name in logan.BizEnvs
func cleanEnv(envs []corev1.EnvVar) []corev1.EnvVar {
	ret := make([]corev1.EnvVar, 0)
	if envs != nil {
		for _, env := range envs {
			if _, found := logan.BizEnvs[env.Name]; !found {
				ret = append(ret, env)
			}
		}
	}
	return ret
}

// SelectLatestRevision will return the latest revision
func (in *BootRevisionList) SelectLatestRevision() *BootRevision {
	if in.Items == nil || len(in.Items) == 0 {
		return nil
	}

	max := 1
	index := 0
	for current, item := range in.Items {
		revisionId := (&item).GetRevisionId()
		if revisionId >= max {
			max = revisionId
			index = current
		}
	}

	return &in.Items[index]
}
