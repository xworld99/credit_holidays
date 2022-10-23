package consts

import "net/http"

// consts used in whole application

const (
	Host              = "0.0.0.0:8080"
	ReportLinkPattern = "http://%s/reports/%s"

	ProofedNow  = "now()"
	ProofedNull = "null"

	OrderSuccess    = "success"
	OrderDeclined   = "declined"
	OrderInProgress = "in_progress"
	OrderProof      = "proof"
	OrderDecline    = "decline"

	OperationAccrual  = "accrual"
	OperationWithdraw = "withdraw"
)

var (
	OrderActions      = map[string]bool{OrderProof: true, OrderDecline: true}
	SortingType       = map[string]bool{"created_at": true, "amount": true}
	ErrorDescriptions = map[int]string{
		http.StatusBadRequest:          "error in user request",
		http.StatusInternalServerError: "cant handle request, internal error",
		http.StatusNotFound:            "not found",
	}
)
