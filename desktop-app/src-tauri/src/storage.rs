use tauri::AppHandle;
use tauri_plugin_store::StoreExt;

const STORE_FILE: &str = "livecode_secure.json";
const ACCESS_TOKEN_KEY: &str = "access_token";
const REFRESH_TOKEN_KEY: &str = "refresh_token";

#[tauri::command]
pub async fn save_tokens(
    app: AppHandle,
    access_token: String,
    refresh_token: String,
) -> Result<(), String> {
    let store = app.store(STORE_FILE).map_err(|e| e.to_string())?;

    store.set(ACCESS_TOKEN_KEY, serde_json::json!(access_token));
    store.set(REFRESH_TOKEN_KEY, serde_json::json!(refresh_token));

    store.save().map_err(|e| e.to_string())?;

    Ok(())
}

#[tauri::command]
pub async fn get_access_token(app: AppHandle) -> Result<Option<String>, String> {
    let store = app.store(STORE_FILE).map_err(|e| e.to_string())?;

    match store.get(ACCESS_TOKEN_KEY) {
        Some(value) => {
            let token = value.as_str().map(|s| s.to_string());
            Ok(token)
        }
        None => Ok(None),
    }
}

#[tauri::command]
pub async fn get_refresh_token(app: AppHandle) -> Result<Option<String>, String> {
    let store = app.store(STORE_FILE).map_err(|e| e.to_string())?;

    match store.get(REFRESH_TOKEN_KEY) {
        Some(value) => {
            let token = value.as_str().map(|s| s.to_string());
            Ok(token)
        }
        None => Ok(None),
    }
}

#[tauri::command]
pub async fn clear_tokens(app: AppHandle) -> Result<(), String> {
    let store = app.store(STORE_FILE).map_err(|e| e.to_string())?;

    store.delete(ACCESS_TOKEN_KEY);
    store.delete(REFRESH_TOKEN_KEY);

    store.save().map_err(|e| e.to_string())?;

    Ok(())
}
