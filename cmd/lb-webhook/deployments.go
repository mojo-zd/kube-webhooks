package main

import (
	"encoding/json"
	"fmt"

	"k8s.io/api/admission/v1beta1"
	"k8s.io/api/apps/v1"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
	"gitee.com/wisecloud/kube-webhooks/pkg"
)

const (
	WISE2C_LB_LABEL = "io.wise2c.service.type"
)

// 针对lb类型的deployment做resource quota默认配置
func mutateDeploy(ar v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	klog.V(2).Info("mutating pods")
	deployResource := metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "deployments"}

	if ar.Request.Resource.Resource != deployResource.Resource {
		klog.Errorf("expect resource to be %s, but resource is %s", deployResource, ar.Request.Resource)
		return nil
	}

	raw := ar.Request.Object.Raw
	deployment := v1.Deployment{}
	deserializer := pkg.Codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(raw, nil, &deployment); err != nil {
		klog.Error(err)
		return toAdmissionResponse(err)
	}

	reviewResponse := v1beta1.AdmissionResponse{}
	if _, ok := deployment.Labels[WISE2C_LB_LABEL]; !ok {
		reviewResponse.Allowed = true
		return &reviewResponse
	}

	quotaOpt, err := getQuotaOptions(deployment.Namespace)
	if err != nil {
		klog.Error(err)
		return toAdmissionResponse(err)
	}

	// pass if quota not set
	if len(quotaOpt) == 0 {
		reviewResponse.Allowed = true
		return &reviewResponse
	}

	bytes, needPatch := generateResourcePatch(deployment, quotaOpt)
	if !needPatch {
		reviewResponse.Allowed = true
		return &reviewResponse
	}

	reviewResponse.Allowed = true
	reviewResponse.Patch = bytes
	pt := v1beta1.PatchTypeJSONPatch
	reviewResponse.PatchType = &pt

	return &reviewResponse
}

func generateResourcePatch(deploy v1.Deployment, opt map[string]bool) ([]byte, bool) {
	patches := []*pkg.Patch{}
	needPatch := false
	for i, container := range deploy.Spec.Template.Spec.Containers {
		if container.Resources.Limits != nil {
			continue
		}

		needPatch = true
		limitOpt := map[string]interface{}{}
		for k := range opt {
			switch k {
			case string(v12.ResourceMemory):
				limitOpt[k] = memory
			case string(v12.ResourceCPU):
				limitOpt[k] = cpu
			}
		}

		patches = append(patches, &pkg.Patch{
			OP:   "add",
			Path: fmt.Sprintf("/spec/template/spec/containers/%d/resources", i),
			Value: map[string]interface{}{
				"limits": limitOpt,
			},
		})
	}
	bytes, _ := json.Marshal(patches)
	return bytes, needPatch
}

func getQuotaOptions(namespace string) (map[string]bool, error) {
	quotaOpt := map[string]bool{}
	quotas, err := clientSet.CoreV1().ResourceQuotas(namespace).List(metav1.ListOptions{})
	if err != nil {
		klog.Errorf("resource quota list failed, err:%s", err.Error())
		return nil, err
	}

	if len(quotas.Items) == 0 {
		return map[string]bool{}, nil
	}

	for _, quota := range quotas.Items {
		cpu := quota.Spec.Hard[v12.ResourceLimitsCPU]
		if cpu.String() != "0" && len(cpu.String()) > 0 {
			quotaOpt[string(v12.ResourceCPU)] = true
		}
		memory := quota.Spec.Hard[v12.ResourceLimitsMemory]
		if memory.String() != "0" && len(memory.String()) > 0 {
			quotaOpt[string(v12.ResourceMemory)] = true
		}
	}

	return quotaOpt, nil
}
