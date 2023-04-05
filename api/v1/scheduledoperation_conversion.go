package v1

import (
	v2 "github.com/Algatux/k8s-reconcyle-tests/api/v2"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

func (src *ScheduledOperation) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v2.ScheduledOperation)

	dst.Status.Test = src.Spec.Test

	return nil
}

func (dst *ScheduledOperation) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v2.ScheduledOperation)

	dst.Spec.Test = src.Status.Test

	return nil
}
