package msg

import (
	"GoWinrm/internal/utils"
	"GoWinrm/internal/wsmv"
	"fmt"
	"time"
)

// CreateShell representa un mensaje WSMV para crear un shell remoto
type CreateShell struct {
	sessionOpts      wsmv.SessionOptions
	shellURI         string
	iStream          string
	oStream          string
	codepage         int
	noprofile        string
	workingDirectory string
	idleTimeout      any
	envVars          map[string]string
}

// UTF8CodePage representa la codificación UTF-8
const UTF8CodePage = 65001

// NewCreateShell crea una nueva instancia de CreateShell
func NewCreateShell(sessionOpts wsmv.SessionOptions, shellOpts map[string]any) *CreateShell {
	cs := &CreateShell{
		sessionOpts: sessionOpts,
		shellURI:    optOrDefault(shellOpts, "shell_uri", wsmv.RESOURCEURICMD).(string),
		iStream:     optOrDefault(shellOpts, "i_stream", "stdin").(string),
		oStream:     optOrDefault(shellOpts, "o_stream", "stdout stderr").(string),
		codepage:    optOrDefault(shellOpts, "codepage", UTF8CodePage).(int),
		noprofile:   optOrDefault(shellOpts, "noprofile", "FALSE").(string),
	}

	if workingDir, ok := shellOpts["working_directory"]; ok {
		cs.workingDirectory = workingDir.(string)
	}

	if idleTimeout, exists := shellOpts["idle_timeout"]; exists {
		cs.idleTimeout = idleTimeout
	}

	if envVars, ok := shellOpts["env_vars"]; ok {
		cs.envVars = envVars.(map[string]string)
	}

	return cs
}

// Headers retorna los encabezados para el mensaje CreateShell
func (cs *CreateShell) Headers() map[string]any {
	return wsmv.MergeHeaders(
		wsmv.SharedHeaders(cs.sessionOpts),
		wsmv.ResourceURIShell(cs.shellURI),
		wsmv.ActionCreate(),
		cs.headerOpts(),
	)
}

// Body retorna el cuerpo para el mensaje CreateShell
func (cs *CreateShell) Body() map[string]any {
	return map[string]any{
		fmt.Sprintf("%s:Shell", wsmv.NS_WIN_SHELL): cs.shellBody(),
	}
}

// shellBody construye el cuerpo del mensaje Shell
func (cs *CreateShell) shellBody() map[string]any {
	body := map[string]any{
		fmt.Sprintf("%s:InputStreams", wsmv.NS_WIN_SHELL):  cs.iStream,
		fmt.Sprintf("%s:OutputStreams", wsmv.NS_WIN_SHELL): cs.oStream,
	}

	if cs.workingDirectory != "" {
		body[fmt.Sprintf("%s:WorkingDirectory", wsmv.NS_WIN_SHELL)] = cs.workingDirectory
	}

	if cs.idleTimeout != nil {
		body[fmt.Sprintf("%s:IdleTimeOut", wsmv.NS_WIN_SHELL)] = formatIdleTimeout(cs.idleTimeout)
	}

	if len(cs.envVars) > 0 {
		body[fmt.Sprintf("%s:Environment", wsmv.NS_WIN_SHELL)] = cs.environmentVarsBody()
	}

	return body
}

// environmentVarsBody construye el cuerpo para las variables de entorno
func (cs *CreateShell) environmentVarsBody() map[string]any {
	variables := make([]map[string]any, 0, len(cs.envVars))

	for name, value := range cs.envVars {
		variables = append(variables, map[string]any{
			"_": value,
			":attributes!": map[string]any{
				"Name": name,
			},
		})
	}

	return map[string]any{
		fmt.Sprintf("%s:Variable", wsmv.NS_WIN_SHELL): variables,
	}
}

// headerOpts construye las opciones del encabezado
func (cs *CreateShell) headerOpts() map[string]any {
	options := []map[string]any{
		{
			"_": cs.noprofile,
			":attributes!": map[string]any{
				"Name": "WINRS_NOPROFILE",
			},
		},
		{
			"_": fmt.Sprintf("%d", cs.codepage),
			":attributes!": map[string]any{
				"Name": "WINRS_CODEPAGE",
			},
		},
	}

	return map[string]any{
		fmt.Sprintf("%s:OptionSet", wsmv.NS_WSMAN_DMTF): map[string]any{
			fmt.Sprintf("%s:Option", wsmv.NS_WSMAN_DMTF): options,
		},
	}
}

func formatIdleTimeout(timeout any) string {
	switch t := timeout.(type) {
	case string:
		return t
	case int:
		return utils.Iso8601Duration(t)
	case time.Duration:
		return utils.Iso8601Duration(int(t.Seconds()))
	default:
		return ""
	}
}

func optOrDefault(options map[string]any, key string, defaultValue any) any {
	if value, ok := options[key]; ok {
		return value
	}
	return defaultValue
}
