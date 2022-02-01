package ctxkeys

type ContextKey string
type logger string

var (
	// CtxXKtbsRequestID context key for X-Ktbs-Request-ID
	CtxXKtbsRequestID ContextKey = "X-Ktbs-Request-ID"

	// CtxKtbsDonationIdentifier context key for X-Ktbs-Donation-Identifier
	CtxKtbsDonationIdentifier ContextKey = "X-Ktbs-Donation-Identifier"

	// CtxLogger context key for logger
	CtxLogger logger = "Ktbs-Logger"

	CtxWatermillProcessID ContextKey = "Ktbs-watermill-process-id"
)

func (c ContextKey) String() string {
	return string(c)
}

func (c logger) String() string {
	return string(c)
}
