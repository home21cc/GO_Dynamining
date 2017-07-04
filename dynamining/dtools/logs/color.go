package logs

import "io"

type ansiColorWriter struct {
	writer io.Writer
	mode outputMode
}

func (color *ansiColorWriter) Write(p []byte) (int, error) {
	return color.writer.Write(p)
}
