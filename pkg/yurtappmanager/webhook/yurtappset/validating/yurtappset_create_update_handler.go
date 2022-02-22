/*
Copyright 2020 The OpenYurt Authors.
Copyright 2019 The Kruise Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package validating

import (
	"context"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	unitv1alpha1 "github.com/openyurtio/yurt-app-manager/pkg/yurtappmanager/apis/apps/v1alpha1"
	webhookutil "github.com/openyurtio/yurt-app-manager/pkg/yurtappmanager/webhook/util"
)

// YurtAppSetCreateUpdateHandler handles YurtAppSet
type YurtAppSetCreateUpdateHandler struct {
	// To use the client, you need to do the following:
	// - uncomment it
	// - import sigs.k8s.io/controller-runtime/pkg/client
	// - uncomment the InjectClient method at the bottom of this file.
	Client client.Client

	// Decoder decodes objects
	Decoder *admission.Decoder
}

var _ webhookutil.Handler = &YurtAppSetCreateUpdateHandler{}

func (h *YurtAppSetCreateUpdateHandler) SetOptions(options webhookutil.Options) {
	h.Client = options.Client
	return
}

// Handle handles admission requests.
func (h *YurtAppSetCreateUpdateHandler) Handle(ctx context.Context, req admission.Request) admission.Response {
	obj := &unitv1alpha1.YurtAppSet{}

	err := h.Decoder.Decode(req, obj)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	switch req.AdmissionRequest.Operation {
	case admissionv1.Create:
		if allErrs := validateYurtAppSet(h.Client, obj); len(allErrs) > 0 {
			return admission.Errored(http.StatusUnprocessableEntity, allErrs.ToAggregate())
		}
	case admissionv1.Update:
		oldObj := &unitv1alpha1.YurtAppSet{}
		if err := h.Decoder.DecodeRaw(req.AdmissionRequest.OldObject, oldObj); err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		}

		validationErrorList := validateYurtAppSet(h.Client, obj)
		updateErrorList := ValidateYurtAppSetUpdate(obj, oldObj)
		if allErrs := append(validationErrorList, updateErrorList...); len(allErrs) > 0 {
			return admission.Errored(http.StatusUnprocessableEntity, allErrs.ToAggregate())
		}
	}

	return admission.ValidationResponse(true, "")
}

var _ admission.DecoderInjector = &YurtAppSetCreateUpdateHandler{}

// InjectDecoder injects the decoder into the YurtAppSetCreateUpdateHandler
func (h *YurtAppSetCreateUpdateHandler) InjectDecoder(d *admission.Decoder) error {
	h.Decoder = d
	return nil
}
