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
	"fmt"

	"camlistore.org/pkg/blob"
	"camlistore.org/pkg/syncutil"
	"camlistore.org/third_party/github.com/ncw/swift"
)

var statGate = syncutil.NewGate(20) // arbitrary

func (sto *swiftStorage) StatBlobs(dest chan<- blob.SizedRef, blobs []blob.Ref) error {
	var wg syncutil.Group

	for _, br := range blobs {
		br := br
		statGate.Start()
		wg.Go(func() error {
			defer statGate.Done()

			obj, _, err := sto.client.Object(sto.container, br.String())
			if err == nil {
				dest <- blob.SizedRef{Ref: br, Size: uint32(obj.Bytes)}
				return nil
			}
			if err == swift.ObjectNotFound {
				return nil
			}
			return fmt.Errorf("error statting %v: %v", br, err)
		})
	}
	return wg.Err()
}
