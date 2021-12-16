package aws

type streamPayloadReader struct {
	data []byte
}

func newStreamReader(data []byte) awsPayloadReader {
	return &streamPayloadReader{data: data}
}

func (rdr *streamPayloadReader) Read() ([]byte, error) {

	return rdr.data, nil
}
