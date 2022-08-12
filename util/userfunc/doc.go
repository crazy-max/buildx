// Copyright 2022 Docker Buildx authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package userfunc implements a HCL extension that allows user-defined
// functions in HCL configuration.
//
// Using this extension requires some integration effort on the part of the
// calling application, to pass any declared functions into a HCL evaluation
// context after processing.
//
// The function declaration syntax looks like this:
//
//     function "foo" {
//       params = ["name"]
//       result = "Hello, ${name}!"
//     }
//
// When a user-defined function is called, the expression given for the "result"
// attribute is evaluated in an isolated evaluation context that defines variables
// named after the given parameter names.
//
// The block name "function" may be overridden by the calling application, if
// that default name conflicts with an existing block or attribute name in
// the application.
package userfunc
