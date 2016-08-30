// Copyright 2015-present Oursky Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package asset

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// FileStore implements Store by storing files on file system
type FileStore struct {
	dir        string
	prefix     string
	postPrefix string
	secret     string
	public     bool
}

// NewFileStore creates a new FileStore
func NewFileStore(dir, prefix, postPrefix, secret string, public bool) *FileStore {
	return &FileStore{dir, prefix, postPrefix, secret, public}
}

// GetFileReader returns a reader for reading files
func (s *FileStore) GetFileReader(name string) (io.ReadCloser, error) {
	path := filepath.Join(s.dir, name)
	return os.Open(path)
}

// PutFileReader stores a file from reader onto file system
func (s *FileStore) PutFileReader(name string, src io.Reader, length int64, contentType string) error {
	path := filepath.Join(s.dir, name)

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	written, err := io.Copy(f, src)
	if err != nil {
		return err
	}

	if written != length {
		return fmt.Errorf("got written %d bytes, expect %d", written, length)
	}

	return nil
}

// GeneratePostFileRequest return a PostFileRequest for uploading asset
func (s *FileStore) GeneratePostFileRequest(name string) (*PostFileRequest, error) {
	return &PostFileRequest{
		Action: strings.Join(
			[]string{s.postPrefix, "files", name},
			"/",
		),
	}, nil
}

// SignedURL returns a signed url with expiry date
func (s *FileStore) SignedURL(name string) (string, error) {
	if !s.IsSignatureRequired() {
		return fmt.Sprintf("%s/%s", s.prefix, name), nil
	}

	expiredAt := time.Now().Add(time.Minute * time.Duration(15))
	expiredAtStr := strconv.FormatInt(expiredAt.Unix(), 10)

	h := hmac.New(sha256.New, []byte(s.secret))
	io.WriteString(h, name)
	io.WriteString(h, expiredAtStr)

	buf := bytes.Buffer{}
	base64Encoder := base64.NewEncoder(base64.URLEncoding, &buf)
	base64Encoder.Write(h.Sum(nil))

	return fmt.Sprintf(
		"%s/%s?expiredAt=%s&signature=%s",
		s.prefix, name, expiredAtStr, buf.String(),
	), nil
}

// ParseSignature tries to parse the asset signature
func (s *FileStore) ParseSignature(signed string, name string, expiredAt time.Time) (valid bool, err error) {
	base64Decoder := base64.NewDecoder(base64.URLEncoding, strings.NewReader(signed))
	remoteSignature, err := ioutil.ReadAll(base64Decoder)
	if err != nil {
		log.Errorf("failed to decode asset url signature: %v", err)

		return false, errors.New("invalid signature")
	}

	h := hmac.New(sha256.New, []byte(s.secret))
	io.WriteString(h, name)
	io.WriteString(h, strconv.FormatInt(expiredAt.Unix(), 10))

	return !hmac.Equal(remoteSignature, h.Sum(nil)), nil
}

// IsSignatureRequired indicates whether a signature is required
func (s *FileStore) IsSignatureRequired() bool {
	return !s.public
}
