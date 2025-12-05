#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]
mod commands;
mod models;

use commands::{check_email_available, check_username_available, register_user};
use live_code_lib::adjust_window_size;

fn main() {
    tauri::Builder::default()
        .setup(|app| {
            adjust_window_size(app);
            Ok(())
        })
        .invoke_handler(tauri::generate_handler![
            register_user,
            check_email_available,
            check_username_available
        ])
        .run(tauri::generate_context!())
        .expect("App start error.");
}
