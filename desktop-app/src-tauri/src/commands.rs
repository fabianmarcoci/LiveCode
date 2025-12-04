use tauri::State;
use sqlx::PgPool;
use live_code_lib::auth::register::models::{RegisterRequest, RegisterResponse};

#[tauri::command]
pub async fn register_user(
    payload: RegisterRequest,
    pool: State<'_, PgPool>,
) -> Result<RegisterResponse, String> {
    live_code_lib::auth::register::handlers::register_user_internal(payload, pool.inner()).await
}

#[tauri::command]
pub async fn check_email_available(
    email: String,
    pool: State<'_, PgPool>,
) -> Result<Option<bool>, ()> {
    Ok(live_code_lib::auth::register::handlers::check_field_available_internal(
        "email",
        email,
        pool.inner()
    ).await)
}

#[tauri::command]
pub async fn check_username_available(
    username: String,
    pool: State<'_, PgPool>,
) -> Result<Option<bool>, ()> {
    Ok(live_code_lib::auth::register::handlers::check_field_available_internal(
        "username",
        username,
        pool.inner()
    ).await)
}