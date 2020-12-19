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

typedef struct {
  Persistent<Context> ptr;
  Isolate* iso;
} m_ctx;

typedef struct {
  Persistent<Value> ptr;
  m_ctx* ctx_ptr;
} m_value;

const char* CopyString(std::string str) {
  int len = str.length();
  char *mem = (char*)malloc(len+1);
  memcpy(mem, str.data(), len);
  mem[len] = 0;
  return mem;
}

const char* CopyString(String::Utf8Value& value) {
  if (value.length() == 0) {
    return nullptr;
  }
  return CopyString(*value);
}

RtnError ExceptionError(TryCatch& try_catch, Isolate* iso, Local<Context> ctx) {
  Locker locker(iso);
  Isolate::Scope isolate_scope(iso);
  HandleScope handle_scope(iso);

  RtnError rtn = {nullptr, nullptr, nullptr};

  if (try_catch.HasTerminated()) {
    rtn.msg = CopyString("ExecutionTerminated: script execution has been terminated");
    return rtn;
  }

  String::Utf8Value exception(iso, try_catch.Exception());
  rtn.msg = CopyString(exception);

  Local<Message> msg = try_catch.Message();
  if (!msg.IsEmpty()) {
    String::Utf8Value origin(iso, msg->GetScriptOrigin().ResourceName());
    std::ostringstream sb;
    sb << *origin;
    Maybe<int> line = try_catch.Message()->GetLineNumber(ctx);
    if (line.IsJust()) {
      sb << ":" << line.ToChecked();
    }
    Maybe<int> start = try_catch.Message()->GetStartColumn(ctx);
    if (start.IsJust()) {
      sb << ":" << start.ToChecked() + 1; // + 1 to match output from stack trace
    }
    rtn.location = CopyString(sb.str());
  }
 
  MaybeLocal<Value> mstack = try_catch.StackTrace(ctx);
  if (!mstack.IsEmpty()) {
    String::Utf8Value stack(iso, mstack.ToLocalChecked());
    rtn.stack = CopyString(stack);
  }
  
  return rtn;
}

extern "C"
{

/********** Isolate **********/

void Init() {
#ifdef _WIN32
    V8::InitializeExternalStartupData(".");
#endif
    V8::InitializePlatform(default_platform.get());
    V8::Initialize();
    return;
}

IsolatePtr NewIsolate() {
    Isolate::CreateParams params;
    params.array_buffer_allocator = default_allocator;
    return static_cast<IsolatePtr>(Isolate::New(params));
}

void IsolateDispose(IsolatePtr ptr) {
    if (ptr == nullptr) {
        return;
    }
    Isolate* iso = static_cast<Isolate*>(ptr);
    iso->Dispose();
}

void IsolateTerminateExecution(IsolatePtr ptr) {
    Isolate* iso = static_cast<Isolate*>(ptr);
    iso->TerminateExecution();
}

IsolateHStatistics IsolationGetHeapStatistics(IsolatePtr ptr) {
  if (ptr == nullptr) {
    return IsolateHStatistics{0};
  }
  Isolate* iso = static_cast<Isolate*>(ptr);
  v8::HeapStatistics hs;
  iso->GetHeapStatistics(&hs);
  
  return IsolateHStatistics{
    hs.total_heap_size(),
    hs.total_heap_size_executable(),
    hs.total_physical_size(),
    hs.total_available_size(),
    hs.used_heap_size(),
    hs.heap_size_limit(),
    hs.malloced_memory(),
    hs.external_memory(),
    hs.peak_malloced_memory(),
    hs.number_of_native_contexts(),
    hs.number_of_detached_contexts()
  };
}

/********** Context **********/

ContextPtr NewContext(IsolatePtr ptr) {
    Isolate* iso = static_cast<Isolate*>(ptr);
    Locker locker(iso);
    Isolate::Scope isolate_scope(iso);
    HandleScope handle_scope(iso);
  
    iso->SetCaptureStackTraceForUncaughtExceptions(true);
    
    m_ctx* ctx = new m_ctx;
    ctx->ptr.Reset(iso, Context::New(iso));
    ctx->iso = iso;
    return static_cast<ContextPtr>(ctx);
}

RtnValue RunScript(ContextPtr ctx_ptr, const char* source, const char* origin) {
    m_ctx* ctx = static_cast<m_ctx*>(ctx_ptr);
    Isolate* iso = ctx->iso;
    Locker locker(iso);
    Isolate::Scope isolate_scope(iso);
    HandleScope handle_scope(iso);
    TryCatch try_catch(iso);

    Local<Context> local_ctx = ctx->ptr.Get(iso);
    Context::Scope context_scope(local_ctx);

    Local<String> src = String::NewFromUtf8(iso, source, NewStringType::kNormal).ToLocalChecked();
    Local<String> ogn = String::NewFromUtf8(iso, origin, NewStringType::kNormal).ToLocalChecked();

    RtnValue rtn = { nullptr, nullptr };

    ScriptOrigin script_origin(ogn);
    MaybeLocal<Script> script = Script::Compile(local_ctx, src, &script_origin);
    if (script.IsEmpty()) {
      rtn.error = ExceptionError(try_catch, iso, local_ctx);
      return rtn;
    } 
    MaybeLocal<v8::Value> result = script.ToLocalChecked()->Run(local_ctx);
    if (result.IsEmpty()) {
      rtn.error = ExceptionError(try_catch, iso, local_ctx);
      return rtn;
    }
    m_value* val = new m_value;
    val->ctx_ptr = ctx;
    val->ptr.Reset(iso, Persistent<Value>(iso, result.ToLocalChecked()));

    rtn.value = static_cast<ValuePtr>(val);
    return rtn;
}

void ContextDispose(ContextPtr ptr) {
    if (ptr == nullptr) {
        return;
    }
    m_ctx* ctx = static_cast<m_ctx*>(ptr);
    if (ctx == nullptr) {
        return;
    }
    ctx->ptr.Reset(); 
    delete ctx;
} 

/********** Value **********/

#define LOCAL_VALUE(ptr) \
    m_value* val = static_cast<m_value*>(ptr); \
    m_ctx* ctx = val->ctx_ptr; \
    Isolate* iso = ctx->iso; \
    Locker locker(iso); \
    Isolate::Scope isolate_scope(iso); \
    HandleScope handle_scope(iso); \
    Context::Scope context_scope(ctx->ptr.Get(iso)); \
    Local<Value> value = val->ptr.Get(iso);

void ValueDispose(ValuePtr ptr) {
    delete static_cast<m_value*>(ptr);
}

const char* ValueToString(ValuePtr ptr) {
    LOCAL_VALUE(ptr);
    String::Utf8Value utf8(iso, value);

    return CopyString(utf8);
} 

int ValueIsUndefined(ValuePtr ptr) {
    LOCAL_VALUE(ptr);
    return value->IsUndefined() ? 1 : 0;
}

int ValueIsNull(ValuePtr ptr) {
    LOCAL_VALUE(ptr);
    return value->IsNull() ? 1 : 0;
}

int ValueIsNullOrUndefined(ValuePtr ptr) {
    LOCAL_VALUE(ptr);
    return value->IsNullOrUndefined() ? 1 : 0;
}

int ValueIsTrue(ValuePtr ptr) {
    LOCAL_VALUE(ptr);
    return value->IsTrue() ? 1 : 0;
}

int ValueIsFalse(ValuePtr ptr) {
    LOCAL_VALUE(ptr);
    return value->IsFalse() ? 1 : 0;
}

int ValueIsString(ValuePtr ptr) {
    LOCAL_VALUE(ptr);

    return value->IsString() ? 1 : 0;
}


/********** Version **********/
  
const char* Version() {
    return V8::GetVersion();
}

}

