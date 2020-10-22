/*
Copyright 2020 CaiCloud, Inc. All rights reserved.

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

package filters

import (
	"strings"

	v1 "k8s.io/api/core/v1"
)

type EventFilter interface {
	Filter(event *v1.Event) bool
}

type EventTypeFilter struct {
	AllowedTypes []string
}

func NewEventTypeFilter(allowedTypes []string) *EventTypeFilter {
	return &EventTypeFilter{
		AllowedTypes: allowedTypes,
	}
}

func (e *EventTypeFilter) Filter(event *v1.Event) bool {
	for _, allowedType := range e.AllowedTypes {
		if strings.EqualFold(event.Type, allowedType) {
			return true
		}
	}
	return false
}
