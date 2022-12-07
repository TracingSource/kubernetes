package certificates

import certificates "k8s.io/api/certificates/v1beta1"

// IsCertificateRequestApproved returns true if a certificate request has the
// "Approved" condition and no "Denied" conditions; false otherwise.
func IsCertificateRequestApproved(csr *certificates.CertificateSigningRequest) bool {
	approved, denied := GetCertApprovalCondition(&csr.Status)
	return approved && !denied
}

func GetCertApprovalCondition(status *certificates.CertificateSigningRequestStatus) (approved bool, denied bool) {
	for _, c := range status.Conditions {
		if c.Type == certificates.CertificateApproved {
			approved = true
		}
		if c.Type == certificates.CertificateDenied {
			denied = true
		}
	}
	return
}
