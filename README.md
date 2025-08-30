# 🖧 GoWinrm

This module provides a basic connection and shell setup for a **WinRM** (Windows Remote Management) client in Go, using **NTLM Negotiate SSP Authentication**.

---

## ⚙️ Basic Usage

```go
args := utils.Arguments{ //map[string]any
    "endpoint": "http://host:5985/wsman",
    "username": "Administrador",
    "password": "miPassword",//Also NTLMHASH 00000000000000000000000000000000:b4b9b02e6f09a9bd760f388b67351e2b
}

conf, err := NewConf(args)
if err != nil {
    log.Fatal(err)
}

conn := NewConnection(conf)
oc, err := conn.Shell.RunCommand("whoami")
	if err != nil {
		Log.Error(err)
	} else {
		Log("STDOUT: ", oc.Stdout.String())
		Log("STDERR: ", oc.Stderr.String())
		Log("EXITCODE: ", oc.ExitCode)
	}

	err = conn.Shell.Close()
	if err != nil {
		Log.Error(err)
	}

```

---

## 🔑 Supported Arguments

| Key               | Type       | Required | Default Value      |
| ----------------- | ---------- | -------- | ------------------ |
| `endpoint`        | `string`   | ✅ Yes   | -                  |
| `username`        | `string`   | ✅ Yes   | -                  |
| `password`        | `string`   | ✅ Yes   | -                  |
| `workstation`     | `string`   | ❌ No    | `""`               |
| `domain`          | `string`   | ❌ No    | `""`               |
| `locale`          | `string`   | ❌ No    | `"en-US"`          |
| `maxEnvelopeSize` | `int`      | ❌ No    | `153600`           |
| `opTimeout`       | `Duration` | ❌ No    | `60 * time.Second` |

---
