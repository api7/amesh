-- Copyright 2022 The Amesh Authors
--
-- Licensed under the Apache License, Version 2.0 (the "License");
-- you may not use this file except in compliance with the License.
-- You may obtain a copy of the License at
--
--    http://www.apache.org/licenses/LICENSE-2.0
--
-- Unless required by applicable law or agreed to in writing, software
-- distributed under the License is distributed on an "AS IS" BASIS,
-- WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
-- See the License for the specific language governing permissions and
-- limitations under the License.
--
local ffi = require "ffi"

local amesh = ffi.load("/amesh/libxds.so")
-- local amesh = ffi.load("./libxds.so")

ffi.cdef[[
typedef signed char GoInt8;
typedef unsigned char GoUint8;
typedef short GoInt16;
typedef unsigned short GoUint16;
typedef int GoInt32;
typedef unsigned int GoUint32;
typedef long long GoInt64;
typedef unsigned long long GoUint64;
typedef GoInt64 GoInt;
typedef GoUint64 GoUint;
typedef float GoFloat32;
typedef double GoFloat64;
typedef float _Complex GoComplex64;
typedef double _Complex GoComplex128;

typedef struct { const char *p; ptrdiff_t n; } _GoString_;
typedef _GoString_ GoString;

typedef void *GoMap;
typedef void *GoChan;
typedef struct { void *t; void *v; } GoInterface;
typedef struct { void *data; GoInt len; GoInt cap; } GoSlice;

extern void Log(GoString msg);
extern void StartAmesh(GoString src);

]]


xdsSrc = ffi.new("GoString")
local s = "grpc://istiod.istio-system.svc.cluster.local:15010"
xdsSrc.p = s;
xdsSrc.n = #s;

amesh.StartAmesh(xdsSrc)

startedStr = ffi.new("GoString")
local msg = "amesh started"
startedStr.p = msg;
startedStr.n = #msg;
amesh.Log(startedStr)
while true do

end