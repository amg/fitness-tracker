package main

import (
	"fmt"
	"testing"
	"time"
)

func TestJWTGenerationAndSigning(t *testing.T) {
	var testTime time.Time
	testTime, _ = time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	fmt.Println(testTime)
	jwt, _ := jwtWithCustomClaims("./jwtRSA256-private.pem", testTime)
	expected := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJpc3MiOiJ0ZXN0Iiwic3ViIjoic29tZWJvZHkiLCJhdWQiOlsic29tZWJvZHlfZWxzZSJdLCJleHAiOjExMzYyMTYwNDUsIm5iZiI6MTEzNjIxNDI0NSwiaWF0IjoxMTM2MjE0MjQ1LCJqdGkiOiIxIn0.XGazhG6L_5xRh3qG6iR7AcHhxM1OyeR5sqbbFXh9-nlBxSVEZIlaFioQYTx1n4cfx_pRKFdd3SjeezhMYuv5qp04qCvGYbM5d_0ndYEt_nvYRWe-weFcsUy1EKPI1fEwfOfNTX0hVVlb7MwXJf6c1QvFsJYUdegU5fGmMwK2lYW0xac6t0tdZ-W8ZXj9925QziKy3R91JwzaBhuAksQp5KUxMb-GG3qO448AeXRBts0j7ptmeQtYoBvCvx9IlU-AoTTATYJt9ovM4gl1fHqlsO5d7lomgnfQCSu53iRnD0KQiJmAYNHFbXdzZNoIO5vvuWIpbIUftGe87xqVtKHhLA"
	if expected != jwt {
		t.Error(jwt)
	}
}
