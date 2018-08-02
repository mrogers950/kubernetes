#!/bin/bash
# ensure env variables make it into the plugin
if [ -z "${test_subject}" ]; then
  exit 1
else
  if [ "${test_subject}" != "CN=system:node:localhost,O=system:nodes" ]; then
    exit 1
  fi
fi
if [ -z "${test_dns}" ]; then
  exit 1
else
  if [ "${test_dns}" != "localhost" ]; then
    exit 1
  fi
fi
if [ -z "${test_ip}" ]; then
  exit 1
else
  if [ "${test_ip}" != "127.0.0.1" ]; then
    exit 1
  fi
fi
if [ -z "${test_other}" ]; then
  exit 1
else
  if [ "${test_other}" != "foo" ]; then
    exit 1
  fi
fi

echo '{"kind":"ExecCredential","apiVersion":"client.authentication.k8s.io/v1alpha1","spec":{},"status":{"clientCertificateData":"-----BEGIN CERTIFICATE-----\nMIIDYzCCAkugAwIBAgIJAM0ln+O2Fm+gMA0GCSqGSIb3DQEBCwUAMBExDzANBgNVBAMMBnRlc3RDQTAeFw0xODA3MjMxOTQ5MDhaFw0xOTA3MjMxOTQ5MDhaMBQxEjAQBgNVBAMMCWxvY2FsaG9zdDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMqcBHPJIMwlh/Udx6AurXVHW2pReRa1hsr/OA6ZW/6ncjt/vq6od356hlFmhKwMtE2gtb+ZWIoV+AEpQj21r002QYM2hZcQXQJtzT2K/jUb+IgbzLjwzgNLPY1Tit77nah25e5mF+Onh28qIF/CLtZYM0veykS43ttJRlqxIYjA44znAYo5ZO7e7oFt4nefETudhl0K5Yg0VZo4ll/FwVO06W0mv5nlrrd5622SHekGMs1m0lVu+KxHI0cQ/L4UPMFUxYbEidhqfyG0U6GHiKMyMJuKdEhgcas9ysEJoyM8RXIKhv2PskC5+SHqA7116zlWxvpwbRaNU46cPTTstccCAwEAAaOBujCBtzAdBgNVHQ4EFgQUXucfmHk8ejT2fqMH3FQqod0tEn8wHwYDVR0jBBgwFoAU59gqfpGfjF7uLILP7J91Azzbh2swCQYDVR0TBAIwADALBgNVHQ8EBAMCBaAwGgYDVR0RBBMwEYIJbG9jYWxob3N0hwR/AAABMCwGCWCGSAGG+EIBDQQfFh1PcGVuU1NMIEdlbmVyYXRlZCBDZXJ0aWZpY2F0ZTATBgNVHSUEDDAKBggrBgEFBQcDATANBgkqhkiG9w0BAQsFAAOCAQEAKYDpEtlHFJbpp5c8aAS57fES7Cd4ZH0uaast2t4OeU1uPfuCI3HsNwys3Jc7/U+vQiUAMF/wASD8yFNzTEef5s6P26Ir+fLr2mXpXWJ1GUosSrTjWlmB9gFHVfJObEaWtBRSwXHWrP/Qz019nYNazH7mb7an4LlCfxkczJv6JTV4jq800kL1PIeRQmeIyl67fzl8lbrDQJUulCyt8GLEsMILglKBO82EGvs9Z8RN7gaN9jAJCkgInfurCHn2dRmEbwjnrFHnLnpqb5r0tcCTR5axAqXQOgILVT52OWq6dkuzUvhFJAF/WhjHIFhAtLzjKx6x24FY5/hq/JcvptcmEg==\n-----END CERTIFICATE-----\n","clientKeyData":"-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAypwEc8kgzCWH9R3HoC6tdUdbalF5FrWGyv84Dplb/qdyO3++rqh3fnqGUWaErAy0TaC1v5lYihX4ASlCPbWvTTZBgzaFlxBdAm3NPYr+NRv4iBvMuPDOA0s9jVOK3vudqHbl7mYX46eHbyogX8Iu1lgzS97KRLje20lGWrEhiMDjjOcBijlk7t7ugW3id58RO52GXQrliDRVmjiWX8XBU7TpbSa/meWut3nrbZId6QYyzWbSVW74rEcjRxD8vhQ8wVTFhsSJ2Gp/IbRToYeIozIwm4p0SGBxqz3KwQmjIzxFcgqG/Y+yQLn5IeoDvXXrOVbG+nBtFo1Tjpw9NOy1xwIDAQABAoIBACiFSXKqr++ENgu0t/72NuS0r7i0sKX1Cg9BOcHZtAdbD8KMiuM9eCCIeqJ/AVuzcr/vu0mlboq3WBFLYh8bXgLwLewDFHag5CkfMqPzT2HpxSvbe3clWd5Yxuej5Ksx4VcW6GdkbbSvBawa3bypBlsB6shqt0NFQfTTU8nBkTZbGk3pv1VBKKEQnrYsV03tSO5Ph5jppdUusYFm2R8mgDwOR5PxcpN0uzA9f2XbKIPIUnhiNNcIss5jASTi5IPojOPgQurTq9aDR2XAnsJB+h+AqZbwQUpzrCMAlK8Oz/xwUs6BqSptVFDY0yVNS1qVfap4xd5oNWxGksQVo3P0kyECgYEA5WNy7gqrZ+3Kuk3vqWghDWmZ+jhKjeuGy45YnMYs7WRFK5ckBo6riucbj/1GwqIfmFAXUIDTlG5tJHL3s0Nmv5Q0Tjq5umjXpf/iiiQMLK/E8FgMpH0adMGqcFW49VBYa/CD2/HT6YjyDqy13Dj9ptYnMbMe3z/BG8puG3hoPJECgYEA4h1DxRQ9HEUVOVKMFLoCvAzBeVh3UZmA5W7Bl904+qLN3wkGy3bIBiVlQkZtaY5AspPgNkccJYsuPoDkx0Qp8qt0M6hWOmgSPBQMpSnzWoQF42YfuUemqFjfM0j4EEIEvGhDO2vbryDMn4NHAk/u1pnxmlVwOM8LNM3hJb+QWNcCgYEAvAC7BHAINcDF49XWdCDc3hJL2bFjIVgE/TZoV+1wiwwgSO6x3xH1dH2fsG6kHQclH/+cbCV5w3CR0UrMysaWIrRD/k3RRP+CpxHGyPNsav+QSG/RxMqn8UN8/l6znZNBNQ5F8/EKfp/3y6Ev2BN5iNCCBRDKX6zwB2fswGT6AZECgYAoB2pB721KHei99yEZYjytscxmgQTOi1BITa00B1PY+w1bGKv9RQ/wFpqweutProFBm/Ara7dN5i/PnN3jcOvELBosMvbg7B+eRyZd7ulH8utf8GpZUJfuYZ1R5O8VYbqY6BRO5q9Dd5kB/CmL/T6Y+zPMUKfHRtADDxd2qU0SjQKBgDrOCBhxDn61cRyMyUMo7zXxsC+YgZYwtDyYtI/s37Px2nWO+lfGVX/qZhM/GfLa20Uqe9a9Khw7tJ1VSDvFayyv/c2WuszIjLXgE1wDLnOVcetQSn3eDag1IFBW9JlHHZ5RgEDdy3/lmWMi4nLDsa+/UNXs++AXOSAB0o+yFUQP\n-----END RSA PRIVATE KEY-----\n"}}'