#include <string>
#include <iostream>

#include "v8go-internal.h"
#include "v8go-profiler.h"

InspectorFrontend::InspectorFrontend(v8::Local<v8::Context> context) {
  isolate_ = context->GetIsolate();
  context_.Reset(isolate_, context);
}

void InspectorFrontend::sendResponse(int callId, std::unique_ptr<v8_inspector::StringBuffer> buffer) {
  response_ = std::move(buffer);
}

void InspectorFrontend::sendNotification(std::unique_ptr<v8_inspector::StringBuffer> buffer) {
  response_ = std::move(buffer);
}

void InspectorFrontend::clearResponse() {
  response_ = nullptr;
}

const char* InspectorFrontend::cStringResponse(int *length_ptr) {
  v8_inspector::StringView view(response_->string());
  v8::HandleScope handle_scope(isolate_);
  v8::Local<v8::String> encoded =
      (view.is8Bit()
           ? v8::String::NewFromOneByte(
                 isolate_,
                 reinterpret_cast<const uint8_t*>(view.characters8()),
                 v8::NewStringType::kNormal, view.length())
           : v8::String::NewFromTwoByte(
                 isolate_,
                 reinterpret_cast<const uint16_t*>(view.characters16()),
                 v8::NewStringType::kNormal, view.length()))
          .ToLocalChecked();
  v8::String::Utf8Value utf8(isolate_, encoded);
  int length = utf8.length();
  char *response_string = (char *)malloc(length);
  if (!response_string)  {
    *length_ptr = 0;
  } else {
    *length_ptr = length;
    memcpy((void*)response_string, *utf8, length);
  }

  return response_string;
}

Profiler::Profiler(v8::Local<v8::Context> context) {
  isolate = context->GetIsolate();
  inspector_ = v8_inspector::V8Inspector::create(isolate, new InspectorClient());
  frontend_ = std::unique_ptr<InspectorFrontend>(new InspectorFrontend(context));
  session_ = inspector_->connect(1, frontend_.get(), v8_inspector::StringView());
  inspector_->contextCreated(v8_inspector::V8ContextInfo(context, 1, v8_inspector::StringView()));
}

void Profiler::start() {
  send_message("{\"id\":0,\"method\":\"Profiler.enable\"}");
  send_message("{\"id\":1,\"method\":\"Profiler.start\"}");
}

const char *Profiler::stop(int *length) {
  send_message("{\"id\":2,\"method\":\"Profiler.stop\"}");
  return frontend_->cStringResponse(length);
}

void Profiler::send_message(std::string msg) {
  v8_inspector::StringView message_view(reinterpret_cast<const uint8_t *>(msg.c_str()), msg.length());
  frontend_->clearResponse();
  session_->dispatchProtocolMessage(message_view);
}
