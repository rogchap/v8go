#include "v8go.h"

#include "v8.h"
#include "libplatform/libplatform.h"

#include <cstdlib>
#include <cstring>
#include <string>
#include <sstream>
#include <stdio.h>

using namespace v8;

auto default_platform = platform::NewDefaultPlatform();
auto default_allocator = ArrayBuffer::Allocator::NewDefaultAllocator();

extern "C" {

/********** Isolate **********/

void Init() {
    V8::InitializePlatform(default_platform.get());
    V8::Initialize();
    return;
}

IsolatePtr NewIsolate() {
    Isolate::CreateParams params;
    params.array_buffer_allocator = default_allocator;
    Isolate::New(params);
    return static_cast<IsolatePtr>(Isolate::New(params));
}

void IsolateDispose(IsolatePtr ptr) {
    if (ptr == nullptr) {
        return;
    }
    Isolate* iso = static_cast<Isolate*>(ptr);
    iso->Dispose();
}

void TerminateExecution(IsolatePtr ptr) {
    Isolate* iso = static_cast<Isolate*>(ptr);
    iso->TerminateExecution();
}

/********** Context **********/

ContextPtr NewContext(IsolatePtr ptr) {
    Isolate* iso = static_cast<Isolate*>(ptr);
    Locker locker(iso);
    Isolate::Scope isolate_scope(iso);
    HandleScope handle_scope(iso);
    iso->SetCaptureStackTraceForUncaughtExceptions(true);
    Local<Context> ctx = Context::New(iso);
    return static_cast<ContextPtr>(std::move(&ctx));
}

void RunScript(ContextPtr ctx_ptr, const char* source, const char* origin) {
    Local<Context> ctx = *(static_cast<Local<Context>*>(ctx_ptr));
    Isolate* iso = ctx->GetIsolate();
    Locker locker(iso);
    Isolate::Scope isolate_scope(iso);
    HandleScope handle_scope(iso);
    TryCatch try_catch(iso);

    Local<Script> script = Script::Compile(ctx, String::NewFromUtf8(iso, "1 + 1").ToLocalChecked()).ToLocalChecked();
    v8::Local<v8::Value> result = script->Run(ctx).ToLocalChecked();
    // Convert the result to an UTF8 string and print it.
    v8::String::Utf8Value utf8(iso, result);
    printf("%s\n", *utf8);
}


/********** Version **********/
  
const char* Version() {
    return V8::GetVersion();
}

}


int main(int argc, char* argv[]) {
    Init();
    auto i = NewIsolate();
    auto c = NewContext(i);
    RunScript(c, "", "");
}
