/*
Copyright 2011 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package swift

import (
	"io"
	"log"
	"strconv"

	"camlistore.org/pkg/blob"
)

func (sto *swiftStorage) FetchStreaming(blob blob.Ref) (file io.ReadCloser, size uint32, err error) {
    log.Println("[SWIFT] FetchStreaming")
    obj, h, err := sto.client.ObjectOpen(sto.container, blob.String(), false, nil);
    if err != nil {
        return nil, size, err
    }
    log.Println("[SWIFT] FetchStreaming", h)
    sz, err := strconv.ParseInt(h["Content-Length"], 10, 64)
    if err != nil {
        obj.Close()
        return nil, size, err
    }
    return obj, uint32(sz), err
}
