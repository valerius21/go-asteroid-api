package user

import (
	"context"
	"github.com/google/uuid"
	"github.com/pastoapp/astroid-api/internal/orbitdb"
	"testing"
)

// WARNING: DO NOT USE THOSE KEYS

const PrivateKey = "-----BEGIN PRIVATE KEY-----\n" +
	"MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDW+0l1hNQeIHNQ\n" +
	"v9fElNVILJs873vq0LQYNTQD5+Md9lk9XmvpQoQwMN5OrQS1rmCxXR5gl1w6WElg\n" +
	"sEC1JaQoJFD6uyoOVonCCj7oozv5sUjrpzkRdrWfI0YJ4cm0VK8C2yZgE9uMRNvF\n" +
	"fnRrvxBRBD/cKlD19C3Eh2gjJVex5k6x3UhGmMAOi5Y7Fh6f9hZaE1cjbLn/Xuyu\n" +
	"/fAv30p+EUpxyUw9bMGTSaVeQPUgJDlWbnQ27vuaNZj94ZMJxsnqgxZ67e9q0v2n\n" +
	"qgzeb1Vdk9dz4xoEMSKkneZcRB97/N5dnGb2XkBqIA2FM05AAcyAJD6PzY1cNk1b\n" +
	"SLSwzd/XAgMBAAECggEANQa0XtVNkMTGTMOVD/r6BJxkCgq0EP47fcEZvGDlB03t\n" +
	"QxLYISfz/QMAwn0wrTN+zCDAzmhnKkj+iGfQ/h1zuScEBmRGAlVcE019qNoif/s3\n" +
	"llCf5evW7pK/ws47d0xwAo9Vjts3CwktADCpvK3a83DN45tYbuI0q8YObKTYVZua\n" +
	"F0xaEt34wIr00+sFGp+EIEbdaIYsA+Xw2xrnoXJ7RZPqab1flY2R3Ywtjj1X7xs+\n" +
	"36pRb32IcHChTFOvsPEquwzzRtSXrsGWh7sC/T16oK2D3zKhPO4cNXeKr34y7RA/\n" +
	"rjvN5NrRQmnVcxTJN9AVKQk3GEGrpDaCXI4xKJuYgQKBgQD5puq0JePzOWye9fVn\n" +
	"3SaAZR9Ad34lYXZbQ0FqFryy4F89pAnlXDEO9IcxEDQ/yO6rtl72JrXCMLAz89rd\n" +
	"lgG/RAxPWJHU66vW6do5QPzaQNol6QZNMAH03pccFueiSxvn8pUjSFclA/i1CR2D\n" +
	"3LfPEVlijWokVE5lRnwcc2WtwQKBgQDccq/KGICLi6jBSOY6iy7Am7rP8FZVfZ5A\n" +
	"GUdSGcxp2NexpzhrsgF5t5MyLhvQVcfyphtqdelQP8ZEEV/qo0AU/PgQOVdD42UV\n" +
	"hZ1cC0hz0ylllBJOLIwnXcsddeL/yrq4lYKxh0er35qPR+5iWOEq2btKNZeP5V0a\n" +
	"0ipeQ9MjlwKBgQCT73A27hkEwOt3ysHwGSE+hi2dsvFCT12mONVjezO4P7LQQQtI\n" +
	"KXg8p3QRA8wiW/vqLrhY3ploCBfxwTfZhJFmeHle71oDNN/wC0UwFi1/W1ypiEse\n" +
	"pCsUom4izAPYrO9NY/cyICTbn13aURt3TuxwvJQeDTRvcLG/26Oh7FO9QQKBgQCn\n" +
	"QezWbRi6x8/WNrVXndpvELkAo90dKB41i/hs9l39wSSGCdXpADnWW842uH0p6Ku4\n" +
	"EweiEqxe0sRfSL95iztHsVeVF6x6OcGsPqBSAc2K85/pDhuFbI86HxPI2nTMlshs\n" +
	"3CM5Bub+JL1kELNxsvXESUas21SPPPsTKnoy40sDKQKBgBdPYIzkOP2p7ACTHyoe\n" +
	"rzIBXvkfKHJPn5rI9FZawF/Tjw3HocQg6Jp7+a43VGWBbhdKKWTjDS7m1wln9BZw\n" +
	"fL6nT6xqtetsSLmtZ6VTU2JwWjNFAJbEIeDsvf1GSPj4sSS6YgszS/2yRiX8F1Fs\n" +
	"oSVs9U8dbbPt5paggGlxOc9e\n" +
	"-----END PRIVATE KEY-----\n"

const PublicKey = "-----BEGIN PUBLIC KEY-----\n" +
	"MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1vtJdYTUHiBzUL/XxJTV\n" +
	"SCybPO976tC0GDU0A+fjHfZZPV5r6UKEMDDeTq0Eta5gsV0eYJdcOlhJYLBAtSWk\n" +
	"KCRQ+rsqDlaJwgo+6KM7+bFI66c5EXa1nyNGCeHJtFSvAtsmYBPbjETbxX50a78Q\n" +
	"UQQ/3CpQ9fQtxIdoIyVXseZOsd1IRpjADouWOxYen/YWWhNXI2y5/17srv3wL99K\n" +
	"fhFKcclMPWzBk0mlXkD1ICQ5Vm50Nu77mjWY/eGTCcbJ6oMWeu3vatL9p6oM3m9V\n" +
	"XZPXc+MaBDEipJ3mXEQfe/zeXZxm9l5AaiANhTNOQAHMgCQ+j82NXDZNW0i0sM3f\n" +
	"1wIDAQAB\n" +
	"-----END PUBLIC KEY-----\n"

func TestNewUser(t *testing.T) {
	cancelFunc, err := orbitdb.InitializeOrbitDB("http://localhost:5001", t.TempDir())
	if err != nil {
		t.Fatalf("Error initializing OrbitDB: %v", err)
	}
	defer cancelFunc()
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("should create a new user", func(t *testing.T) {
		user, err := NewUser(PublicKey, false)

		if err != nil {
			t.Errorf("an error occurred %v\n", err)
		}

		if user.PublicKey != PublicKey {
			t.Errorf("public key does not match")
		}

		if user.Nonce == "" {
			t.Errorf("no valid Nonce found")
		}

		if user.IsAdmin {
			t.Errorf("user should not be admin")
		}
	})

	t.Run("should create a user and find it", func(t *testing.T) {
		user, err := NewUser(PublicKey, false)
		if err != nil {
			t.Errorf("error creating the user, %v\n", err)
		}

		resp, err := Find(user.ID.String())
		if err != nil {
			t.Errorf("error finding the user %v - %v\n", user, resp)
		}
	})

	t.Run("should create multiple users and find them", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			user, err := NewUser(PublicKey, false)
			if err != nil {
				t.Errorf("error creating user %v\n", err)
			}
			resp, err := Find(user.ID.String())

			_, err = uuid.Parse(resp.ID.String())
			if err != nil {
				t.Errorf("user could not be queried")
			}
		}

	})
}
