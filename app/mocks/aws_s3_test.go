package mocks

import "testing"

func TestUploadImage(t *testing.T) {
	s3Mock := NewAWSS3Repository()
	_, err := s3Mock.UploadImage(nil, nil)
	if err != nil {
		t.Fatalf("An error occurred: %v", err)
	}
}
