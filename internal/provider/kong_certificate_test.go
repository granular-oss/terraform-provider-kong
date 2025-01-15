package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccKongCertificateDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "kong_certificate" "test" {
  id = "56c71339-eda3-4c34-9aff-22ac90539119"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					assertStringValue("data.kong_certificate.test", "id", "56c71339-eda3-4c34-9aff-22ac90539119"),
					assertStringArrayValue("data.kong_certificate.test", "snis", []string{}),
					assertNullValue("data.kong_certificate.test", "tags"),
					statecheck.ExpectKnownValue("data.kong_certificate.test", tfjsonpath.New("cert"), knownvalue.StringRegexp(regexp.MustCompile("--BEGIN\\sCERTIFICATE--"))),
					statecheck.ExpectKnownValue("data.kong_certificate.test", tfjsonpath.New("key"), knownvalue.StringRegexp(regexp.MustCompile("--BEGIN\\sPRIVATE\\sKEY--"))),
				},
			},
		},
	})
}

func TestAccKongCertificateResource(t *testing.T) {
	idStaysSame := statecheck.CompareValue(compare.ValuesSame())
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "kong_certificate" "test" {
  cert = <<EOF
%s
EOF
  key = <<EOF
%s
EOF
}
`, testCert1, testKey1),
				ConfigStateChecks: []statecheck.StateCheck{
					assertNotNull("kong_certificate.test", "id"),
					idStaysSame.AddStateValue("kong_certificate.test", tfjsonpath.New("id")),
					assertStringArrayValue("kong_certificate.test", "snis", []string{}),
					assertNullValue("kong_certificate.test", "tags"),
					statecheck.ExpectKnownValue("kong_certificate.test", tfjsonpath.New("cert"), knownvalue.StringRegexp(regexp.MustCompile("--BEGIN\\sCERTIFICATE--"))),
					statecheck.ExpectKnownValue("kong_certificate.test", tfjsonpath.New("key"), knownvalue.StringRegexp(regexp.MustCompile("--BEGIN\\sPRIVATE\\sKEY--"))),
				},
			},
			{
				ResourceName:      "kong_certificate.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: providerConfig + fmt.Sprintf(`
resource "kong_certificate" "test" {
  cert = <<EOF
%s
EOF
  key = <<EOF
%s
EOF
  tags = ["foo"]
  snis = ["test.com"]
}
`, testCert1, testKey1),
				ConfigStateChecks: []statecheck.StateCheck{
					idStaysSame.AddStateValue("kong_certificate.test", tfjsonpath.New("id")),
					assertStringArrayValue("kong_certificate.test", "snis", []string{"test.com"}),
					assertStringArrayValue("kong_certificate.test", "tags", []string{"foo"}),
					statecheck.ExpectKnownValue("kong_certificate.test", tfjsonpath.New("cert"), knownvalue.StringRegexp(regexp.MustCompile("--BEGIN\\sCERTIFICATE--"))),
					statecheck.ExpectKnownValue("kong_certificate.test", tfjsonpath.New("key"), knownvalue.StringRegexp(regexp.MustCompile("--BEGIN\\sPRIVATE\\sKEY--"))),
				},
			},
			{
				Config: providerConfig + fmt.Sprintf(`
resource "kong_certificate" "test" {
  cert = <<EOF
%s
EOF
  key = <<EOF
%s
EOF
}
`, testCert1, testKey1),
				ConfigStateChecks: []statecheck.StateCheck{
					idStaysSame.AddStateValue("kong_certificate.test", tfjsonpath.New("id")),
					assertStringArrayValue("kong_certificate.test", "snis", []string{}),
					assertNullValue("kong_certificate.test", "tags"),
					statecheck.ExpectKnownValue("kong_certificate.test", tfjsonpath.New("cert"), knownvalue.StringRegexp(regexp.MustCompile("--BEGIN\\sCERTIFICATE--"))),
					statecheck.ExpectKnownValue("kong_certificate.test", tfjsonpath.New("key"), knownvalue.StringRegexp(regexp.MustCompile("--BEGIN\\sPRIVATE\\sKEY--"))),
				},
			},
		},
	})

}

const (
	testKey1 = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDXi8zFDavAN7fl
RJO2G4oLj6NIT86BJnzM3XqtGl6pvfp0bo9so+h/0HhGtnIh7Je4BL7PGsv5BSdg
6EDDZDXZn/ZDe3jje+Ee1sfn98H+1mjTDlm0U2it/cWaZ+a8GEhPwidWyI1AeS4O
XrM5VcbcmVoIPRr5z0iJBJ2LRY0L/rPVsJOGbT1WPsFjFMZrc7GAixrjjk0jovKr
X7hK/Sj3vRGjZsP6CDEVIeEJQOMDvup9YCbgPTofn3gH9PoKFQdA0CbuGsdKiUOF
bfrp2MFP8933WsCGhUbpsdDd/Vt9JvpFlR0aCShv0wIWhvzv8WQ5MUWFQvKTmXDv
tdnRqXTvAgMBAAECggEARP648zKnEYZEVR0Ycyhpjb3StGjnXyvksucKR7KzLn5j
VzW0rz/gQlmGxovMCNPk1MCgG0cml3Vw33I4mNLQ8fJkL8GsNpUGwIpbvwLtlBcp
wrVLPY+daGRdBknP79GOBAnP8dWMcWDYvzzM/cNZPm/QA+cbZW9WdpWFoHkI5xdw
LIuI4q7jsjbRHEK5Mc055EDC5hwkmZrjVi9hsqk2De5AYaXhX1kqPopOuKvxJjff
H0QGvE285HFi0qdDYcEzpFEgzBUMjU144bqyl+s7VLhuJ302VVsUJZbNq0m1h9xN
/KSqK6Q3d3ZgsP9fWpq/H9N1NcxpmJbxyET+3P4lEQKBgQDvHzYRpn2bCgPrZaF7
oCbFcdBPdRc+GSTuNzZLYl5e/dGLD/+soxOobZzH9HmhpMkOwpZX4ajz7/5y3GTH
OJ9nyuh52b1weOAWOodmUrMqwRcxNNtkUE6LFMVdatgGpC8ihvoMfDM3GK91Hdbr
bs1Ws7tfIkftad3nwM4YO62oQwKBgQDmwpZ1sXgDUxswWq9vvTy4TDOkRdcEO2J/
8Aphut8RJ/itT2cm5WKWeNkQvyC1Gs4PMGnqU2LWG+31YPB1GAk3a+OP7GN7WewS
rMBnmgdoCjDqXOBg15KOn2m23Ac6a6dBGI9pvtQfhX0e59izcSdA3ZwboFfYdmMI
hYn49uO75QKBgQCXYQnouJ7R1NBQaKGHUwbYfkni03yoWmCv0hI0PQ0DU+ohADra
/s5GFUZoq5OIynpiNrvY3MoJzAgojO/b0zPPEHyGD1tHZa5vRBRNqdM1INJe21h8
s/5VPAwKLMafxbb1Q7/uwX3mxmDlYsOZfibOWbAn9NrWKOxLeBrA6p7wYwKBgF/x
J31ne+5l7zf7fFWI6GX3yMDUCMHJrvpiYu6fM39+jvX/vXN+i67kL9u2m3Kw4luO
VXsHkGBU3GrZEyCcDbjtMn/0WKhAitZ43MY2VD39frjyRJf/CQAjZ2CPurGfcLqv
63Cb1rYEWjEvU/nHYfqmKPGTiPKGxkYUv3izrZvBAoGAVyi8N8LtYhG+T52tVhXM
0wopUsU5NkZjuTof8lG0oAUv6B0I6osacsB9AFv1Ai8Y6hBpZlLRALKDLu/pLVrj
GGRb4/tQrQz7/cn73RiefOha1TzDSEEfq237/eBy7dqk75C3nfLNtu6vEFuAgQbS
dTLBTUGXDSVpcySM5XwoVuM=
-----END PRIVATE KEY-----`
	testCert1 = `-----BEGIN CERTIFICATE-----
MIIDKjCCAhICCQDV2QFIWXZy0DANBgkqhkiG9w0BAQsFADBXMQswCQYDVQQGEwJH
QjENMAsGA1UECAwEQ0FNQjESMBAGA1UEBwwJQ2FtYnJpZGdlMRQwEgYDVQQKDAtr
ZXZob2xkaXRjaDEPMA0GA1UEAwwGZ29rb25nMB4XDTE5MDExNDIxMTgxMVoXDTI5
MDExMTIxMTgxMVowVzELMAkGA1UEBhMCR0IxDTALBgNVBAgMBENBTUIxEjAQBgNV
BAcMCUNhbWJyaWRnZTEUMBIGA1UECgwLa2V2aG9sZGl0Y2gxDzANBgNVBAMMBmdv
a29uZzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANeLzMUNq8A3t+VE
k7YbiguPo0hPzoEmfMzdeq0aXqm9+nRuj2yj6H/QeEa2ciHsl7gEvs8ay/kFJ2Do
QMNkNdmf9kN7eON74R7Wx+f3wf7WaNMOWbRTaK39xZpn5rwYSE/CJ1bIjUB5Lg5e
szlVxtyZWgg9GvnPSIkEnYtFjQv+s9Wwk4ZtPVY+wWMUxmtzsYCLGuOOTSOi8qtf
uEr9KPe9EaNmw/oIMRUh4QlA4wO+6n1gJuA9Oh+feAf0+goVB0DQJu4ax0qJQ4Vt
+unYwU/z3fdawIaFRumx0N39W30m+kWVHRoJKG/TAhaG/O/xZDkxRYVC8pOZcO+1
2dGpdO8CAwEAATANBgkqhkiG9w0BAQsFAAOCAQEAP6xjv2nqMb9NmyUPz6bGlLNq
8lqUE4zWK61YS6P3BinRIswwDfUg42eMcafebOBgyc34yLBSbKF9paDupuI/xcyk
ySQk48vSGYAuo0wlN8YAmf6SC7tkfk7PL8uVl8bblDREk+D28UzEMNMA4ScCoYtQ
21G2HUhMonRI+MGRtbaVmc14XXjpPww29W6s5nxuG5MaGWd6wkIL7pmHmVBSN2QK
RQzGLmfi0TxOiCNCb9fArIaxlXYfR/yBoV/NdEKrFdpQg3pxKNKu0+IYJHomJDpZ
+Hr3Nf7YNDiX/eCuG//beQaE2H4A9/K7i15szIbv/inpIkcx7z5eIGULR7Hykw==
-----END CERTIFICATE-----`
)
