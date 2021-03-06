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
	"camlistore.org/pkg/blob"
	"camlistore.org/pkg/syncutil"
)

var removeGate = syncutil.NewGate(20) // arbitrary

func (sto *swiftStorage) RemoveBlobs(blobs []blob.Ref) error {
	objs := make([]string, len(blobs))
	for i := range blobs {
		objs[i] = blobs[i].String()
	}
	removeGate.Start()
	defer removeGate.Done()
	_, err := sto.client.BulkDelete(sto.container, objs)
	return err
}
