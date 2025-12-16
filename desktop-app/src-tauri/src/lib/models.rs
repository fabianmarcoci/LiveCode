use serde::Serialize;

#[derive(Debug, Clone, Serialize)]
pub struct ClientErrorLog {
    pub timestamp: String,
    pub error_type: String,
    pub error_message: String,
    pub app_version: String,
    pub os: String,
}
