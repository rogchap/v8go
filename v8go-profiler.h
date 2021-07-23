#ifndef V8GO_PROFILER_H
#define V8GO_PROFILER_H

#include <stddef.h>
#include <stdint.h>
#include <v8-inspector.h>
#include <v8.h>

class InspectorFrontend final : public v8_inspector::V8Inspector::Channel {
 public:
  InspectorFrontend(v8::Local<v8::Context> context);
  void sendResponse(int callId, std::unique_ptr<v8_inspector::StringBuffer> buffer) override;
  void sendNotification(std::unique_ptr<v8_inspector::StringBuffer> buffer) override;
  void flushProtocolNotifications() override {};

  void clearResponse();
  const char* cStringResponse(int *length);

 private:

  std::unique_ptr<v8_inspector::StringBuffer> response_;
  v8::Isolate* isolate_;
  v8::Global<v8::Context> context_;
};

class InspectorClient : public v8_inspector::V8InspectorClient {
};

class Profiler final {
 public:
  Profiler(v8::Local<v8::Context> context);

  void start();
  const char *stop(int *length);

  v8::Isolate* isolate;

 private:
  void send_message(std::string msg);

  std::unique_ptr<v8_inspector::V8Inspector> inspector_;
  std::unique_ptr<v8_inspector::V8InspectorSession> session_;
  std::unique_ptr<InspectorFrontend> frontend_;
};

#endif  // V8GO_PROFILER_H
