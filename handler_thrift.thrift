# thrift IDL for thrift handler
# A example of writing network rpc handler in thrift.

namespace go logging

// tailored logging fields to report.
struct ThriftLogRecord {
    1: string name,
    2: i32 level,
    3: string path_name,
    4: string file_name,
    5: i32 line_no,
    6: string func_name,
    7: string message,
}

service ThriftLoggingService {
    oneway void report(1: ThriftLogRecord record),
}
