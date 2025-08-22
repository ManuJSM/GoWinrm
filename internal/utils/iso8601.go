package utils

import "fmt"

// convierte segundos a formato ISO8601 duración.
func Iso8601Duration(seconds int) string {
	isoStr := "P"

	if seconds > 604800 { // más de una semana
		weeks := seconds / 604800
		seconds -= weeks * 604800
		isoStr += fmt.Sprintf("%dW", weeks)
	}
	if seconds > 86400 { // más de un día
		days := seconds / 86400
		seconds -= days * 86400
		isoStr += fmt.Sprintf("%dD", days)
	}
	if seconds > 0 {
		isoStr += "T"
		if seconds > 3600 { // más de una hora
			hours := seconds / 3600
			seconds -= hours * 3600
			isoStr += fmt.Sprintf("%dH", hours)
		}
		if seconds > 60 { // más de un minuto
			minutes := seconds / 60
			seconds -= minutes * 60
			isoStr += fmt.Sprintf("%dM", minutes)
		}
		isoStr += fmt.Sprintf("%dS", seconds)
	}

	return isoStr
}
