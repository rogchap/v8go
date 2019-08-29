#include "v8go.h"

#include "v8.h"
#include "libplatform/libplatform.h"

#include <cstdlib>
#include <cstring>
#include <string>
#include <sstream>
#include <stdio.h>

using namespace v8;

auto allocator = ArrayBuffer::Allocator::NewDefaultAllocator();

extern "C" {

void Init() {
  std::unique_ptr<Platform> plt = platform::NewDefaultPlatform();
  V8::InitializePlatform(plt.get());
  V8::Initialize();
  return;
}

IsolatePtr NewIsolate() {
//  Isolate::CreateParams params;
//  params.array_buffer_allocator = allocator;
//  Isolate::New(params);
//  return static_cast<IsolatePtr>(Isolate::New(params));
v8::Isolate::CreateParams create_params;
  create_params.array_buffer_allocator =
      v8::ArrayBuffer::Allocator::NewDefaultAllocator();
  v8::Isolate* isolate = v8::Isolate::New(create_params);
  isolate->Dispose();
  v8::V8::Dispose();
  v8::V8::ShutdownPlatform();
  delete create_params.array_buffer_allocator;
return nullptr;
}

void IsolateRelease(IsolatePtr ptr) {
  if (ptr == nullptr) {
    return;
  }
  Isolate* iso = static_cast<Isolate*>(ptr);
  iso->Dispose();
}
  
const char* Version() {
  return V8::GetVersion();
}

}

