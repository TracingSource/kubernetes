package apps

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	imageutils "k8s.io/kubernetes/test/utils/image"
)

// NOTE(claudiub): These constants should NOT be used as Pod Container Images.
const (
	WebserverImageName = "httpd"
	AgnhostImageName   = "agnhost"
)

var (
	// CronJobGroupVersionResourceAlpha unambiguously identifies a resource of cronjob with alpha status
	CronJobGroupVersionResourceAlpha = schema.GroupVersionResource{Group: "batch", Version: "v2alpha1", Resource: "cronjobs"}

	// CronJobGroupVersionResourceBeta unambiguously identifies a resource of cronjob with beta status
	CronJobGroupVersionResourceBeta = schema.GroupVersionResource{Group: "batch", Version: "v1beta1", Resource: "cronjobs"}

	// NautilusImage is the fully qualified URI to the Nautilus image
	NautilusImage = imageutils.GetE2EImage(imageutils.Nautilus)

	// KittenImage is the fully qualified URI to the Kitten image
	KittenImage = imageutils.GetE2EImage(imageutils.Kitten)

	// WebserverImage is the fully qualified URI to the Httpd image
	WebserverImage = imageutils.GetE2EImage(imageutils.Httpd)

	// NewWebserverImage is the fully qualified URI to the HttpdNew image
	NewWebserverImage = imageutils.GetE2EImage(imageutils.HttpdNew)

	// AgnhostImage is the fully qualified URI to the Agnhost image
	AgnhostImage = imageutils.GetE2EImage(imageutils.Agnhost)
)
