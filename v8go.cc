// Copyright 2019 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

#include "v8go.h"

#include <stdio.h>

#include <cstdlib>
#include <cstring>
#include <iostream>
#include <sstream>
#include <string>
#include <vector>

#include "_cgo_export.h"

using namespace v8;

auto default_platform = platform::NewDefaultPlatform();
auto default_allocator = ArrayBuffer::Allocator::NewDefaultAllocator();

const int ScriptCompilerNoCompileOptions = ScriptCompiler::kNoCompileOptions;
const int ScriptCompilerConsumeCodeCache = ScriptCompiler::kConsumeCodeCache;
const int ScriptCompilerEagerCompile = ScriptCompiler::kEagerCompile;

struct m_ctx {
  Isolate* iso;
  StartupData* startup_data;
  std::vector<m_value*> vals;
  std::vector<m_unboundScript*> unboundScripts;
  Persistent<Context> ptr;
};

struct m_value {
  Isolate* iso;
  m_ctx* ctx;
  Persistent<Value, CopyablePersistentTraits<Value>> ptr;
};

struct m_template {
  Isolate* iso;
  Persistent<Template> ptr;
};

struct m_unboundScript {
  Persistent<UnboundScript> ptr;
};

const char* CopyString(std::string str) {
  int len = str.length();
  char* mem = (char*)malloc(len + 1);
  memcpy(mem, str.data(), len);
  mem[len] = 0;
  return mem;
}

const char* CopyString(String::Utf8Value& value) {
  if (value.length() == 0) {
    return nullptr;
  }
  return CopyString(std::string(*value, value.length()));
}

static RtnError ExceptionError(TryCatch& try_catch,
                               Isolate* iso,
                               Local<Context> ctx) {
  HandleScope handle_scope(iso);

  RtnError rtn = {nullptr, nullptr, nullptr};

  if (try_catch.HasTerminated()) {
    rtn.msg =
        CopyString("ExecutionTerminated: script execution has been terminated");
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
      sb << ":"
         << start.ToChecked() + 1;  // + 1 to match output from stack trace
    }
    rtn.location = CopyString(sb.str());
  }

  Local<Value> mstack;
  if (try_catch.StackTrace(ctx).ToLocal(&mstack)) {
    String::Utf8Value stack(iso, mstack);
    rtn.stack = CopyString(stack);
  }

  return rtn;
}

m_value* tracked_value(m_ctx* ctx, m_value* val) {
  // (rogchap) we track values against a context so that when the context is
  // closed (either manually or GC'd by Go) we can also release all the
  // values associated with the context; previously the Go GC would not run
  // quickly enough, as it has no understanding of the C memory allocation size.
  // By doing so we hold pointers to all values that are created/returned to Go
  // until the context is released; this is a compromise.
  // Ideally we would be able to delete the value object and cancel the
  // finalizer on the Go side, but we currently don't pass the Go ptr, but
  // rather the C ptr. A potential future iteration would be to use an
  // unordered_map, where we could do O(1) lookups for the value, but then know
  // if the object has been finalized or not by being in the map or not. This
  // would require some ref id for the value rather than passing the ptr between
  // Go <--> C, which would be a significant change, as there are places where
  // we get the context from the value, but if we then need the context to get
  // the value, we would be in a circular bind.
  ctx->vals.push_back(val);

  return val;
}

m_unboundScript* tracked_unbound_script(m_ctx* ctx, m_unboundScript* us) {
  ctx->unboundScripts.push_back(us);

  return us;
}

extern "C" {

/********** Isolate **********/

#define ISOLATE_SCOPE(iso)           \
  Locker locker(iso);                \
  Isolate::Scope isolate_scope(iso); \
  HandleScope handle_scope(iso);

#define ISOLATE_SCOPE_INTERNAL_CONTEXT(iso) \
  ISOLATE_SCOPE(iso);                       \
  m_ctx* ctx = isolateInternalContext(iso);

void Init() {
#ifdef _WIN32
  V8::InitializeExternalStartupData(".");
#endif
  V8::InitializePlatform(default_platform.get());
  V8::Initialize();
  return;
}

IsolatePtr NewIsolate(IsolateOptions options) {
  Isolate::CreateParams params;
  params.array_buffer_allocator = default_allocator;

  StartupData* startup_data;
  if (options.snapshot_blob_data != nullptr) {
    startup_data = new StartupData{options.snapshot_blob_data,
                                   options.snapshot_blob_raw_size};
    params.snapshot_blob = startup_data;
  } else {
    startup_data = nullptr;
  }

  Isolate* iso = Isolate::New(params);
  Locker locker(iso);
  Isolate::Scope isolate_scope(iso);
  HandleScope handle_scope(iso);

  iso->SetCaptureStackTraceForUncaughtExceptions(true);

  // Create a Context for internal use
  m_ctx* ctx = new m_ctx;
  ctx->ptr.Reset(iso, Context::New(iso));
  ctx->iso = iso;
  ctx->startup_data = startup_data;
  iso->SetData(0, ctx);

  return iso;
}

static inline m_ctx* isolateInternalContext(Isolate* iso) {
  return static_cast<m_ctx*>(iso->GetData(0));
}

void IsolatePerformMicrotaskCheckpoint(IsolatePtr iso) {
  ISOLATE_SCOPE(iso)
  iso->PerformMicrotaskCheckpoint();
}

void IsolateDispose(IsolatePtr iso) {
  if (iso == nullptr) {
    return;
  }
  ContextFree(isolateInternalContext(iso));

  iso->Dispose();
}

void IsolateTerminateExecution(IsolatePtr iso) {
  iso->TerminateExecution();
}

int IsolateIsExecutionTerminating(IsolatePtr iso) {
  return iso->IsExecutionTerminating();
}

IsolateHStatistics IsolationGetHeapStatistics(IsolatePtr iso) {
  if (iso == nullptr) {
    return IsolateHStatistics{0};
  }
  v8::HeapStatistics hs;
  iso->GetHeapStatistics(&hs);

  return IsolateHStatistics{hs.total_heap_size(),
                            hs.total_heap_size_executable(),
                            hs.total_physical_size(),
                            hs.total_available_size(),
                            hs.used_heap_size(),
                            hs.heap_size_limit(),
                            hs.malloced_memory(),
                            hs.external_memory(),
                            hs.peak_malloced_memory(),
                            hs.number_of_native_contexts(),
                            hs.number_of_detached_contexts()};
}

RtnUnboundScript IsolateCompileUnboundScript(IsolatePtr iso,
                                             const char* s,
                                             const char* o,
                                             CompileOptions opts) {
  ISOLATE_SCOPE_INTERNAL_CONTEXT(iso);
  TryCatch try_catch(iso);
  Local<Context> local_ctx = ctx->ptr.Get(iso);
  Context::Scope context_scope(local_ctx);

  RtnUnboundScript rtn = {};

  Local<String> src =
      String::NewFromUtf8(iso, s, NewStringType::kNormal).ToLocalChecked();
  Local<String> ogn =
      String::NewFromUtf8(iso, o, NewStringType::kNormal).ToLocalChecked();

  ScriptCompiler::CompileOptions option =
      static_cast<ScriptCompiler::CompileOptions>(opts.compileOption);

  ScriptCompiler::CachedData* cached_data = nullptr;

  if (opts.cachedData.data) {
    cached_data = new ScriptCompiler::CachedData(opts.cachedData.data,
                                                 opts.cachedData.length);
  }

  ScriptOrigin script_origin(ogn);

  ScriptCompiler::Source source(src, script_origin, cached_data);

  Local<UnboundScript> unbound_script;
  if (!ScriptCompiler::CompileUnboundScript(iso, &source, option)
           .ToLocal(&unbound_script)) {
    rtn.error = ExceptionError(try_catch, iso, local_ctx);
    return rtn;
  };

  if (cached_data) {
    rtn.cachedDataRejected = cached_data->rejected;
  }

  m_unboundScript* us = new m_unboundScript;
  us->ptr.Reset(iso, unbound_script);
  rtn.ptr = tracked_unbound_script(ctx, us);
  return rtn;
}

/********** SnapshotCreator **********/

RtnSnapshotCreator NewSnapshotCreator() {
  RtnSnapshotCreator rtn = {};
  SnapshotCreator* creator = new SnapshotCreator;
  Isolate* iso = creator->GetIsolate();
  rtn.creator = creator;
  rtn.iso = iso;

  return rtn;
}

void DeleteSnapshotCreator(SnapshotCreatorPtr snapshotCreator) {
  delete snapshotCreator;
}

void SetDefaultContext(SnapshotCreatorPtr snapshotCreator, ContextPtr ctx) {
  Isolate* iso = ctx->iso;
  Locker locker(iso);
  Isolate::Scope isolate_scope(iso);
  HandleScope handle_scope(iso);
  Local<Context> local_ctx = ctx->ptr.Get(iso);
  Context::Scope context_scope(local_ctx);

  ContextFree(ctx);

  snapshotCreator->SetDefaultContext(local_ctx);
}

size_t AddContext(SnapshotCreatorPtr snapshotCreator, ContextPtr ctx) {
  Isolate* iso = ctx->iso;
  Locker locker(iso);
  Isolate::Scope isolate_scope(iso);
  HandleScope handle_scope(iso);
  Local<Context> local_ctx = ctx->ptr.Get(iso);
  Context::Scope context_scope(local_ctx);

  ContextFree(ctx);

  return snapshotCreator->AddContext(local_ctx);
}

RtnSnapshotBlob* CreateBlob(SnapshotCreatorPtr snapshotCreator,
                            int function_code_handling) {
  //  kKeep - keeps any compiled functions
  //  kClear - does not keep any compiled functions
  StartupData startup_data = snapshotCreator->CreateBlob(
      SnapshotCreator::FunctionCodeHandling(function_code_handling));

  RtnSnapshotBlob* rtn = new RtnSnapshotBlob;
  rtn->data = startup_data.data;
  rtn->raw_size = startup_data.raw_size;
  delete snapshotCreator;
  return rtn;
}

void SnapshotBlobDelete(RtnSnapshotBlob* ptr) {
  delete[] ptr->data;
  delete ptr;
}

/********** Exceptions & Errors **********/

ValuePtr IsolateThrowException(IsolatePtr iso, ValuePtr value) {
  ISOLATE_SCOPE(iso);
  m_ctx* ctx = value->ctx;

  Local<Value> throw_ret_val = iso->ThrowException(value->ptr.Get(iso));

  m_value* new_val = new m_value;
  new_val->iso = iso;
  new_val->ctx = ctx;
  new_val->ptr =
      Persistent<Value, CopyablePersistentTraits<Value>>(iso, throw_ret_val);

  return tracked_value(ctx, new_val);
}

/********** CpuProfiler **********/

CPUProfiler* NewCPUProfiler(IsolatePtr iso_ptr) {
  Isolate* iso = static_cast<Isolate*>(iso_ptr);
  Locker locker(iso);
  Isolate::Scope isolate_scope(iso);
  HandleScope handle_scope(iso);

  CPUProfiler* c = new CPUProfiler;
  c->iso = iso;
  c->ptr = CpuProfiler::New(iso);
  return c;
}

void CPUProfilerDispose(CPUProfiler* profiler) {
  if (profiler->ptr == nullptr) {
    return;
  }
  profiler->ptr->Dispose();

  delete profiler;
}

void CPUProfilerStartProfiling(CPUProfiler* profiler, const char* title) {
  if (profiler->iso == nullptr) {
    return;
  }

  Locker locker(profiler->iso);
  Isolate::Scope isolate_scope(profiler->iso);
  HandleScope handle_scope(profiler->iso);

  Local<String> title_str =
      String::NewFromUtf8(profiler->iso, title, NewStringType::kNormal)
          .ToLocalChecked();
  profiler->ptr->StartProfiling(title_str);
}

CPUProfileNode* NewCPUProfileNode(const CpuProfileNode* ptr_) {
  int count = ptr_->GetChildrenCount();
  CPUProfileNode** children = new CPUProfileNode*[count];
  for (int i = 0; i < count; ++i) {
    children[i] = NewCPUProfileNode(ptr_->GetChild(i));
  }

  CPUProfileNode* root = new CPUProfileNode{
      ptr_,
      ptr_->GetScriptResourceNameStr(),
      ptr_->GetFunctionNameStr(),
      ptr_->GetLineNumber(),
      ptr_->GetColumnNumber(),
      count,
      children,
  };
  return root;
}

CPUProfile* CPUProfilerStopProfiling(CPUProfiler* profiler, const char* title) {
  if (profiler->iso == nullptr) {
    return nullptr;
  }

  Locker locker(profiler->iso);
  Isolate::Scope isolate_scope(profiler->iso);
  HandleScope handle_scope(profiler->iso);

  Local<String> title_str =
      String::NewFromUtf8(profiler->iso, title, NewStringType::kNormal)
          .ToLocalChecked();

  CPUProfile* profile = new CPUProfile;
  profile->ptr = profiler->ptr->StopProfiling(title_str);

  Local<String> str = profile->ptr->GetTitle();
  String::Utf8Value t(profiler->iso, str);
  profile->title = CopyString(t);

  CPUProfileNode* root = NewCPUProfileNode(profile->ptr->GetTopDownRoot());
  profile->root = root;

  profile->startTime = profile->ptr->GetStartTime();
  profile->endTime = profile->ptr->GetEndTime();

  return profile;
}

void CPUProfileNodeDelete(CPUProfileNode* node) {
  for (int i = 0; i < node->childrenCount; ++i) {
    CPUProfileNodeDelete(node->children[i]);
  }

  delete[] node->children;
  delete node;
}

void CPUProfileDelete(CPUProfile* profile) {
  if (profile->ptr == nullptr) {
    return;
  }
  profile->ptr->Delete();
  free((void*)profile->title);

  CPUProfileNodeDelete(profile->root);

  delete profile;
}

/********** Template **********/

#define LOCAL_TEMPLATE(tmpl_ptr)     \
  Isolate* iso = tmpl_ptr->iso;      \
  Locker locker(iso);                \
  Isolate::Scope isolate_scope(iso); \
  HandleScope handle_scope(iso);     \
  Local<Template> tmpl = tmpl_ptr->ptr.Get(iso);

void TemplateFreeWrapper(TemplatePtr tmpl) {
  tmpl->ptr.Empty();  // Just does `val_ = 0;` without calling V8::DisposeGlobal
  delete tmpl;
}

void TemplateSetValue(TemplatePtr ptr,
                      const char* name,
                      ValuePtr val,
                      int attributes) {
  LOCAL_TEMPLATE(ptr);

  Local<String> prop_name =
      String::NewFromUtf8(iso, name, NewStringType::kNormal).ToLocalChecked();
  tmpl->Set(prop_name, val->ptr.Get(iso), (PropertyAttribute)attributes);
}

void TemplateSetTemplate(TemplatePtr ptr,
                         const char* name,
                         TemplatePtr obj,
                         int attributes) {
  LOCAL_TEMPLATE(ptr);

  Local<String> prop_name =
      String::NewFromUtf8(iso, name, NewStringType::kNormal).ToLocalChecked();
  tmpl->Set(prop_name, obj->ptr.Get(iso), (PropertyAttribute)attributes);
}

/********** ObjectTemplate **********/

TemplatePtr NewObjectTemplate(IsolatePtr iso) {
  Locker locker(iso);
  Isolate::Scope isolate_scope(iso);
  HandleScope handle_scope(iso);

  m_template* ot = new m_template;
  ot->iso = iso;
  ot->ptr.Reset(iso, ObjectTemplate::New(iso));
  return ot;
}

RtnValue ObjectTemplateNewInstance(TemplatePtr ptr, ContextPtr ctx) {
  LOCAL_TEMPLATE(ptr);
  TryCatch try_catch(iso);
  Local<Context> local_ctx = ctx->ptr.Get(iso);
  Context::Scope context_scope(local_ctx);

  RtnValue rtn = {};

  Local<ObjectTemplate> obj_tmpl = tmpl.As<ObjectTemplate>();
  Local<Object> obj;
  if (!obj_tmpl->NewInstance(local_ctx).ToLocal(&obj)) {
    rtn.error = ExceptionError(try_catch, iso, local_ctx);
    return rtn;
  }

  m_value* val = new m_value;
  val->iso = iso;
  val->ctx = ctx;
  val->ptr = Persistent<Value, CopyablePersistentTraits<Value>>(iso, obj);
  rtn.value = tracked_value(ctx, val);
  return rtn;
}

void ObjectTemplateSetInternalFieldCount(TemplatePtr ptr, int field_count) {
  LOCAL_TEMPLATE(ptr);

  Local<ObjectTemplate> obj_tmpl = tmpl.As<ObjectTemplate>();
  obj_tmpl->SetInternalFieldCount(field_count);
}

int ObjectTemplateInternalFieldCount(TemplatePtr ptr) {
  LOCAL_TEMPLATE(ptr);

  Local<ObjectTemplate> obj_tmpl = tmpl.As<ObjectTemplate>();
  return obj_tmpl->InternalFieldCount();
}

/********** FunctionTemplate **********/

static void FunctionTemplateCallback(const FunctionCallbackInfo<Value>& info) {
  Isolate* iso = info.GetIsolate();
  ISOLATE_SCOPE(iso);

  // This callback function can be called from any Context, which we only know
  // at runtime. We extract the Context reference from the embedder data so that
  // we can use the context registry to match the Context on the Go side
  Local<Context> local_ctx = iso->GetCurrentContext();
  int ctx_ref = local_ctx->GetEmbedderData(1).As<Integer>()->Value();
  m_ctx* ctx = goContext(ctx_ref);

  int callback_ref = info.Data().As<Integer>()->Value();

  m_value* _this = new m_value;
  _this->iso = iso;
  _this->ctx = ctx;
  _this->ptr.Reset(iso, Persistent<Value, CopyablePersistentTraits<Value>>(
                            iso, info.This()));

  int args_count = info.Length();
  ValuePtr thisAndArgs[args_count + 1];
  thisAndArgs[0] = tracked_value(ctx, _this);
  ValuePtr* args = thisAndArgs + 1;
  for (int i = 0; i < args_count; i++) {
    m_value* val = new m_value;
    val->iso = iso;
    val->ctx = ctx;
    val->ptr.Reset(
        iso, Persistent<Value, CopyablePersistentTraits<Value>>(iso, info[i]));
    args[i] = tracked_value(ctx, val);
  }

  ValuePtr val =
      goFunctionCallback(ctx_ref, callback_ref, thisAndArgs, args_count);
  if (val != nullptr) {
    info.GetReturnValue().Set(val->ptr.Get(iso));
  } else {
    info.GetReturnValue().SetUndefined();
  }
}

TemplatePtr NewFunctionTemplate(IsolatePtr iso, int callback_ref) {
  Locker locker(iso);
  Isolate::Scope isolate_scope(iso);
  HandleScope handle_scope(iso);

  // (rogchap) We only need to store one value, callback_ref, into the
  // C++ callback function data, but if we needed to store more items we could
  // use an V8::Array; this would require the internal context from
  // iso->GetData(0)
  Local<Integer> cbData = Integer::New(iso, callback_ref);

  m_template* ot = new m_template;
  ot->iso = iso;
  ot->ptr.Reset(iso,
                FunctionTemplate::New(iso, FunctionTemplateCallback, cbData));
  return ot;
}

RtnValue FunctionTemplateGetFunction(TemplatePtr ptr, ContextPtr ctx) {
  LOCAL_TEMPLATE(ptr);
  TryCatch try_catch(iso);
  Local<Context> local_ctx = ctx->ptr.Get(iso);
  Context::Scope context_scope(local_ctx);

  Local<FunctionTemplate> fn_tmpl = tmpl.As<FunctionTemplate>();
  RtnValue rtn = {};
  Local<Function> fn;
  if (!fn_tmpl->GetFunction(local_ctx).ToLocal(&fn)) {
    rtn.error = ExceptionError(try_catch, iso, local_ctx);
    return rtn;
  }

  m_value* val = new m_value;
  val->iso = iso;
  val->ctx = ctx;
  val->ptr = Persistent<Value, CopyablePersistentTraits<Value>>(iso, fn);
  rtn.value = tracked_value(ctx, val);
  return rtn;
}

/********** Context **********/

#define LOCAL_CONTEXT(ctx)                      \
  Isolate* iso = ctx->iso;                      \
  Locker locker(iso);                           \
  Isolate::Scope isolate_scope(iso);            \
  HandleScope handle_scope(iso);                \
  TryCatch try_catch(iso);                      \
  Local<Context> local_ctx = ctx->ptr.Get(iso); \
  Context::Scope context_scope(local_ctx);

ContextPtr NewContext(IsolatePtr iso,
                      TemplatePtr global_template_ptr,
                      int ref) {
  Locker locker(iso);
  Isolate::Scope isolate_scope(iso);
  HandleScope handle_scope(iso);

  Local<ObjectTemplate> global_template;
  if (global_template_ptr != nullptr) {
    global_template = global_template_ptr->ptr.Get(iso).As<ObjectTemplate>();
  } else {
    global_template = ObjectTemplate::New(iso);
  }

  // For function callbacks we need a reference to the context, but because of
  // the complexities of C -> Go function pointers, we store a reference to the
  // context as a simple integer identifier; this can then be used on the Go
  // side to lookup the context in the context registry. We use slot 1 as slot 0
  // has special meaning for the Chrome debugger.
  Local<Context> local_ctx = Context::New(iso, nullptr, global_template);
  local_ctx->SetEmbedderData(1, Integer::New(iso, ref));

  m_ctx* ctx = new m_ctx;
  ctx->ptr.Reset(iso, local_ctx);
  ctx->iso = iso;
  ctx->startup_data = nullptr;
  return ctx;
}

ContextPtr NewContextFromSnapshot(IsolatePtr iso,
                                  size_t snapshot_blob_index,
                                  int ref) {
  Locker locker(iso);
  Isolate::Scope isolate_scope(iso);
  HandleScope handle_scope(iso);

  // For function callbacks we need a reference to the context, but because of
  // the complexities of C -> Go function pointers, we store a reference to the
  // context as a simple integer identifier; this can then be used on the Go
  // side to lookup the context in the context registry. We use slot 1 as slot 0
  // has special meaning for the Chrome debugger.

  Local<Context> local_ctx =
      Context::FromSnapshot(iso, snapshot_blob_index).ToLocalChecked();
  local_ctx->SetEmbedderData(1, Integer::New(iso, ref));

  m_ctx* ctx = new m_ctx;
  ctx->ptr.Reset(iso, local_ctx);
  ctx->iso = iso;
  ctx->startup_data = nullptr;
  return ctx;
}

void ContextFree(ContextPtr ctx) {
  if (ctx == nullptr) {
    return;
  }
  ctx->ptr.Reset();

  for (m_value* val : ctx->vals) {
    val->ptr.Reset();
    delete val;
  }

  for (m_unboundScript* us : ctx->unboundScripts) {
    us->ptr.Reset();
    delete us;
  }

  if (ctx->startup_data) {
    delete ctx->startup_data;
  }

  delete ctx;
}

RtnValue RunScript(ContextPtr ctx, const char* source, const char* origin) {
  LOCAL_CONTEXT(ctx);

  RtnValue rtn = {};

  MaybeLocal<String> maybeSrc =
      String::NewFromUtf8(iso, source, NewStringType::kNormal);
  MaybeLocal<String> maybeOgn =
      String::NewFromUtf8(iso, origin, NewStringType::kNormal);
  Local<String> src, ogn;
  if (!maybeSrc.ToLocal(&src) || !maybeOgn.ToLocal(&ogn)) {
    rtn.error = ExceptionError(try_catch, iso, local_ctx);
    return rtn;
  }

  ScriptOrigin script_origin(ogn);
  Local<Script> script;
  if (!Script::Compile(local_ctx, src, &script_origin).ToLocal(&script)) {
    rtn.error = ExceptionError(try_catch, iso, local_ctx);
    return rtn;
  }
  Local<Value> result;
  if (!script->Run(local_ctx).ToLocal(&result)) {
    rtn.error = ExceptionError(try_catch, iso, local_ctx);
    return rtn;
  }
  m_value* val = new m_value;
  val->iso = iso;
  val->ctx = ctx;
  val->ptr = Persistent<Value, CopyablePersistentTraits<Value>>(iso, result);

  rtn.value = tracked_value(ctx, val);
  return rtn;
}

/********** UnboundScript & ScriptCompilerCachedData **********/

ScriptCompilerCachedData* UnboundScriptCreateCodeCache(
    IsolatePtr iso,
    UnboundScriptPtr us_ptr) {
  ISOLATE_SCOPE(iso);

  Local<UnboundScript> unbound_script = us_ptr->ptr.Get(iso);

  ScriptCompiler::CachedData* cached_data =
      ScriptCompiler::CreateCodeCache(unbound_script);

  ScriptCompilerCachedData* cd = new ScriptCompilerCachedData;
  cd->ptr = cached_data;
  cd->data = cached_data->data;
  cd->length = cached_data->length;
  cd->rejected = cached_data->rejected;
  return cd;
}

void ScriptCompilerCachedDataDelete(ScriptCompilerCachedData* cached_data) {
  delete cached_data->ptr;
  delete cached_data;
}

// This can only run in contexts that belong to the same isolate
// the script was compiled in
RtnValue UnboundScriptRun(ContextPtr ctx, UnboundScriptPtr us_ptr) {
  LOCAL_CONTEXT(ctx)

  RtnValue rtn = {};

  Local<UnboundScript> unbound_script = us_ptr->ptr.Get(iso);

  Local<Script> script = unbound_script->BindToCurrentContext();
  Local<Value> result;
  if (!script->Run(local_ctx).ToLocal(&result)) {
    rtn.error = ExceptionError(try_catch, iso, local_ctx);
    return rtn;
  }
  m_value* val = new m_value;
  val->iso = iso;
  val->ctx = ctx;
  val->ptr = Persistent<Value, CopyablePersistentTraits<Value>>(iso, result);

  rtn.value = tracked_value(ctx, val);
  return rtn;
}

RtnValue JSONParse(ContextPtr ctx, const char* str) {
  LOCAL_CONTEXT(ctx);
  RtnValue rtn = {};

  Local<String> v8Str;
  if (!String::NewFromUtf8(iso, str, NewStringType::kNormal).ToLocal(&v8Str)) {
    rtn.error = ExceptionError(try_catch, iso, local_ctx);
  }

  Local<Value> result;
  if (!JSON::Parse(local_ctx, v8Str).ToLocal(&result)) {
    rtn.error = ExceptionError(try_catch, iso, local_ctx);
    return rtn;
  }
  m_value* val = new m_value;
  val->iso = iso;
  val->ctx = ctx;
  val->ptr = Persistent<Value, CopyablePersistentTraits<Value>>(iso, result);

  rtn.value = tracked_value(ctx, val);
  return rtn;
}

const char* JSONStringify(ContextPtr ctx, ValuePtr val) {
  Isolate* iso;
  Local<Context> local_ctx;

  if (ctx != nullptr) {
    iso = ctx->iso;
  } else {
    iso = val->iso;
  }

  Locker locker(iso);
  Isolate::Scope isolate_scope(iso);
  HandleScope handle_scope(iso);

  if (ctx != nullptr) {
    local_ctx = ctx->ptr.Get(iso);
  } else {
    if (val->ctx != nullptr) {
      local_ctx = val->ctx->ptr.Get(iso);
    } else {
      m_ctx* ctx = isolateInternalContext(iso);
      local_ctx = ctx->ptr.Get(iso);
    }
  }

  Context::Scope context_scope(local_ctx);

  Local<String> str;
  if (!JSON::Stringify(local_ctx, val->ptr.Get(iso)).ToLocal(&str)) {
    return nullptr;
  }
  String::Utf8Value json(iso, str);
  return CopyString(json);
}

ValuePtr ContextGlobal(ContextPtr ctx) {
  LOCAL_CONTEXT(ctx);
  m_value* val = new m_value;

  val->iso = iso;
  val->ctx = ctx;
  val->ptr = Persistent<Value, CopyablePersistentTraits<Value>>(
      iso, local_ctx->Global());

  return tracked_value(ctx, val);
}

/********** Value **********/

#define LOCAL_VALUE(val)                   \
  Isolate* iso = val->iso;                 \
  Locker locker(iso);                      \
  Isolate::Scope isolate_scope(iso);       \
  HandleScope handle_scope(iso);           \
  TryCatch try_catch(iso);                 \
  m_ctx* ctx = val->ctx;                   \
  Local<Context> local_ctx;                \
  if (ctx != nullptr) {                    \
    local_ctx = ctx->ptr.Get(iso);         \
  } else {                                 \
    ctx = isolateInternalContext(iso);     \
    local_ctx = ctx->ptr.Get(iso);         \
  }                                        \
  Context::Scope context_scope(local_ctx); \
  Local<Value> value = val->ptr.Get(iso);

ValuePtr NewValueInteger(IsolatePtr iso, int32_t v) {
  ISOLATE_SCOPE_INTERNAL_CONTEXT(iso);
  m_value* val = new m_value;
  val->iso = iso;
  val->ctx = ctx;
  val->ptr = Persistent<Value, CopyablePersistentTraits<Value>>(
      iso, Integer::New(iso, v));
  return tracked_value(ctx, val);
}

ValuePtr NewValueIntegerFromUnsigned(IsolatePtr iso, uint32_t v) {
  ISOLATE_SCOPE_INTERNAL_CONTEXT(iso);
  m_value* val = new m_value;
  val->iso = iso;
  val->ctx = ctx;
  val->ptr = Persistent<Value, CopyablePersistentTraits<Value>>(
      iso, Integer::NewFromUnsigned(iso, v));
  return tracked_value(ctx, val);
}

RtnValue NewValueString(IsolatePtr iso, const char* v, int v_length) {
  ISOLATE_SCOPE_INTERNAL_CONTEXT(iso);
  TryCatch try_catch(iso);
  RtnValue rtn = {};
  Local<String> str;
  if (!String::NewFromUtf8(iso, v, NewStringType::kNormal, v_length)
           .ToLocal(&str)) {
    rtn.error = ExceptionError(try_catch, iso, ctx->ptr.Get(iso));
    return rtn;
  }
  m_value* val = new m_value;
  val->iso = iso;
  val->ctx = ctx;
  val->ptr = Persistent<Value, CopyablePersistentTraits<Value>>(iso, str);
  rtn.value = tracked_value(ctx, val);
  return rtn;
}

ValuePtr NewValueNull(IsolatePtr iso) {
  ISOLATE_SCOPE_INTERNAL_CONTEXT(iso);
  m_value* val = new m_value;
  val->iso = iso;
  val->ctx = ctx;
  val->ptr = Persistent<Value, CopyablePersistentTraits<Value>>(iso, Null(iso));
  return tracked_value(ctx, val);
}

ValuePtr NewValueUndefined(IsolatePtr iso) {
  ISOLATE_SCOPE_INTERNAL_CONTEXT(iso);
  m_value* val = new m_value;
  val->iso = iso;
  val->ctx = ctx;
  val->ptr =
      Persistent<Value, CopyablePersistentTraits<Value>>(iso, Undefined(iso));
  return tracked_value(ctx, val);
}

ValuePtr NewValueBoolean(IsolatePtr iso, int v) {
  ISOLATE_SCOPE_INTERNAL_CONTEXT(iso);
  m_value* val = new m_value;
  val->iso = iso;
  val->ctx = ctx;
  val->ptr = Persistent<Value, CopyablePersistentTraits<Value>>(
      iso, Boolean::New(iso, v));
  return tracked_value(ctx, val);
}

ValuePtr NewValueNumber(IsolatePtr iso, double v) {
  ISOLATE_SCOPE_INTERNAL_CONTEXT(iso);
  m_value* val = new m_value;
  val->iso = iso;
  val->ctx = ctx;
  val->ptr = Persistent<Value, CopyablePersistentTraits<Value>>(
      iso, Number::New(iso, v));
  return tracked_value(ctx, val);
}

ValuePtr NewValueBigInt(IsolatePtr iso, int64_t v) {
  ISOLATE_SCOPE_INTERNAL_CONTEXT(iso);
  m_value* val = new m_value;
  val->iso = iso;
  val->ctx = ctx;
  val->ptr = Persistent<Value, CopyablePersistentTraits<Value>>(
      iso, BigInt::New(iso, v));
  return tracked_value(ctx, val);
}

ValuePtr NewValueBigIntFromUnsigned(IsolatePtr iso, uint64_t v) {
  ISOLATE_SCOPE_INTERNAL_CONTEXT(iso);
  m_value* val = new m_value;
  val->iso = iso;
  val->ctx = ctx;
  val->ptr = Persistent<Value, CopyablePersistentTraits<Value>>(
      iso, BigInt::NewFromUnsigned(iso, v));
  return tracked_value(ctx, val);
}

RtnValue NewValueBigIntFromWords(IsolatePtr iso,
                                 int sign_bit,
                                 int word_count,
                                 const uint64_t* words) {
  ISOLATE_SCOPE_INTERNAL_CONTEXT(iso);
  TryCatch try_catch(iso);
  Local<Context> local_ctx = ctx->ptr.Get(iso);

  RtnValue rtn = {};
  Local<BigInt> bigint;
  if (!BigInt::NewFromWords(local_ctx, sign_bit, word_count, words)
           .ToLocal(&bigint)) {
    rtn.error = ExceptionError(try_catch, iso, local_ctx);
    return rtn;
  }
  m_value* val = new m_value;
  val->iso = iso;
  val->ctx = ctx;
  val->ptr = Persistent<Value, CopyablePersistentTraits<Value>>(iso, bigint);
  rtn.value = tracked_value(ctx, val);
  return rtn;
}

const uint32_t* ValueToArrayIndex(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  Local<Uint32> array_index;
  if (!value->ToArrayIndex(local_ctx).ToLocal(&array_index)) {
    return nullptr;
  }

  uint32_t* idx = (uint32_t*)malloc(sizeof(uint32_t));
  *idx = array_index->Value();
  return idx;
}

int ValueToBoolean(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->BooleanValue(iso);
}

int32_t ValueToInt32(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->Int32Value(local_ctx).ToChecked();
}

int64_t ValueToInteger(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IntegerValue(local_ctx).ToChecked();
}

double ValueToNumber(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->NumberValue(local_ctx).ToChecked();
}

RtnString ValueToDetailString(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  RtnString rtn = {0};
  Local<String> str;
  if (!value->ToDetailString(local_ctx).ToLocal(&str)) {
    rtn.error = ExceptionError(try_catch, iso, local_ctx);
    return rtn;
  }
  String::Utf8Value ds(iso, str);
  rtn.data = CopyString(ds);
  rtn.length = ds.length();
  return rtn;
}

RtnString ValueToString(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  RtnString rtn = {0};
  // String::Utf8Value will result in an empty string if conversion to a string
  // fails
  // TODO: Consider propagating the JS error. A fallback value could be returned
  // in Value.String()
  String::Utf8Value src(iso, value);
  char* data = static_cast<char*>(malloc(src.length()));
  memcpy(data, *src, src.length());
  rtn.data = data;
  rtn.length = src.length();
  return rtn;
}

uint32_t ValueToUint32(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->Uint32Value(local_ctx).ToChecked();
}

ValueBigInt ValueToBigInt(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  Local<BigInt> bint;
  if (!value->ToBigInt(local_ctx).ToLocal(&bint)) {
    return {nullptr, 0};
  }

  int word_count = bint->WordCount();
  int sign_bit = 0;
  uint64_t* words = (uint64_t*)malloc(sizeof(uint64_t) * word_count);
  bint->ToWordsArray(&sign_bit, &word_count, words);
  ValueBigInt rtn = {words, word_count, sign_bit};
  return rtn;
}

RtnValue ValueToObject(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  RtnValue rtn = {};
  Local<Object> obj;
  if (!value->ToObject(local_ctx).ToLocal(&obj)) {
    rtn.error = ExceptionError(try_catch, iso, local_ctx);
    return rtn;
  }
  m_value* new_val = new m_value;
  new_val->iso = iso;
  new_val->ctx = ctx;
  new_val->ptr = Persistent<Value, CopyablePersistentTraits<Value>>(iso, obj);
  rtn.value = tracked_value(ctx, new_val);
  return rtn;
}

int ValueSameValue(ValuePtr val1, ValuePtr val2) {
  Isolate* iso = val1->iso;
  ISOLATE_SCOPE(iso);
  Local<Value> value1 = val1->ptr.Get(iso);
  Local<Value> value2 = val2->ptr.Get(iso);

  return value1->SameValue(value2);
}

int ValueIsUndefined(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsUndefined();
}

int ValueIsNull(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsNull();
}

int ValueIsNullOrUndefined(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsNullOrUndefined();
}

int ValueIsTrue(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsTrue();
}

int ValueIsFalse(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsFalse();
}

int ValueIsName(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsName();
}

int ValueIsString(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsString();
}

int ValueIsSymbol(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsSymbol();
}

int ValueIsFunction(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsFunction();
}

int ValueIsObject(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsObject();
}

int ValueIsBigInt(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsBigInt();
}

int ValueIsBoolean(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsBoolean();
}

int ValueIsNumber(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsNumber();
}

int ValueIsExternal(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsExternal();
}

int ValueIsInt32(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsInt32();
}

int ValueIsUint32(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsUint32();
}

int ValueIsDate(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsDate();
}

int ValueIsArgumentsObject(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsArgumentsObject();
}

int ValueIsBigIntObject(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsBigIntObject();
}

int ValueIsNumberObject(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsNumberObject();
}

int ValueIsStringObject(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsStringObject();
}

int ValueIsSymbolObject(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsSymbolObject();
}

int ValueIsNativeError(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsNativeError();
}

int ValueIsRegExp(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsRegExp();
}

int ValueIsAsyncFunction(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsAsyncFunction();
}

int ValueIsGeneratorFunction(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsGeneratorFunction();
}

int ValueIsGeneratorObject(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsGeneratorObject();
}

int ValueIsPromise(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsPromise();
}

int ValueIsMap(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsMap();
}

int ValueIsSet(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsSet();
}

int ValueIsMapIterator(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsMapIterator();
}

int ValueIsSetIterator(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsSetIterator();
}

int ValueIsWeakMap(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsWeakMap();
}

int ValueIsWeakSet(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsWeakSet();
}

int ValueIsArray(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsArray();
}

int ValueIsArrayBuffer(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsArrayBuffer();
}

int ValueIsArrayBufferView(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsArrayBufferView();
}

int ValueIsTypedArray(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsTypedArray();
}

int ValueIsUint8Array(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsUint8Array();
}

int ValueIsUint8ClampedArray(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsUint8ClampedArray();
}

int ValueIsInt8Array(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsInt8Array();
}

int ValueIsUint16Array(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsUint16Array();
}

int ValueIsInt16Array(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsInt16Array();
}

int ValueIsUint32Array(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsUint32Array();
}

int ValueIsInt32Array(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsInt32Array();
}

int ValueIsFloat32Array(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsFloat32Array();
}

int ValueIsFloat64Array(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsFloat64Array();
}

int ValueIsBigInt64Array(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsBigInt64Array();
}

int ValueIsBigUint64Array(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsBigUint64Array();
}

int ValueIsDataView(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsDataView();
}

int ValueIsSharedArrayBuffer(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsSharedArrayBuffer();
}

int ValueIsProxy(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsProxy();
}

int ValueIsWasmModuleObject(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsWasmModuleObject();
}

int ValueIsModuleNamespaceObject(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  return value->IsModuleNamespaceObject();
}

/********** Object **********/

#define LOCAL_OBJECT(ptr) \
  LOCAL_VALUE(ptr)        \
  Local<Object> obj = value.As<Object>()

void ObjectSet(ValuePtr ptr, const char* key, ValuePtr prop_val) {
  LOCAL_OBJECT(ptr);
  Local<String> key_val =
      String::NewFromUtf8(iso, key, NewStringType::kNormal).ToLocalChecked();
  obj->Set(local_ctx, key_val, prop_val->ptr.Get(iso)).Check();
}

void ObjectSetIdx(ValuePtr ptr, uint32_t idx, ValuePtr prop_val) {
  LOCAL_OBJECT(ptr);
  obj->Set(local_ctx, idx, prop_val->ptr.Get(iso)).Check();
}

int ObjectSetInternalField(ValuePtr ptr, int idx, ValuePtr val_ptr) {
  LOCAL_OBJECT(ptr);
  m_value* prop_val = static_cast<m_value*>(val_ptr);

  if (idx >= obj->InternalFieldCount()) {
    return 0;
  }

  obj->SetInternalField(idx, prop_val->ptr.Get(iso));

  return 1;
}

int ObjectInternalFieldCount(ValuePtr ptr) {
  LOCAL_OBJECT(ptr);
  return obj->InternalFieldCount();
}

RtnValue ObjectGet(ValuePtr ptr, const char* key) {
  LOCAL_OBJECT(ptr);
  RtnValue rtn = {};

  Local<String> key_val;
  if (!String::NewFromUtf8(iso, key, NewStringType::kNormal)
           .ToLocal(&key_val)) {
    rtn.error = ExceptionError(try_catch, iso, local_ctx);
    return rtn;
  }
  Local<Value> result;
  if (!obj->Get(local_ctx, key_val).ToLocal(&result)) {
    rtn.error = ExceptionError(try_catch, iso, local_ctx);
    return rtn;
  }
  m_value* new_val = new m_value;
  new_val->iso = iso;
  new_val->ctx = ctx;
  new_val->ptr =
      Persistent<Value, CopyablePersistentTraits<Value>>(iso, result);

  rtn.value = tracked_value(ctx, new_val);
  return rtn;
}

ValuePtr ObjectGetInternalField(ValuePtr ptr, int idx) {
  LOCAL_OBJECT(ptr);

  if (idx >= obj->InternalFieldCount()) {
    return nullptr;
  }

  Local<Value> result = obj->GetInternalField(idx);

  m_value* new_val = new m_value;
  new_val->iso = iso;
  new_val->ctx = ctx;
  new_val->ptr =
      Persistent<Value, CopyablePersistentTraits<Value>>(iso, result);

  return tracked_value(ctx, new_val);
}

RtnValue ObjectGetIdx(ValuePtr ptr, uint32_t idx) {
  LOCAL_OBJECT(ptr);
  RtnValue rtn = {};

  Local<Value> result;
  if (!obj->Get(local_ctx, idx).ToLocal(&result)) {
    rtn.error = ExceptionError(try_catch, iso, local_ctx);
    return rtn;
  }
  m_value* new_val = new m_value;
  new_val->iso = iso;
  new_val->ctx = ctx;
  new_val->ptr =
      Persistent<Value, CopyablePersistentTraits<Value>>(iso, result);

  rtn.value = tracked_value(ctx, new_val);
  return rtn;
}

int ObjectHas(ValuePtr ptr, const char* key) {
  LOCAL_OBJECT(ptr);
  Local<String> key_val =
      String::NewFromUtf8(iso, key, NewStringType::kNormal).ToLocalChecked();
  return obj->Has(local_ctx, key_val).ToChecked();
}

int ObjectHasIdx(ValuePtr ptr, uint32_t idx) {
  LOCAL_OBJECT(ptr);
  return obj->Has(local_ctx, idx).ToChecked();
}

int ObjectDelete(ValuePtr ptr, const char* key) {
  LOCAL_OBJECT(ptr);
  Local<String> key_val =
      String::NewFromUtf8(iso, key, NewStringType::kNormal).ToLocalChecked();
  return obj->Delete(local_ctx, key_val).ToChecked();
}

int ObjectDeleteIdx(ValuePtr ptr, uint32_t idx) {
  LOCAL_OBJECT(ptr);
  return obj->Delete(local_ctx, idx).ToChecked();
}

/********** Promise **********/

RtnValue NewPromiseResolver(ContextPtr ctx) {
  LOCAL_CONTEXT(ctx);
  RtnValue rtn = {};
  Local<Promise::Resolver> resolver;
  if (!Promise::Resolver::New(local_ctx).ToLocal(&resolver)) {
    rtn.error = ExceptionError(try_catch, iso, local_ctx);
    return rtn;
  }
  m_value* val = new m_value;
  val->iso = iso;
  val->ctx = ctx;
  val->ptr = Persistent<Value, CopyablePersistentTraits<Value>>(iso, resolver);
  rtn.value = tracked_value(ctx, val);
  return rtn;
}

ValuePtr PromiseResolverGetPromise(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  Local<Promise::Resolver> resolver = value.As<Promise::Resolver>();
  Local<Promise> promise = resolver->GetPromise();
  m_value* promise_val = new m_value;
  promise_val->iso = iso;
  promise_val->ctx = ctx;
  promise_val->ptr =
      Persistent<Value, CopyablePersistentTraits<Value>>(iso, promise);
  return tracked_value(ctx, promise_val);
}

int PromiseResolverResolve(ValuePtr ptr, ValuePtr resolve_val) {
  LOCAL_VALUE(ptr);
  Local<Promise::Resolver> resolver = value.As<Promise::Resolver>();
  return resolver->Resolve(local_ctx, resolve_val->ptr.Get(iso)).ToChecked();
}

int PromiseResolverReject(ValuePtr ptr, ValuePtr reject_val) {
  LOCAL_VALUE(ptr);
  Local<Promise::Resolver> resolver = value.As<Promise::Resolver>();
  return resolver->Reject(local_ctx, reject_val->ptr.Get(iso)).ToChecked();
}

int PromiseState(ValuePtr ptr) {
  LOCAL_VALUE(ptr)
  Local<Promise> promise = value.As<Promise>();
  return promise->State();
}

RtnValue PromiseThen(ValuePtr ptr, int callback_ref) {
  LOCAL_VALUE(ptr)
  RtnValue rtn = {};
  Local<Promise> promise = value.As<Promise>();
  Local<Integer> cbData = Integer::New(iso, callback_ref);
  Local<Function> func;
  if (!Function::New(local_ctx, FunctionTemplateCallback, cbData)
           .ToLocal(&func)) {
    rtn.error = ExceptionError(try_catch, iso, local_ctx);
    return rtn;
  }
  Local<Promise> result;
  if (!promise->Then(local_ctx, func).ToLocal(&result)) {
    rtn.error = ExceptionError(try_catch, iso, local_ctx);
    return rtn;
  }
  m_value* result_val = new m_value;
  result_val->iso = iso;
  result_val->ctx = ctx;
  result_val->ptr =
      Persistent<Value, CopyablePersistentTraits<Value>>(iso, result);
  rtn.value = tracked_value(ctx, result_val);
  return rtn;
}

RtnValue PromiseThen2(ValuePtr ptr, int on_fulfilled_ref, int on_rejected_ref) {
  LOCAL_VALUE(ptr)
  RtnValue rtn = {};
  Local<Promise> promise = value.As<Promise>();
  Local<Integer> onFulfilledData = Integer::New(iso, on_fulfilled_ref);
  Local<Function> onFulfilledFunc;
  if (!Function::New(local_ctx, FunctionTemplateCallback, onFulfilledData)
           .ToLocal(&onFulfilledFunc)) {
    rtn.error = ExceptionError(try_catch, iso, local_ctx);
    return rtn;
  }
  Local<Integer> onRejectedData = Integer::New(iso, on_rejected_ref);
  Local<Function> onRejectedFunc;
  if (!Function::New(local_ctx, FunctionTemplateCallback, onRejectedData)
           .ToLocal(&onRejectedFunc)) {
    rtn.error = ExceptionError(try_catch, iso, local_ctx);
    return rtn;
  }
  Local<Promise> result;
  if (!promise->Then(local_ctx, onFulfilledFunc, onRejectedFunc)
           .ToLocal(&result)) {
    rtn.error = ExceptionError(try_catch, iso, local_ctx);
    return rtn;
  }
  m_value* result_val = new m_value;
  result_val->iso = iso;
  result_val->ctx = ctx;
  result_val->ptr =
      Persistent<Value, CopyablePersistentTraits<Value>>(iso, result);
  rtn.value = tracked_value(ctx, result_val);
  return rtn;
}

RtnValue PromiseCatch(ValuePtr ptr, int callback_ref) {
  LOCAL_VALUE(ptr)
  RtnValue rtn = {};
  Local<Promise> promise = value.As<Promise>();
  Local<Integer> cbData = Integer::New(iso, callback_ref);
  Local<Function> func;
  if (!Function::New(local_ctx, FunctionTemplateCallback, cbData)
           .ToLocal(&func)) {
    rtn.error = ExceptionError(try_catch, iso, local_ctx);
    return rtn;
  }
  Local<Promise> result;
  if (!promise->Catch(local_ctx, func).ToLocal(&result)) {
    rtn.error = ExceptionError(try_catch, iso, local_ctx);
    return rtn;
  }
  m_value* result_val = new m_value;
  result_val->iso = iso;
  result_val->ctx = ctx;
  result_val->ptr =
      Persistent<Value, CopyablePersistentTraits<Value>>(iso, result);
  rtn.value = tracked_value(ctx, result_val);
  return rtn;
}

ValuePtr PromiseResult(ValuePtr ptr) {
  LOCAL_VALUE(ptr)
  Local<Promise> promise = value.As<Promise>();
  Local<Value> result = promise->Result();
  m_value* result_val = new m_value;
  result_val->iso = iso;
  result_val->ctx = ctx;
  result_val->ptr =
      Persistent<Value, CopyablePersistentTraits<Value>>(iso, result);
  return tracked_value(ctx, result_val);
}

/********** Function **********/

static void buildCallArguments(Isolate* iso,
                               Local<Value>* argv,
                               int argc,
                               ValuePtr args[]) {
  for (int i = 0; i < argc; i++) {
    argv[i] = args[i]->ptr.Get(iso);
  }
}

RtnValue FunctionCall(ValuePtr ptr, ValuePtr recv, int argc, ValuePtr args[]) {
  LOCAL_VALUE(ptr)

  RtnValue rtn = {};
  Local<Function> fn = Local<Function>::Cast(value);
  Local<Value> argv[argc];
  buildCallArguments(iso, argv, argc, args);

  Local<Value> local_recv = recv->ptr.Get(iso);

  Local<Value> result;
  if (!fn->Call(local_ctx, local_recv, argc, argv).ToLocal(&result)) {
    rtn.error = ExceptionError(try_catch, iso, local_ctx);
    return rtn;
  }
  m_value* rtnval = new m_value;
  rtnval->iso = iso;
  rtnval->ctx = ctx;
  rtnval->ptr = Persistent<Value, CopyablePersistentTraits<Value>>(iso, result);
  rtn.value = tracked_value(ctx, rtnval);
  return rtn;
}

RtnValue FunctionNewInstance(ValuePtr ptr, int argc, ValuePtr args[]) {
  LOCAL_VALUE(ptr)
  RtnValue rtn = {};
  Local<Function> fn = Local<Function>::Cast(value);
  Local<Value> argv[argc];
  buildCallArguments(iso, argv, argc, args);
  Local<Object> result;
  if (!fn->NewInstance(local_ctx, argc, argv).ToLocal(&result)) {
    rtn.error = ExceptionError(try_catch, iso, local_ctx);
    return rtn;
  }
  m_value* rtnval = new m_value;
  rtnval->iso = iso;
  rtnval->ctx = ctx;
  rtnval->ptr = Persistent<Value, CopyablePersistentTraits<Value>>(iso, result);
  rtn.value = tracked_value(ctx, rtnval);
  return rtn;
}

ValuePtr FunctionSourceMapUrl(ValuePtr ptr) {
  LOCAL_VALUE(ptr)
  Local<Function> fn = Local<Function>::Cast(value);
  Local<Value> result = fn->GetScriptOrigin().SourceMapUrl();
  m_value* rtnval = new m_value;
  rtnval->iso = iso;
  rtnval->ctx = ctx;
  rtnval->ptr = Persistent<Value, CopyablePersistentTraits<Value>>(iso, result);
  return tracked_value(ctx, rtnval);
}

/********** v8::V8 **********/

const char* Version() {
  return V8::GetVersion();
}

void SetFlags(const char* flags) {
  V8::SetFlagsFromString(flags);
}
}
