package models

type (
	HTTPProxyParams struct {
		URL string `json:"url" query:"url"`
	}
)

func (p *HTTPProxyParams) IsValid() bool {
	return p.URL != ""
}
