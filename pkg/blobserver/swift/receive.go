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

	"camlistore.org/pkg/blob"
	"camlistore.org/pkg/syncutil"
)

var receiveGate = syncutil.NewGate(5) // arbitrary

func (sto *swiftStorage) ReceiveBlob(b blob.Ref, source io.Reader) (sr blob.SizedRef, err error) {
    log.Println("[SWIFT] ReceiveBlob")
	removeGate.Start()
    defer removeGate.Done()
    obj, err := sto.client.ObjectCreate(sto.container, b.String(), true, "", "", nil)
	if err != nil {
		return sr, err
	}
    size, err := io.Copy(obj, source)
    if err != nil {
        return sr, err
    }
	return blob.SizedRef{Ref: b, Size: uint32(size)}, nil
}
