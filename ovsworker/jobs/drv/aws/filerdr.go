package aws

import "os"

type payloadFileReader struct {
	path string
}

func newFileReader(path string) awsPayloadReader {

	return &payloadFileReader{path: path}
}
func (rdr *payloadFileReader) Read() ([]byte, error) {

	return os.ReadFile(rdr.path)
}
