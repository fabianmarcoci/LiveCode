#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]
mod auth;
mod models;
mod storage;

use auth::{check_email_available, check_username_available, login_user, register_user};
use live_code_lib::adjust_window_size;
use storage::{clear_tokens, get_access_token, get_refresh_token, save_tokens};

fn main() {
    tauri::Builder::default()
        .plugin(tauri_plugin_store::Builder::new().build())
        .setup(|app| {
            adjust_window_size(app);
            Ok(())
        })
        .invoke_handler(tauri::generate_handler![
            register_user,
            login_user,
            check_email_available,
            check_username_available,
            save_tokens,
            get_access_token,
            get_refresh_token,
            clear_tokens
        ])
        .run(tauri::generate_context!())
        .expect("App start error.");
}
