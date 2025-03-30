package logger

import "bufio"

// BufferedWriteSyncer untuk buffer log ke stdout
type BufferedWriteSyncer struct {
	writer *bufio.Writer
}

func (b *BufferedWriteSyncer) Write(p []byte) (n int, err error) {
	n, err = b.writer.Write(p)
	b.writer.Flush() // Flush setiap kali ada log baru
	return n, err
}

func (b *BufferedWriteSyncer) Sync() error {
	return b.writer.Flush()
}
