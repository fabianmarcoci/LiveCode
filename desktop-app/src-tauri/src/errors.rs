use crate::config::ApiConfig;
use crate::models::ClientErrorLog;
use reqwest;
use std::collections::HashMap;
use std::sync::Mutex;
use std::time::{SystemTime, UNIX_EPOCH};
use tauri_plugin_store::Error as StoreError;

lazy_static::lazy_static! {
    static ref LAST_ERROR_SENT: Mutex<HashMap<String, u64>> = Mutex::new(HashMap::new());
}

const ERROR_COOLDOWN_SECS: u64 = 60;

fn should_send_error(error_type: &str) -> bool {
    let now = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .unwrap()
        .as_secs();

    let mut last_sent = LAST_ERROR_SENT.lock().unwrap();

    if let Some(&last_time) = last_sent.get(error_type) {
        if now - last_time < ERROR_COOLDOWN_SECS {
            return false;
        }
    }

    last_sent.insert(error_type.to_string(), now);
    true
}

async fn send_client_error(error_type: &str, error_message: &str) {
    if !should_send_error(error_type) {
        return;
    }

    let config = ApiConfig::new();
    let client = reqwest::Client::new();

    let payload = ClientErrorLog {
        timestamp: chrono::Utc::now().to_rfc3339(),
        error_type: error_type.to_string(),
        error_message: error_message.to_string(),
        app_version: env!("CARGO_PKG_VERSION").to_string(),
        os: std::env::consts::OS.to_string(),
    };

    let _ = client
        .post(format!("{}/monitoring/client-errors", config.base_url))
        .json(&payload)
        .send()
        .await;
}

pub fn network_error_to_string(error: reqwest::Error) -> String {
    if error.is_decode() {
        let error_msg = format!("JSON decode failed: {}", error);
        tokio::spawn(async move {
            send_client_error("json_decode_network", &error_msg).await;
        });
        return "Received invalid response from server. Please try again.".to_string();
    }

    if error.is_timeout() {
        return "Connection timed out. Please try again.".to_string();
    }

    if error.is_connect() {
        return "Could not connect to server. Please check your internet connection.".to_string();
    }

    if error.is_body() {
        return "Received invalid response from server. Please try again.".to_string();
    }

    "Network error occurred. Please try again.".to_string()
}

pub fn storage_error_to_string(error: StoreError) -> String {
    let error_msg = format!("Storage error: {:?}", error);
    tokio::spawn(async move {
        send_client_error("storage_error", &error_msg).await;
    });

    "Could not access local storage. Please check your permissions.".to_string()
}
