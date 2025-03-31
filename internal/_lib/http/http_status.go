package httpstatus

// HttpStatus represents HTTP status codes as constants
type HttpStatus int

const (
	Continue                     HttpStatus = 100
	SwitchingProtocols           HttpStatus = 101
	Processing                   HttpStatus = 102
	OK                           HttpStatus = 200
	Created                      HttpStatus = 201
	Accepted                     HttpStatus = 202
	NonAuthoritativeInformation  HttpStatus = 203
	NoContent                    HttpStatus = 204
	ResetContent                 HttpStatus = 205
	PartialContent               HttpStatus = 206
	Ambiguous                    HttpStatus = 300
	MovedPermanently             HttpStatus = 301
	Found                        HttpStatus = 302
	SeeOther                     HttpStatus = 303
	NotModified                  HttpStatus = 304
	TemporaryRedirect            HttpStatus = 307
	PermanentRedirect            HttpStatus = 308
	BadRequest                   HttpStatus = 400
	Unauthorized                 HttpStatus = 401
	PaymentRequired              HttpStatus = 402
	Forbidden                    HttpStatus = 403
	NotFound                     HttpStatus = 404
	MethodNotAllowed             HttpStatus = 405
	NotAcceptable                HttpStatus = 406
	ProxyAuthenticationRequired  HttpStatus = 407
	RequestTimeout               HttpStatus = 408
	Conflict                     HttpStatus = 409
	Gone                         HttpStatus = 410
	LengthRequired               HttpStatus = 411
	PreconditionFailed           HttpStatus = 412
	PayloadTooLarge              HttpStatus = 413
	URITooLong                   HttpStatus = 414
	UnsupportedMediaType         HttpStatus = 415
	RequestedRangeNotSatisfiable HttpStatus = 416
	ExpectationFailed            HttpStatus = 417
	UnprocessableEntity          HttpStatus = 422
	TooManyRequests              HttpStatus = 429
	InternalServerError          HttpStatus = 500
	NotImplemented               HttpStatus = 501
	BadGateway                   HttpStatus = 502
	ServiceUnavailable           HttpStatus = 503
	GatewayTimeout               HttpStatus = 504
	HTTPVersionNotSupported      HttpStatus = 505
)
