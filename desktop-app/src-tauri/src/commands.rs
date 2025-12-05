use crate::models::{RegisterRequest, RegisterResponse};
use serde::Deserialize;

#[tauri::command]
pub async fn register_user(payload: RegisterRequest) -> Result<RegisterResponse, String> {
    let client = reqwest::Client::new();

    let response = client
        .post("http://localhost:3000/api/auth/register")
        .json(&payload)
        .send()
        .await
        .map_err(|e| e.to_string())?;

    let data: RegisterResponse = response.json().await.map_err(|e| e.to_string())?;

    Ok(data)
}

#[tauri::command]
pub async fn check_email_available(email: String) -> Result<Option<bool>, String> {
    let client = reqwest::Client::new();

    let response = client
        .get("http://localhost:3000/api/auth/check-field")
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
    let client = reqwest::Client::new();

    let response = client
        .get("http://localhost:3000/api/auth/check-field")
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
