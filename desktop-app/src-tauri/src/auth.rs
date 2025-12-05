use crate::models::{LoginRequest, LoginResponse, RegisterRequest, RegisterResponse};
use live_code_lib::config::ApiConfig;
use serde::Deserialize;

#[tauri::command]
pub async fn register_user(payload: RegisterRequest) -> Result<RegisterResponse, String> {
    let config = ApiConfig::new();
    let client = reqwest::Client::new();

    let response = client
        .post(config.auth_url("register"))
        .json(&payload)
        .send()
        .await
        .map_err(|e| e.to_string())?;

    if !response.status().is_success() {
        if let Ok(error_response) = response.json::<RegisterResponse>().await {
            return Ok(error_response);
        }
        return Err("An unexpected error occurred. Please try again.".to_string());
    }

    let data: RegisterResponse = response.json().await.map_err(|e| e.to_string())?;

    Ok(data)
}

#[tauri::command]
pub async fn check_email_available(email: String) -> Result<Option<bool>, String> {
    let config = ApiConfig::new();
    let client = reqwest::Client::new();

    let response = client
        .get(config.auth_url("check-field"))
        .query(&[("field", "email"), ("value", &email)])
        .send()
        .await
        .map_err(|e| e.to_string())?;

    #[derive(Deserialize)]
    struct CheckResponse {
        available: Option<bool>,
    }

    let data: CheckResponse = response.json().await.map_err(|e| e.to_string())?;

    Ok(data.available)
}

#[tauri::command]
pub async fn check_username_available(username: String) -> Result<Option<bool>, String> {
    let config = ApiConfig::new();
    let client = reqwest::Client::new();

    let response = client
        .get(config.auth_url("check-field"))
        .query(&[("field", "username"), ("value", &username)])
        .send()
        .await
        .map_err(|e| e.to_string())?;

    #[derive(Deserialize)]
    struct CheckResponse {
        available: Option<bool>,
    }

    let data: CheckResponse = response.json().await.map_err(|e| e.to_string())?;

    Ok(data.available)
}

#[tauri::command]
pub async fn login_user(payload: LoginRequest) -> Result<LoginResponse, String> {
    let config = ApiConfig::new();
    let client = reqwest::Client::new();

    let response = client
        .post(config.auth_url("login"))
        .json(&payload)
        .send()
        .await
        .map_err(|e| e.to_string())?;

    if !response.status().is_success() {
        if let Ok(error_response) = response.json::<LoginResponse>().await {
            return Ok(error_response);
        }
        return Err("An unexpected error occurred. Please try again.".to_string());
    }

    let data: LoginResponse = response.json().await.map_err(|e| e.to_string())?;

    Ok(data)
}
