package transport

import (
	"GoWinrm/internal/log"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/ManuJSM/GoNtlm/client"
)

type ntlmNego struct {
	httpcli  *http.Client
	endpoint string
	ntlmcli  *client.Client
}
type NegotiateOpts = client.ClientOpts

const userAgent = "TEST-USER-AGENT"

func NewNtlmNego(endpoint string, opts *NegotiateOpts) *ntlmNego {
	return &ntlmNego{
		httpcli:  &http.Client{},
		endpoint: endpoint,
		ntlmcli:  client.NewClient(opts),
	}
}

func createAuthRequest(authToken string, endpoint string) *http.Request {

	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer([]byte("")))
	req.Header.Set("Authorization", "Negotiate "+authToken)
	req.Header.Set("Connection", "Keep-Alive")
	req.Header.Set("Content-Type", "application/soap+xml;charset=UTF-8")
	req.Header.Set("User-Agent", userAgent)

	return req
}

func createSendRequest(endpoint string, message []byte) *http.Request {
	req, _ := http.NewRequest("POST", endpoint, bytes.NewReader(message))
	req.Header.Set("Content-Type", `multipart/encrypted;protocol="application/HTTP-SPNEGO-session-encrypted";boundary="Encrypted Boundary"`)
	req.Header.Set("User-Agent", userAgent)

	return req
}

func (nn *ntlmNego) issueChallengeResponse(auth3Token string) error {

	type3req := createAuthRequest(auth3Token, nn.endpoint)

	resp, err := nn.httpcli.Do(type3req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return nil
	}

	return fmt.Errorf("error: %d", resp.StatusCode)
}

func (nn *ntlmNego) initAuth() error {

	// Paso 1: mensaje inicial
	auth1Token := base64.StdEncoding.EncodeToString(nn.ntlmcli.Type1Request())

	type1req := createAuthRequest(auth1Token, nn.endpoint)

	resp1, err := nn.httpcli.Do(type1req)
	if err != nil {
		return err
	}
	defer resp1.Body.Close()

	authHeader := resp1.Header.Get("WWW-Authenticate")
	if authHeader == "" {
		return fmt.Errorf("missing WWW-Authenticate header")
	}

	parts := strings.Split(authHeader, " ")
	itok := parts[len(parts)-1]

	// Paso 2: respuesta al desafío
	challengeBytes, _ := base64.StdEncoding.DecodeString(itok)
	type3, err := nn.ntlmcli.Type3Request(challengeBytes)
	if err != nil {
		return err
	}
	auth3Encoded := base64.StdEncoding.EncodeToString(type3)

	return nn.issueChallengeResponse(auth3Encoded)
}

func (nn *ntlmNego) seal(message []byte) ([]byte, error) {
	sealmsg, err := nn.ntlmcli.SealMessage(message)
	if err != nil {
		return nil, err
	}
	signature, err := nn.ntlmcli.SignMessage(message)
	if err != nil {
		return nil, err
	}

	prefix := []byte{0x10, 0x00, 0x00, 0x00}
	result := append(prefix, signature...)
	result = append(result, sealmsg...)

	return result, nil
}

func (nn *ntlmNego) winrmDecrypt(resp *http.Response) ([]byte, error) {

	contentType := resp.Header.Get("Content-Type")
	if matched, _ := regexp.MatchString(`(?i)^application/soap\+xml`, contentType); matched {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return body, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if len(body) == 0 {
		return []byte{}, nil
	}

	re := regexp.MustCompile(`(?s)^.*Content-Type: application/octet-stream\r\n(.*)--Encrypted.*$`)
	matches := re.FindSubmatch(body)
	if len(matches) < 2 {
		return nil, errors.New("could not extract encrypted message")
	}
	encrypted := matches[1]

	if len(encrypted) < 20 {
		return nil, errors.New("invalid encrypted message format")
	}
	signature := encrypted[4:20]
	encryptedMessage := encrypted[20:]

	nn.ntlmcli.UnSealMessage(encryptedMessage)

	if ok, err := nn.ntlmcli.VerifySignature(signature, encryptedMessage); !ok {
		return nil, fmt.Errorf("could not decrypt NTLM message, signature verification fail: %v", err)
	}
	return encryptedMessage, nil
}

func body(message string, length int, contentType ...string) string {
	t := "application/HTTP-SPNEGO-session-encrypted"
	if len(contentType) > 0 {
		t = contentType[0]
	}

	return strings.Join([]string{
		"--Encrypted Boundary",
		fmt.Sprintf("Content-Type: %s", t),
		fmt.Sprintf("OriginalContent: type=application/soap+xml;charset=UTF-8;Length=%d", length),
		"--Encrypted Boundary",
		"Content-Type: application/octet-stream",
		fmt.Sprintf("%s--Encrypted Boundary--", message),
	}, "\r\n") + "\r\n"
}

func (nn *ntlmNego) SendRequest(message []byte) ([]byte, error) {

	if !nn.ntlmcli.IsSession() {
		err := nn.initAuth()
		if err != nil {
			return nil, err
		}
	}

	log.Debug(string(message) + "\n")

	sealed, err := nn.seal(message)
	if err != nil {
		return nil, err
	}
	post := body(string(sealed), len(message))

	req := createSendRequest(nn.endpoint, []byte(post))

	resp, _ := nn.httpcli.Do(req)

	decrypted, err := nn.winrmDecrypt(resp)
	if err != nil {
		return nil, fmt.Errorf("decrypt failed: %v", err)
	}

	log.Debug(string(decrypted) + "\n")

	return RespHandler(decrypted, resp.StatusCode)
}
