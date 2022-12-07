package v2alpha1

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
)

func addConversionFuncs(scheme *runtime.Scheme) error {
	var err error
	// Add field label conversions for kinds having selectable nothing but ObjectMeta fields.
	for _, k := range []string{"Job", "JobTemplate", "CronJob"} {
		kind := k // don't close over range variables
		err = scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind(kind),
			func(label, value string) (string, string, error) {
				switch label {
				case "metadata.name", "metadata.namespace", "status.successful":
					return label, value, nil
				default:
					return "", "", fmt.Errorf("field label %q not supported for %q", label, kind)
				}
			})
		if err != nil {
			return err
		}
	}
	return nil
}
