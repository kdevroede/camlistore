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

/*
Package swift registers the "swift" blobserver storage type, storing
blobs in an Openstack Swift storage container.

Example low-level config:

     "/r1/": {
         "handler": "swift",
         "handlerArgs": {
            "container": "...",
            "user_name": "...",
            "api_key": "...",
            "auth_url": "...",
            "tenant": "...",
            "skipStartupCheck": false
          }
     },

*/
package swift

import (
	"fmt"

	"camlistore.org/pkg/blobserver"
	"camlistore.org/pkg/jsonconfig"

	"camlistore.org/third_party/github.com/ncw/swift"
)

type swiftStorage struct {
	client    *swift.Connection
	container string
}

func newFromConfig(_ blobserver.Loader, config jsonconfig.Obj) (blobserver.Storage, error) {
	client := &swift.Connection{
		UserName: config.RequiredString("user_name"),
		Tenant:   config.RequiredString("tenant"),
		ApiKey:   config.RequiredString("secret"),
		AuthUrl:  config.RequiredString("auth_url"),
	}
	if err := client.Authenticate(); err != nil {
		return nil, fmt.Errorf("Authentication failure on swift: %v", err)
	}
	sto := &swiftStorage{
		client:    client,
		container: config.OptionalString("container", "default"),
	}
	skipStartupCheck := config.OptionalBool("skipStartupCheck", false)
	if err := config.Validate(); err != nil {
		return nil, err
	}
	if !skipStartupCheck {
		containers, err := client.ContainerNames(nil)
		if err != nil {
			return nil, fmt.Errorf("Failed to get container list from swift: %v", err)
		}
		haveContainer := make(map[string]bool)
		for _, c := range containers {
			haveContainer[c] = true
		}
		if !haveContainer[sto.container] {
			if err := client.ContainerCreate(sto.container, nil); err != nil {
				return nil, fmt.Errorf("Swift bucket %s doesn't exist and it's impossible to create it: ", sto.container, err)
			}
		}
	}
	return sto, nil
}

func init() {
	blobserver.RegisterStorageConstructor("swift", blobserver.StorageConstructor(newFromConfig))
}
