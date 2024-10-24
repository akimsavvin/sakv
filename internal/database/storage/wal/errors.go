package wal

type CreationError struct {
	innerErr error
}

func newCreationError(innerErr error) CreationError {
	return CreationError{
		innerErr: innerErr,
	}
}

func (err CreationError) Error() string {
	return "could not create wal"
}

func (err CreationError) Unwrap() error {
	return err.innerErr
}

type SegmentsProcessingError struct {
	innerErr    error
	segmentName *string
}

func newSegmentsProcessingError(innerErr error, segmentName ...string) SegmentsProcessingError {
	err := SegmentsProcessingError{
		innerErr: innerErr,
	}

	if len(segmentName) > 0 {
		err.segmentName = &segmentName[0]
	}

	return err
}

func (err SegmentsProcessingError) Error() string {
	return "could not process segment"
}

func (err SegmentsProcessingError) Unwrap() error {
	return err.innerErr
}

func (err SegmentsProcessingError) SegmentName() (string, bool) {
	if err.segmentName == nil {
		return "", false
	}

	return *err.segmentName, true
}
